# GitHub PR Dashboard

A terminal-based dashboard for GitHub pull requests built with Go and Bubbletea TUI framework.

## Features (Minimal Prototype)

- Basic TUI interface with repository list
- **Organization support** - shows repositories from all organizations you're a member of
- **Repository navigation** - browse repositories by organization
- **Pull Request list view** - view active pull requests for each repository
- **Pagination** - displays 10 items per page with navigation
- **Repository grouping** - organizes repos by user/organization with visual headers
- **PR status indicators** - visual indicators for draft, approved, and review states
- Navigation with arrow keys
- Clean, modern terminal UI

## Prerequisites

- Go 1.21 or later
- GitHub CLI (`gh`) installed and authenticated

## Setup

1. Clone the repository:
```bash
git clone <your-repo-url>
cd ghprs
```

2. Ensure you're authenticated with GitHub CLI:
```bash
gh auth login
```

The application will automatically use your GitHub CLI token. If you prefer to use an environment variable instead, you can set:
```bash
export GITHUB_TOKEN="your-github-token-here"
```

3. Install dependencies:
```bash
go mod tidy
```

4. Run the application:
```bash
go run cmd/main.go
```

## Development

### Pre-commit Setup

Install pre-commit hooks:
```bash
# Install pre-commit
pip install pre-commit

# Install the git hook scripts
pre-commit install
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o bin/ghprs cmd/main.go
```

## Usage

### Three-Level Navigation

The application uses a hierarchical navigation system:

1. **Organization Selection**: First, select an organization or user
2. **Repository Selection**: Then, select a repository from that organization
3. **Pull Request List**: View active pull requests for the selected repository

### Navigation Controls

- **↑/↓ or j/k**: Navigate through items on current page
- **←/→ or h/l**: Navigate between pages
- **g**: Go to first page
- **G**: Go to last page
- **Enter**: Select organization (first level), repository (second level), or view PR details (third level)
- **b or Backspace**: Go back to previous level
- **d**: Toggle debug mode (shows token info and timestamps)
- **r**: Reload repositories
- **q**: Quit the application

### Repository Organization

The application fetches repositories from:
- **Your personal repositories**
- **All organizations you're a member of**

**Three-Level Navigation:**
- **Level 1**: Select organization/user (shows repo count)
- **Level 2**: Select repository from that organization
- **Level 3**: View pull requests with status indicators
- **Visual indicators**: 📁 for organizations, 📦 for repositories, 🔵🟢🔴🟡 for PR status
- **Easy navigation**: Use Enter to select, Backspace to go back

### Pagination

The application displays 10 items per page to handle large lists efficiently:
- **Organization view**: Shows current page and total organization count
- **Repository view**: Shows current page and total repository count for selected organization
- Navigate between pages with arrow keys or h/l
- Jump to first/last page with g/G
- Cursor resets to top when changing pages

### Debug Mode

Press `d` to toggle debug mode, which will show:
- Last update timestamp
- Masked GitHub token (for verification)
- Loading status and error messages

## Architecture

This is a minimal prototype that will be incrementally enhanced with:

- GitHub API integration
- Pull request data fetching
- Caching layer
- Detailed PR views
- Review and comment display

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run pre-commit hooks: `pre-commit run --all-files`
5. Submit a pull request

## License

MIT
