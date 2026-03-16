package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestContextBar_GradientStyle(t *testing.T) {
	cfg := config.Default()
	cfg.ContextBar.Style = config.BarGradient
	pct := 80.0
	data := &schema.Input{}
	data.ContextWindow.UsedPercentage = &pct
	result := Get("context_bar").Render(data, cfg, theme.Get("default"))
	if result == "" {
		t.Error("expected non-empty result for gradient bar")
	}
	// percentage should still appear
	if !strings.Contains(result, "80%") {
		t.Errorf("expected percentage in output, got %q", result)
	}
}

func TestContextBar_GradientNilPct(t *testing.T) {
	cfg := config.Default()
	cfg.ContextBar.Style = config.BarGradient
	data := &schema.Input{}
	result := Get("context_bar").Render(data, cfg, theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty for nil percentage, got %q", result)
	}
}
