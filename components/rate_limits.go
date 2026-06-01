package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

// rateLimitsNow is overridable from tests for deterministic countdown output.
var rateLimitsNow = time.Now

type rateLimitsComponent struct {
	key      ComponentKey
	withTime bool
}

func init() {
	Register(&rateLimitsComponent{key: "rate_limits", withTime: false})
	Register(&rateLimitsComponent{key: "rate_limits_reset", withTime: true})
}

func (r *rateLimitsComponent) Key() ComponentKey { return r.key }

func (r *rateLimitsComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	if data.RateLimits == nil {
		return ""
	}

	parts := []string{}
	now := rateLimitsNow()

	if w := data.RateLimits.FiveHour; w != nil {
		parts = append(parts, formatRateWindow("5h", w, r.withTime, now, cfg.ContextBar.Thresholds, th))
	}
	if w := data.RateLimits.SevenDay; w != nil {
		parts = append(parts, formatRateWindow("7d", w, r.withTime, now, cfg.ContextBar.Thresholds, th))
	}
	if len(parts) == 0 {
		return ""
	}

	return GetMeta(r.Key()).Prefix(cfg) + strings.Join(parts, " · ")
}

func formatRateWindow(label string, w *schema.RateLimitWindow, withTime bool, now time.Time, thresholds []int, th *theme.Theme) string {
	style := ContextStyle(th, w.UsedPercentage, thresholds)
	body := fmt.Sprintf("%s %.0f%%", label, w.UsedPercentage)
	if withTime && w.ResetsAt > 0 {
		body += " in " + formatResetCountdown(time.Unix(w.ResetsAt, 0).Sub(now))
	}
	return style.Render(body)
}

// formatResetCountdown renders a time.Duration as a compact countdown.
// Examples: "2h14m", "45m", "3d", "23h", "now".
func formatResetCountdown(d time.Duration) string {
	if d <= 0 {
		return "now"
	}
	secs := int(d / time.Second)
	days := secs / 86400
	hours := (secs % 86400) / 3600
	mins := (secs % 3600) / 60

	if days > 0 {
		if hours > 0 {
			return fmt.Sprintf("%dd%dh", days, hours)
		}
		return fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		if mins > 0 {
			return fmt.Sprintf("%dh%dm", hours, mins)
		}
		return fmt.Sprintf("%dh", hours)
	}
	if mins > 0 {
		return fmt.Sprintf("%dm", mins)
	}
	return fmt.Sprintf("%ds", secs)
}
