package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestPR_RendersNumber(t *testing.T) {
	data := &schema.Input{PR: &schema.PR{Number: 1234, URL: "https://example/pr/1234"}}
	result := Get("pr").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "#1234") {
		t.Errorf("expected '#1234' in output, got %q", result)
	}
}

func TestPR_EmptyWhenNil(t *testing.T) {
	data := &schema.Input{}
	result := Get("pr").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestPR_EmptyWhenNumberZero(t *testing.T) {
	data := &schema.Input{PR: &schema.PR{}}
	result := Get("pr").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty when number=0, got %q", result)
	}
}

func TestPR_ReviewStateGlyphs(t *testing.T) {
	cases := []struct {
		state string
		glyph string
	}{
		{"approved", "✓"},
		{"changes_requested", "✗"},
		{"pending", "⏳"},
		{"draft", "·"},
	}
	for _, tc := range cases {
		data := &schema.Input{PR: &schema.PR{Number: 1, ReviewState: tc.state}}
		result := Get("pr").Render(data, config.Default(), theme.Get("default"))
		if !strings.Contains(result, tc.glyph) {
			t.Errorf("state %q: expected %q in output, got %q", tc.state, tc.glyph, result)
		}
	}
}

func TestPR_WrapsHyperlinkWhenURLPresent(t *testing.T) {
	data := &schema.Input{PR: &schema.PR{Number: 7, URL: "https://example.com/7"}}
	result := Get("pr").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "\x1b]8;;https://example.com/7\x07") {
		t.Errorf("expected OSC 8 link sequence in output, got %q", result)
	}
}

func TestPR_NoHyperlinkWhenURLEmpty(t *testing.T) {
	data := &schema.Input{PR: &schema.PR{Number: 7}}
	result := Get("pr").Render(data, config.Default(), theme.Get("default"))
	if strings.Contains(result, "\x1b]8;;") {
		t.Errorf("expected no OSC 8 link without URL, got %q", result)
	}
	if !strings.Contains(result, "#7") {
		t.Errorf("expected '#7', got %q", result)
	}
}
