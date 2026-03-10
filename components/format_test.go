package components

import "testing"

func TestHumanizeTokens(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1.0k"},
		{1234, "1.2k"},
		{9999, "10.0k"},
		{10000, "10k"},
		{12345, "12k"},
		{999999, "999k"},
		{1000000, "1.0M"},
		{1234567, "1.2M"},
		{9999999, "10.0M"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := HumanizeTokens(tt.input)
			if got != tt.want {
				t.Errorf("HumanizeTokens(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
