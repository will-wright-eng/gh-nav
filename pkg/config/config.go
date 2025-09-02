package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	GitHub GitHubConfig `yaml:"github"`
	UI     UIConfig     `yaml:"ui"`
}

// GitHubConfig holds GitHub API configuration
type GitHubConfig struct {
	Token   string `yaml:"token"`
	BaseURL string `yaml:"base_url"`
}

// UIConfig holds UI-specific configuration
type UIConfig struct {
	Theme       string        `yaml:"theme"`
	RefreshRate time.Duration `yaml:"refresh_rate"`
}

// Load loads configuration from environment and defaults
func Load() (*Config, error) {
	// Try to get token from gh CLI first, fallback to environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		// Try to get token from gh CLI
		if ghToken, err := getGitHubCLIToken(); err == nil {
			token = ghToken
		}
	}

	cfg := &Config{
		GitHub: GitHubConfig{
			Token:   token,
			BaseURL: "https://api.github.com",
		},
		UI: UIConfig{
			Theme:       "dark",
			RefreshRate: time.Second,
		},
	}

	return cfg, nil
}

// getGitHubCLIToken retrieves the GitHub token from the gh CLI
func getGitHubCLIToken() (string, error) {
	// Execute gh auth token command
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get gh token: %w", err)
	}

	// Trim whitespace and newlines
	token := strings.TrimSpace(string(output))
	if token == "" {
		return "", fmt.Errorf("gh token is empty")
	}

	return token, nil
}
