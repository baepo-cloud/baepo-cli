package config

import (
	"fmt"
	"os"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Contexts map[string]*Context `yaml:"contexts" env-upd:""`
	Context  string              `yaml:"context" env:"BAEPO_CONTEXT"`

	CurrentContext *Context `yaml:"-"` // Not saved to config file

	ConfigVersion string `yaml:"version"`
}

type Context struct {
	SecretKey   *string `yaml:"secret_key" env:"BAEPO_SECRET_KEY"`
	WorkspaceID *string `yaml:"workspace_id" env:"BAEPO_WORKSPACE_ID"`
	UserID      *string `yaml:"user_id" env:"BAEPO_USER_ID"`
	URL         string  `yaml:"url" env:"BAEPO_URL"`
}

var DefaultContext = &Context{
	SecretKey:   nil,
	WorkspaceID: nil,
	UserID:      nil,
	URL:         "http://138.201.222.180:3000/",
}

// LoadConfig loads the configuration from the specified file path.
// First it tries to load from the file in $HOME/.baepo/config.yaml
// If that file does not exist, it will create it and fill with default values.
// If the file exists, it will load the configuration from it.
//
// Then it will check if the current context is set.
// If it is not set, it will set it to the default context.
// Then the current context will be set to the context specified in the config file.
// It can be surcharged by the environment variable BAEPO_CURRENT_CONTEXT.
// Also SecretKey and WorkspaceID can be surcharged by the environment variables BAEPO_SECRET_KEY and BAEPO_WORKSPACE_ID.
func LoadConfig(currentContext string) (*Config, error) {
	var configuration *Config
	// Get config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := path.Join(homeDir, ".baepo")
	configPath := path.Join(configDir, "config.yaml")

	// Ensure config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	// Initialize with default values
	configuration = &Config{
		Contexts: map[string]*Context{
			"default": DefaultContext,
		},
		Context:       "default",
		ConfigVersion: "0.1",
	}

	// Create config file with defaults if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := SaveConfig(configuration); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
	}

	// Read configuration file (cleanenv will use default values from the struct if file doesn't exist)
	if err := cleanenv.ReadConfig(configPath, configuration); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Load environment variables into the config (cleanenv automatically handles this with the env tags)
	if err := cleanenv.ReadEnv(configuration); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}

	// Determine which context to use
	contextName := configuration.Context

	// Override with parameter if provided
	if currentContext != "" && currentContext != "default" {
		contextName = currentContext
	}

	// Ensure the context exists
	selectedContext, exists := configuration.Contexts[contextName]
	if !exists {
		return nil, fmt.Errorf("context '%s' does not exist in configuration", contextName)
	}

	// Set the current context
	configuration.CurrentContext = selectedContext

	cleanenv.UpdateEnv(&configuration.CurrentContext)

	return configuration, nil
}

func SaveConfig(cfg *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := path.Join(homeDir, ".baepo")
	configPath := path.Join(configDir, "config.yaml")

	// Ensure config directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}
	if err := os.WriteFile(configPath, out, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
