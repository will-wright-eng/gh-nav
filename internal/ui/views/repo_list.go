package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/will/ghprs/internal/constants"
)

// RepoListModel represents the repository selection view
type RepoListModel struct {
	BaseView
	owner string
	repos []string
}

// NewRepoList creates a new repository list view
func NewRepoList(pageSize int) *RepoListModel {
	return &RepoListModel{
		BaseView: NewBaseView(pageSize),
		owner:    "",
		repos:    []string{},
	}
}

// SetData sets the repository data for the view
func (r *RepoListModel) SetData(owner string, repos []string) {
	r.owner = owner
	r.repos = repos
	r.page = 0
	r.cursor = 0
}

// GetVisibleRepos returns the repositories visible on the current page
func (r *RepoListModel) GetVisibleRepos() []string {
	start, end := r.GetVisibleRange(len(r.repos))
	if start >= len(r.repos) {
		return []string{}
	}
	return r.repos[start:end]
}

// GetSelectedRepo returns the currently selected repository
func (r *RepoListModel) GetSelectedRepo() string {
	visibleRepos := r.GetVisibleRepos()
	if r.cursor < len(visibleRepos) {
		return visibleRepos[r.cursor]
	}
	return ""
}

// GetRepoName extracts the repository name from the full name
func (r *RepoListModel) GetRepoName(fullName string) string {
	parts := strings.Split(fullName, "/")
	if len(parts) == 2 {
		return parts[1]
	}
	return fullName
}

// Update handles messages and updates the view
func (r *RepoListModel) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			r.MoveCursorUp(len(r.GetVisibleRepos()))
		case "down", "j":
			r.MoveCursorDown(len(r.GetVisibleRepos()))
		case "left", "h":
			r.PreviousPage()
		case "right", "l":
			r.NextPage(len(r.repos))
		case "g":
			r.GoToFirstPage()
		case "G":
			r.GoToLastPage(len(r.repos))
		}
	}
	return r, nil
}

// View renders the repository list
func (r *RepoListModel) View() string {
	if r.width == 0 {
		return "Loading..."
	}

	list := ""
	visibleRepos := r.GetVisibleRepos()

	for i, repo := range visibleRepos {
		cursor := " "
		if r.cursor == i {
			cursor = ">"
		}

		style := lipgloss.NewStyle().MarginLeft(2)
		if r.cursor == i {
			style = style.Foreground(lipgloss.Color("#00FF00"))
		}

		repoName := r.GetRepoName(repo)
		list += style.Render(fmt.Sprintf("%s %s %s", cursor, constants.IconRepository, repoName)) + "\n"
	}

	return list
}

// GetPageInfo returns pagination information
func (r *RepoListModel) GetPageInfo() string {
	return r.BaseView.GetPageInfo(len(r.repos), "repositories")
}

// GetOwner returns the current owner
func (r *RepoListModel) GetOwner() string {
	return r.owner
}
