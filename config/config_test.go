package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.toml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Theme != "default" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "default")
	}
	if cfg.Emojis != EmojiAll {
		t.Errorf("Emojis = %q, want %q", cfg.Emojis, EmojiAll)
	}
	if cfg.ContextBar.Style != BarBlock {
		t.Errorf("ContextBar.Style = %q, want %q", cfg.ContextBar.Style, BarBlock)
	}
	if cfg.ContextBar.Width != 10 {
		t.Errorf("ContextBar.Width = %d, want 10", cfg.ContextBar.Width)
	}
	if len(cfg.ContextBar.Thresholds) != 2 || cfg.ContextBar.Thresholds[0] != 70 || cfg.ContextBar.Thresholds[1] != 90 {
		t.Errorf("ContextBar.Thresholds = %v, want [70, 90]", cfg.ContextBar.Thresholds)
	}
	if cfg.Separator.Character != "|" {
		t.Errorf("Separator.Character = %q, want %q", cfg.Separator.Character, "|")
	}
	if len(cfg.Lines) != 2 {
		t.Fatalf("Lines = %d, want 2", len(cfg.Lines))
	}
	if len(cfg.Lines[0].Components) == 0 {
		t.Error("Lines[0] has no components")
	}
	if len(cfg.Lines[1].Components) == 0 {
		t.Error("Lines[1] has no components")
	}
}

func TestLoadFile_Valid(t *testing.T) {
	path := writeTemp(t, `
theme = "default"
emojis = "none"

[context_bar]
style = "ascii"
width = 20
thresholds = [60, 85]

[separator]
character = "•"

[[line]]
components = ["model", "cost"]

[[line]]
components = ["context_pct"]
`)

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Emojis != EmojiNone {
		t.Errorf("Emojis = %q, want %q", cfg.Emojis, EmojiNone)
	}
	if cfg.ContextBar.Style != BarASCII {
		t.Errorf("ContextBar.Style = %q, want %q", cfg.ContextBar.Style, BarASCII)
	}
	if cfg.ContextBar.Width != 20 {
		t.Errorf("ContextBar.Width = %d, want 20", cfg.ContextBar.Width)
	}
	if cfg.ContextBar.Thresholds[0] != 60 || cfg.ContextBar.Thresholds[1] != 85 {
		t.Errorf("ContextBar.Thresholds = %v, want [60, 85]", cfg.ContextBar.Thresholds)
	}
	if cfg.Separator.Character != "•" {
		t.Errorf("Separator.Character = %q, want %q", cfg.Separator.Character, "•")
	}
	if len(cfg.Lines) != 2 {
		t.Fatalf("Lines = %d, want 2", len(cfg.Lines))
	}
	if cfg.Lines[0].Components[0] != "model" {
		t.Errorf("Lines[0].Components[0] = %q, want %q", cfg.Lines[0].Components[0], "model")
	}
}

func TestLoadFile_MissingFile_ReturnsDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.toml")

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}

	// Should be identical to Default()
	def := Default()
	if cfg.Theme != def.Theme {
		t.Errorf("Theme = %q, want %q", cfg.Theme, def.Theme)
	}
	if cfg.Emojis != def.Emojis {
		t.Errorf("Emojis = %q, want %q", cfg.Emojis, def.Emojis)
	}
	if len(cfg.Lines) != len(def.Lines) {
		t.Errorf("Lines count = %d, want %d", len(cfg.Lines), len(def.Lines))
	}
}

func TestLoadFile_MissingFields_FilledWithDefaults(t *testing.T) {
	path := writeTemp(t, `theme = "default"`)

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Fields not in TOML should be defaults
	if cfg.Emojis != EmojiAll {
		t.Errorf("Emojis = %q, want %q (default)", cfg.Emojis, EmojiAll)
	}
	if cfg.ContextBar.Style != BarBlock {
		t.Errorf("ContextBar.Style = %q, want %q (default)", cfg.ContextBar.Style, BarBlock)
	}
	if cfg.ContextBar.Width != 10 {
		t.Errorf("ContextBar.Width = %d, want 10 (default)", cfg.ContextBar.Width)
	}
	if len(cfg.Lines) == 0 {
		t.Error("Lines should be non-empty (default layout)")
	}
}

func TestLoadFile_InvalidEmojiMode(t *testing.T) {
	path := writeTemp(t, `emojis = "banana"`)

	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for invalid emoji mode, got nil")
	}
}

func TestLoadFile_InvalidBarStyle(t *testing.T) {
	path := writeTemp(t, `
[context_bar]
style = "invalid"
`)

	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for invalid bar style, got nil")
	}
}

func TestLoadFile_InvalidTOML(t *testing.T) {
	path := writeTemp(t, `not valid toml = = =`)

	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for invalid TOML, got nil")
	}
}

func TestLoadFile_ThresholdsReversed(t *testing.T) {
	path := writeTemp(t, `
[context_bar]
thresholds = [90, 70]
`)

	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for reversed thresholds, got nil")
	}
}

func TestLoadFile_ThresholdsWrongCount(t *testing.T) {
	path := writeTemp(t, `
[context_bar]
thresholds = [70]
`)

	_, err := LoadFile(path)
	if err == nil {
		t.Error("expected error for wrong threshold count, got nil")
	}
}

func TestEmojiMode_Constants(t *testing.T) {
	if EmojiAll != "all" {
		t.Errorf("EmojiAll = %q, want %q", EmojiAll, "all")
	}
	if EmojiNone != "none" {
		t.Errorf("EmojiNone = %q, want %q", EmojiNone, "none")
	}
	if EmojiCustom != "custom" {
		t.Errorf("EmojiCustom = %q, want %q", EmojiCustom, "custom")
	}
}

func TestBarStyle_Constants(t *testing.T) {
	if BarBlock != "block" {
		t.Errorf("BarBlock = %q, want %q", BarBlock, "block")
	}
	if BarSolid != "solid" {
		t.Errorf("BarSolid = %q, want %q", BarSolid, "solid")
	}
	if BarASCII != "ascii" {
		t.Errorf("BarASCII = %q, want %q", BarASCII, "ascii")
	}
	if BarPercent != "percent" {
		t.Errorf("BarPercent = %q, want %q", BarPercent, "percent")
	}
}
