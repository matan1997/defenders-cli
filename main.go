package main

import (
	"fmt"
	"os"

	"defenders-cli/cmd"
	"defenders-cli/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(utils.HELPER)
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	// Check for global help flag
	if command == "-h" || command == "--help" || command == "help" {
		fmt.Println(utils.HELPER)
		os.Exit(0)
	}

	switch command {
	case "conf":
		confCmd := cmd.ParseConfArgs(args)
		confCmd.Run()

	case "get-token":
		getTokenCmd := cmd.ParseGetTokenArgs(args)
		getTokenCmd.Run()

	case "cado":
		cadoCmd := cmd.ParseCadoArgs(args)
		cadoCmd.Run()

	case "prme":
		prmeCmd := cmd.ParsePrmeArgs(args)
		prmeCmd.Run()

	case "release":
		releaseCmd := cmd.ParsePiperunArgs(args)
		releaseCmd.Run()

	case "pr":
		prCmd := cmd.ParsePrhandlerArgs(args)
		prCmd.Run()

	default:
		fmt.Printf("Unknown command: %s\nSee 'defenders --help'\n", command)
		os.Exit(1)
	}
}
