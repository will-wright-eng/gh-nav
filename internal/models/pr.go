package models

import (
	"time"

	"github.com/google/go-github/v58/github"
)

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID        int64     `json:"id"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"` // open, closed, merged
	RepoID    int64     `json:"repo_id"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Additions int       `json:"additions"`
	Deletions int       `json:"deletions"`

	// Additional fields for list view
	Labels       []string `json:"labels"`
	Assignees    []string `json:"assignees"`
	Reviewers    []string `json:"reviewers"`
	ReviewStatus string   `json:"review_status"` // approved, changes_requested, pending
	IsDraft      bool     `json:"is_draft"`
	Mergeable    *bool    `json:"mergeable"`
	Comments     int      `json:"comments"`
	Commits      int      `json:"commits"`
}

// FromGitHubPR converts a GitHub PR to our model
func FromGitHubPR(pr *github.PullRequest, repoName string) *PullRequest {
	// Safe dereferencing helper
	safeString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	safeInt := func(i *int) int {
		if i == nil {
			return 0
		}
		return *i
	}

	safeInt64 := func(i *int64) int64 {
		if i == nil {
			return 0
		}
		return *i
	}

	safeBool := func(b *bool) bool {
		if b == nil {
			return false
		}
		return *b
	}

	safeTime := func(t *github.Timestamp) time.Time {
		if t == nil {
			return time.Time{}
		}
		return t.Time
	}

	// Extract labels safely
	labels := make([]string, 0, len(pr.Labels))
	for _, label := range pr.Labels {
		if label != nil && label.Name != nil {
			labels = append(labels, *label.Name)
		}
	}

	// Extract assignees safely
	assignees := make([]string, 0, len(pr.Assignees))
	for _, assignee := range pr.Assignees {
		if assignee != nil && assignee.Login != nil {
			assignees = append(assignees, *assignee.Login)
		}
	}

	// Extract reviewers safely
	reviewers := make([]string, 0, len(pr.RequestedReviewers))
	for _, reviewer := range pr.RequestedReviewers {
		if reviewer != nil && reviewer.Login != nil {
			reviewers = append(reviewers, *reviewer.Login)
		}
	}

	// Determine review status based on mergeable state and reviews
	reviewStatus := "pending"
	if pr.Mergeable != nil && *pr.Mergeable {
		reviewStatus = "approved"
	} else if pr.Mergeable != nil && !*pr.Mergeable {
		reviewStatus = "changes_requested"
	}

	// Extract author safely
	author := ""
	if pr.User != nil && pr.User.Login != nil {
		author = *pr.User.Login
	}

	return &PullRequest{
		ID:           safeInt64(pr.ID),
		Number:       safeInt(pr.Number),
		Title:        safeString(pr.Title),
		State:        safeString(pr.State),
		RepoID:       0, // Will be set by caller if needed
		Author:       author,
		CreatedAt:    safeTime(pr.CreatedAt),
		UpdatedAt:    safeTime(pr.UpdatedAt),
		Additions:    safeInt(pr.Additions),
		Deletions:    safeInt(pr.Deletions),
		Labels:       labels,
		Assignees:    assignees,
		Reviewers:    reviewers,
		ReviewStatus: reviewStatus,
		IsDraft:      safeBool(pr.Draft),
		Mergeable:    pr.Mergeable,
		Comments:     safeInt(pr.Comments),
		Commits:      safeInt(pr.Commits),
	}
}
