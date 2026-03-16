package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func contextBarInput(pct float64) *schema.Input {
	return &schema.Input{ContextWindow: schema.Context{UsedPercentage: ptr(pct)}}
}

func cfgWithBarStyle(style config.BarStyle, width int) *config.Config {
	cfg := config.Default()
	cfg.ContextBar.Style = style
	cfg.ContextBar.Width = width
	return cfg
}

func TestContextBar_NilPercentage(t *testing.T) {
	c := Get("context_bar")
	th := theme.Get("default")

	result := c.Render(&schema.Input{}, config.Default(), th)

	if result != "" {
		t.Errorf("expected empty string for nil percentage, got %q", result)
	}
}

func TestContextBar_BlockStyle(t *testing.T) {
	c := Get("context_bar")
	th := theme.Get("default")

	result := c.Render(contextBarInput(28), cfgWithBarStyle(config.BarBlock, 10), th)

	if !strings.Contains(result, "28%") {
		t.Errorf("expected '28%%' in output, got %q", result)
	}
	if !strings.Contains(result, "▓") {
		t.Errorf("expected block char '▓' in output, got %q", result)
	}
	if !strings.Contains(result, "░") {
		t.Errorf("expected empty char '░' in output, got %q", result)
	}
}

func TestContextBar_SolidStyle(t *testing.T) {
	c := Get("context_bar")
	th := theme.Get("default")

	result := c.Render(contextBarInput(28), cfgWithBarStyle(config.BarSolid, 10), th)

	if !strings.Contains(result, "█") {
		t.Errorf("expected solid char '█' in output, got %q", result)
	}
	if !strings.Contains(result, "░") {
		t.Errorf("expected empty char '░' in output, got %q", result)
	}
}

func TestContextBar_ASCIIStyle(t *testing.T) {
	c := Get("context_bar")
	th := theme.Get("default")

	result := c.Render(contextBarInput(28), cfgWithBarStyle(config.BarASCII, 10), th)

	if !strings.Contains(result, "[") || !strings.Contains(result, "]") {
		t.Errorf("expected brackets in ASCII bar, got %q", result)
	}
	if !strings.Contains(result, "=") {
		t.Errorf("expected '=' in ASCII bar, got %q", result)
	}
	if !strings.Contains(result, "-") {
		t.Errorf("expected '-' in ASCII bar, got %q", result)
	}
}

func TestContextBar_PercentStyle(t *testing.T) {
	c := Get("context_bar")
	th := theme.Get("default")

	result := c.Render(contextBarInput(28), cfgWithBarStyle(config.BarPercent, 10), th)

	if !strings.Contains(result, "28%") {
		t.Errorf("expected '28%%' in output, got %q", result)
	}
	stripped := stripANSI(result)
	if strings.Contains(stripped, "▓") || strings.Contains(stripped, "█") || strings.Contains(stripped, "[") {
		t.Errorf("percent style should not contain bar chars, got %q", stripped)
	}
}

func TestContextBar_Width(t *testing.T) {
	c := Get("context_bar")
	th := theme.Get("default")

	r5 := c.Render(contextBarInput(50), cfgWithBarStyle(config.BarBlock, 5), th)
	r20 := c.Render(contextBarInput(50), cfgWithBarStyle(config.BarBlock, 20), th)

	if len([]rune(stripANSI(r5))) >= len([]rune(stripANSI(r20))) {
		t.Errorf("width=5 bar should be shorter than width=20 bar")
	}
}

func TestContextBar_ColorThresholds(t *testing.T) {
	c := Get("context_bar")
	th := theme.Get("default")
	cfg := config.Default()

	low := c.Render(contextBarInput(28), cfg, th)
	warn := c.Render(contextBarInput(75), cfg, th)
	danger := c.Render(contextBarInput(95), cfg, th)

	if !strings.Contains(low, "\033[32m") {
		t.Errorf("low context should be green, got %q", low)
	}
	if !strings.Contains(warn, "\033[33m") {
		t.Errorf("warning context should be yellow, got %q", warn)
	}
	if !strings.Contains(danger, "\033[31m") {
		t.Errorf("danger context should be red, got %q", danger)
	}
}

func TestContextBar_NoEmoji(t *testing.T) {
	c := Get("context_bar")
	th := theme.Get("default")
	cfg := config.Default()
	cfg.Emojis = config.EmojiNone

	result := c.Render(contextBarInput(28), cfg, th)

	if strings.Contains(result, "📊") {
		t.Errorf("expected no emoji, got %q", result)
	}
	if !strings.Contains(result, "28%") {
		t.Errorf("expected percentage in output, got %q", result)
	}
}

func stripANSI(s string) string {
	result := []rune{}
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result = append(result, r)
	}
	return string(result)
}
