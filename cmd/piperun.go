package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"defenders-cli/internal/utils"
)

const piperunHelp = `release - Azure DevOps Pipeline Runner

USAGE:
  defenders release <subcommand> [flags]

SUBCOMMANDS:
  run              Run a pipeline directly
  monitor-trigger  Monitor a pipeline and trigger another when it completes

FLAGS:
  -t, --token      Personal Access Token (overrides config/env)
  -i, --interval   Check interval in seconds (default: 30, for monitor-trigger)
  -h, --help       Show this help message

EXAMPLES:
  defenders release run <pipeline-definition-url>
  defenders release run <pipeline-definition-url> -t <token>
  defenders release monitor-trigger <wait-for-build-url> <trigger-pipeline-url>
  defenders release monitor-trigger <wait-url> <trigger-url> --interval 60

URL FORMATS:
  <wait-for-build-url>:      https://dev.azure.com/org/proj/_build/results?buildId=123
  <pipeline-definition-url>: https://dev.azure.com/org/proj/_build?definitionId=456

AUTHENTICATION:
  PAT with 'Build (Read & Execute)' permissions required.
  Create at: https://msazure.visualstudio.com/_usersSettings/tokens
`

const piperunRunHelp = `release run - Run an Azure DevOps pipeline

USAGE:
  defenders release run <pipeline-url> [-t <token>]

ARGUMENTS:
  <pipeline-url>  URL of the pipeline definition to run

FLAGS:
  -t, --token     Personal Access Token (overrides config/env)

EXAMPLE:
  defenders release run https://dev.azure.com/org/project/_build?definitionId=456
`

const piperunMonitorHelp = `release monitor-trigger - Monitor a pipeline and trigger another when complete

USAGE:
  defenders release monitor-trigger <wait-for-url> <trigger-url> [flags]

ARGUMENTS:
  <wait-for-url>  URL of the pipeline run to wait for (buildId URL)
  <trigger-url>   URL of the pipeline definition to trigger

FLAGS:
  -t, --token      Personal Access Token (overrides config/env)
  -i, --interval   Check interval in seconds (default: 30)

EXAMPLE:
  defenders release monitor-trigger \
    https://dev.azure.com/org/proj/_build/results?buildId=123 \
    https://dev.azure.com/org/proj/_build?definitionId=456 \
    --interval 60
`

type PiperunCmd struct {
	Subcommand  string
	PipelineURL string
	WaitForURL  string
	TriggerURL  string
	PAT         string
	Interval    int
}

// parseADOUrl parses Azure DevOps URL and extracts org, project, and query params
// Supports URLs like:
// - https://dev.azure.com/{org}/{project}/_build?definitionId=123
// - https://dev.azure.com/{org}/{project}/_build/results?buildId=123
// - https://{org}.visualstudio.com/{project}/_build?definitionId=123
func parseADOUrl(rawURL string) (orgURL string, project string, queryParams url.Values, err error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", "", nil, err
	}

	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")

	// Check if it's dev.azure.com format (org is in path)
	if strings.Contains(parsed.Host, "dev.azure.com") {
		// Format: https://dev.azure.com/{org}/{project}/...
		if len(pathParts) < 2 {
			return "", "", nil, fmt.Errorf("invalid ADO URL - expected org and project in path")
		}
		org := pathParts[0]
		project = pathParts[1]
		orgURL = fmt.Sprintf("https://dev.azure.com/%s", org)
	} else if strings.Contains(parsed.Host, "visualstudio.com") {
		// Format: https://{org}.visualstudio.com/{project}/...
		// Org is part of hostname
		orgURL = fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
		if len(pathParts) < 1 {
			return "", "", nil, fmt.Errorf("invalid ADO URL - expected project in path")
		}
		project = pathParts[0]
	} else {
		return "", "", nil, fmt.Errorf("unrecognized ADO URL format")
	}

	queryParams = parsed.Query()
	return orgURL, project, queryParams, nil
}

func (p *PiperunCmd) Run() {
	switch p.Subcommand {
	case "run":
		p.runPipeline()
	case "monitor-trigger":
		p.monitorAndTrigger()
	default:
		fmt.Println(piperunHelp)
		os.Exit(1)
	}
}

func (p *PiperunCmd) runPipeline() {
	if p.PipelineURL == "" {
		fmt.Println(piperunRunHelp)
		os.Exit(1)
	}

	pat := utils.GetPAT(p.PAT)
	if pat == "" {
		fmt.Fprintln(os.Stderr, "Error: PAT is required. Run 'defenders conf' or set ADO_PAT environment variable.")
		fmt.Fprintln(os.Stderr, "Create a PAT at: https://dev.azure.com/{your-org}/_usersSettings/tokens")
		os.Exit(1)
	}

	orgURL, project, queryParams, err := parseADOUrl(p.PipelineURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URL: %s\n", err)
		os.Exit(1)
	}

	definitionID := queryParams.Get("definitionId")
	if definitionID == "" {
		fmt.Fprintln(os.Stderr, "Error: Could not extract definitionId from URL")
		os.Exit(1)
	}

	fmt.Printf("Triggering pipeline: %s\n", p.PipelineURL)
	fmt.Printf("Project: %s, Definition ID: %s\n", project, definitionID)

	// Run pipeline using az CLI
	stdout, stderr, err := utils.RunCommand("az", "pipelines", "run",
		"--id", definitionID,
		"--org", orgURL,
		"--project", project,
		"-o", "json",
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to trigger pipeline:")
		fmt.Fprintln(os.Stderr, stderr)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %s\n", err)
		os.Exit(1)
	}

	buildID := ""
	if id, ok := result["id"].(float64); ok {
		buildID = fmt.Sprintf("%.0f", id)
	}

	fmt.Println("\nSuccessfully triggered pipeline!")
	fmt.Printf("Build ID: %s\n", buildID)
	fmt.Printf("URL: %s/%s/_build/results?buildId=%s&view=results\n", orgURL, project, buildID)
}

func (p *PiperunCmd) monitorAndTrigger() {
	if p.WaitForURL == "" || p.TriggerURL == "" {
		fmt.Println(piperunMonitorHelp)
		os.Exit(1)
	}

	pat := utils.GetPAT(p.PAT)
	if pat == "" {
		fmt.Fprintln(os.Stderr, "Error: PAT is required. Run 'defenders conf' or set ADO_PAT environment variable.")
		os.Exit(1)
	}

	interval := p.Interval
	if interval <= 0 {
		interval = 30
	}

	waitOrgURL, project, queryParams, err := parseADOUrl(p.WaitForURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing wait URL: %s\n", err)
		os.Exit(1)
	}

	buildID := queryParams.Get("buildId")
	if buildID == "" {
		fmt.Fprintln(os.Stderr, "Error: Could not extract buildId from wait URL")
		os.Exit(1)
	}

	triggerOrgURL, triggerProject, triggerParams, err := parseADOUrl(p.TriggerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing trigger URL: %s\n", err)
		os.Exit(1)
	}

	definitionID := triggerParams.Get("definitionId")
	if definitionID == "" {
		fmt.Fprintln(os.Stderr, "Error: Could not extract definitionId from trigger URL")
		os.Exit(1)
	}

	fmt.Println("Starting pipeline monitor...")
	fmt.Printf("Monitoring: %s\n", p.WaitForURL)
	fmt.Printf("Will trigger: %s\n", p.TriggerURL)
	fmt.Printf("Check interval: %d seconds\n\n", interval)

	for {
		// Get build status
		stdout, stderr, err := utils.RunCommand("az", "pipelines", "runs", "show",
			"--id", buildID,
			"--org", waitOrgURL,
			"--project", project,
			"-o", "json",
		)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking pipeline status: %s\n", stderr)
			fmt.Printf("Retrying in %d seconds...\n", interval)
			time.Sleep(time.Duration(interval) * time.Second)
			continue
		}

		var buildResult map[string]interface{}
		if err := json.Unmarshal([]byte(stdout), &buildResult); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %s\n", err)
			time.Sleep(time.Duration(interval) * time.Second)
			continue
		}

		status := ""
		if s, ok := buildResult["status"].(string); ok {
			status = s
		}

		result := ""
		if r, ok := buildResult["result"].(string); ok {
			result = r
		}

		fmt.Printf("Pipeline %s status: %s", buildID, status)
		if result != "" {
			fmt.Printf(" (result: %s)", result)
		}
		fmt.Println()

		if status == "completed" {
			fmt.Printf("\nPipeline %s completed with result: %s\n", buildID, result)

			if result == "succeeded" {
				fmt.Println("Triggering second pipeline...")

				stdout, stderr, err := utils.RunCommand("az", "pipelines", "run",
					"--id", definitionID,
					"--org", triggerOrgURL,
					"--project", triggerProject,
					"-o", "json",
				)

				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to trigger pipeline: %s\n", stderr)
					os.Exit(1)
				}

				var triggerResult map[string]interface{}
				json.Unmarshal([]byte(stdout), &triggerResult)

				newBuildID := ""
				if id, ok := triggerResult["id"].(float64); ok {
					newBuildID = fmt.Sprintf("%.0f", id)
				}

				fmt.Printf("Successfully triggered pipeline %s\n", newBuildID)
				fmt.Printf("URL: %s/%s/_build/results?buildId=%s&view=results\n", triggerOrgURL, triggerProject, newBuildID)
			} else {
				fmt.Printf("Pipeline %s failed with result: %s\n", buildID, result)
			}
			break
		}

		currTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] Pipeline still running. Checking again in %d seconds...\n", currTime, interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// ParsePiperunArgs parses command line arguments for piperun command
func ParsePiperunArgs(args []string) *PiperunCmd {
	cmd := &PiperunCmd{
		Interval: 30,
	}

	if len(args) == 0 {
		return cmd
	}

	// First argument is subcommand
	cmd.Subcommand = args[0]

	if cmd.Subcommand == "-h" || cmd.Subcommand == "--help" {
		fmt.Println(piperunHelp)
		os.Exit(0)
	}

	positionalArgs := []string{}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-t" || arg == "--token":
			if i+1 < len(args) {
				i++
				cmd.PAT = args[i]
			}
		case strings.HasPrefix(arg, "--token="):
			cmd.PAT = strings.TrimPrefix(arg, "--token=")
		case strings.HasPrefix(arg, "-t="):
			cmd.PAT = strings.TrimPrefix(arg, "-t=")
		case arg == "-i" || arg == "--interval":
			if i+1 < len(args) {
				i++
				if val, err := strconv.Atoi(args[i]); err == nil {
					cmd.Interval = val
				}
			}
		case strings.HasPrefix(arg, "--interval="):
			if val, err := strconv.Atoi(strings.TrimPrefix(arg, "--interval=")); err == nil {
				cmd.Interval = val
			}
		case arg == "-h" || arg == "--help":
			switch cmd.Subcommand {
			case "run":
				fmt.Println(piperunRunHelp)
			case "monitor-trigger":
				fmt.Println(piperunMonitorHelp)
			default:
				fmt.Println(piperunHelp)
			}
			os.Exit(0)
		default:
			if !strings.HasPrefix(arg, "-") {
				positionalArgs = append(positionalArgs, arg)
			}
		}
	}

	// Assign positional args based on subcommand
	switch cmd.Subcommand {
	case "run":
		if len(positionalArgs) >= 1 {
			cmd.PipelineURL = positionalArgs[0]
		}
	case "monitor-trigger":
		if len(positionalArgs) >= 1 {
			cmd.WaitForURL = positionalArgs[0]
		}
		if len(positionalArgs) >= 2 {
			cmd.TriggerURL = positionalArgs[1]
		}
	}

	return cmd
}
