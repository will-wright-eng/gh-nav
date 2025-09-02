package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/will-wright-eng/gh-nav/internal/api"
	"github.com/will-wright-eng/gh-nav/internal/constants"
	"github.com/will-wright-eng/gh-nav/internal/models"
	"github.com/will-wright-eng/gh-nav/internal/ui/theme"
	"github.com/will-wright-eng/gh-nav/internal/ui/views"
	"github.com/will-wright-eng/gh-nav/pkg/config"
)

// Message types for the UI
type tickMsg time.Time
type reposLoadedMsg struct {
	repos []string
	err   error
}
type prsLoadedMsg struct {
	prs []*models.PullRequest
	err error
}

// ViewMode represents the current view state
type ViewMode int

const (
	OwnerSelection ViewMode = iota
	RepoSelection
	PRList
)

// AppModel represents the main application state
type AppModel struct {
	config *config.Config
	theme  *theme.Theme

	// View management
	currentView ViewMode
	views       map[ViewMode]views.View

	// Repository grouping
	repoGroups map[string][]string // owner -> repos

	// Navigation state
	selectedOwner string
	selectedRepo  string

	// UI state
	width   int
	height  int
	loading bool
	error   string

	// Debug information
	debugMode bool
	debugInfo string
}

// NewApp creates a new application model
func NewApp(cfg *config.Config) *AppModel {
	// Initialize views
	ownerList := views.NewOwnerList(constants.DefaultPageSize)
	repoList := views.NewRepoList(constants.DefaultPageSize)
	prList := views.NewPRList(constants.DefaultPageSize)

	viewsMap := map[ViewMode]views.View{
		OwnerSelection: ownerList,
		RepoSelection:  repoList,
		PRList:         prList,
	}

	return &AppModel{
		config:        cfg,
		theme:         &theme.DefaultTheme,
		currentView:   OwnerSelection,
		views:         viewsMap,
		repoGroups:    make(map[string][]string),
		selectedOwner: "",
		selectedRepo:  "",
		width:         0,
		height:        0,
		loading:       true,
		error:         "",
		debugMode:     false,
		debugInfo:     "Initializing...",
	}
}

// Init initializes the application
func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		loadRepositories(m.config),
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
	)
}

// Update handles messages and updates the model
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case constants.KeyQuitAlt, constants.KeyQuit:
			return m, tea.Quit
		case constants.KeyUp, "k", constants.KeyDown, "j", constants.KeyLeft, "h", constants.KeyRight, "l", constants.KeyFirst, constants.KeyLast:
			// Delegate navigation to current view
			if view, exists := m.views[m.currentView]; exists {
				updatedView, cmd := view.Update(msg)
				m.views[m.currentView] = updatedView
				return m, cmd
			}
		case constants.KeyEnter:
			return m.handleEnterKey()
		case constants.KeyBack, constants.KeyBackAlt:
			return m.handleBackKey()
		case constants.KeyDebug:
			// Toggle debug mode
			m.debugMode = !m.debugMode
		case constants.KeyReload:
			// Reload repositories
			m.loading = true
			m.error = ""
			m.currentView = OwnerSelection
			m.selectedOwner = ""
			m.selectedRepo = ""
			return m, loadRepositories(m.config)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Set size on all views
		for _, view := range m.views {
			view.SetSize(msg.Width, msg.Height)
		}
	case tickMsg:
		// Update debug info every second
		m.debugInfo = fmt.Sprintf("Last update: %s | Token: %s",
			time.Now().Format("15:04:05"),
			maskToken(m.config.GitHub.Token))
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case reposLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.error = msg.err.Error()
		} else {
			m.groupRepositories(msg.repos) // Group repositories by owner
			m.error = ""

			// Update owner list with data
			if ownerList, ok := m.views[OwnerSelection].(*views.OwnerListModel); ok {
				ownerList.SetData(m.repoGroups)
				m.views[OwnerSelection] = ownerList
			}
		}
	case prsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.error = msg.err.Error()
		} else {
			m.error = ""

			// Update PR list with data
			if prList, ok := m.views[PRList].(*views.PRListModel); ok {
				prList.SetData(m.selectedRepo, msg.prs)
				m.views[PRList] = prList
			}
		}
	}

	return m, nil
}

// handleEnterKey handles the enter key press for navigation
func (m *AppModel) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case OwnerSelection:
		if ownerList, ok := m.views[OwnerSelection].(*views.OwnerListModel); ok {
			selectedOwner := ownerList.GetSelectedOwner()
			if selectedOwner != "" {
				m.selectedOwner = selectedOwner
				m.currentView = RepoSelection

				// Update repo list with data
				if repoList, ok := m.views[RepoSelection].(*views.RepoListModel); ok {
					repos := m.getReposForOwner(selectedOwner)
					repoList.SetData(selectedOwner, repos)
					m.views[RepoSelection] = repoList
				}
			}
		}
	case RepoSelection:
		if repoList, ok := m.views[RepoSelection].(*views.RepoListModel); ok {
			selectedRepo := repoList.GetSelectedRepo()
			if selectedRepo != "" {
				m.selectedRepo = selectedRepo
				m.currentView = PRList
				m.loading = true
				m.error = ""
				return m, loadPullRequests(m.config, m.selectedOwner, m.selectedRepo)
			}
		}
	}
	return m, nil
}

// handleBackKey handles the back key press for navigation
func (m *AppModel) handleBackKey() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case RepoSelection:
		m.currentView = OwnerSelection
		m.selectedOwner = ""
		m.selectedRepo = ""
	case PRList:
		m.currentView = RepoSelection
		m.selectedRepo = ""
		// Clear PR list data
		if prList, ok := m.views[PRList].(*views.PRListModel); ok {
			prList.SetData("", []*models.PullRequest{})
			m.views[PRList] = prList
		}
	}
	return m, nil
}

// View renders the UI
func (m AppModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	title := m.theme.Styles.Title.Render("GitHub PR Dashboard")

	// Status line
	status := ""
	if m.loading {
		status = m.theme.Styles.Warning.Render(fmt.Sprintf("%s Loading repositories...", m.theme.Icons.Loading))
	} else if m.error != "" {
		status = m.theme.Styles.Error.Render(fmt.Sprintf("%s %s", m.theme.Icons.Error, m.error))
	} else {
		switch m.currentView {
		case OwnerSelection:
			orgCount := len(m.repoGroups)
			status = m.theme.Styles.Success.Render(fmt.Sprintf("%s %s (%d organizations)",
				m.theme.Icons.Success, m.getPageInfo(), orgCount))
		case RepoSelection:
			status = m.theme.Styles.Success.Render(fmt.Sprintf("%s %s - %s",
				m.theme.Icons.Success, m.getPageInfo(), m.selectedOwner))
		case PRList:
			status = m.theme.Styles.Success.Render(fmt.Sprintf("%s %s - %s",
				m.theme.Icons.Success, m.getPageInfo(), m.selectedRepo))
		}
	}

	// Get content from current view
	list := ""
	if view, exists := m.views[m.currentView]; exists {
		list = view.View()
	}

	help := m.theme.Styles.Help.Render(constants.HelpNavigation)

	// Debug information
	debug := ""
	if m.debugMode {
		debug = m.theme.Styles.Debug.Render(fmt.Sprintf("Debug: %s", m.debugInfo))
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", title, status, list, help, debug)
}

// loadPullRequests fetches pull requests from GitHub API
func loadPullRequests(cfg *config.Config, owner, repo string) tea.Cmd {
	return func() tea.Msg {
		// Create GitHub API client
		client := api.NewClient(cfg)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
		defer cancel()

		// Extract repo name from full name
		parts := strings.Split(repo, "/")
		repoName := repo
		if len(parts) == 2 {
			repoName = parts[1]
		}

		// Fetch pull requests
		githubPRs, err := client.GetPullRequests(ctx, owner, repoName)
		if err != nil {
			return prsLoadedMsg{
				prs: nil,
				err: err,
			}
		}

		// Convert to our model
		var prs []*models.PullRequest
		for _, githubPR := range githubPRs {
			pr := models.FromGitHubPR(githubPR, repo)
			prs = append(prs, pr)
		}

		return prsLoadedMsg{
			prs: prs,
			err: nil,
		}
	}
}

// maskToken masks most of the token for security
func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// getPageInfo returns pagination information from the current view
func (m AppModel) getPageInfo() string {
	if view, exists := m.views[m.currentView]; exists {
		switch m.currentView {
		case OwnerSelection:
			if ownerList, ok := view.(*views.OwnerListModel); ok {
				return ownerList.GetPageInfo()
			}
		case RepoSelection:
			if repoList, ok := view.(*views.RepoListModel); ok {
				return repoList.GetPageInfo()
			}
		case PRList:
			if prList, ok := view.(*views.PRListModel); ok {
				return prList.GetPageInfo()
			}
		}
	}
	return "Loading..."
}

// groupRepositories groups repositories by owner (user/org)
func (m *AppModel) groupRepositories(repos []string) {
	m.repoGroups = make(map[string][]string)
	for _, repo := range repos {
		parts := strings.Split(repo, "/")
		if len(parts) == 2 {
			owner := parts[0]
			m.repoGroups[owner] = append(m.repoGroups[owner], repo)
		}
	}
}

// getReposForOwner returns repositories for a specific owner
func (m AppModel) getReposForOwner(owner string) []string {
	if repos, exists := m.repoGroups[owner]; exists {
		return repos
	}
	return []string{}
}

// loadRepositories fetches repositories from GitHub API
func loadRepositories(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		// Create GitHub API client
		client := api.NewClient(cfg)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultTimeout)
		defer cancel()

		// Fetch repositories
		repos, err := client.GetUserRepositories(ctx)
		if err != nil {
			return reposLoadedMsg{
				repos: nil,
				err:   err,
			}
		}

		return reposLoadedMsg{
			repos: repos,
			err:   nil,
		}
	}
}
