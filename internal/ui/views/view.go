package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// View represents a UI view component
type View interface {
	// Update handles messages and updates the view
	Update(msg tea.Msg) (View, tea.Cmd)

	// View renders the view content
	View() string

	// SetSize updates the view dimensions
	SetSize(width, height int)

	// GetCursor returns the current cursor position
	GetCursor() int

	// SetCursor sets the cursor position
	SetCursor(cursor int)

	// GetPage returns the current page
	GetPage() int

	// SetPage sets the current page
	SetPage(page int)
}

// BaseView provides common functionality for all views
type BaseView struct {
	width    int
	height   int
	cursor   int
	page     int
	pageSize int
}

// NewBaseView creates a new base view
func NewBaseView(pageSize int) BaseView {
	return BaseView{
		width:    0,
		height:   0,
		cursor:   0,
		page:     0,
		pageSize: pageSize,
	}
}

// SetSize updates the view dimensions
func (b *BaseView) SetSize(width, height int) {
	b.width = width
	b.height = height
}

// GetCursor returns the current cursor position
func (b *BaseView) GetCursor() int {
	return b.cursor
}

// SetCursor sets the cursor position
func (b *BaseView) SetCursor(cursor int) {
	b.cursor = cursor
}

// GetPage returns the current page
func (b *BaseView) GetPage() int {
	return b.page
}

// SetPage sets the current page
func (b *BaseView) SetPage(page int) {
	b.page = page
}

// GetPageSize returns the page size
func (b *BaseView) GetPageSize() int {
	return b.pageSize
}

// MoveCursorUp moves the cursor up
func (b *BaseView) MoveCursorUp(maxItems int) {
	if b.cursor > 0 {
		b.cursor--
	}
}

// MoveCursorDown moves the cursor down
func (b *BaseView) MoveCursorDown(maxItems int) {
	if b.cursor < maxItems-1 {
		b.cursor++
	}
}

// NextPage goes to the next page
func (b *BaseView) NextPage(maxItems int) bool {
	totalPages := (maxItems - 1) / b.pageSize
	if b.page < totalPages {
		b.page++
		b.cursor = 0
		return true
	}
	return false
}

// PreviousPage goes to the previous page
func (b *BaseView) PreviousPage() bool {
	if b.page > 0 {
		b.page--
		b.cursor = 0
		return true
	}
	return false
}

// GoToFirstPage goes to the first page
func (b *BaseView) GoToFirstPage() {
	if b.page != 0 {
		b.page = 0
		b.cursor = 0
	}
}

// GoToLastPage goes to the last page
func (b *BaseView) GoToLastPage(maxItems int) {
	totalPages := (maxItems - 1) / b.pageSize
	if b.page != totalPages {
		b.page = totalPages
		b.cursor = 0
	}
}

// GetVisibleRange returns the visible range for pagination
func (b *BaseView) GetVisibleRange(maxItems int) (start, end int) {
	start = b.page * b.pageSize
	end = start + b.pageSize
	if end > maxItems {
		end = maxItems
	}
	if start >= maxItems {
		return 0, 0
	}
	return start, end
}

// GetPageInfo returns pagination information
func (b *BaseView) GetPageInfo(maxItems int, itemType string) string {
	if maxItems == 0 {
		return fmt.Sprintf("No %s found", itemType)
	}

	totalPages := (maxItems - 1) / b.pageSize
	start, end := b.GetVisibleRange(maxItems)

	if totalPages == 0 {
		return fmt.Sprintf("Showing %d %s", maxItems, itemType)
	}

	return fmt.Sprintf("Page %d/%d (%s %d-%d of %d)",
		b.page+1, totalPages+1, itemType, start+1, end, maxItems)
}
