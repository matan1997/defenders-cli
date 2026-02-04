package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"defenders-cli/internal/utils"
)

const prmeHelp = `prme - Create Azure DevOps PR from current branch to default branch

USAGE:
  defenders prme [flags]

FLAGS:
  -i, --work-item  Work item ID to link to the PR
  -t, --title      Custom PR title (default: branch name after last /)
  -h, --help       Show this help message

EXAMPLES:
  defenders prme
  defenders prme -i 12345
  defenders prme -t "My PR Title"
  defenders prme -i 12345 -t "My PR Title"
`

type PrmeCmd struct {
	WorkItem string
	Title    string
}

func (p *PrmeCmd) Run() {
	// Get current branch
	branch, err := utils.GetCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if branch == "" {
		fmt.Fprintln(os.Stderr, "Error: Not in a git repository or no branch checked out")
		os.Exit(1)
	}

	// Determine default branch
	defaultBranch, err := utils.GetDefaultBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Set title (custom or default from branch name)
	title := p.Title
	if title == "" {
		title = utils.GetBranchTitle(branch)
	}

	fmt.Printf("Creating PR: %s -> %s\n", branch, defaultBranch)
	fmt.Printf("Title: %s\n", title)
	if p.WorkItem != "" {
		fmt.Printf("Work Item: %s\n", p.WorkItem)
	}

	// Build az command arguments
	args := []string{"repos", "pr", "create",
		"-s", branch,
		"-t", defaultBranch,
		"--title", title,
		"-o", "json",
	}

	if p.WorkItem != "" {
		args = append(args, "--work-items", p.WorkItem)
	}

	// Execute az command
	stdout, stderr, err := utils.RunCommand("az", args...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create PR:")
		fmt.Fprintln(os.Stderr, stderr)
		os.Exit(1)
	}

	// Parse JSON response
	var prJSON map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &prJSON); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %s\n", err)
		os.Exit(1)
	}

	// Extract repository name and PR ID
	repoName := ""
	if repo, ok := prJSON["repository"].(map[string]interface{}); ok {
		repoName = repo["name"].(string)
	}

	prID := ""
	if id, ok := prJSON["pullRequestId"].(float64); ok {
		prID = fmt.Sprintf("%.0f", id)
	}

	if repoName != "" && prID != "" {
		org := utils.GetOrganization("")
		project := utils.GetProject("")
		fmt.Printf("%s/%s/_git/%s/pullrequest/%s\n", org, project, repoName, prID)
	} else {
		fmt.Println("PR created successfully")
		fmt.Println(stdout)
	}
}

// ParsePrmeArgs parses command line arguments for prme command
func ParsePrmeArgs(args []string) *PrmeCmd {
	cmd := &PrmeCmd{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-i" || arg == "--work-item":
			if i+1 < len(args) {
				i++
				cmd.WorkItem = args[i]
			}
		case arg == "-t" || arg == "--title":
			if i+1 < len(args) {
				i++
				cmd.Title = args[i]
			}
		case strings.HasPrefix(arg, "--work-item="):
			cmd.WorkItem = strings.TrimPrefix(arg, "--work-item=")
		case strings.HasPrefix(arg, "--title="):
			cmd.Title = strings.TrimPrefix(arg, "--title=")
		case arg == "-h" || arg == "--help":
			fmt.Println(prmeHelp)
			os.Exit(0)
		}
	}

	return cmd
}
