package wizard

import (
	"reflect"
	"sort"
	"strings"
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

// reasonsContain reports whether reasons contains a substring match.
func reasonsContain(reasons []string, want string) bool {
	for _, r := range reasons {
		if strings.Contains(r, want) {
			return true
		}
	}
	return false
}

func TestStateFromConfig_DefaultRoundTrip(t *testing.T) {
	want := DefaultState()
	got, reasons := StateFromConfig(want.ToConfig())

	if len(reasons) > 0 {
		t.Fatalf("default-state round trip should be lossless, got reasons: %v", reasons)
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
		{"solid", "context_bar", config.BarSolid, "solid"},
		{"block_to_solid", "context_bar", config.BarBlock, "solid"},
		{"ascii", "context_bar", config.BarASCII, "ascii"},
		{"gradient", "context_bar", config.BarGradient, "gradient"},
		{"pct", "context_pct", "", "pct"},
		{"tokens", "context_tokens", "", "tokens"},
		{"tokens_bar", "context_tokens_bar", config.BarGradient, "tokens_bar"},
		{"tokens_bar_pct", "context_tokens_bar_pct", config.BarGradient, "tokens_bar_pct"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := defaultConfigWithLines([]string{tc.component})
			if tc.barStyle != "" {
				cfg.ContextBar.Style = tc.barStyle
			}
			got, reasons := StateFromConfig(cfg)
			if len(reasons) > 0 {
				t.Fatalf("unexpected reasons: %v", reasons)
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
			got, reasons := StateFromConfig(cfg)
			if len(reasons) > 0 {
				t.Fatalf("unexpected reasons: %v", reasons)
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

	got, reasons := StateFromConfig(cfg)
	if len(reasons) > 0 {
		t.Fatalf("unexpected reasons: %v", reasons)
	}
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
	got, reasons := StateFromConfig(cfg)
	if len(reasons) > 0 {
		t.Fatalf("unexpected reasons: %v", reasons)
	}
	if got.BarWidth != 25 {
		t.Errorf("BarWidth: got %d want 25", got.BarWidth)
	}
}

func TestStateFromConfig_LossyOnUnknownComponent(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"model", "totally_made_up_component"})
	got, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "unknown") {
		t.Errorf("expected reasons to mention unknown components, got %v", reasons)
	}
	if !equalFeatures(got.Features, []string{"model"}) {
		t.Errorf("Features should skip unknown: got %v want [model]", got.Features)
	}
}

func TestStateFromConfig_LossyOnCustomSeparator(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"model", "cost"})
	cfg.Separator.Character = "::"
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "separator") {
		t.Errorf("expected reasons to mention separator, got %v", reasons)
	}
}

func TestStateFromConfig_LossyOnCustomLayout(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"cost", "model"})
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "layout") {
		t.Errorf("expected reasons to mention layout, got %v", reasons)
	}
}

func TestStateFromConfig_NotLossyForCanonicalLayout(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"model", "cost"})
	_, reasons := StateFromConfig(cfg)
	if len(reasons) > 0 {
		t.Errorf("expected no reasons for canonical small layout, got %v", reasons)
	}
}

func TestStateFromConfig_FeatureOrderIsCanonical(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"cost", "model", "git_status"})
	got, _ := StateFromConfig(cfg)

	want := []string{"model", "git", "cost"}
	if !reflect.DeepEqual(got.Features, want) {
		t.Errorf("Features order: got %v want %v", got.Features, want)
	}
}

func TestStateFromConfig_LossyOnDuplicateFeatureComponents(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"cache", "cache_hit"})
	got, reasons := StateFromConfig(cfg)

	if !reasonsContain(reasons, "layout") {
		t.Errorf("expected reasons to mention layout (duplicate caught by component mismatch), got %v", reasons)
	}
	if !equalFeatures(got.Features, []string{"cache"}) {
		t.Errorf("Features should dedupe: got %v", got.Features)
	}
	if got.CacheStyle != "counts" {
		t.Errorf("CacheStyle: got %q want %q (first occurrence wins)", got.CacheStyle, "counts")
	}
}

func TestStateFromConfig_ConfigDefaultRoundTrip(t *testing.T) {
	cfg := config.Default()
	got, reasons := StateFromConfig(cfg)

	if len(reasons) > 0 {
		t.Fatalf("config.Default() should be lossless, got reasons: %v", reasons)
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

func TestStateFromConfig_LossyOnUnknownTheme(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"model"})
	cfg.Theme = "my_custom_fork_theme"
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "theme") {
		t.Errorf("expected reasons to mention theme, got %v", reasons)
	}
}

func TestStateFromConfig_LossyOnEmojiCustom(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"model"})
	cfg.Emojis = config.EmojiCustom
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "emoji") {
		t.Errorf("expected reasons to mention emoji, got %v", reasons)
	}
}

func TestStateFromConfig_LossyOnCustomThresholds(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"context_bar"})
	cfg.ContextBar.Thresholds = []int{50, 80}
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "threshold") {
		t.Errorf("expected reasons to mention thresholds, got %v", reasons)
	}
}

func TestStateFromConfig_LossyOnBarPercent(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"context_bar"})
	cfg.ContextBar.Style = config.BarPercent
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "bar style") {
		t.Errorf("expected reasons to mention bar style for BarPercent, got %v", reasons)
	}
}

func TestStateFromConfig_LossyOnTokensBarWithNonGradientStyle(t *testing.T) {
	// context_tokens_bar always renders gradient; any other ContextBar.Style
	// is silently re-coerced on save, so flag it.
	cfg := defaultConfigWithLines([]string{"context_tokens_bar"})
	cfg.ContextBar.Style = config.BarBlock
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "bar style") {
		t.Errorf("expected reasons to mention bar style for tokens_bar+BarBlock, got %v", reasons)
	}
}

func TestStateFromConfig_LossyOnWidthZero(t *testing.T) {
	cfg := defaultConfigWithLines([]string{"context_bar"})
	cfg.ContextBar.Width = 0
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "width") {
		t.Errorf("expected reasons to mention width for width=0, got %v", reasons)
	}
}

func TestStateFromConfig_FeatureKeyAsComponent_IsRejected(t *testing.T) {
	// "git" is a feature key but not a real component (the real ones are
	// git_branch / git_status). Should be treated as unknown.
	cfg := defaultConfigWithLines([]string{"model", "git"})
	_, reasons := StateFromConfig(cfg)
	if !reasonsContain(reasons, "unknown") {
		t.Errorf("expected reasons to mention unknown for feature-key-as-component, got %v", reasons)
	}
}

func TestSaveConfirmDescription(t *testing.T) {
	cases := []struct {
		name              string
		replacingExisting bool
		reasons           []string
		want              string
	}{
		{"fresh-no-existing", false, nil, ""},
		{"fresh-replacing", true, nil, "This will replace your existing config."},
		{"patch-lossless", false, nil, ""},
		{"patch-lossy", true, []string{"custom theme name"}, "Replacing existing config. The following won't be preserved: custom theme name."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := saveConfirmDescription(tc.replacingExisting, tc.reasons); got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestFormatLossyReasons(t *testing.T) {
	if got := FormatLossyReasons(nil); got != "" {
		t.Errorf("empty reasons should format to empty string, got %q", got)
	}
	got := FormatLossyReasons([]string{"custom theme name", "hand-edited layout"})
	want := "The following won't be preserved: custom theme name, hand-edited layout."
	if got != want {
		t.Errorf("FormatLossyReasons: got %q want %q", got, want)
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
