package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Force disables interactive prompts when true
var Force bool

const HELPER = `Hi There!
This CLI for Azure DevOps operations

SYNTAX:
  defenders <command> [flags]

COMMANDS:
  conf        Configure CLI settings (PAT, org, project, etc.)
  get-token   Open browser to create PAT with required permissions
  cado        Create ADO Feature work item with parent link
  prme        Create PR from current branch to default branch
  release     Pipeline operations (run, monitor-trigger)
  pr          PR approval operations (approve, reset)

GLOBAL FLAGS:
  -h, --help  Show this help message

EXAMPLES:
  defenders conf                              # Interactive setup
  defenders conf show                         # Show current config
  defenders cado --title "My Feature"
  defenders cado --title "My Feature" --parent 12345
  defenders prme
  defenders prme -i 12345 -t "My PR Title"
  defenders release run <pipeline-url>
  defenders release monitor-trigger <wait-url> <trigger-url>
  defenders pr --approve <pr-url>
  defenders pr --reset <pr-url>

AUTHENTICATION:
  Run 'defenders conf' to set up your configuration.
  Alternatively, set ADO_PAT environment variable.
  Create a PAT at: https://dev.azure.com/{org}/_usersSettings/tokens
`

// AskUser prompts user for confirmation. Returns true if user confirms.
func AskUser(message string, args ...any) bool {
	if Force {
		return true
	}

	fmt.Printf(message, args...)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

// RunCommand executes a shell command and returns stdout, stderr, and error
func RunCommand(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RunCommandWithOutput executes a command and prints output in real-time
func RunCommandWithOutput(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommandWithPAT executes a command with AZURE_DEVOPS_EXT_PAT environment variable set
// This is used for Azure CLI commands that need PAT authentication
// It also isolates the Azure config to prevent fallback to az login credentials
func RunCommandWithPAT(pat string, name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)

	// Create isolated Azure config directory to prevent using az login credentials
	// Uses os.TempDir() for cross-platform compatibility (Linux/Mac: /tmp, Windows: %TEMP%)
	isolatedConfigDir := filepath.Join(os.TempDir(), "defenders-az-isolated")
	os.MkdirAll(isolatedConfigDir, 0755)

	// Inherit current environment and add/override PAT + isolated config
	cmd.Env = append(os.Environ(),
		"AZURE_DEVOPS_EXT_PAT="+pat,
		"AZURE_CONFIG_DIR="+isolatedConfigDir,
	)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RunCommandWithOptionalPAT executes a command, using PAT authentication if provided
// If pat is empty, it falls back to az login identity
func RunCommandWithOptionalPAT(pat string, name string, args ...string) (string, string, error) {
	if pat != "" {
		return RunCommandWithPAT(pat, name, args...)
	}
	return RunCommand(name, args...)
}

// GetEnvOrDefault returns environment variable value or default if not set
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// RequireEnv returns environment variable or exits with error message
func RequireEnv(key, errorMessage string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Fprintf(os.Stderr, "Error: %s\n", errorMessage)
		fmt.Fprintf(os.Stderr, "Set %s environment variable or use appropriate flag.\n", key)
		os.Exit(1)
	}
	return value
}
