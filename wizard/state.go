package wizard

import (
	"bytes"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/saars/claude-code-statusline/config"
)

// WizardState holds all choices collected across wizard steps.
type WizardState struct {
	// Features is the set of high-level data categories the user wants to see.
	// Valid values: "model", "context", "tokens", "cache", "cost", "duration",
	// "git_branch", "git_status", "lines_changed", "directory".
	Features []string

	// ContextStyle controls how the context window is displayed.
	// "pct" → context_pct; "block"/"solid"/"ascii" → context_bar.
	// Only used when "context" is in Features.
	ContextStyle string

	// CacheStyle controls how cache stats are displayed.
	// "hit" → cache_hit (efficiency %); "counts" → cache (raw token counts).
	// Only used when "cache" is in Features.
	CacheStyle string

	// LinesStyle controls how lines changed are displayed.
	// "summary" → lines_summary (±total); "detail" → lines_changed (+N -M).
	// Only used when "lines_changed" is in Features.
	LinesStyle string

	// GitStyle controls how git information is displayed.
	// "branch" → git_branch (branch name only); "status" → git_status (branch + file counts).
	// Only used when "git" is in Features.
	GitStyle string

	// Emojis is "all" or "none".
	Emojis string

	// BarWidth is the character width of progress bars (default 10).
	BarWidth int
}

// identityFeatures are shown in the first row (who/where am I).
var identityFeatures = []string{"model", "git", "lines_changed", "directory", "agent", "worktree"}

// statsFeatures are shown in the second row (numbers/metrics).
var statsFeatures = []string{"context", "tokens", "cache", "cost", "duration"}

// DefaultState returns a WizardState that matches config.Default().
func DefaultState() *WizardState {
	return &WizardState{
		Features:     []string{"model", "git", "context", "tokens", "cache", "cost"},
		ContextStyle: "solid",
		CacheStyle:   "counts",
		LinesStyle:   "detail",
		GitStyle:     "status",
		Emojis:       "all",
		BarWidth:     10,
	}
}

// HasContext reports whether the user selected the context window feature.
func (s *WizardState) HasContext() bool { return s.hasFeature("context") }

// HasCache reports whether the user selected the cache feature.
func (s *WizardState) HasCache() bool { return s.hasFeature("cache") }

// HasLines reports whether the user selected the lines_changed feature.
func (s *WizardState) HasLines() bool { return s.hasFeature("lines_changed") }

// HasGit reports whether the user selected the git feature.
func (s *WizardState) HasGit() bool { return s.hasFeature("git") }

func (s *WizardState) hasFeature(key string) bool {
	for _, f := range s.Features {
		if f == key {
			return true
		}
	}
	return false
}

// featureToComponent maps a feature key to its component key, resolving style
// choices (e.g. "context" + ContextStyle="block" → "context_bar").
func (s *WizardState) featureToComponent(feature string) string {
	switch feature {
	case "context":
		switch s.ContextStyle {
		case "pct":
			return "context_pct"
		case "tokens":
			return "context_tokens"
		case "tokens_bar":
			return "context_tokens_bar"
		default:
			return "context_bar"
		}
	case "cache":
		if s.CacheStyle == "hit" {
			return "cache_hit"
		}
		return "cache"
	case "lines_changed":
		if s.LinesStyle == "summary" {
			return "lines_summary"
		}
		return "lines_changed"
	case "git":
		if s.GitStyle == "branch" {
			return "git_branch"
		}
		return "git_status"
	default:
		return feature // all other feature keys match component keys 1:1
	}
}

// contextBarStyle maps ContextStyle to a config.BarStyle.
func (s *WizardState) contextBarStyle() config.BarStyle {
	switch s.ContextStyle {
	case "solid":
		return config.BarSolid
	case "ascii":
		return config.BarASCII
	case "gradient":
		return config.BarGradient
	default:
		return config.BarBlock
	}
}

// InferLayout distributes the selected features into rows:
//   - Row 1 (identity): model, git, lines_changed, directory
//   - Row 2 (metrics):  context, tokens, cache, cost, duration
//
// If only one category has selections, a single row is returned.
func (s *WizardState) InferLayout() [][]string {
	featureSet := make(map[string]bool, len(s.Features))
	for _, f := range s.Features {
		featureSet[f] = true
	}

	var row1, row2 []string
	for _, f := range identityFeatures {
		if featureSet[f] {
			row1 = append(row1, s.featureToComponent(f))
		}
	}
	for _, f := range statsFeatures {
		if featureSet[f] {
			row2 = append(row2, s.featureToComponent(f))
		}
	}

	switch {
	case len(row1) > 0 && len(row2) > 0:
		return [][]string{row1, row2}
	case len(row1) > 0:
		return [][]string{row1}
	case len(row2) > 0:
		return [][]string{row2}
	default:
		return nil
	}
}

// ToConfig converts the wizard state to a *config.Config.
func (s *WizardState) ToConfig() *config.Config {
	cfg := config.Default()
	cfg.Emojis = config.EmojiMode(s.Emojis)
	if s.HasContext() && s.ContextStyle != "pct" {
		cfg.ContextBar.Style = s.contextBarStyle()
	}
	if s.BarWidth > 0 {
		cfg.ContextBar.Width = s.BarWidth
	}

	layout := s.InferLayout()
	cfg.Lines = make([]config.LineConfig, len(layout))
	for i, comps := range layout {
		cfg.Lines[i] = config.LineConfig{Components: comps}
	}
	return cfg
}

// tomlConfig is an encodable mirror of config.Config.
type tomlConfig struct {
	Theme      string                  `toml:"theme"`
	Emojis     string                  `toml:"emojis"`
	ContextBar config.ContextBarConfig `toml:"context_bar"`
	Separator  config.SeparatorConfig  `toml:"separator"`
	Lines      []config.LineConfig     `toml:"line"`
}

// ToTOML encodes the wizard state as a TOML string for writing to the config
// file.
func (s *WizardState) ToTOML() (string, error) {
	cfg := s.ToConfig()
	tc := tomlConfig{
		Theme:      cfg.Theme,
		Emojis:     string(cfg.Emojis),
		ContextBar: cfg.ContextBar,
		Separator:  cfg.Separator,
		Lines:      cfg.Lines,
	}

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(tc); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()) + "\n", nil
}
