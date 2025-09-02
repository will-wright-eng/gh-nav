package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/will/ghprs/internal/constants"
)

// ColorPalette defines the color scheme
type ColorPalette struct {
	Primary    string
	Secondary  string
	Success    string
	Warning    string
	Error      string
	Info       string
	Muted      string
	Background string
	Text       string
}

// StylePalette defines common styles
type StylePalette struct {
	Title      lipgloss.Style
	Status     lipgloss.Style
	Error      lipgloss.Style
	Success    lipgloss.Style
	Warning    lipgloss.Style
	Info       lipgloss.Style
	Help       lipgloss.Style
	Debug      lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Margin     lipgloss.Style
}

// IconSet defines the icons used throughout the UI
type IconSet struct {
	Organization string
	Repository   string
	PullRequest  string
	Open         string
	Draft        string
	Approved     string
	Changes      string
	Loading      string
	Error        string
	Success      string
	Info         string
}

// Theme represents the complete UI theme
type Theme struct {
	Colors ColorPalette
	Styles StylePalette
	Icons  IconSet
}

// DefaultTheme provides a dark theme optimized for terminal use
var DefaultTheme = Theme{
	Colors: ColorPalette{
		Primary:    constants.ColorPrimary,   // Bright green
		Secondary:  constants.ColorSecondary, // Gray
		Success:    constants.ColorSuccess,   // Green
		Warning:    constants.ColorWarning,   // Yellow
		Error:      constants.ColorError,     // Red
		Info:       constants.ColorInfo,      // Cyan
		Muted:      constants.ColorMuted,     // Dark gray
		Background: "#000000",                // Black
		Text:       constants.ColorText,      // Light gray
	},
	Styles: StylePalette{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			MarginLeft(2),
		Status: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			MarginLeft(2),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			MarginLeft(2),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			MarginLeft(2),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			MarginLeft(2),
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			MarginLeft(2),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginTop(2).
			MarginLeft(2),
		Debug: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginTop(1).
			MarginLeft(2),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			MarginLeft(2),
		Unselected: lipgloss.NewStyle().
			MarginLeft(2),
		Margin: lipgloss.NewStyle().
			MarginLeft(2),
	},
	Icons: IconSet{
		Organization: constants.IconOrganization,
		Repository:   constants.IconRepository,
		PullRequest:  constants.IconPullRequest,
		Open:         constants.IconOpen,
		Draft:        constants.IconDraft,
		Approved:     constants.IconApproved,
		Changes:      constants.IconChanges,
		Loading:      constants.IconLoading,
		Error:        constants.IconError,
		Success:      constants.IconSuccess,
		Info:         constants.IconInfo,
	},
}

// GetStatusStyle returns the appropriate style for a status
func (t *Theme) GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "loading":
		return t.Styles.Warning
	case "error":
		return t.Styles.Error
	case "success":
		return t.Styles.Success
	default:
		return t.Styles.Info
	}
}

// GetStatusIcon returns the appropriate icon for a status
func (t *Theme) GetStatusIcon(status string) string {
	switch status {
	case "loading":
		return t.Icons.Loading
	case "error":
		return t.Icons.Error
	case "success":
		return t.Icons.Success
	default:
		return t.Icons.Info
	}
}

// GetPRStatusIcon returns the appropriate icon for a PR status
func (t *Theme) GetPRStatusIcon(isDraft bool, reviewStatus string) string {
	if isDraft {
		return t.Icons.Draft
	}

	switch reviewStatus {
	case "approved":
		return t.Icons.Approved
	case "changes_requested":
		return t.Icons.Changes
	default:
		return t.Icons.Open
	}
}
