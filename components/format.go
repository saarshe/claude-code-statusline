package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/theme"
)

// EmojiPrefix returns "emoji " when emojis are enabled, or textFallback otherwise.
// This eliminates the repeated if-block pattern across all components.
func EmojiPrefix(cfg *config.Config, emoji, textFallback string) string {
	if cfg.Emojis != config.EmojiNone {
		return emoji + " "
	}
	return textFallback
}

// Clamp restricts v to the range [lo, hi].
func Clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ContextStyle returns a lipgloss.Style based on how full the context window is.
// thresholds should be [warn, danger] (e.g. [70, 90]).
func ContextStyle(th *theme.Theme, pct float64, thresholds []int) lipgloss.Style {
	if len(thresholds) == 2 {
		if pct >= float64(thresholds[1]) {
			return th.Danger
		}
		if pct >= float64(thresholds[0]) {
			return th.Warning
		}
	}
	return th.Success
}

// HumanizeTokens formats a token count into a human-readable string.
// < 1000: "999"
// 1000-9999: "1.2k" (one decimal)
// 10000-999999: "12k" (no decimal)
// >= 1000000: "1.2M" (one decimal) or "12M"
func HumanizeTokens(n int) string {
	switch {
	case n < 1000:
		return fmt.Sprintf("%d", n)
	case n < 10000:
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	case n < 1000000:
		return fmt.Sprintf("%dk", n/1000)
	case n < 10000000:
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	default:
		return fmt.Sprintf("%dM", n/1000000)
	}
}
