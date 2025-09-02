package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/will/ghprs/pkg/config"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	config *config.Config
}

// NewClient creates a new GitHub API client
func NewClient(cfg *config.Config) *Client {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	client := github.NewClient(httpClient)
	if cfg.GitHub.Token != "" {
		client = client.WithAuthToken(cfg.GitHub.Token)
	}

	return &Client{
		client: client,
		config: cfg,
	}
}

// GetUserRepositories fetches repositories for the authenticated user and their organizations
func (c *Client) GetUserRepositories(ctx context.Context) ([]string, error) {
	// Get the authenticated user first
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}

	var allRepos []*github.Repository

	// 1. Fetch user's personal repositories
	userRepos, err := c.getRepositoriesForUser(ctx, *user.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to get user repositories: %w", err)
	}
	allRepos = append(allRepos, userRepos...)

	// 2. Fetch user's organizations
	orgs, err := c.getUserOrganizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	// 3. Fetch repositories from each organization
	for _, org := range orgs {
		orgRepos, err := c.getRepositoriesForOrganization(ctx, *org.Login)
		if err != nil {
			// Log error but continue with other orgs
			fmt.Printf("Warning: failed to get repos for org %s: %v\n", *org.Login, err)
			continue
		}
		allRepos = append(allRepos, orgRepos...)
	}

	// Convert to string slice
	var repoNames []string
	for _, repo := range allRepos {
		repoNames = append(repoNames, *repo.FullName)
	}

	return repoNames, nil
}

// getRepositoriesForUser fetches repositories for a specific user
func (c *Client) getRepositoriesForUser(ctx context.Context, username string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByUserOptions{
		Type:        "all", // all, owner, public, private, member
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := c.client.Repositories.ListByUser(ctx, username, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories for user %s: %w", username, err)
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

// getUserOrganizations fetches organizations the user belongs to
func (c *Client) getUserOrganizations(ctx context.Context) ([]*github.Organization, error) {
	opt := &github.ListOptions{PerPage: 100}
	var allOrgs []*github.Organization

	for {
		orgs, resp, err := c.client.Organizations.List(ctx, "", opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list organizations: %w", err)
		}

		allOrgs = append(allOrgs, orgs...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allOrgs, nil
}

// getRepositoriesForOrganization fetches repositories for a specific organization
func (c *Client) getRepositoriesForOrganization(ctx context.Context, orgName string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		Type:        "all", // all, public, private, forks, sources, member
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := c.client.Repositories.ListByOrg(ctx, orgName, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories for org %s: %w", orgName, err)
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

// GetRepository fetches a specific repository
func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*github.Repository, error) {
	repository, _, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository %s/%s: %w", owner, repo, err)
	}
	return repository, nil
}

// GetPullRequests fetches pull requests for a specific repository
func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		State:       "open", // open, closed, all
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allPRs []*github.PullRequest
	for {
		prs, resp, err := c.client.PullRequests.List(ctx, owner, repo, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests for %s/%s: %w", owner, repo, err)
		}

		allPRs = append(allPRs, prs...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allPRs, nil
}
