package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func ptr[T any](v T) *T { return &v }

func TestContextPct_Low(t *testing.T) {
	c := Get("context_pct")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{ContextWindow: schema.Context{UsedPercentage: ptr(28.0)}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "28%") {
		t.Errorf("expected '28%%', got %q", result)
	}
	if !strings.Contains(result, "\033[32m") {
		t.Errorf("expected green color for low context, got %q", result)
	}
}

func TestContextPct_Warning(t *testing.T) {
	c := Get("context_pct")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{ContextWindow: schema.Context{UsedPercentage: ptr(75.0)}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "75%") {
		t.Errorf("expected '75%%', got %q", result)
	}
	if !strings.Contains(result, "\033[33m") {
		t.Errorf("expected yellow color for warning context, got %q", result)
	}
}

func TestContextPct_Danger(t *testing.T) {
	c := Get("context_pct")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{ContextWindow: schema.Context{UsedPercentage: ptr(95.0)}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "95%") {
		t.Errorf("expected '95%%', got %q", result)
	}
	if !strings.Contains(result, "\033[31m") {
		t.Errorf("expected red color for danger context, got %q", result)
	}
}

func TestContextPct_NilPercentage(t *testing.T) {
	c := Get("context_pct")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{ContextWindow: schema.Context{UsedPercentage: nil}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "--") {
		t.Errorf("expected '--' for nil percentage, got %q", result)
	}
}

func TestContextPct_Threshold_ExactlyAt70(t *testing.T) {
	c := Get("context_pct")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{ContextWindow: schema.Context{UsedPercentage: ptr(70.0)}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "\033[33m") {
		t.Errorf("expected yellow at exactly 70%%, got %q", result)
	}
}

func TestContextPct_Threshold_ExactlyAt90(t *testing.T) {
	c := Get("context_pct")
	th := theme.Get("default")
	cfg := config.Default()
	input := &schema.Input{ContextWindow: schema.Context{UsedPercentage: ptr(90.0)}}

	result := c.Render(input, cfg, th)

	if !strings.Contains(result, "\033[31m") {
		t.Errorf("expected red at exactly 90%%, got %q", result)
	}
}
