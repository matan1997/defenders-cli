package utils

import (
	"fmt"
	"strings"
)

// GetCurrentBranch returns the current git branch name
func GetCurrentBranch() (string, error) {
	stdout, _, err := RunCommand("git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("not in a git repository or no branch checked out")
	}
	return strings.TrimSpace(stdout), nil
}

// GetDefaultBranch determines the default branch (develop, main, or master)
func GetDefaultBranch() (string, error) {
	// Check for develop
	_, _, err := RunCommand("git", "show-ref", "--verify", "--quiet", "refs/remotes/origin/develop")
	if err == nil {
		return "develop", nil
	}

	// Check for main
	_, _, err = RunCommand("git", "show-ref", "--verify", "--quiet", "refs/remotes/origin/main")
	if err == nil {
		return "main", nil
	}

	// Check for master
	_, _, err = RunCommand("git", "show-ref", "--verify", "--quiet", "refs/remotes/origin/master")
	if err == nil {
		return "master", nil
	}

	// Try to get from remote HEAD
	stdout, _, err := RunCommand("git", "remote", "show", "origin")
	if err == nil {
		lines := strings.Split(stdout, "\n")
		for _, line := range lines {
			if strings.Contains(line, "HEAD branch") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not determine default branch")
}

// GetBranchTitle extracts title from branch name (part after last /)
func GetBranchTitle(branch string) string {
	parts := strings.Split(branch, "/")
	return parts[len(parts)-1]
}
