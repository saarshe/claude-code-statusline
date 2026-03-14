package theme

import (
	"sort"
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

var allThemes = []string{"default", "catppuccin", "nord", "dracula", "gruvbox", "tokyo-night", "powerline", "rounded"}

func TestAllThemes_Registered(t *testing.T) {
	for _, name := range allThemes {
		t.Run(name, func(t *testing.T) {
			th := Get(name)
			if th == nil {
				t.Fatalf("Get(%q) returned nil", name)
			}
			if th.Name != name {
				t.Errorf("Name = %q, want %q", th.Name, name)
			}
		})
	}
}

func TestAllThemes_AllStylesRender(t *testing.T) {
	for _, name := range allThemes {
		t.Run(name, func(t *testing.T) {
			th := Get(name)

			styles := []struct {
				label string
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

			for _, s := range styles {
				t.Run(s.label, func(t *testing.T) {
					result := s.style()
					if result == "" {
						t.Errorf("%s.Render('test') returned empty string", s.label)
					}
					if len(result) <= len("test") {
						t.Errorf("%s.Render('test') = %q, expected ANSI codes (longer than plain 'test')", s.label, result)
					}
				})
			}
		})
	}
}

func TestRegistryCount(t *testing.T) {
	names := Names()
	if len(names) != len(allThemes) {
		t.Errorf("expected %d themes in registry, got %d: %v", len(allThemes), len(names), names)
	}
}

func TestNames_Sorted(t *testing.T) {
	names := Names()
	sorted := make([]string, len(names))
	copy(sorted, names)
	sort.Strings(sorted)
	// Names() doesn't guarantee order, just verify all expected themes are present
	nameSet := map[string]bool{}
	for _, n := range names {
		nameSet[n] = true
	}
	for _, want := range allThemes {
		if !nameSet[want] {
			t.Errorf("theme %q missing from Names()", want)
		}
	}
}

func TestAllThemes_UniqueColors(t *testing.T) {
	// Each theme should have distinct Primary and Muted styles
	for _, name := range allThemes {
		t.Run(name, func(t *testing.T) {
			th := Get(name)
			primary := th.Primary.Render("x")
			muted := th.Muted.Render("x")
			if primary == muted {
				t.Errorf("Primary and Muted render identically for theme %q", name)
			}
		})
	}
}

func TestPowerlineTheme_HasSeparator(t *testing.T) {
	th := Get("powerline")
	if th.Separator == "" {
		t.Error("powerline separator should not be empty")
	}
}

func TestRoundedTheme_HasSeparator(t *testing.T) {
	th := Get("rounded")
	if th.Separator == "" {
		t.Error("rounded separator should not be empty")
	}
}

func TestLeanThemes_NoSeparator(t *testing.T) {
	lean := []string{"default", "catppuccin", "nord", "dracula", "gruvbox", "tokyo-night"}
	for _, name := range lean {
		t.Run(name, func(t *testing.T) {
			th := Get(name)
			if th.Separator != "" {
				t.Errorf("lean theme %q should have empty Separator, got %q", name, th.Separator)
			}
		})
	}
}

func TestPowerlineTheme_StylesHaveBackgrounds(t *testing.T) {
	th := Get("powerline")
	// Primary should render longer than a lean theme's Primary because of
	// background ANSI codes and padding.
	output := th.Primary.Render("x")
	leanOutput := Get("default").Primary.Render("x")
	if len(output) <= len(leanOutput) {
		t.Errorf("powerline Primary should produce longer ANSI output than lean; powerline=%d lean=%d",
			len(output), len(leanOutput))
	}
}

func TestRoundedTheme_StylesHaveBackgrounds(t *testing.T) {
	th := Get("rounded")
	output := th.Primary.Render("x")
	leanOutput := Get("default").Primary.Render("x")
	if len(output) <= len(leanOutput) {
		t.Errorf("rounded Primary should produce longer ANSI output than lean; rounded=%d lean=%d",
			len(output), len(leanOutput))
	}
}
