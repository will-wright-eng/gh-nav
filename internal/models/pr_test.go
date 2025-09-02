package models

import (
	"testing"
	"time"

	"github.com/google/go-github/v58/github"
)

func TestFromGitHubPR(t *testing.T) {
	// Create a mock GitHub PR
	now := time.Now()
	trueVal := true
	falseVal := false
	prID := int64(123)
	prNumber := 456
	prTitle := "Test PR"
	prState := "open"
	prAuthor := "testuser"
	prAdditions := 10
	prDeletions := 5
	prComments := 3
	prCommits := 2

	githubPR := &github.PullRequest{
		ID:                 &prID,
		Number:             &prNumber,
		Title:              &prTitle,
		State:              &prState,
		User:               &github.User{Login: &prAuthor},
		CreatedAt:          &github.Timestamp{Time: now},
		UpdatedAt:          &github.Timestamp{Time: now},
		Additions:          &prAdditions,
		Deletions:          &prDeletions,
		Comments:           &prComments,
		Commits:            &prCommits,
		Draft:              &falseVal,
		Mergeable:          &trueVal,
		Labels:             []*github.Label{},
		Assignees:          []*github.User{},
		RequestedReviewers: []*github.User{},
	}

	// Convert to our model
	pr := FromGitHubPR(githubPR, "test/repo")

	// Verify the conversion
	if pr.ID != prID {
		t.Errorf("Expected ID %d, got %d", prID, pr.ID)
	}
	if pr.Number != prNumber {
		t.Errorf("Expected Number %d, got %d", prNumber, pr.Number)
	}
	if pr.Title != prTitle {
		t.Errorf("Expected Title %s, got %s", prTitle, pr.Title)
	}
	if pr.State != prState {
		t.Errorf("Expected State %s, got %s", prState, pr.State)
	}
	if pr.Author != prAuthor {
		t.Errorf("Expected Author %s, got %s", prAuthor, pr.Author)
	}
	if pr.Additions != prAdditions {
		t.Errorf("Expected Additions %d, got %d", prAdditions, pr.Additions)
	}
	if pr.Deletions != prDeletions {
		t.Errorf("Expected Deletions %d, got %d", prDeletions, pr.Deletions)
	}
	if pr.Comments != prComments {
		t.Errorf("Expected Comments %d, got %d", prComments, pr.Comments)
	}
	if pr.Commits != prCommits {
		t.Errorf("Expected Commits %d, got %d", prCommits, pr.Commits)
	}
	if pr.IsDraft != false {
		t.Errorf("Expected IsDraft false, got %v", pr.IsDraft)
	}
	if pr.ReviewStatus != "approved" {
		t.Errorf("Expected ReviewStatus 'approved', got %s", pr.ReviewStatus)
	}

	// Test with nil values to ensure safe dereferencing
	nilPR := &github.PullRequest{}
	nilResult := FromGitHubPR(nilPR, "test/repo")

	// Should not panic and should have default values
	if nilResult.ID != 0 {
		t.Errorf("Expected ID 0 for nil PR, got %d", nilResult.ID)
	}
	if nilResult.Title != "" {
		t.Errorf("Expected empty title for nil PR, got %s", nilResult.Title)
	}
	if nilResult.Author != "" {
		t.Errorf("Expected empty author for nil PR, got %s", nilResult.Author)
	}
}
