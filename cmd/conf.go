package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"defenders-cli/internal/utils"
)

const confHelp = `conf - Configure defenders CLI

USAGE:
  defenders conf [subcommand]

SUBCOMMANDS:
  (none)    Interactive configuration wizard
  show      Show current configuration
  path      Show configuration file path
  reset     Reset configuration to defaults

EXAMPLES:
  defenders conf           # Interactive setup
  defenders conf show      # Show current config
  defenders conf path      # Show config file location
  defenders conf reset     # Reset to defaults

CONFIG FILE LOCATION:
  Linux/macOS: ~/.config/defenders/config.json
  Windows:     %%APPDATA%%\defenders\config.json
`

type ConfCmd struct {
	Subcommand string
}

func (c *ConfCmd) Run() {
	switch c.Subcommand {
	case "show":
		c.showConfig()
	case "path":
		c.showPath()
	case "reset":
		c.resetConfig()
	case "", "setup":
		c.interactiveSetup()
	default:
		fmt.Printf("Unknown subcommand: %s\n", c.Subcommand)
		fmt.Println(confHelp)
		os.Exit(1)
	}
}

func (c *ConfCmd) showConfig() {
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
		os.Exit(1)
	}

	if config == nil {
		fmt.Println("No configuration file found.")
		fmt.Println("Run 'defenders conf' to create one.")
		return
	}

	fmt.Println("Current Configuration:")
	fmt.Println("─────────────────────────────────────")
	fmt.Printf("  PAT:          %s\n", maskPAT(config.PAT))
	fmt.Printf("  Organization: %s\n", config.Organization)
	fmt.Printf("  Project:      %s\n", config.Project)
	fmt.Printf("  Team:         %s\n", config.Team)
	fmt.Printf("  Area:         %s\n", config.Area)
	fmt.Printf("  Assigned To:  %s\n", config.AssignedTo)
	fmt.Println("─────────────────────────────────────")

	configPath, _ := utils.GetConfigPath()
	fmt.Printf("\nConfig file: %s\n", configPath)
}

func (c *ConfCmd) showPath() {
	configPath, err := utils.GetConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(configPath)

	if utils.ConfigExists() {
		fmt.Println("(file exists)")
	} else {
		fmt.Println("(file does not exist)")
	}
}

func (c *ConfCmd) resetConfig() {
	if !utils.ConfigExists() {
		fmt.Println("No configuration file exists.")
		return
	}

	fmt.Print("Are you sure you want to reset configuration? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input != "y" && input != "yes" {
		fmt.Println("Cancelled.")
		return
	}

	config := utils.DefaultConfig()
	if err := utils.SaveConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration reset to defaults.")
}

func (c *ConfCmd) interactiveSetup() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║           Defenders CLI Configuration Wizard               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Load existing config or use defaults
	existingConfig, _ := utils.LoadConfig()
	defaults := utils.DefaultConfig()

	if existingConfig != nil {
		defaults = existingConfig
		fmt.Println("Existing configuration found. Press Enter to keep current values.")
	} else {
		fmt.Println("No existing configuration. Setting up new config.")
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	config := &utils.Config{}

	// PAT Token
	fmt.Println("1. Personal Access Token (PAT)")
	fmt.Println("   Create at: https://dev.azure.com/<your-org>/_usersSettings/tokens")
	fmt.Println("   Required scopes: Work Items (Read/Write), Code (Read/Write), Build (Read/Execute)")
	if defaults.PAT != "" {
		fmt.Printf("   Current: %s\n", maskPAT(defaults.PAT))
	}
	fmt.Print("   PAT Token: ")
	pat, _ := reader.ReadString('\n')
	pat = strings.TrimSpace(pat)
	if pat == "" {
		config.PAT = defaults.PAT
	} else {
		config.PAT = pat
	}
	fmt.Println()

	// Organization
	fmt.Println("2. Azure DevOps Organization URL")
	fmt.Printf("   Default: %s\n", defaults.Organization)
	fmt.Print("   Organization URL: ")
	org, _ := reader.ReadString('\n')
	org = strings.TrimSpace(org)
	if org == "" {
		config.Organization = defaults.Organization
	} else {
		config.Organization = org
	}
	fmt.Println()

	// Project
	fmt.Println("3. Project Name")
	fmt.Printf("   Default: %s\n", defaults.Project)
	fmt.Print("   Project: ")
	project, _ := reader.ReadString('\n')
	project = strings.TrimSpace(project)
	if project == "" {
		config.Project = defaults.Project
	} else {
		config.Project = project
	}
	fmt.Println()

	// Team
	fmt.Println("4. Team Name")
	fmt.Printf("   Default: %s\n", defaults.Team)
	fmt.Print("   Team: ")
	team, _ := reader.ReadString('\n')
	team = strings.TrimSpace(team)
	if team == "" {
		config.Team = defaults.Team
	} else {
		config.Team = team
	}
	fmt.Println()

	// Area Path
	fmt.Println("5. Area Path (for work items)")
	fmt.Printf("   Default: %s\n", defaults.Area)
	fmt.Print("   Area Path: ")
	area, _ := reader.ReadString('\n')
	area = strings.TrimSpace(area)
	if area == "" {
		config.Area = defaults.Area
	} else {
		config.Area = area
	}
	fmt.Println()

	// Assigned To
	fmt.Println("6. Default Assigned To (email for work items)")
	if defaults.AssignedTo != "" {
		fmt.Printf("   Current: %s\n", defaults.AssignedTo)
	} else {
		fmt.Println("   (optional - leave empty to skip)")
	}
	fmt.Print("   Assigned To: ")
	assignedTo, _ := reader.ReadString('\n')
	assignedTo = strings.TrimSpace(assignedTo)
	if assignedTo == "" {
		config.AssignedTo = defaults.AssignedTo
	} else {
		config.AssignedTo = assignedTo
	}
	fmt.Println()

	// Summary
	fmt.Println("─────────────────────────────────────")
	fmt.Println("Configuration Summary:")
	fmt.Println("─────────────────────────────────────")
	fmt.Printf("  PAT:          %s\n", maskPAT(config.PAT))
	fmt.Printf("  Organization: %s\n", config.Organization)
	fmt.Printf("  Project:      %s\n", config.Project)
	fmt.Printf("  Team:         %s\n", config.Team)
	fmt.Printf("  Area:         %s\n", config.Area)
	fmt.Printf("  Assigned To:  %s\n", config.AssignedTo)
	fmt.Println("─────────────────────────────────────")
	fmt.Println()

	fmt.Print("Save this configuration? [Y/n]: ")
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "n" || confirm == "no" {
		fmt.Println("Configuration cancelled.")
		return
	}

	// Save config
	if err := utils.SaveConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving configuration: %s\n", err)
		os.Exit(1)
	}

	configPath, _ := utils.GetConfigPath()
	fmt.Printf("\n✓ Configuration saved to: %s\n", configPath)
	fmt.Println("\nYou can now use defenders CLI commands!")
}

// maskPAT masks the PAT token for display, showing only first and last 4 chars
func maskPAT(pat string) string {
	if pat == "" {
		return "(not set)"
	}
	if len(pat) <= 8 {
		return "****"
	}
	return pat[:4] + "..." + pat[len(pat)-4:]
}

// ParseConfArgs parses command line arguments for conf command
func ParseConfArgs(args []string) *ConfCmd {
	cmd := &ConfCmd{}

	for _, arg := range args {
		switch {
		case arg == "-h" || arg == "--help":
			fmt.Println(confHelp)
			os.Exit(0)
		default:
			if !strings.HasPrefix(arg, "-") && cmd.Subcommand == "" {
				cmd.Subcommand = arg
			}
		}
	}

	return cmd
}
