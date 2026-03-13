package wizard

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestThemeSelectorModel_DefaultsToCurrentTheme(t *testing.T) {
	m := newThemeSelectorModel("nord")
	if m.names[m.cursor] != "nord" {
		t.Errorf("cursor at %q, want %q", m.names[m.cursor], "nord")
	}
}

func TestThemeSelectorModel_DefaultsToFirstIfUnknown(t *testing.T) {
	m := newThemeSelectorModel("nonexistent")
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 for unknown theme", m.cursor)
	}
}

func TestThemeSelectorModel_EnterReturnsSelection(t *testing.T) {
	m := newThemeSelectorModel("default")
	// Move down to second theme
	m, _ = applyThemeKey(m, tea.KeyDown)
	m, _ = applyThemeKey(m, tea.KeyEnter)
	if !m.done {
		t.Fatal("expected done after enter")
	}
	if m.result != m.names[1] {
		t.Errorf("result = %q, want %q", m.result, m.names[1])
	}
}

func TestThemeSelectorModel_CtrlCReturnsEmpty(t *testing.T) {
	m := newThemeSelectorModel("default")
	m, _ = applyThemeKey(m, tea.KeyCtrlC)
	if !m.done {
		t.Fatal("expected done after ctrl+c")
	}
	if m.result != "" {
		t.Errorf("result = %q, want empty on cancel", m.result)
	}
}

func TestThemeSelectorModel_ArrowNavigation(t *testing.T) {
	m := newThemeSelectorModel("default")
	if m.cursor != 0 {
		t.Fatalf("initial cursor = %d, want 0", m.cursor)
	}

	m, _ = applyThemeKey(m, tea.KeyDown)
	if m.cursor != 1 {
		t.Errorf("after down: cursor = %d, want 1", m.cursor)
	}

	m, _ = applyThemeKey(m, tea.KeyUp)
	if m.cursor != 0 {
		t.Errorf("after up: cursor = %d, want 0", m.cursor)
	}
}

func TestThemeSelectorModel_ClampsAtBounds(t *testing.T) {
	m := newThemeSelectorModel("default")
	// Up at top should stay at 0
	m, _ = applyThemeKey(m, tea.KeyUp)
	if m.cursor != 0 {
		t.Errorf("up at top: cursor = %d, want 0", m.cursor)
	}

	// Go to bottom
	for i := 0; i < len(m.names)+5; i++ {
		m, _ = applyThemeKey(m, tea.KeyDown)
	}
	if m.cursor != len(m.names)-1 {
		t.Errorf("past bottom: cursor = %d, want %d", m.cursor, len(m.names)-1)
	}
}

func TestThemeSelectorModel_ViewShowsPreview(t *testing.T) {
	m := newThemeSelectorModel("default")
	view := m.View()
	if view == "" {
		t.Fatal("View() returned empty string")
	}
	if !containsAny(view, "preview:") {
		t.Error("View() should contain preview label")
	}
}

func TestThemeSelectorModel_ViewEmptyWhenDone(t *testing.T) {
	m := newThemeSelectorModel("default")
	m, _ = applyThemeKey(m, tea.KeyEnter)
	if m.View() != "" {
		t.Error("View() should be empty when done")
	}
}

func TestThemeSelectorModel_ViewChangesWithCursor(t *testing.T) {
	m := newThemeSelectorModel("default")
	view1 := m.View()
	m, _ = applyThemeKey(m, tea.KeyDown)
	view2 := m.View()
	if view1 == view2 {
		t.Error("View() should change when cursor moves to a different theme")
	}
}

func TestThemePreview_AllThemes(t *testing.T) {
	for _, name := range sortedThemeNames() {
		t.Run(name, func(t *testing.T) {
			output := themePreview(name)
			if output == "" {
				t.Errorf("themePreview(%q) returned empty string", name)
			}
			if len(output) <= 10 {
				t.Errorf("themePreview(%q) too short, expected styled output: %q", name, output)
			}
		})
	}
}

func TestThemePreview_DifferentThemesProduceDifferentOutput(t *testing.T) {
	defaultOut := themePreview("default")
	draculaOut := themePreview("dracula")
	if defaultOut == draculaOut {
		t.Error("expected different output for default vs dracula themes")
	}
}

// containsAny checks if s contains substr.
func containsAny(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && contains(s, substr)
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func applyThemeKey(m themeSelectorModel, keyType tea.KeyType) (themeSelectorModel, tea.Cmd) {
	model, cmd := m.Update(tea.KeyMsg{Type: keyType})
	return model.(themeSelectorModel), cmd
}
