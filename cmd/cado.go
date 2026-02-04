package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"defenders-cli/internal/utils"
)

const cadoHelp = `cado - Create ADO Feature work item with parent link and current iteration

USAGE:
  defenders cado --title="Your Title" [--parent=12345]

FLAGS:
  --title       (required) Title of the Feature work item
  --parent      (optional) Parent work item ID to link
  --assigned-to (optional) Override assigned-to from config

EXAMPLES:
  defenders cado --title "Implement new feature"
  defenders cado --title "My Task" --parent 12345

NOTE:
  Uses configuration from 'defenders conf' for org, project, team, and area.
`

type CadoCmd struct {
	Title      string
	Parent     string
	AssignedTo string
}

func (c *CadoCmd) Run() {
	if c.Title == "" {
		fmt.Println(cadoHelp)
		os.Exit(1)
	}

	// Get config values
	org := utils.GetOrganization("")
	project := utils.GetProject("")
	team := utils.GetTeam("")
	area := utils.GetArea("")
	assignedTo := utils.GetAssignedTo(c.AssignedTo)

	fmt.Printf("Creating Feature: %s\n", c.Title)
	if c.Parent != "" {
		fmt.Printf("Parent: %s\n", c.Parent)
	}

	// Get the current iteration from ADO
	stdout, stderr, err := utils.RunCommand("az", "boards", "iteration", "team", "list",
		"--org", org,
		"--project", project,
		"--team", team,
		"--timeframe", "current",
		"--query", "[0].path",
		"-o", "tsv",
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not get current iteration: %s\n", stderr)
		os.Exit(1)
	}

	iteration := strings.TrimSpace(stdout)
	if iteration == "" {
		fmt.Fprintln(os.Stderr, "Error: Could not get current iteration")
		os.Exit(1)
	}

	fmt.Printf("Iteration: %s\n", iteration)

	// Build work item create arguments
	args := []string{"boards", "work-item", "create",
		"--org", org,
		"--project", project,
		"--type", "Feature",
		"--title", c.Title,
		"--iteration", iteration,
		"--area", area,
		"--query", "id",
		"-o", "tsv",
	}

	// Add assigned-to if specified
	if assignedTo != "" {
		args = append(args, "--assigned-to", assignedTo)
	}

	// Create the work item
	stdout, stderr, err = utils.RunCommand("az", args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create work item: %s\n", stderr)
		os.Exit(1)
	}

	newID := strings.TrimSpace(stdout)
	if newID == "" {
		fmt.Fprintln(os.Stderr, "Error: Failed to create work item")
		os.Exit(1)
	}

	// Add parent link if provided
	if c.Parent != "" {
		_, stderr, err = utils.RunCommand("az", "boards", "work-item", "relation", "add",
			"--org", org,
			"--id", newID,
			"--relation-type", "parent",
			"--target-id", c.Parent,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to add parent link: %s\n", stderr)
		}
	}

	fmt.Printf("%s/%s/_workitems/edit/%s\n", org, project, newID)
}

// ParseCadoArgs parses command line arguments for cado command
func ParseCadoArgs(args []string) *CadoCmd {
	cmd := &CadoCmd{}

	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--title="):
			cmd.Title = strings.TrimPrefix(arg, "--title=")
		case strings.HasPrefix(arg, "--parent="):
			cmd.Parent = strings.TrimPrefix(arg, "--parent=")
		case strings.HasPrefix(arg, "--assigned-to="):
			cmd.AssignedTo = strings.TrimPrefix(arg, "--assigned-to=")
		case arg == "-h" || arg == "--help":
			fmt.Println(cadoHelp)
			os.Exit(0)
		}
	}

	return cmd
}

// Helper function to parse JSON from az command output
func parseJSON(data string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	return result
}
