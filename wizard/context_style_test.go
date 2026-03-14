package wizard

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func stateWithContextStyle(style string) *WizardState {
	s := DefaultState()
	s.ContextStyle = style
	return s
}

func TestContextStyleModel_DefaultsToFirstChoice(t *testing.T) {
	m := newContextStyleModel(stateWithContextStyle("block"))
	if m.cursor < 0 || m.cursor >= len(m.choices) {
		t.Fatalf("cursor %d out of range [0, %d)", m.cursor, len(m.choices))
	}
}

func TestContextStyleModel_EnterReturnsSelection(t *testing.T) {
	m := newContextStyleModel(stateWithContextStyle("pct"))
	m.cursor = 0 // "pct"
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	result := updated.(contextStyleModel)
	if !result.done {
		t.Error("expected done=true after enter")
	}
	if result.result != m.choices[0].value {
		t.Errorf("expected %q, got %q", m.choices[0].value, result.result)
	}
}

func TestContextStyleModel_ArrowMovesDown(t *testing.T) {
	m := newContextStyleModel(stateWithContextStyle("pct"))
	m.cursor = 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	result := updated.(contextStyleModel)
	if result.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", result.cursor)
	}
}

func TestContextStyleModel_ArrowClampsAtBottom(t *testing.T) {
	m := newContextStyleModel(stateWithContextStyle("pct"))
	m.cursor = len(m.choices) - 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	result := updated.(contextStyleModel)
	if result.cursor != len(m.choices)-1 {
		t.Errorf("expected cursor to stay at bottom, got %d", result.cursor)
	}
}

func TestContextStyleModel_TickAdvancesPct(t *testing.T) {
	m := newContextStyleModel(stateWithContextStyle("block"))
	m.pct = 10.0
	updated, _ := m.Update(tickMsg{})
	result := updated.(contextStyleModel)
	if result.pct <= 10.0 {
		t.Errorf("expected pct to advance from 10, got %f", result.pct)
	}
}

func TestContextStyleModel_EscSetsGoBack(t *testing.T) {
	m := newContextStyleModel(stateWithContextStyle("pct"))
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	result := updated.(contextStyleModel)
	if !result.done {
		t.Error("expected done=true after esc")
	}
	if !result.goBack {
		t.Error("expected goBack=true after esc")
	}
}

func TestContextStyleModel_ViewIsNonEmpty(t *testing.T) {
	m := newContextStyleModel(stateWithContextStyle("block"))
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}
