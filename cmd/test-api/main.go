package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/will/ghprs/internal/api"
	"github.com/will/ghprs/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.GitHub.Token == "" {
		log.Fatal("No GitHub token available. Please run 'gh auth login' first.")
	}

	fmt.Printf("Using token: %s\n", maskToken(cfg.GitHub.Token))

	// Create API client
	client := api.NewClient(cfg)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Fetching repositories...")

	// Fetch repositories
	repos, err := client.GetUserRepositories(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch repositories: %v", err)
	}

	fmt.Printf("Found %d repositories:\n", len(repos))

	// Group repositories by organization/user
	repoGroups := make(map[string][]string)
	for _, repo := range repos {
		parts := strings.Split(repo, "/")
		if len(parts) == 2 {
			owner := parts[0]
			repoGroups[owner] = append(repoGroups[owner], repo)
		}
	}

	// Display grouped repositories
	for owner, ownerRepos := range repoGroups {
		fmt.Printf("\n%s (%d repos):\n", owner, len(ownerRepos))
		for i, repo := range ownerRepos {
			if i < 5 { // Show first 5 repos per org
				fmt.Printf("  - %s\n", repo)
			} else if i == 5 {
				fmt.Printf("  ... and %d more\n", len(ownerRepos)-5)
				break
			}
		}
	}
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
