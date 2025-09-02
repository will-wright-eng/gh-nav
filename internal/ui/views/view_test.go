package views

import "testing"

func TestBaseView(t *testing.T) {
	view := NewBaseView(5)

	// Test initial state
	if view.GetCursor() != 0 {
		t.Errorf("Expected cursor 0, got %d", view.GetCursor())
	}
	if view.GetPage() != 0 {
		t.Errorf("Expected page 0, got %d", view.GetPage())
	}
	if view.GetPageSize() != 5 {
		t.Errorf("Expected page size 5, got %d", view.GetPageSize())
	}

	// Test cursor movement
	view.MoveCursorUp(10)
	if view.GetCursor() != 0 {
		t.Errorf("Expected cursor to stay at 0 when moving up from 0, got %d", view.GetCursor())
	}

	view.MoveCursorDown(10)
	if view.GetCursor() != 1 {
		t.Errorf("Expected cursor 1, got %d", view.GetCursor())
	}

	view.MoveCursorDown(10)
	if view.GetCursor() != 2 {
		t.Errorf("Expected cursor 2, got %d", view.GetCursor())
	}

	// Test pagination
	if !view.NextPage(15) {
		t.Error("Expected NextPage to return true for 15 items")
	}
	if view.GetPage() != 1 {
		t.Errorf("Expected page 1, got %d", view.GetPage())
	}
	if view.GetCursor() != 0 {
		t.Errorf("Expected cursor to reset to 0, got %d", view.GetCursor())
	}

	if !view.PreviousPage() {
		t.Error("Expected PreviousPage to return true")
	}
	if view.GetPage() != 0 {
		t.Errorf("Expected page 0, got %d", view.GetPage())
	}

	// Test page info
	info := view.GetPageInfo(15, "items")
	expected := "Page 1/3 (items 1-5 of 15)"
	if info != expected {
		t.Errorf("Expected '%s', got '%s'", expected, info)
	}
}

func TestBaseViewEdgeCases(t *testing.T) {
	view := NewBaseView(5)

	// Test with no items
	info := view.GetPageInfo(0, "items")
	expected := "No items found"
	if info != expected {
		t.Errorf("Expected '%s', got '%s'", expected, info)
	}

	// Test with items less than page size
	info = view.GetPageInfo(3, "items")
	expected = "Showing 3 items"
	if info != expected {
		t.Errorf("Expected '%s', got '%s'", expected, info)
	}

	// Test cursor bounds - MoveCursorDown should respect maxItems
	view.SetCursor(0)
	view.MoveCursorDown(5)
	if view.GetCursor() != 1 {
		t.Errorf("Expected cursor to be 1 when moving down from 0, got %d", view.GetCursor())
	}

	// Test that cursor doesn't go beyond maxItems-1
	view.SetCursor(3)
	view.MoveCursorDown(5)
	if view.GetCursor() != 4 {
		t.Errorf("Expected cursor to be 4 when moving down from 3, got %d", view.GetCursor())
	}

	// Test that cursor doesn't go beyond maxItems-1
	view.SetCursor(4)
	view.MoveCursorDown(5)
	if view.GetCursor() != 4 {
		t.Errorf("Expected cursor to stay at 4 when already at max, got %d", view.GetCursor())
	}

	// Test cursor bounds when moving down
	view.SetCursor(0)
	view.MoveCursorDown(5)
	if view.GetCursor() != 1 {
		t.Errorf("Expected cursor to be 1 when moving down from 0, got %d", view.GetCursor())
	}
}
