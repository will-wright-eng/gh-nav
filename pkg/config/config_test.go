package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test loading config
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that we have a token (either from gh CLI or env var)
	if cfg.GitHub.Token == "" {
		t.Error("Expected GitHub token to be set")
	}

	// Check other config values
	if cfg.GitHub.BaseURL != "https://api.github.com" {
		t.Errorf("Expected base URL to be https://api.github.com, got %s", cfg.GitHub.BaseURL)
	}

	if cfg.UI.Theme != "dark" {
		t.Errorf("Expected theme to be dark, got %s", cfg.UI.Theme)
	}
}

func TestGitHubCLIToken(t *testing.T) {
	// Test gh CLI token retrieval
	token, err := getGitHubCLIToken()
	if err != nil {
		// It's okay if gh CLI is not available, but if it is, we should get a token
		t.Logf("gh CLI token not available: %v", err)
		return
	}

	if token == "" {
		t.Error("Expected non-empty token from gh CLI")
	}
}

func TestLoadWithEnvVar(t *testing.T) {
	// Test with environment variable
	testToken := "test-token-123"
	os.Setenv("GITHUB_TOKEN", testToken)
	defer os.Unsetenv("GITHUB_TOKEN")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.GitHub.Token != testToken {
		t.Errorf("Expected token to be %s, got %s", testToken, cfg.GitHub.Token)
	}
}
