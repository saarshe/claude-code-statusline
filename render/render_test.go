package render

import (
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func ptr[T any](v T) *T { return &v }

func defaultCfg() *config.Config {
	cfg := config.Default()
	// Use simple single-line layout for render tests
	cfg.Lines = []config.LineConfig{
		{Components: []string{"model", "context_pct", "cost"}},
	}
	return cfg
}

func fullInput() *schema.Input {
	return &schema.Input{
		Model: schema.Model{DisplayName: "Opus"},
		Cost:  schema.Cost{TotalCostUSD: 0.42},
		ContextWindow: schema.Context{
			UsedPercentage: ptr(28.0),
		},
	}
}

func TestRender_FullInput(t *testing.T) {
	output := Render(fullInput(), defaultCfg())

	if !strings.Contains(output, "Opus") {
		t.Errorf("output should contain model name 'Opus', got: %q", output)
	}
	if !strings.Contains(output, "28%") {
		t.Errorf("output should contain context percentage '28%%', got: %q", output)
	}
	if !strings.Contains(output, "$0.42") {
		t.Errorf("output should contain cost '$0.42', got: %q", output)
	}
}

func TestRender_NilUsedPercentage(t *testing.T) {
	input := fullInput()
	input.ContextWindow.UsedPercentage = nil

	output := Render(input, defaultCfg())

	if !strings.Contains(output, "--") {
		t.Errorf("output should contain '--' for nil percentage, got: %q", output)
	}
}

func TestRender_ZeroCost(t *testing.T) {
	input := fullInput()
	input.Cost.TotalCostUSD = 0

	output := Render(input, defaultCfg())

	if !strings.Contains(output, "$0.00") {
		t.Errorf("output should contain '$0.00' for zero cost, got: %q", output)
	}
}

func TestRender_NilCurrentUsage(t *testing.T) {
	input := &schema.Input{
		Model: schema.Model{DisplayName: "Opus"},
		Cost:  schema.Cost{TotalCostUSD: 0.42},
		ContextWindow: schema.Context{
			UsedPercentage: ptr(28.0),
			CurrentUsage:   nil,
		},
	}

	output := Render(input, defaultCfg())

	if output == "" {
		t.Error("output should not be empty")
	}
}

func TestRender_EmptyModel(t *testing.T) {
	input := &schema.Input{
		Model: schema.Model{DisplayName: ""},
		Cost:  schema.Cost{TotalCostUSD: 0.42},
		ContextWindow: schema.Context{
			UsedPercentage: ptr(28.0),
		},
	}

	output := Render(input, defaultCfg())
	if output == "" {
		t.Error("output should not be empty even with empty model")
	}
}

func TestRender_HighContextPercentage(t *testing.T) {
	input := fullInput()
	input.ContextWindow.UsedPercentage = ptr(95.0)

	output := Render(input, defaultCfg())

	if !strings.Contains(output, "95%") {
		t.Errorf("output should contain '95%%', got: %q", output)
	}
}

func TestRender_NoEmojis(t *testing.T) {
	cfg := defaultCfg()
	cfg.Emojis = config.EmojiNone


	output := Render(fullInput(), cfg)

	if strings.Contains(output, "🤖") {
		t.Errorf("output should not contain emoji with EmojiNone, got: %q", output)
	}
	if !strings.Contains(output, "Opus") {
		t.Errorf("output should still contain model name, got: %q", output)
	}
}

func TestRender_CustomSeparator(t *testing.T) {
	cfg := defaultCfg()
	cfg.Separator.Character = "•"

	output := Render(fullInput(), cfg)

	if !strings.Contains(output, "•") {
		t.Errorf("output should contain custom separator '•', got: %q", output)
	}
}

func TestRenderWithTheme_MultiLine(t *testing.T) {
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		{Components: []string{"model"}},
		{Components: []string{"cost"}},
	}
	th := theme.Get("default")

	output := RenderWithTheme(fullInput(), cfg, th)

	if !strings.Contains(output, "\n") {
		t.Errorf("multi-line output should contain newline, got: %q", output)
	}
	lines := strings.Split(output, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %q", len(lines), output)
	}
}

func TestRenderWithTheme_EmptyComponentsSkipped(t *testing.T) {
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		{Components: []string{"model"}},
	}
	th := theme.Get("default")

	// Empty model name — component returns ""
	input := &schema.Input{Model: schema.Model{DisplayName: ""}}
	output := RenderWithTheme(input, cfg, th)

	if output != "" {
		t.Errorf("line with all-empty components should be skipped, got: %q", output)
	}
}

func TestRenderWithTheme_ThemeSeparatorOverridesConfig(t *testing.T) {
	cfg := defaultCfg()
	cfg.Separator.Character = "|"

	th := theme.Get("default")

	// Default theme has no separator set, so config separator is used.
	output1 := RenderWithTheme(fullInput(), cfg, th)
	if !strings.Contains(output1, "|") {
		t.Errorf("expected config separator '|', got: %q", output1)
	}

	// Create a theme with a custom separator.
	custom := *th
	custom.Separator = " · "
	output2 := RenderWithTheme(fullInput(), cfg, &custom)
	if !strings.Contains(output2, "·") {
		t.Errorf("expected theme separator '·', got: %q", output2)
	}
}

func TestRender_NewComponentsComposeTogether(t *testing.T) {
	// Integration check: exercise all four new components through the full
	// render pipeline (registry, separator, theme styling) at once.
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		{Components: []string{"model", "effort", "pr"}},
		{Components: []string{"rate_limits_reset"}},
	}
	input := fullInput()
	input.Effort = &schema.Effort{Level: "high"}
	input.PR = &schema.PR{Number: 99, URL: "https://example.com/99", ReviewState: "approved"}
	input.RateLimits = &schema.RateLimits{
		FiveHour: &schema.RateLimitWindow{UsedPercentage: 25, ResetsAt: 9999999999},
		SevenDay: &schema.RateLimitWindow{UsedPercentage: 50, ResetsAt: 9999999999},
	}

	output := Render(input, cfg)
	lines := strings.Split(output, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), output)
	}
	if !strings.Contains(lines[0], "high") {
		t.Errorf("line 1 should contain effort 'high', got %q", lines[0])
	}
	if !strings.Contains(lines[0], "#99") {
		t.Errorf("line 1 should contain PR '#99', got %q", lines[0])
	}
	// OSC 8 hyperlink survives the render pipeline.
	if !strings.Contains(lines[0], "\x1b]8;;https://example.com/99\x07") {
		t.Errorf("line 1 should contain OSC 8 link, got %q", lines[0])
	}
	if !strings.Contains(lines[1], "5h 25%") || !strings.Contains(lines[1], "7d 50%") {
		t.Errorf("line 2 should contain both rate-limit windows, got %q", lines[1])
	}
	// Separator should appear between effort and pr on line 1.
	if !strings.Contains(lines[0], "|") {
		t.Errorf("line 1 should contain separator between components, got %q", lines[0])
	}
}

func TestRender_NewComponentsAllAbsent_ProducesNoOutput(t *testing.T) {
	// When the JSON omits the new fields entirely (e.g. API-key user, model
	// without effort, no open PR), a line containing only the new components
	// should render as empty and be dropped from the output.
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		{Components: []string{"effort", "pr", "rate_limits_reset"}},
		{Components: []string{"model"}},
	}
	input := &schema.Input{Model: schema.Model{DisplayName: "Opus"}}

	output := Render(input, cfg)
	lines := strings.Split(output, "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line (all-absent line skipped), got %d: %q", len(lines), output)
	}
	if !strings.Contains(lines[0], "Opus") {
		t.Errorf("surviving line should contain 'Opus', got %q", lines[0])
	}
}

func TestRenderWithTheme_UnknownComponentSkipped(t *testing.T) {
	cfg := config.Default()
	cfg.Lines = []config.LineConfig{
		{Components: []string{"unknown_component", "cost"}},
	}
	th := theme.Get("default")

	output := RenderWithTheme(fullInput(), cfg, th)

	if !strings.Contains(output, "$0.42") {
		t.Errorf("known component should still render, got: %q", output)
	}
}
