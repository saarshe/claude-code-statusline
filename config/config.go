package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type EmojiMode string

const (
	EmojiAll    EmojiMode = "all"
	EmojiNone   EmojiMode = "none"
	EmojiCustom EmojiMode = "custom"
)

type BarStyle string

const (
	BarBlock    BarStyle = "block"
	BarSolid    BarStyle = "solid"
	BarASCII    BarStyle = "ascii"
	BarPercent  BarStyle = "percent"
	BarGradient BarStyle = "gradient"
)

type LineConfig struct {
	Components []string `toml:"components"`
}

type ContextBarConfig struct {
	Style      BarStyle `toml:"style"`
	Width      int      `toml:"width"`
	Thresholds []int    `toml:"thresholds"`
}

type SeparatorConfig struct {
	Character string `toml:"character"`
}

type Config struct {
	Theme      string           `toml:"theme"`
	Emojis     EmojiMode        `toml:"emojis"`
	ContextBar ContextBarConfig `toml:"context_bar"`
	Separator  SeparatorConfig  `toml:"separator"`
	Lines      []LineConfig     `toml:"line"`
}

func Default() *Config {
	return &Config{
		Theme:  "default",
		Emojis: EmojiAll,
		ContextBar: ContextBarConfig{
			Style:      BarBlock,
			Width:      10,
			Thresholds: []int{70, 90},
		},
		Separator: SeparatorConfig{
			Character: "|",
		},
		Lines: []LineConfig{
			{Components: []string{"model", "git_status"}},
			{Components: []string{"context_bar", "tokens", "cache", "cost"}},
		},
	}
}

func LoadFile(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return nil, err
	}

	// Decode into a separate struct to detect which fields were explicitly set.
	// We merge explicitly set fields over defaults.
	var raw rawConfig
	if _, err := toml.Decode(string(data), &raw); err != nil {
		return nil, fmt.Errorf("invalid TOML: %w", err)
	}

	if raw.Theme != "" {
		cfg.Theme = raw.Theme
	}
	if raw.Emojis != "" {
		cfg.Emojis = raw.Emojis
	}
	if raw.ContextBar.Style != "" {
		cfg.ContextBar.Style = raw.ContextBar.Style
	}
	if raw.ContextBar.Width != 0 {
		cfg.ContextBar.Width = raw.ContextBar.Width
	}
	if len(raw.ContextBar.Thresholds) != 0 {
		cfg.ContextBar.Thresholds = raw.ContextBar.Thresholds
	}
	if raw.Separator.Character != "" {
		cfg.Separator.Character = raw.Separator.Character
	}
	if len(raw.Lines) != 0 {
		cfg.Lines = raw.Lines
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// rawConfig mirrors Config but with zero values indicating "not set".
type rawConfig struct {
	Theme      string           `toml:"theme"`
	Emojis     EmojiMode        `toml:"emojis"`
	ContextBar ContextBarConfig `toml:"context_bar"`
	Separator  SeparatorConfig  `toml:"separator"`
	Lines      []LineConfig     `toml:"line"`
}

func validate(cfg *Config) error {
	switch cfg.Emojis {
	case EmojiAll, EmojiNone, EmojiCustom:
	default:
		return fmt.Errorf("invalid emojis value %q: must be 'all', 'none', or 'custom'", cfg.Emojis)
	}

	switch cfg.ContextBar.Style {
	case BarBlock, BarSolid, BarASCII, BarPercent, BarGradient:
	default:
		return fmt.Errorf("invalid context_bar.style %q: must be 'block', 'solid', 'ascii', 'percent', or 'gradient'", cfg.ContextBar.Style)
	}

	if len(cfg.ContextBar.Thresholds) != 2 {
		return fmt.Errorf("context_bar.thresholds must have exactly 2 values, got %d", len(cfg.ContextBar.Thresholds))
	}
	if cfg.ContextBar.Thresholds[0] >= cfg.ContextBar.Thresholds[1] {
		return fmt.Errorf("context_bar.thresholds[0] (%d) must be less than thresholds[1] (%d)",
			cfg.ContextBar.Thresholds[0], cfg.ContextBar.Thresholds[1])
	}

	return nil
}
