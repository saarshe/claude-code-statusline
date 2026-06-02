package components

import (
	"strings"
	"testing"
	"time"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

// fixedNow returns a stable time so countdown strings are deterministic.
func withFixedNow(t *testing.T, now time.Time) {
	t.Helper()
	prev := rateLimitsNow
	rateLimitsNow = func() time.Time { return now }
	t.Cleanup(func() { rateLimitsNow = prev })
}

func TestRateLimits_BothWindows(t *testing.T) {
	data := &schema.Input{RateLimits: &schema.RateLimits{
		FiveHour: &schema.RateLimitWindow{UsedPercentage: 23},
		SevenDay: &schema.RateLimitWindow{UsedPercentage: 41},
	}}
	result := Get("rate_limits").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "5h 23%") {
		t.Errorf("expected '5h 23%%', got %q", result)
	}
	if !strings.Contains(result, "7d 41%") {
		t.Errorf("expected '7d 41%%', got %q", result)
	}
}

func TestRateLimits_OnlyFiveHour(t *testing.T) {
	data := &schema.Input{RateLimits: &schema.RateLimits{
		FiveHour: &schema.RateLimitWindow{UsedPercentage: 50},
	}}
	result := Get("rate_limits").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "5h 50%") {
		t.Errorf("expected '5h 50%%', got %q", result)
	}
	if strings.Contains(result, "7d") {
		t.Errorf("did not expect 7d window, got %q", result)
	}
}

func TestRateLimits_EmptyWhenNil(t *testing.T) {
	data := &schema.Input{}
	result := Get("rate_limits").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestRateLimits_EmptyWhenBothWindowsAbsent(t *testing.T) {
	data := &schema.Input{RateLimits: &schema.RateLimits{}}
	result := Get("rate_limits").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty when both windows absent, got %q", result)
	}
}

func TestRateLimits_OnlySevenDay(t *testing.T) {
	data := &schema.Input{RateLimits: &schema.RateLimits{
		SevenDay: &schema.RateLimitWindow{UsedPercentage: 60},
	}}
	result := Get("rate_limits").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "7d 60%") {
		t.Errorf("expected '7d 60%%', got %q", result)
	}
	if strings.Contains(result, "5h") {
		t.Errorf("did not expect 5h window, got %q", result)
	}
}

func TestRateLimits_ThresholdColors(t *testing.T) {
	// Default thresholds are [70, 90]: green < 70 ≤ yellow < 90 ≤ red.
	cases := []struct {
		pct    float64
		ansi   string
		label  string
	}{
		{30, "\033[32m", "green"},
		{75, "\033[33m", "yellow"},
		{95, "\033[31m", "red"},
	}
	for _, tc := range cases {
		data := &schema.Input{RateLimits: &schema.RateLimits{
			FiveHour: &schema.RateLimitWindow{UsedPercentage: tc.pct},
		}}
		result := Get("rate_limits").Render(data, config.Default(), theme.Get("default"))
		if !strings.Contains(result, tc.ansi) {
			t.Errorf("pct=%v: expected %s color (%q) in output, got %q", tc.pct, tc.label, tc.ansi, result)
		}
	}
}

func TestRateLimitsReset_OmitsCountdownWhenResetsAtZero(t *testing.T) {
	// ResetsAt is documented as Unix epoch seconds; treat 0 as "not provided"
	// and render the percentage without "in …".
	data := &schema.Input{RateLimits: &schema.RateLimits{
		FiveHour: &schema.RateLimitWindow{UsedPercentage: 23, ResetsAt: 0},
	}}
	result := Get("rate_limits_reset").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "5h 23%") {
		t.Errorf("expected percentage to still render, got %q", result)
	}
	if strings.Contains(result, " in ") {
		t.Errorf("did not expect countdown when ResetsAt=0, got %q", result)
	}
}

func TestRateLimitsReset_IncludesCountdown(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	withFixedNow(t, now)

	data := &schema.Input{RateLimits: &schema.RateLimits{
		FiveHour: &schema.RateLimitWindow{
			UsedPercentage: 23,
			ResetsAt:       now.Add(2*time.Hour + 14*time.Minute).Unix(),
		},
		SevenDay: &schema.RateLimitWindow{
			UsedPercentage: 41,
			ResetsAt:       now.Add(3 * 24 * time.Hour).Unix(),
		},
	}}
	result := Get("rate_limits_reset").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "2h14m") {
		t.Errorf("expected '2h14m' countdown, got %q", result)
	}
	if !strings.Contains(result, "3d") {
		t.Errorf("expected '3d' countdown, got %q", result)
	}
}

func TestRateLimits_EmojiOff_UsesTextPrefix(t *testing.T) {
	data := &schema.Input{RateLimits: &schema.RateLimits{
		FiveHour: &schema.RateLimitWindow{UsedPercentage: 23},
	}}
	cfg := config.Default()
	result := Get("rate_limits").Render(data, cfg, theme.Get("default"))
	if !strings.Contains(result, "📈") {
		t.Errorf("expected 📈 emoji when emojis on, got %q", result)
	}

	cfg.Emojis = config.EmojiNone
	result = Get("rate_limits").Render(data, cfg, theme.Get("default"))
	if strings.Contains(result, "📈") {
		t.Errorf("expected no emoji when disabled, got %q", result)
	}
	if !strings.Contains(result, "Limits: ") {
		t.Errorf("expected 'Limits: ' text prefix when emojis off, got %q", result)
	}
}

func TestFormatResetCountdown(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{-time.Second, "now"},
		{0, "now"},
		{45 * time.Second, "45s"},
		{45 * time.Minute, "45m"},
		{2*time.Hour + 14*time.Minute, "2h14m"},
		{2 * time.Hour, "2h"},
		{25 * time.Hour, "1d1h"},
		{3 * 24 * time.Hour, "3d"},
		// Rounding boundary: truncates (does not round up). 1h59m59s stays "1h59m".
		{time.Hour + 59*time.Minute + 59*time.Second, "1h59m"},
	}
	for _, tt := range tests {
		got := formatResetCountdown(tt.d)
		if got != tt.want {
			t.Errorf("formatResetCountdown(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}
