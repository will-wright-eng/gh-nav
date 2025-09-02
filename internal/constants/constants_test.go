package constants

import "testing"

func TestConstants(t *testing.T) {
	// Test navigation constants
	if KeyUp != "up" {
		t.Errorf("Expected KeyUp to be 'up', got %s", KeyUp)
	}
	if KeyDown != "down" {
		t.Errorf("Expected KeyDown to be 'down', got %s", KeyDown)
	}
	if KeyEnter != "enter" {
		t.Errorf("Expected KeyEnter to be 'enter', got %s", KeyEnter)
	}
	if KeyQuit != "q" {
		t.Errorf("Expected KeyQuit to be 'q', got %s", KeyQuit)
	}

	// Test UI constants
	if DefaultPageSize != 10 {
		t.Errorf("Expected DefaultPageSize to be 10, got %d", DefaultPageSize)
	}
	if MaxTitleLength != 60 {
		t.Errorf("Expected MaxTitleLength to be 60, got %d", MaxTitleLength)
	}

	// Test icon constants
	if IconOrganization != "📁" {
		t.Errorf("Expected IconOrganization to be '📁', got %s", IconOrganization)
	}
	if IconRepository != "📦" {
		t.Errorf("Expected IconRepository to be '📦', got %s", IconRepository)
	}
	if IconOpen != "🔵" {
		t.Errorf("Expected IconOpen to be '🔵', got %s", IconOpen)
	}

	// Test color constants
	if ColorPrimary != "#00FF00" {
		t.Errorf("Expected ColorPrimary to be '#00FF00', got %s", ColorPrimary)
	}
	if ColorError != "#FF0000" {
		t.Errorf("Expected ColorError to be '#FF0000', got %s", ColorError)
	}

	// Test help text
	if HelpNavigation == "" {
		t.Error("Expected HelpNavigation to not be empty")
	}
}
