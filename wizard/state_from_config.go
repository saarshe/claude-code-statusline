package wizard

import (
	"reflect"

	"github.com/saarshe/claude-code-statusline/components"
	"github.com/saarshe/claude-code-statusline/config"
)

// StateFromConfig builds a WizardState from a loaded Config. The second
// return value is true when the config has fields the wizard cannot fully
// round-trip — a custom layout, a non-default separator, or unknown
// components. Callers warn the user before saving in that case.
func StateFromConfig(cfg *config.Config) (*WizardState, bool) {
	state := DefaultState()
	state.Theme = cfg.Theme
	state.Emojis = string(cfg.Emojis)
	if cfg.ContextBar.Width > 0 {
		state.BarWidth = cfg.ContextBar.Width
	}

	lossy := false
	if cfg.Separator.Character != config.Default().Separator.Character {
		lossy = true
	}

	seen := make(map[string]bool)
	features := make([]string, 0, len(featureOrder))

	for _, line := range cfg.Lines {
		for _, comp := range line.Components {
			feature, style, ok := componentToFeature(comp, cfg.ContextBar.Style)
			if !ok {
				lossy = true
				continue
			}
			if seen[feature] {
				continue
			}
			seen[feature] = true
			features = append(features, feature)
			applyFeatureStyle(state, feature, style)
		}
	}

	// Reorder features by canonical order so layout inference is deterministic.
	state.Features = orderFeatures(features)

	state.InvalidateLayout()
	canonical := state.InferLayout()
	if !sameComponents(canonical, cfg.Lines) {
		lossy = true
	}

	return state, lossy
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

	for _, fm := range components.FeatureMeta {
		if fm.Key == comp {
			return comp, "", true
		}
	}

	return "", "", false
}

// barStyleToContextValue maps config.BarStyle back to the wizard style value
// used for the "context" feature when the component is "context_bar".
func barStyleToContextValue(s config.BarStyle) string {
	switch s {
	case config.BarSolid:
		return "solid"
	case config.BarASCII:
		return "ascii"
	default:
		return "block"
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

// orderFeatures returns features sorted by canonical featureOrder. Features
// not in featureOrder come last in their input order.
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
