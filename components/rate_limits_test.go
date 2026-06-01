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
	}
	for _, tt := range tests {
		got := formatResetCountdown(tt.d)
		if got != tt.want {
			t.Errorf("formatResetCountdown(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}
