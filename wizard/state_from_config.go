package wizard

import (
	"reflect"
	"strings"

	"github.com/saarshe/claude-code-statusline/components"
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/theme"
)

// StateFromConfig builds a WizardState from a loaded Config. The second
// return value is a list of human-readable reasons why the config can't be
// fully round-tripped through the wizard — empty when no information will be
// lost on save.
//
// Callers use len(reasons) > 0 as a "should we warn the user before saving?"
// signal, and show the reasons themselves in the confirmation dialog.
func StateFromConfig(cfg *config.Config) (*WizardState, []string) {
	state := DefaultState()
	state.Theme = cfg.Theme
	state.Emojis = string(cfg.Emojis)
	state.Layout = string(cfg.Layout)
	if state.Layout == "" {
		state.Layout = string(config.LayoutFixed)
	}
	if cfg.ContextBar.Width > 0 {
		state.BarWidth = cfg.ContextBar.Width
	}

	var reasons []string

	if !themeKnown(cfg.Theme) {
		reasons = append(reasons, "custom theme name")
	}
	if cfg.Emojis != config.EmojiAll && cfg.Emojis != config.EmojiNone {
		reasons = append(reasons, "custom emoji settings")
	}
	if !reflect.DeepEqual(cfg.ContextBar.Thresholds, config.Default().ContextBar.Thresholds) {
		reasons = append(reasons, "context bar thresholds")
	}
	if cfg.ContextBar.Width <= 0 {
		reasons = append(reasons, "context bar width")
	}
	if cfg.Separator.Character != config.Default().Separator.Character {
		reasons = append(reasons, "custom separator")
	}

	seen := make(map[string]bool)
	features := make([]string, 0, len(featureOrder))
	hasUnknown := false
	hasBarStyleMismatch := false

	for _, line := range cfg.Lines {
		for _, comp := range line.Components {
			feature, style, ok := componentToFeature(comp, cfg.ContextBar.Style)
			if !ok {
				hasUnknown = true
				continue
			}
			if isContextComponent(comp) && !barStyleRepresentable(comp, cfg.ContextBar.Style) {
				hasBarStyleMismatch = true
			}
			if seen[feature] {
				continue
			}
			seen[feature] = true
			features = append(features, feature)
			applyFeatureStyle(state, feature, style)
		}
	}

	if hasUnknown {
		reasons = append(reasons, "unknown components")
	}
	if hasBarStyleMismatch {
		reasons = append(reasons, "context bar style")
	}

	state.Features = orderFeatures(features)

	// In auto mode the renderer reflows at runtime, so the config stores a
	// single ordered list rather than pre-split lines — comparing it against
	// the inferred split would always (falsely) look hand-edited.
	if state.Layout != string(config.LayoutAuto) {
		state.InvalidateLayout()
		canonical := state.InferLayout()
		if !sameComponents(canonical, cfg.Lines) {
			reasons = append(reasons, "hand-edited layout")
		}
	}

	return state, reasons
}

// themeKnown reports whether the given theme name is registered.
func themeKnown(name string) bool {
	for _, n := range theme.Names() {
		if n == name {
			return true
		}
	}
	return false
}

// isContextComponent reports whether the component belongs to the context
// family — i.e. its rendering is influenced by ContextBar.Style.
func isContextComponent(comp string) bool {
	switch comp {
	case "context_bar", "context_pct", "context_tokens", "context_tokens_bar", "context_tokens_bar_pct":
		return true
	}
	return false
}

// barStyleRepresentable reports whether the wizard can preserve the given
// ContextBar.Style when paired with the given context component. Bar styles
// that don't survive the round-trip add to the lossy reasons.
func barStyleRepresentable(comp string, s config.BarStyle) bool {
	switch comp {
	case "context_bar":
		// All bar styles except BarPercent are representable as a wizard
		// ContextStyle. BarBlock and BarSolid render identically (see
		// components/context_bar.go), so we treat them as equivalent.
		return s == config.BarBlock || s == config.BarSolid || s == config.BarASCII || s == config.BarGradient
	case "context_tokens_bar", "context_tokens_bar_pct":
		// These components always render their bar as a gradient.
		return s == config.BarGradient
	case "context_pct", "context_tokens":
		// These don't render a bar at all, so style is irrelevant.
		return true
	}
	return true
}

// componentToFeature reverses featureToComponent: given a component key, find
// the wizard feature it belongs to and the style value that selected it.
// For context_bar, the style value depends on the ContextBar.Style field.
func componentToFeature(comp string, barStyle config.BarStyle) (feature, style string, ok bool) {
	switch comp {
	case "context_bar":
		return "context", barStyleToContextValue(barStyle), true
	case "context_pct":
		return "context", "pct", true
	case "context_tokens":
		return "context", "tokens", true
	case "context_tokens_bar":
		return "context", "tokens_bar", true
	case "context_tokens_bar_pct":
		return "context", "tokens_bar_pct", true
	}

	for f, styles := range components.FeatureStyles {
		for _, so := range styles {
			if string(so.ComponentKey) == comp {
				return f, so.Value, true
			}
		}
	}

	// 1:1 fallback: feature key == component key, for components that have
	// no style variants (model, cost, duration, etc.). We restrict the match
	// to a known whitelist of such keys so that a literal "git" or "context"
	// (which are feature keys but not real components) doesn't slip through.
	if isOneToOneComponent(comp) {
		return comp, "", true
	}

	return "", "", false
}

// isOneToOneComponent returns true for components whose feature key equals
// their component key (no style variants). These are valid component names.
func isOneToOneComponent(comp string) bool {
	switch comp {
	case "model", "cost", "duration", "directory", "agent", "worktree", "session_id", "effort", "plan", "pr":
		return true
	}
	return false
}

// barStyleToContextValue maps config.BarStyle back to the wizard style value
// used for the "context" feature when the component is "context_bar". The
// wizard has no "block" option — BarBlock renders identically to BarSolid in
// components/context_bar.go, so we collapse both to "solid".
func barStyleToContextValue(s config.BarStyle) string {
	switch s {
	case config.BarASCII:
		return "ascii"
	case config.BarGradient:
		return "gradient"
	case config.BarSolid, config.BarBlock:
		return "solid"
	default:
		// BarPercent or unknown — pick the closest visually similar choice.
		// barStyleRepresentable will independently flag this as lossy.
		return "solid"
	}
}

// applyFeatureStyle writes a style value into the WizardState field that
// corresponds to the feature.
func applyFeatureStyle(s *WizardState, feature, style string) {
	if style == "" {
		return
	}
	switch feature {
	case "context":
		s.ContextStyle = style
	case "tokens":
		s.TokenStyle = style
	case "cache":
		s.CacheStyle = style
	case "lines_changed":
		s.LinesStyle = style
	case "git":
		s.GitStyle = style
	case "rate_limits":
		s.RateLimitsStyle = style
	}
}

// orderFeatures returns features sorted by canonical featureOrder.
func orderFeatures(features []string) []string {
	known := make(map[string]bool, len(features))
	for _, f := range features {
		known[f] = true
	}
	out := make([]string, 0, len(features))
	for _, f := range featureOrder {
		if known[f] {
			out = append(out, f)
		}
	}
	return out
}

// sameComponents reports whether the flat component sequences match. Line
// wrapping is terminal-width-adaptive and the wizard always re-infers it on
// save, so we ignore where rows break and only compare the underlying order.
func sameComponents(canonical [][]string, cfgLines []config.LineConfig) bool {
	var a, b []string
	for _, row := range canonical {
		a = append(a, row...)
	}
	for _, l := range cfgLines {
		b = append(b, l.Components...)
	}
	return reflect.DeepEqual(a, b)
}

// FormatLossyReasons turns a list of lossy reasons into a single-line
// description suitable for a confirmation dialog.
func FormatLossyReasons(reasons []string) string {
	if len(reasons) == 0 {
		return ""
	}
	return "The following won't be preserved: " + strings.Join(reasons, ", ") + "."
}
