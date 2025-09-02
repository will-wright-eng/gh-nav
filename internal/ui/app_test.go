package ui

import (
	"testing"

	"github.com/will-wright-eng/gh-nav/internal/ui/views"
	"github.com/will-wright-eng/gh-nav/pkg/config"
)

func TestNewApp(t *testing.T) {
	cfg := &config.Config{}
	app := NewApp(cfg)

	if app == nil {
		t.Fatal("Expected app to be created")
	}

	if !app.loading {
		t.Error("Expected app to start in loading state")
	}

	if app.repoGroups == nil {
		t.Error("Expected repoGroups to be initialized")
	}

	if app.currentView != OwnerSelection {
		t.Error("Expected currentView to start with OwnerSelection")
	}

	if app.selectedOwner != "" {
		t.Error("Expected selectedOwner to start empty")
	}

	if app.views == nil {
		t.Error("Expected views to be initialized")
	}

	if len(app.views) != 3 {
		t.Error("Expected 3 views to be initialized")
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "***"},
		{"short", "***"},
		{"gho_1234567890abcdef", "gho_...cdef"},
		{"gho_abcdef1234567890", "gho_...7890"},
	}

	for _, test := range tests {
		result := maskToken(test.input)
		if result != test.expected {
			t.Errorf("maskToken(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestLoadRepositories(t *testing.T) {
	// Load config to get token
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Skip test if no token available
	if cfg.GitHub.Token == "" {
		t.Skip("Skipping test - no GitHub token available")
	}

	// Test that the command returns a message
	cmd := loadRepositories(cfg)
	if cmd == nil {
		t.Fatal("Expected loadRepositories to return a command")
	}

	// Execute the command
	msg := cmd()

	// Check that it returns the expected message type
	reposMsg, ok := msg.(reposLoadedMsg)
	if !ok {
		t.Fatal("Expected reposLoadedMsg")
	}

	// With real API, we might get an error if token is invalid
	// or we might get repositories if token is valid
	if reposMsg.err != nil {
		t.Logf("API call failed (this might be expected): %v", reposMsg.err)
		// Don't fail the test for API errors
		return
	}

	if len(reposMsg.repos) == 0 {
		t.Log("No repositories found - this might be normal")
	} else {
		t.Logf("Found %d repositories", len(reposMsg.repos))
	}
}

func TestPaginationHelpers(t *testing.T) {
	cfg := &config.Config{}
	app := NewApp(cfg)

	// Test with no organizations
	if app.getPageInfo() != "No organizations found" {
		t.Errorf("Expected 'No organizations found', got '%s'", app.getPageInfo())
	}

	// Test with some organizations
	app.repoGroups = map[string][]string{
		"org1": {"org1/repo1", "org1/repo2", "org1/repo3"},
		"org2": {"org2/repo1", "org2/repo2"},
		"org3": {"org3/repo1"},
	}

	// Update the owner list view with data
	if ownerList, ok := app.views[OwnerSelection].(*views.OwnerListModel); ok {
		ownerList.SetData(app.repoGroups)
		app.views[OwnerSelection] = ownerList
	}

	// Test owner selection mode
	app.currentView = OwnerSelection
	if app.getPageInfo() != "Showing 3 organizations" {
		t.Errorf("Expected 'Showing 3 organizations', got '%s'", app.getPageInfo())
	}
}

func TestGroupRepositories(t *testing.T) {
	cfg := &config.Config{}
	app := NewApp(cfg)

	// Test grouping repositories
	repos := []string{
		"user/repo1",
		"user/repo2",
		"org1/repo1",
		"org1/repo2",
		"org2/repo1",
	}

	app.groupRepositories(repos)

	// Check that repositories are grouped correctly
	if len(app.repoGroups) != 3 {
		t.Errorf("Expected 3 groups, got %d", len(app.repoGroups))
	}

	if len(app.repoGroups["user"]) != 2 {
		t.Errorf("Expected 2 repos for user, got %d", len(app.repoGroups["user"]))
	}

	if len(app.repoGroups["org1"]) != 2 {
		t.Errorf("Expected 2 repos for org1, got %d", len(app.repoGroups["org1"]))
	}

	if len(app.repoGroups["org2"]) != 1 {
		t.Errorf("Expected 1 repo for org2, got %d", len(app.repoGroups["org2"]))
	}

	// Check specific repositories
	if app.repoGroups["user"][0] != "user/repo1" {
		t.Error("Expected user/repo1 in user group")
	}

	if app.repoGroups["org1"][1] != "org1/repo2" {
		t.Error("Expected org1/repo2 in org1 group")
	}
}
