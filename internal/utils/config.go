package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the defenders CLI configuration
type Config struct {
	PAT          string `json:"pat"`
	Organization string `json:"organization"`
	Project      string `json:"project"`
	Team         string `json:"team"`
	Area         string `json:"area"`
	AssignedTo   string `json:"assigned_to"`
}

// DefaultConfig returns default configuration values
func DefaultConfig() *Config {
	return &Config{
		Organization: "https://dev.azure.com/msazure",
		Project:      "One",
		Team:         "Rome",
		Area:         `One\Rome\CNAPP\Defenders\BarTeam`,
		AssignedTo:   "",
	}
}

// GetConfigDir returns the configuration directory path based on OS
func GetConfigDir() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		configDir = filepath.Join(appData, "defenders")
	default: // linux, darwin, etc.
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config", "defenders")
	}

	return configDir, nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig loads configuration from the config file
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No config file exists
		}
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to the config file
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("could not serialize config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

// GetConfigValue returns a config value with priority: flag > env > config file > default
func GetConfigValue(flagValue, envKey, configValue, defaultValue string) string {
	// Priority 1: Flag value
	if flagValue != "" {
		return flagValue
	}

	// Priority 2: Environment variable
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}

	// Priority 3: Config file value
	if configValue != "" {
		return configValue
	}

	// Priority 4: Default value
	return defaultValue
}

// GetPAT returns the PAT token with priority: flag > env > config
func GetPAT(flagValue string) string {
	config, _ := LoadConfig()
	configPAT := ""
	if config != nil {
		configPAT = config.PAT
	}
	return GetConfigValue(flagValue, "ADO_PAT", configPAT, "")
}

// GetOrganization returns the organization with priority: flag > env > config
func GetOrganization(flagValue string) string {
	config, _ := LoadConfig()
	configOrg := ""
	if config != nil {
		configOrg = config.Organization
	}
	return GetConfigValue(flagValue, "ADO_ORG", configOrg, "https://dev.azure.com/msazure")
}

// GetProject returns the project with priority: flag > env > config
func GetProject(flagValue string) string {
	config, _ := LoadConfig()
	configProject := ""
	if config != nil {
		configProject = config.Project
	}
	return GetConfigValue(flagValue, "ADO_PROJECT", configProject, "One")
}

// GetTeam returns the team with priority: flag > env > config
func GetTeam(flagValue string) string {
	config, _ := LoadConfig()
	configTeam := ""
	if config != nil {
		configTeam = config.Team
	}
	return GetConfigValue(flagValue, "ADO_TEAM", configTeam, "Rome")
}

// GetArea returns the area with priority: flag > env > config
func GetArea(flagValue string) string {
	config, _ := LoadConfig()
	configArea := ""
	if config != nil {
		configArea = config.Area
	}
	return GetConfigValue(flagValue, "ADO_AREA", configArea, `One\Rome\CNAPP\Defenders\BarTeam`)
}

// GetAssignedTo returns the assigned-to email with priority: flag > env > config
func GetAssignedTo(flagValue string) string {
	config, _ := LoadConfig()
	configAssignedTo := ""
	if config != nil {
		configAssignedTo = config.AssignedTo
	}
	return GetConfigValue(flagValue, "ADO_ASSIGNED_TO", configAssignedTo, "")
}

// ConfigExists checks if a config file exists
func ConfigExists() bool {
	configPath, err := GetConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(configPath)
	return err == nil
}
