package constants

import "time"

// Navigation constants
const (
	KeyUp      = "up"
	KeyDown    = "down"
	KeyLeft    = "left"
	KeyRight   = "right"
	KeyEnter   = "enter"
	KeyBack    = "backspace"
	KeyBackAlt = "b"
	KeyQuit    = "q"
	KeyQuitAlt = "ctrl+c"
	KeyDebug   = "d"
	KeyReload  = "r"
	KeyFirst   = "g"
	KeyLast    = "G"
)

// View mode constants
const (
	ViewModeOwnerSelection = iota
	ViewModeRepoSelection
	ViewModePRList
)

// UI constants
const (
	DefaultPageSize = 10
	MaxTitleLength  = 60
	DefaultMargin   = 2
	RefreshInterval = time.Second
	DefaultTimeout  = 30 * time.Second
)

// Status constants
const (
	StatusLoading = "loading"
	StatusError   = "error"
	StatusSuccess = "success"
	StatusInfo    = "info"
)

// PR status constants
const (
	PRStatusOpen             = "open"
	PRStatusDraft            = "draft"
	PRStatusApproved         = "approved"
	PRStatusChangesRequested = "changes_requested"
	PRStatusPending          = "pending"
)

// Icon constants
const (
	IconOrganization = "📁"
	IconRepository   = "📦"
	IconPullRequest  = "🔀"
	IconOpen         = "🔵"
	IconDraft        = "🟡"
	IconApproved     = "🟢"
	IconChanges      = "🔴"
	IconLoading      = "🔄"
	IconError        = "❌"
	IconSuccess      = "✅"
	IconInfo         = "ℹ️"
)

// Color constants
const (
	ColorPrimary   = "#00FF00"
	ColorSecondary = "#666666"
	ColorSuccess   = "#00FF00"
	ColorWarning   = "#FFFF00"
	ColorError     = "#FF0000"
	ColorInfo      = "#00FFFF"
	ColorMuted     = "#888888"
	ColorText      = "#FAFAFA"
)

// Help text constants
const (
	HelpNavigation = "↑/↓: Navigate • ←/→: Page • g/G: First/Last • Enter: Select • b: Back • d: Debug • r: Reload • q: Quit"
	HelpLoading    = "Loading..."
	HelpNoData     = "No data found"
)

// Error messages
const (
	ErrLoadingRepos   = "Error loading repositories"
	ErrLoadingPRs     = "Error loading pull requests"
	ErrInvalidToken   = "Invalid GitHub token"
	ErrNetworkTimeout = "Network timeout"
	ErrUnauthorized   = "Unauthorized access"
)

// Success messages
const (
	MsgReposLoaded   = "Repositories loaded successfully"
	MsgPRsLoaded     = "Pull requests loaded successfully"
	MsgDataRefreshed = "Data refreshed successfully"
)
