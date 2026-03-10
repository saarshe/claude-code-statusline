package wizard

import "testing"

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
