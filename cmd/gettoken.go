package cmd

import (
	"fmt"
	"os"
	"runtime"

	"defenders-cli/internal/utils"
)

const gettokenHelp = `get-token - Open browser to create a PAT token with required permissions

USAGE:
  defenders get-token

This command opens your browser to the Azure DevOps PAT creation page with
instructions on which permissions to select.

REQUIRED PERMISSIONS:
  - Work Items: Read & Write (for 'cado' command)
  - Code: Read & Write (for 'prme' and 'pr' commands)
  - Build: Read & Execute (for 'release' command)

After creating the token, run 'defenders conf' to save it.
`

type GetTokenCmd struct{}

func (g *GetTokenCmd) Run() {
	org := utils.GetOrganization("")

	// Build the PAT creation URL
	patURL := fmt.Sprintf("%s/_usersSettings/tokens", org)

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║            Create Azure DevOps PAT Token                   ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Opening browser to create a new PAT token...")
	fmt.Println()
	fmt.Printf("URL: %s\n", patURL)
	fmt.Println()
	fmt.Println("─────────────────────────────────────")
	fmt.Println("REQUIRED PERMISSIONS:")
	fmt.Println("─────────────────────────────────────")
	fmt.Println()
	fmt.Println("  1. Work Items")
	fmt.Println("     ☑ Read & Write")
	fmt.Println()
	fmt.Println("  2. Code")
	fmt.Println("     ☑ Read & Write")
	fmt.Println()
	fmt.Println("  3. Build")
	fmt.Println("     ☑ Read & Execute")
	fmt.Println()
	fmt.Println("─────────────────────────────────────")
	fmt.Println()
	fmt.Println("After creating the token, run:")
	fmt.Println("  defenders conf")
	fmt.Println()
	fmt.Println("to save it to your configuration.")
	fmt.Println()

	// Open browser
	err := openBrowser(patURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open browser: %s\n", err)
		fmt.Println("Please open the URL manually in your browser.")
	}
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	_, _, err := utils.RunCommand(cmd, args...)
	return err
}

// ParseGetTokenArgs parses command line arguments for get-token command
func ParseGetTokenArgs(args []string) *GetTokenCmd {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			fmt.Println(gettokenHelp)
			os.Exit(0)
		}
	}
	return &GetTokenCmd{}
}
