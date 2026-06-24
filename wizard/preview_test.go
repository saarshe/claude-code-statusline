package wizard

import (
	"strings"
	"testing"
)

func TestMockInput_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MockInput() panicked: %v", r)
		}
	}()
	input := MockInput()
	if input == nil {
		t.Error("MockInput() returned nil")
	}
}

func TestPreview_NonEmpty(t *testing.T) {
	state := DefaultState()
	output := Preview(state)
	if output == "" {
		t.Error("Preview() returned empty string")
	}
}

func TestPreview_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Preview() panicked: %v", r)
		}
	}()
	Preview(DefaultState())
}

func TestPreview_ReflectsState(t *testing.T) {
	// Only model selected — should produce non-empty output
	state := DefaultState()
	state.Features = []string{"model"}
	output := Preview(state)
	if output == "" {
		t.Error("Preview with model feature should produce output")
	}
}

func TestPreview_ShowsSessionID(t *testing.T) {
	// The mock must carry a session id, otherwise the session_id component
	// renders empty and is invisible (and unmeasured) in the wizard preview.
	id := MockInput().SessionID
	if id == "" {
		t.Fatal("MockInput() must provide a non-empty SessionID so the wizard preview can show the session_id component")
	}

	state := DefaultState()
	state.Features = []string{"session_id"}
	output := Preview(state)
	if !strings.Contains(output, id) {
		t.Errorf("expected preview to contain session id %q, got:\n%s", id, output)
	}
}
