package components

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func TestPlanLabel(t *testing.T) {
	cases := []struct {
		orgType string
		tier    string
		want    string
	}{
		{"claude_max", "default_claude_max_20x", "Max 20x"},
		{"claude_max", "default_claude_max_5x", "Max 5x"},
		{"claude_pro", "", "Pro"},
		{"claude_team", "", "Team"},
		{"claude_enterprise", "default_claude_enterprise", "Enterprise"},
		{"", "default_claude_max_20x", ""}, // no family -> nothing
		{"claude_", "", ""},                // empty after prefix -> nothing
	}
	for _, c := range cases {
		if got := planLabel(c.orgType, c.tier); got != c.want {
			t.Errorf("planLabel(%q, %q) = %q, want %q", c.orgType, c.tier, got, c.want)
		}
	}
}

func TestPlan_RendersLabel(t *testing.T) {
	data := &schema.Input{Plan: schema.Plan{OrgType: "claude_max", RateLimitTier: "default_claude_max_20x"}}
	result := Get("plan").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "Max 20x") {
		t.Errorf("expected 'Max 20x' in output, got %q", result)
	}
}

func TestPlan_EmptyWhenNoPlan(t *testing.T) {
	data := &schema.Input{}
	result := Get("plan").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty when no plan, got %q", result)
	}
}

func TestPlan_EmojiPrefix(t *testing.T) {
	data := &schema.Input{Plan: schema.Plan{OrgType: "claude_max", RateLimitTier: "default_claude_max_20x"}}
	cfg := config.Default()
	result := Get("plan").Render(data, cfg, theme.Get("default"))
	if !strings.Contains(result, "🎫") {
		t.Errorf("expected 🎫 emoji, got %q", result)
	}

	cfg.Emojis = config.EmojiNone
	result = Get("plan").Render(data, cfg, theme.Get("default"))
	if strings.Contains(result, "🎫") {
		t.Errorf("expected no emoji when disabled, got %q", result)
	}
	if !strings.Contains(result, "Plan:") {
		t.Errorf("expected 'Plan:' text prefix when emojis off, got %q", result)
	}
}
