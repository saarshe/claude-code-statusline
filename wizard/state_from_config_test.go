package wizard

import (
	"reflect"
	"sort"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
)

// equalFeatures compares two feature slices regardless of order.
func equalFeatures(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	ac := append([]string(nil), a...)
	bc := append([]string(nil), b...)
	sort.Strings(ac)
	sort.Strings(bc)
	return reflect.DeepEqual(ac, bc)
}

func TestStateFromConfig_DefaultRoundTrip(t *testing.T) {
	want := DefaultState()
	got, lossy := StateFromConfig(want.ToConfig())

	if lossy {
		t.Fatalf("default-state round trip should not be lossy")
	}
	if got.Theme != want.Theme {
		t.Errorf("Theme: got %q want %q", got.Theme, want.Theme)
	}
	if !equalFeatures(got.Features, want.Features) {
		t.Errorf("Features: got %v want %v", got.Features, want.Features)
	}
	if got.ContextStyle != want.ContextStyle {
		t.Errorf("ContextStyle: got %q want %q", got.ContextStyle, want.ContextStyle)
	}
	if got.TokenStyle != want.TokenStyle {
		t.Errorf("TokenStyle: got %q want %q", got.TokenStyle, want.TokenStyle)
	}
	if got.CacheStyle != want.CacheStyle {
		t.Errorf("CacheStyle: got %q want %q", got.CacheStyle, want.CacheStyle)
	}
	if got.LinesStyle != want.LinesStyle {
		t.Errorf("LinesStyle: got %q want %q", got.LinesStyle, want.LinesStyle)
	}
	if got.GitStyle != want.GitStyle {
		t.Errorf("GitStyle: got %q want %q", got.GitStyle, want.GitStyle)
	}
	if got.RateLimitsStyle != want.RateLimitsStyle {
		t.Errorf("RateLimitsStyle: got %q want %q", got.RateLimitsStyle, want.RateLimitsStyle)
	}
	if got.Emojis != want.Emojis {
		t.Errorf("Emojis: got %q want %q", got.Emojis, want.Emojis)
	}
	if got.BarWidth != want.BarWidth {
		t.Errorf("BarWidth: got %d want %d", got.BarWidth, want.BarWidth)
	}
}

func TestStateFromConfig_ContextStyles(t *testing.T) {
	cases := []struct {
		name      string
		component string
		barStyle  config.BarStyle
		want      string
	}{
		{"block", "context_bar", config.BarBlock, "block"},
		{"solid", "context_bar", config.BarSolid, "solid"},
		{"ascii", "context_bar", config.BarASCII, "ascii"},
		{"pct", "context_pct", "", "pct"},
		{"tokens", "context_tokens", "", "tokens"},
		{"tokens_bar", "context_tokens_bar", config.BarGradient, "tokens_bar"},
		{"tokens_bar_pct", "context_tokens_bar_pct", config.BarGradient, "tokens_bar_pct"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				Theme:  "default",
				Emojis: config.EmojiAll,
				ContextBar: config.ContextBarConfig{
					Style:      tc.barStyle,
					Width:      10,
					Thresholds: []int{70, 90},
				},
				Separator: config.SeparatorConfig{Character: "|"},
				Lines:     []config.LineConfig{{Components: []string{tc.component}}},
			}
			got, lossy := StateFromConfig(cfg)
			if lossy {
				t.Fatalf("unexpected lossy=true")
			}
			if !equalFeatures(got.Features, []string{"context"}) {
				t.Errorf("Features: got %v want [context]", got.Features)
			}
			if got.ContextStyle != tc.want {
				t.Errorf("ContextStyle: got %q want %q", got.ContextStyle, tc.want)
			}
		})
	}
}

func TestStateFromConfig_FeatureStyleMapping(t *testing.T) {
	cases := []struct {
		name      string
		component string
		feature   string
		style     string
		field     func(*WizardState) string
	}{
		{"tokens-turn", "tokens", "tokens", "turn", func(s *WizardState) string { return s.TokenStyle }},
		{"tokens-cache", "tokens_cache", "tokens", "turn_cache", func(s *WizardState) string { return s.TokenStyle }},
		{"tokens-session", "tokens_session", "tokens", "session", func(s *WizardState) string { return s.TokenStyle }},
		{"tokens-full", "tokens_full", "tokens", "full", func(s *WizardState) string { return s.TokenStyle }},
		{"cache-hit", "cache_hit", "cache", "hit", func(s *WizardState) string { return s.CacheStyle }},
		{"cache-counts", "cache", "cache", "counts", func(s *WizardState) string { return s.CacheStyle }},
		{"git-branch", "git_branch", "git", "branch", func(s *WizardState) string { return s.GitStyle }},
		{"git-status", "git_status", "git", "status", func(s *WizardState) string { return s.GitStyle }},
		{"lines-summary", "lines_summary", "lines_changed", "summary", func(s *WizardState) string { return s.LinesStyle }},
		{"lines-detail", "lines_changed", "lines_changed", "detail", func(s *WizardState) string { return s.LinesStyle }},
		{"rate-pct", "rate_limits", "rate_limits", "pct", func(s *WizardState) string { return s.RateLimitsStyle }},
		{"rate-reset", "rate_limits_reset", "rate_limits", "reset", func(s *WizardState) string { return s.RateLimitsStyle }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultConfigWithLines([]string{tc.component})
			got, lossy := StateFromConfig(cfg)
			if lossy {
				t.Fatalf("unexpected lossy=true")
			}
			if !equalFeatures(got.Features, []string{tc.feature}) {
				t.Errorf("Features: got %v want [%s]", got.Features, tc.feature)
			}
			if v := tc.field(got); v != tc.style {
				t.Errorf("style: got %q want %q", v, tc.style)
			}
		})
	}
}

func TestStateFromConfig_ThemeAndEmojisCopied(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"model"})
	cfg.Theme = "catppuccin"
	cfg.Emojis = config.EmojiNone

	got, _ := StateFromConfig(cfg)
	if got.Theme != "catppuccin" {
		t.Errorf("Theme: got %q want catppuccin", got.Theme)
	}
	if got.Emojis != "none" {
		t.Errorf("Emojis: got %q want none", got.Emojis)
	}
}

func TestStateFromConfig_BarWidthCopied(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"context_bar"})
	cfg.ContextBar.Width = 25
	got, _ := StateFromConfig(cfg)
	if got.BarWidth != 25 {
		t.Errorf("BarWidth: got %d want 25", got.BarWidth)
	}
}

func TestStateFromConfig_LossyOnUnknownComponent(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"model", "totally_made_up_component"})
	got, lossy := StateFromConfig(cfg)
	if !lossy {
		t.Errorf("expected lossy=true for unknown component")
	}
	// Unknown component should be skipped from features.
	if !equalFeatures(got.Features, []string{"model"}) {
		t.Errorf("Features should skip unknown: got %v want [model]", got.Features)
	}
}

func TestStateFromConfig_LossyOnCustomSeparator(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"model", "cost"})
	cfg.Separator.Character = "::"
	_, lossy := StateFromConfig(cfg)
	if !lossy {
		t.Errorf("expected lossy=true for non-default separator")
	}
}

func TestStateFromConfig_LossyOnCustomLayout(t *testing.T) {
	// model first then cost on a single hand-arranged line — wizard would
	// canonically place cost on a separate stats row, but here it's mixed.
	cfg := defaultConfigWithLines([]string{"cost", "model"})
	_, lossy := StateFromConfig(cfg)
	if !lossy {
		t.Errorf("expected lossy=true for non-canonical line order")
	}
}

func TestStateFromConfig_NotLossyForCanonicalLayout(t *testing.T) {
	// A natural single-row layout from a small feature set should NOT be lossy.
	cfg := defaultConfigWithLines([]string{"model", "cost"})
	_, lossy := StateFromConfig(cfg)
	if lossy {
		t.Errorf("expected lossy=false for canonical small layout")
	}
}

func TestStateFromConfig_FeatureOrderIsCanonical(t *testing.T) {
	// Hand-arranged line: stats before identity. Features should still come
	// out in canonical order (identity before stats).
	cfg := defaultConfigWithLines([]string{"cost", "model", "git_status"})
	got, _ := StateFromConfig(cfg)

	want := []string{"model", "git", "cost"}
	if !reflect.DeepEqual(got.Features, want) {
		t.Errorf("Features order: got %v want %v", got.Features, want)
	}
}

func TestStateFromConfig_LossyOnDuplicateFeatureComponents(t *testing.T) {
	// Two components mapping to the same feature is a hand-edit the wizard
	// can't represent — keep the first, mark lossy. (Layout comparison
	// catches it because the canonical layout collapses to one component.)
	cfg := defaultConfigWithLines([]string{"cache", "cache_hit"})
	got, lossy := StateFromConfig(cfg)

	if !lossy {
		t.Errorf("expected lossy=true for duplicate feature components")
	}
	if !equalFeatures(got.Features, []string{"cache"}) {
		t.Errorf("Features should dedupe to single 'cache': got %v", got.Features)
	}
	// The first occurrence wins.
	if got.CacheStyle != "counts" {
		t.Errorf("CacheStyle: got %q want %q (first occurrence wins)", got.CacheStyle, "counts")
	}
}

func TestStateFromConfig_ConfigDefaultRoundTrip(t *testing.T) {
	// The on-disk default (config.Default()) should also round-trip cleanly
	// — this is the path a user with an unmodified config takes.
	cfg := config.Default()
	got, lossy := StateFromConfig(cfg)

	if lossy {
		t.Fatalf("config.Default() should not be lossy")
	}

	want := DefaultState().Features
	if !equalFeatures(got.Features, want) {
		t.Errorf("Features: got %v want %v", got.Features, want)
	}
	if got.Theme != "default" {
		t.Errorf("Theme: got %q want default", got.Theme)
	}
	if got.Emojis != "all" {
		t.Errorf("Emojis: got %q want all", got.Emojis)
	}
}

// defaultConfigWithLines builds a config with default everything except the
// component lines, which are placed on a single line.
func defaultConfigWithLines(comps []string) *config.Config {
	return &config.Config{
		Theme:  "default",
		Emojis: config.EmojiAll,
		ContextBar: config.ContextBarConfig{
			Style:      config.BarBlock,
			Width:      10,
			Thresholds: []int{70, 90},
		},
		Separator: config.SeparatorConfig{Character: "|"},
		Lines:     []config.LineConfig{{Components: comps}},
	}
}
