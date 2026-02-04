package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"defenders-cli/internal/utils"
)

const prhandlerHelp = `pr - Azure DevOps Pull Request CLI tool

USAGE:
  defenders pr <action> <pr-url> [-t <token>]

ACTIONS:
  --approve  Approve the Pull Request
  --reset    Reset your vote on the Pull Request

FLAGS:
  -t, --token  Personal Access Token (overrides config/env - use another user's PAT)

ARGUMENTS:
  <pr-url>   Azure DevOps Pull Request URL

EXAMPLES:
  defenders pr --approve https://dev.azure.com/org/project/_git/repo/pullrequest/123
  defenders pr --reset https://dev.azure.com/org/project/_git/repo/pullrequest/123
  defenders pr --approve <url> -t <other-user-pat>

AUTHENTICATION:
  PAT with PR approval permissions required.
  Use -t to provide another user's PAT token to approve/vote on their behalf.
  Create at: https://msazure.visualstudio.com/_usersSettings/tokens
`

type PrhandlerCmd struct {
	Approve bool
	Reset   bool
	PRURL   string
	PAT     string
}

// parsePRUrl parses Azure DevOps PR URL and extracts components
// Supports:
// - https://dev.azure.com/{org}/{project}/_git/{repo}/pullrequest/{pr_id}
// - https://{org}.visualstudio.com/{project}/_git/{repo}/pullrequest/{pr_id}
func parsePRUrl(rawURL string) (orgURL, project, repository string, prID string, err error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", "", err
	}

	orgURL = fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")

	// Find indices for _git and pullrequest
	gitIndex := -1
	prIndex := -1
	for i, part := range pathParts {
		if part == "_git" {
			gitIndex = i
		}
		if part == "pullrequest" {
			prIndex = i
		}
	}

	if gitIndex == -1 || prIndex == -1 {
		return "", "", "", "", fmt.Errorf("invalid PR URL format")
	}

	// Project is the part right before _git
	project = pathParts[gitIndex-1]
	repository = pathParts[gitIndex+1]
	prID = pathParts[prIndex+1]

	return orgURL, project, repository, prID, nil
}

func (p *PrhandlerCmd) Run() {
	// Validate that exactly one action is specified
	if p.Approve == p.Reset {
		fmt.Fprintln(os.Stderr, "Error: You must specify either --approve or --reset (but not both)")
		fmt.Println(prhandlerHelp)
		os.Exit(1)
	}

	if p.PRURL == "" {
		fmt.Fprintln(os.Stderr, "Error: PR URL is required")
		fmt.Println(prhandlerHelp)
		os.Exit(1)
	}

	// Get PAT
	pat := utils.GetPAT(p.PAT)
	if pat == "" {
		fmt.Fprintln(os.Stderr, "Error: PAT is required. Run 'defenders conf' or set ADO_PAT environment variable.")
		fmt.Fprintln(os.Stderr, "Get yours from https://msazure.visualstudio.com/_usersSettings/tokens")
		os.Exit(1)
	}

	orgURL, project, repository, prID, err := parsePRUrl(p.PRURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing PR URL: %s\n", err)
		os.Exit(1)
	}

	// Determine vote value
	// Vote values: 10 = approved, 5 = approved with suggestions, 0 = no vote, -5 = waiting for author, -10 = rejected
	var vote string
	var action string
	if p.Approve {
		vote = "approve"
		action = "approved"
	} else {
		vote = "reset"
		action = "vote reset"
	}

	fmt.Printf("Processing PR #%s...\n", prID)
	fmt.Printf("Organization: %s\n", orgURL)
	fmt.Printf("Project: %s\n", project)
	fmt.Printf("Repository: %s\n", repository)
	fmt.Printf("Action: %s\n", action)

	// Use az repos pr set-vote command
	stdout, stderr, err := utils.RunCommand("az", "repos", "pr", "set-vote",
		"--id", prID,
		"--vote", vote,
		"--org", orgURL,
		"-o", "json",
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", stderr)
		os.Exit(1)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		// Command succeeded but response parsing failed - still a success
		fmt.Printf("✓ PR #%s %s successfully!\n", prID, action)
		fmt.Printf("  Repository: %s\n", repository)
		fmt.Printf("  Project: %s\n", project)
		return
	}

	fmt.Printf("✓ PR #%s %s successfully!\n", prID, action)
	fmt.Printf("  Repository: %s\n", repository)
	fmt.Printf("  Project: %s\n", project)
}

// ParsePrhandlerArgs parses command line arguments for prhandler command
func ParsePrhandlerArgs(args []string) *PrhandlerCmd {
	cmd := &PrhandlerCmd{}

	positionalArgs := []string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--approve":
			cmd.Approve = true
		case arg == "--reset":
			cmd.Reset = true
		case arg == "-t" || arg == "--token":
			if i+1 < len(args) {
				i++
				cmd.PAT = args[i]
			}
		case strings.HasPrefix(arg, "--token="):
			cmd.PAT = strings.TrimPrefix(arg, "--token=")
		case strings.HasPrefix(arg, "-t="):
			cmd.PAT = strings.TrimPrefix(arg, "-t=")
		case arg == "-h" || arg == "--help":
			fmt.Println(prhandlerHelp)
			os.Exit(0)
		default:
			if !strings.HasPrefix(arg, "-") {
				positionalArgs = append(positionalArgs, arg)
			}
		}
	}

	// Assign positional args
	if len(positionalArgs) >= 1 {
		cmd.PRURL = positionalArgs[0]
	}

	return cmd
}
