package wizard

import (
	"bytes"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
	"github.com/saarshe/claude-code-statusline/components"
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/render"
	"golang.org/x/term"
)

// WizardState holds all choices collected across wizard steps.
type WizardState struct {
	// Theme is the color theme name (e.g. "default", "catppuccin", "nord").
	Theme string

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

	// TokenStyle controls how token usage is displayed.
	// "turn" → tokens; "turn_cache" → tokens_cache; "session" → tokens_session; "full" → tokens_full.
	// Only used when "tokens" is in Features.
	TokenStyle string

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
		Theme:        "default",
		Features:     []string{"model", "git", "lines_changed", "directory", "agent", "worktree", "context", "tokens", "cache", "cost", "duration"},
		ContextStyle: "solid",
		TokenStyle:   "full",
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

// HasTokens reports whether the user selected the tokens feature.
func (s *WizardState) HasTokens() bool { return s.hasFeature("tokens") }

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

// featureStyleValue returns the user's chosen style value for a feature.
func (s *WizardState) featureStyleValue(feature string) string {
	switch feature {
	case "context":
		return s.ContextStyle
	case "tokens":
		return s.TokenStyle
	case "cache":
		return s.CacheStyle
	case "lines_changed":
		return s.LinesStyle
	case "git":
		return s.GitStyle
	default:
		return ""
	}
}

// featureToComponent maps a feature key to its component key, resolving style
// choices (e.g. "context" + ContextStyle="block" → "context_bar").
func (s *WizardState) featureToComponent(feature string) string {
	// Context is special: multiple style values map to the same component key.
	if feature == "context" {
		switch s.ContextStyle {
		case "pct":
			return "context_pct"
		case "tokens":
			return "context_tokens"
		case "tokens_bar":
			return "context_tokens_bar"
		case "tokens_bar_pct":
			return "context_tokens_bar_pct"
		default:
			return "context_bar"
		}
	}

	// For features with style options, look up via FeatureStyles.
	if styles, ok := components.FeatureStyles[feature]; ok {
		styleVal := s.featureStyleValue(feature)
		for _, so := range styles {
			if so.Value == styleVal {
				return string(so.ComponentKey)
			}
		}
		// Default to first option if style value doesn't match.
		return string(styles[0].ComponentKey)
	}

	return feature // all other feature keys match component keys 1:1
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

// featureOrder defines the canonical display order for all features.
// Identity features first, then stats — so the split point is natural
// when we need to wrap to two lines.
var featureOrder = append(append([]string{}, identityFeatures...), statsFeatures...)

// InferLayout places all selected components on a single line, then
// progressively splits the widest line in half until every line fits within
// the terminal width (or each line has a single component).
func (s *WizardState) InferLayout() [][]string {
	featureSet := make(map[string]bool, len(s.Features))
	for _, f := range s.Features {
		featureSet[f] = true
	}

	// Collect all components in canonical order.
	var all []string
	for _, f := range featureOrder {
		if featureSet[f] {
			all = append(all, s.featureToComponent(f))
		}
	}
	if len(all) == 0 {
		return nil
	}

	layout := [][]string{all}

	// Progressively split the widest overflowing line until everything fits.
	for layoutExceedsWidth(layout, s) {
		idx := widestLine(layout, s)
		if len(layout[idx]) <= 1 {
			break // can't split a single-component line
		}
		mid := len(layout[idx]) / 2
		left := layout[idx][:mid]
		right := layout[idx][mid:]
		// Replace the line with two halves.
		newLayout := make([][]string, 0, len(layout)+1)
		newLayout = append(newLayout, layout[:idx]...)
		newLayout = append(newLayout, left, right)
		newLayout = append(newLayout, layout[idx+1:]...)
		layout = newLayout
	}

	return layout
}

// widestLine returns the index of the line with the greatest rendered width.
func widestLine(layout [][]string, s *WizardState) int {
	cfg := s.toConfigWithLayout(layout)
	output := render.Render(MockInput(), cfg)
	lines := strings.Split(output, "\n")

	best, bestW := 0, 0
	for i, line := range lines {
		if i >= len(layout) {
			break
		}
		w := lipgloss.Width(line)
		if w > bestW {
			best, bestW = i, w
		}
	}
	return best
}

// termWidth returns the terminal width, defaulting to 80 if unavailable.
func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

// layoutExceedsWidth checks if any line in the layout would be wider than the
// terminal when rendered with mock data and the current wizard state.
func layoutExceedsWidth(layout [][]string, s *WizardState) bool {
	cfg := s.toConfigWithLayout(layout)
	output := render.Render(MockInput(), cfg)
	tw := termWidth()
	for _, line := range strings.Split(output, "\n") {
		if lipgloss.Width(line) > tw {
			return true
		}
	}
	return false
}

// ToConfig converts the wizard state to a *config.Config.
func (s *WizardState) ToConfig() *config.Config {
	return s.toConfigWithLayout(s.InferLayout())
}

// toConfigWithLayout builds a config using the given layout (list of component
// rows) instead of calling InferLayout. This avoids recursion when
// layoutExceedsWidth needs to render a candidate layout.
func (s *WizardState) toConfigWithLayout(layout [][]string) *config.Config {
	cfg := config.Default()
	if s.Theme != "" {
		cfg.Theme = s.Theme
	}
	cfg.Emojis = config.EmojiMode(s.Emojis)
	if s.HasContext() && s.ContextStyle != "pct" {
		cfg.ContextBar.Style = s.contextBarStyle()
	}
	if s.BarWidth > 0 {
		cfg.ContextBar.Width = s.BarWidth
	}

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
