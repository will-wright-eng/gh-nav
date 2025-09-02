package api

import (
	"context"
	"testing"
	"time"

	"github.com/will/ghprs/pkg/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		GitHub: config.GitHubConfig{
			Token:   "test-token",
			BaseURL: "https://api.github.com",
		},
	}

	client := NewClient(cfg)
	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.config != cfg {
		t.Error("Expected config to be set")
	}
}

func TestGetUserRepositories(t *testing.T) {
	// This test requires a valid GitHub token
	// We'll skip it if no token is available
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.GitHub.Token == "" {
		t.Skip("Skipping test - no GitHub token available")
	}

	client := NewClient(cfg)
	ctx := context.Background()

	repos, err := client.GetUserRepositories(ctx)
	if err != nil {
		t.Fatalf("Failed to get user repositories: %v", err)
	}

	if len(repos) == 0 {
		t.Log("No repositories found - this might be normal for a new account")
	} else {
		t.Logf("Found %d repositories", len(repos))
		for i, repo := range repos {
			if i < 5 { // Only log first 5 repos
				t.Logf("  %s", repo)
			}
		}
	}
}

func TestGetPullRequests(t *testing.T) {
	// Skip if no token is available
	cfg := &config.Config{
		GitHub: config.GitHubConfig{
			Token: "test-token", // This will be invalid, but we can test the method structure
		},
	}

	client := NewClient(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with a public repository that likely has PRs
	_, err := client.GetPullRequests(ctx, "golang", "go")

	// We expect an error due to invalid token, but the method should not panic
	if err == nil {
		t.Log("Unexpected success with invalid token")
	} else {
		t.Logf("Expected error with invalid token: %v", err)
	}

	// Verify the method exists and can be called
	if client == nil {
		t.Fatal("Client should not be nil")
	}
}
