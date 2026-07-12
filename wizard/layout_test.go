package wizard

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
)

func sliceHas(s []string, want string) bool {
	for _, v := range s {
		if v == want {
			return true
		}
	}
	return false
}

func TestDefaultState_LayoutIsAuto(t *testing.T) {
	if DefaultState().Layout != "auto" {
		t.Errorf("DefaultState().Layout = %q, want \"auto\"", DefaultState().Layout)
	}
}

func TestToConfig_AutoWritesSingleFlatLine(t *testing.T) {
	s := DefaultState()
	s.Features = []string{"model", "directory", "session_id", "cost"}
	s.Layout = "auto"

	cfg := s.ToConfig()

	if cfg.Layout != config.LayoutAuto {
		t.Fatalf("cfg.Layout = %q, want %q", cfg.Layout, config.LayoutAuto)
	}
	if len(cfg.Lines) != 1 {
		t.Fatalf("auto mode should write a single flat line, got %d lines: %v", len(cfg.Lines), cfg.Lines)
	}
	for _, want := range []string{"model", "directory", "session_id", "cost"} {
		if !sliceHas(cfg.Lines[0].Components, want) {
			t.Errorf("flat line missing %q, got %v", want, cfg.Lines[0].Components)
		}
	}
}

func TestToConfig_FixedSetsLayoutFixed(t *testing.T) {
	s := DefaultState()
	s.Layout = "fixed"

	cfg := s.ToConfig()

	if cfg.Layout != config.LayoutFixed {
		t.Errorf("cfg.Layout = %q, want %q", cfg.Layout, config.LayoutFixed)
	}
}

func TestToTOML_IncludesAutoLayout(t *testing.T) {
	s := DefaultState()
	s.Layout = "auto"

	out, err := s.ToTOML()
	if err != nil {
		t.Fatalf("ToTOML error: %v", err)
	}
	if !strings.Contains(out, `layout = "auto"`) {
		t.Errorf("ToTOML should contain layout = \"auto\", got:\n%s", out)
	}
}

func TestStateFromConfig_RoundTripsAuto(t *testing.T) {
	cfg := config.Default()
	cfg.Layout = config.LayoutAuto
	cfg.Lines = []config.LineConfig{
		{Components: []string{"model", "directory", "session_id"}},
	}

	state, reasons := StateFromConfig(cfg)

	if state.Layout != "auto" {
		t.Errorf("state.Layout = %q, want \"auto\"", state.Layout)
	}
	for _, r := range reasons {
		if r == "hand-edited layout" {
			t.Errorf("auto config's single flat line should not be flagged hand-edited; reasons=%v", reasons)
		}
	}
}
