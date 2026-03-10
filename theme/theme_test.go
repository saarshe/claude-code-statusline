package theme

import (
	"testing"
)

func TestGet_Default(t *testing.T) {
	th := Get("default")
	if th == nil {
		t.Fatal("Get('default') returned nil")
	}
	if th.Name != "default" {
		t.Errorf("Name = %q, want %q", th.Name, "default")
	}
}

func TestGet_UnknownFallsBackToDefault(t *testing.T) {
	th := Get("nonexistent")
	if th == nil {
		t.Fatal("Get('nonexistent') returned nil, want fallback to default")
	}
	if th.Name != "default" {
		t.Errorf("Name = %q, want %q (fallback)", th.Name, "default")
	}
}

func TestGet_EmptyStringFallsBackToDefault(t *testing.T) {
	th := Get("")
	if th == nil {
		t.Fatal("Get('') returned nil, want fallback to default")
	}
	if th.Name != "default" {
		t.Errorf("Name = %q, want %q", th.Name, "default")
	}
}

func TestDefaultTheme_AllStylesRender(t *testing.T) {
	th := Get("default")

	tests := []struct {
		name  string
		style func() string
	}{
		{"Primary", func() string { return th.Primary.Render("test") }},
		{"Secondary", func() string { return th.Secondary.Render("test") }},
		{"Accent", func() string { return th.Accent.Render("test") }},
		{"Success", func() string { return th.Success.Render("test") }},
		{"Warning", func() string { return th.Warning.Render("test") }},
		{"Danger", func() string { return th.Danger.Render("test") }},
		{"Muted", func() string { return th.Muted.Render("test") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style()
			if result == "" {
				t.Errorf("%s.Render('test') returned empty string", tt.name)
			}
			if len(result) <= len("test") {
				t.Errorf("%s.Render('test') = %q, expected ANSI codes (longer than plain 'test')", tt.name, result)
			}
		})
	}
}

func TestRegistryCount(t *testing.T) {
	names := Names()
	if len(names) != 1 {
		t.Errorf("expected 1 theme in registry, got %d: %v", len(names), names)
	}
	if names[0] != "default" {
		t.Errorf("expected 'default' theme, got %q", names[0])
	}
}
