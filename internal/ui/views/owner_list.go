package views

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/will-wright-eng/gh-nav/internal/constants"
)

// OwnerListModel represents the owner selection view
type OwnerListModel struct {
	BaseView
	owners     []string
	repoGroups map[string][]string
}

// NewOwnerList creates a new owner list view
func NewOwnerList(pageSize int) *OwnerListModel {
	return &OwnerListModel{
		BaseView:   NewBaseView(pageSize),
		owners:     []string{},
		repoGroups: make(map[string][]string),
	}
}

// SetData sets the repository data for the view
func (o *OwnerListModel) SetData(repoGroups map[string][]string) {
	o.repoGroups = repoGroups
	o.owners = o.getOwners()
}

// getOwners returns a sorted list of owners
func (o *OwnerListModel) getOwners() []string {
	var owners []string
	for owner := range o.repoGroups {
		owners = append(owners, owner)
	}
	sort.Strings(owners)
	return owners
}

// GetVisibleOwners returns the owners visible on the current page
func (o *OwnerListModel) GetVisibleOwners() []string {
	start, end := o.GetVisibleRange(len(o.owners))
	if start >= len(o.owners) {
		return []string{}
	}
	return o.owners[start:end]
}

// GetSelectedOwner returns the currently selected owner
func (o *OwnerListModel) GetSelectedOwner() string {
	visibleOwners := o.GetVisibleOwners()
	if o.cursor < len(visibleOwners) {
		return visibleOwners[o.cursor]
	}
	return ""
}

// GetRepoCount returns the number of repositories for an owner
func (o *OwnerListModel) GetRepoCount(owner string) int {
	if repos, exists := o.repoGroups[owner]; exists {
		return len(repos)
	}
	return 0
}

// Update handles messages and updates the view
func (o *OwnerListModel) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			o.MoveCursorUp(len(o.GetVisibleOwners()))
		case "down", "j":
			o.MoveCursorDown(len(o.GetVisibleOwners()))
		case "left", "h":
			o.PreviousPage()
		case "right", "l":
			o.NextPage(len(o.owners))
		case "g":
			o.GoToFirstPage()
		case "G":
			o.GoToLastPage(len(o.owners))
		}
	}
	return o, nil
}

// View renders the owner list
func (o *OwnerListModel) View() string {
	if o.width == 0 {
		return "Loading..."
	}

	list := ""
	visibleOwners := o.GetVisibleOwners()

	for i, owner := range visibleOwners {
		cursor := " "
		if o.cursor == i {
			cursor = ">"
		}

		style := lipgloss.NewStyle().MarginLeft(2)
		if o.cursor == i {
			style = style.Foreground(lipgloss.Color("#00FF00"))
		}

		repoCount := o.GetRepoCount(owner)
		list += style.Render(fmt.Sprintf("%s %s %s (%d repos)", cursor, constants.IconOrganization, owner, repoCount)) + "\n"
	}

	return list
}

// GetPageInfo returns pagination information
func (o *OwnerListModel) GetPageInfo() string {
	return o.BaseView.GetPageInfo(len(o.owners), "organizations")
}
