package components

import "fmt"

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
