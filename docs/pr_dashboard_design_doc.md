# GitHub PR Dashboard - Technical Design Document

## Overview
A terminal-based dashboard for GitHub pull requests using Bubbletea TUI framework with a computed views + cache architecture. Focus on hierarchical navigation: Repos → PRs → Details.

## Directory Structure

```
pr-dashboard/
├── cmd/
│   └── main.go                    # Entry point, CLI flags
├── internal/
│   ├── api/                       # GitHub API client layer
│   │   ├── client.go              # HTTP client with rate limiting
│   │   ├── notifications.go       # Notifications endpoint
│   │   ├── pulls.go              # Pull requests endpoint
│   │   └── reviews.go            # Reviews and comments endpoint
│   ├── cache/                     # Data storage and caching
│   │   ├── store.go              # Cache interface and implementations
│   │   ├── memory.go             # In-memory cache (dev/testing)
│   │   └── sqlite.go             # SQLite cache (production)
│   ├── models/                    # Data models
│   │   ├── repo.go               # Repository model
│   │   ├── pr.go                 # Pull request model
│   │   ├── comment.go            # Comments and reviews
│   │   └── views.go              # Computed view models
│   ├── services/                  # Business logic
│   │   ├── sync.go               # Data synchronization
│   │   ├── aggregator.go         # View computation
│   │   └── filter.go             # Data filtering logic
│   └── ui/                        # TUI components
│       ├── app.go                # Main application state
│       ├── views/                # View components
│       │   ├── repos.go          # Repository list view
│       │   ├── prs.go            # PR list view
│       │   └── details.go        # PR details view
│       └── components/           # Reusable UI components
│           ├── table.go          # Data table component
│           ├── stats.go          # Statistics display
│           └── thread.go         # Conversation thread display
├── pkg/
│   └── config/                   # Configuration management
│       ├── config.go
│       └── auth.go               # GitHub token handling
└── README.md
```

## Data Models

### Core Models
```go
type Repository struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    FullName    string    `json:"full_name"`
    LastUpdated time.Time `json:"last_updated"`
}

type PullRequest struct {
    ID          int64     `json:"id"`
    Number      int       `json:"number"`
    Title       string    `json:"title"`
    State       string    `json:"state"` // open, closed, merged
    RepoID      int64     `json:"repo_id"`
    Author      string    `json:"author"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Additions   int       `json:"additions"`
    Deletions   int       `json:"deletions"`
}

type Comment struct {
    ID        int64     `json:"id"`
    PRID      int64     `json:"pr_id"`
    ThreadID  string    `json:"thread_id"`  // For grouping related comments
    Author    string    `json:"author"`
    Body      string    `json:"body"`
    Type      string    `json:"type"`       // review, comment, review_comment
    FilePath  string    `json:"file_path"`  // For review comments
    Line      int       `json:"line"`       // For review comments
    Resolved  bool      `json:"resolved"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Computed Views
```go
type RepoWeeklyStats struct {
    RepoID        int64 `json:"repo_id"`
    RepoName      string `json:"repo_name"`
    Week          time.Time `json:"week"`
    PRsTotal      int `json:"prs_total"`
    PRsMerged     int `json:"prs_merged"`
    PRsOpen       int `json:"prs_open"`
    ThreadsTotal  int `json:"threads_total"`
    ThreadsActive int `json:"threads_active"`
    LastUpdated   time.Time `json:"last_updated"`
}

type PRDetailView struct {
    PullRequest   `json:",inline"`
    RepoName      string      `json:"repo_name"`
    CIStatus      string      `json:"ci_status"`
    ReviewStatus  ReviewState `json:"review_status"`
    Threads       []Thread    `json:"threads"`
    Timeline      []Event     `json:"timeline"`
}

type Thread struct {
    ID        string    `json:"id"`
    FilePath  string    `json:"file_path"`
    StartLine int       `json:"start_line"`
    Resolved  bool      `json:"resolved"`
    Comments  []Comment `json:"comments"`
}
```

## Data Access Patterns

### 1. Cache Layer Interface
```go
type CacheStore interface {
    // Repository operations
    GetRepos() ([]Repository, error)
    SetRepos([]Repository) error

    // Pull request operations
    GetPRsForRepo(repoID int64, since time.Time) ([]PullRequest, error)
    SetPRsForRepo(repoID int64, prs []PullRequest) error

    // Comment operations
    GetCommentsForPR(prID int64) ([]Comment, error)
    SetCommentsForPR(prID int64, comments []Comment) error

    // Computed views
    GetWeeklyStats(week time.Time) ([]RepoWeeklyStats, error)
    SetWeeklyStats(week time.Time, stats []RepoWeeklyStats) error

    GetPRDetails(prID int64) (*PRDetailView, error)
    SetPRDetails(prID int64, details *PRDetailView) error

    // Cache management
    InvalidateRepo(repoID int64) error
    GetLastSync() (time.Time, error)
    SetLastSync(time.Time) error
}
```

### 2. Synchronization Service
```go
type SyncService struct {
    client    *api.Client
    cache     CacheStore
    semaphore chan struct{} // Rate limiting
}

func (s *SyncService) SyncRepositories() error {
    // 1. Fetch user's repos from GitHub API
    // 2. Update cache with repo list
    // 3. Trigger PR sync for each repo
}

func (s *SyncService) SyncRepoData(repoID int64) error {
    // 1. Fetch PRs from last week
    // 2. For each PR, fetch comments/reviews
    // 3. Update cache atomically
    // 4. Trigger view computation
}

func (s *SyncService) ComputeViews() error {
    // 1. Generate weekly statistics
    // 2. Build PR detail views
    // 3. Update computed view cache
}
```

### 3. UI State Management (Bubbletea)
```go
type AppModel struct {
    currentView ViewType

    // View-specific state
    repoList    RepoListModel
    prList      PRListModel
    prDetails   PRDetailsModel

    // Shared state
    cache       CacheStore
    syncService *SyncService

    // Navigation
    viewStack   []ViewState
}

type ViewType int
const (
    RepoListView ViewType = iota
    PRListView
    PRDetailsView
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case SyncCompleteMsg:
        return m.handleSyncComplete(msg)
    }

    // Delegate to current view
    switch m.currentView {
    case RepoListView:
        return m.updateRepoList(msg)
    case PRListView:
        return m.updatePRList(msg)
    case PRDetailsView:
        return m.updatePRDetails(msg)
    }

    return m, nil
}
```

## Data Flow

### Startup Sequence
1. Load configuration (GitHub token, cache location)
2. Initialize cache store (SQLite or memory)
3. Check last sync timestamp
4. If stale (>5min), trigger background sync
5. Load cached data for initial view
6. Start TUI with repo list

### Background Sync
1. Rate-limited GitHub API calls (100 requests/burst)
2. Incremental sync (only changed data since last sync)
3. Atomic cache updates per repository
4. View recomputation after data updates
5. UI refresh notifications

### Navigation Flow
1. **Repo List** → Load weekly stats from cache
2. **Select Repo** → Load PR list for selected repo
3. **Select PR** → Load detailed PR view with threads
4. **Back Navigation** → Pop view stack, restore previous state

## Caching Strategy

### Cache Keys
- `repos:list` → Repository list
- `prs:{repo_id}:{week}` → PRs for repo in specific week
- `comments:{pr_id}` → Comments/reviews for PR
- `stats:{week}` → Weekly aggregated statistics
- `details:{pr_id}` → Computed PR detail view

### TTL Strategy
- Repository list: 1 hour
- PR data: 5 minutes
- Comments: 2 minutes
- Computed views: Invalidate on base data change

### Memory Management
- LRU eviction for detailed views
- Keep current navigation path in memory
- Background cleanup of old week data

## Configuration

```yaml
github:
  token: ${GITHUB_TOKEN}
  base_url: "https://api.github.com"

cache:
  type: "sqlite"  # memory, sqlite
  path: "~/.config/pr-dashboard/cache.db"
  max_size: "100MB"

sync:
  interval: "5m"
  rate_limit: 100  # requests per hour burst

ui:
  theme: "dark"
  refresh_rate: "1s"
```

This architecture provides clear separation of concerns, efficient caching for GitHub's rate limits, and a responsive TUI experience with Bubbletea's reactive model.
