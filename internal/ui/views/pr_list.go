package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/will/ghprs/internal/constants"
	"github.com/will/ghprs/internal/models"
)

// PRListModel represents the pull request list view
type PRListModel struct {
	BaseView
	repoName string
	prs      []*models.PullRequest
}

// NewPRList creates a new pull request list view
func NewPRList(pageSize int) *PRListModel {
	return &PRListModel{
		BaseView: NewBaseView(pageSize),
		repoName: "",
		prs:      []*models.PullRequest{},
	}
}

// SetData sets the pull request data for the view
func (p *PRListModel) SetData(repoName string, prs []*models.PullRequest) {
	p.repoName = repoName
	p.prs = prs
	p.page = 0
	p.cursor = 0
}

// GetVisiblePRs returns the pull requests visible on the current page
func (p *PRListModel) GetVisiblePRs() []*models.PullRequest {
	start, end := p.GetVisibleRange(len(p.prs))
	if start >= len(p.prs) {
		return []*models.PullRequest{}
	}
	return p.prs[start:end]
}

// GetSelectedPR returns the currently selected pull request
func (p *PRListModel) GetSelectedPR() *models.PullRequest {
	visiblePRs := p.GetVisiblePRs()
	if p.cursor < len(visiblePRs) {
		return visiblePRs[p.cursor]
	}
	return nil
}

// GetStatusIcon returns the appropriate status icon for a PR
func (p *PRListModel) GetStatusIcon(pr *models.PullRequest) string {
	if pr.IsDraft {
		return constants.IconDraft // draft
	} else if pr.ReviewStatus == "approved" {
		return constants.IconApproved // approved
	} else if pr.ReviewStatus == "changes_requested" {
		return constants.IconChanges // changes requested
	}
	return constants.IconOpen // open
}

// TruncateTitle truncates the PR title if it's too long
func (p *PRListModel) TruncateTitle(title string, maxLength int) string {
	if len(title) <= maxLength {
		return title
	}
	return title[:maxLength-3] + "..."
}

// Update handles messages and updates the view
func (p *PRListModel) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			p.MoveCursorUp(len(p.GetVisiblePRs()))
		case "down", "j":
			p.MoveCursorDown(len(p.GetVisiblePRs()))
		case "left", "h":
			p.PreviousPage()
		case "right", "l":
			p.NextPage(len(p.prs))
		case "g":
			p.GoToFirstPage()
		case "G":
			p.GoToLastPage(len(p.prs))
		}
	}
	return p, nil
}

// View renders the pull request list
func (p *PRListModel) View() string {
	if p.width == 0 {
		return "Loading..."
	}

	list := ""
	visiblePRs := p.GetVisiblePRs()

	for i, pr := range visiblePRs {
		cursor := " "
		if p.cursor == i {
			cursor = ">"
		}

		style := lipgloss.NewStyle().MarginLeft(2)
		if p.cursor == i {
			style = style.Foreground(lipgloss.Color("#00FF00"))
		}

		statusIcon := p.GetStatusIcon(pr)
		title := p.TruncateTitle(pr.Title, constants.MaxTitleLength)

		list += style.Render(fmt.Sprintf("%s %s #%d %s", cursor, statusIcon, pr.Number, title)) + "\n"
	}

	return list
}

// GetPageInfo returns pagination information
func (p *PRListModel) GetPageInfo() string {
	return p.BaseView.GetPageInfo(len(p.prs), "pull requests")
}

// GetRepoName returns the current repository name
func (p *PRListModel) GetRepoName() string {
	return p.repoName
}

// GetPRCount returns the total number of pull requests
func (p *PRListModel) GetPRCount() int {
	return len(p.prs)
}
