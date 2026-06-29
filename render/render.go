package render

import (
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/saarshe/claude-code-statusline/components"
	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

func Render(input *schema.Input, cfg *config.Config) string {
	th := theme.Get(cfg.Theme)
	return RenderWithTheme(input, cfg, th)
}

func RenderWithTheme(input *schema.Input, cfg *config.Config, th *theme.Theme) string {
	sep := separator(cfg, th)

	rows := rowsForLayout(input, cfg, th)

	lineOutputs := make([]string, 0, len(rows))
	for _, comps := range rows {
		if line := renderRow(input, cfg, th, sep, comps); line != "" {
			lineOutputs = append(lineOutputs, line)
		}
	}

	return strings.Join(lineOutputs, "\n")
}

// rowsForLayout returns the component rows to render. In fixed mode these are
// the configured lines verbatim; in auto mode all components are flattened into
// one ordered list and reflowed to fit the live terminal width.
func rowsForLayout(input *schema.Input, cfg *config.Config, th *theme.Theme) [][]string {
	if cfg.Layout == config.LayoutAuto {
		var flat []string
		for _, line := range cfg.Lines {
			flat = append(flat, line.Components...)
		}
		return FlowComponents(input, cfg, th, flat, UsableWidth(columnsWidth()))
	}

	rows := make([][]string, len(cfg.Lines))
	for i, line := range cfg.Lines {
		rows[i] = line.Components
	}
	return rows
}

// FlowComponents arranges comps across lines so that, when rendered, no line is
// wider than width. It packs greedily in order: each component joins the
// current line if it still fits, otherwise it starts a new line. This fills
// each line to capacity and minimizes line count (a single component wider than
// width sits alone, since it cannot be split). Width is measured by actually
// rendering each candidate line, so the result matches real output.
func FlowComponents(input *schema.Input, cfg *config.Config, th *theme.Theme, comps []string, width int) [][]string {
	if len(comps) == 0 {
		return nil
	}

	sep := separator(cfg, th)
	var rows [][]string
	var cur []string

	for _, key := range comps {
		candidate := append(append([]string{}, cur...), key)
		if len(cur) > 0 && lipgloss.Width(renderRow(input, cfg, th, sep, candidate)) > width {
			rows = append(rows, cur)
			cur = []string{key}
			continue
		}
		cur = candidate
	}
	if len(cur) > 0 {
		rows = append(rows, cur)
	}

	return rows
}

// renderRow renders one line of components joined by the separator. Components
// that produce no output are skipped; an all-empty row yields "".
func renderRow(input *schema.Input, cfg *config.Config, th *theme.Theme, sep string, comps []string) string {
	parts := make([]string, 0, len(comps))
	for _, key := range comps {
		c := components.Get(key)
		if c == nil {
			continue
		}
		if result := c.Render(input, cfg, th); result != "" {
			parts = append(parts, result)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, sep)
}

// separator resolves the styled separator, preferring a theme-provided one.
func separator(cfg *config.Config, th *theme.Theme) string {
	if th.Separator != "" {
		return th.Muted.Render(th.Separator)
	}
	return th.Muted.Render(" " + cfg.Separator.Character + " ")
}

// LayoutSafetyMargin is the number of columns kept free when reflowing, so
// lines never pack to the exact terminal edge. The host renders the status
// line in slightly less room than COLUMNS reports (a reserved final column,
// and emoji glyphs that can occupy marginally more space than measured), which
// would otherwise truncate a line measured at exactly the full width.
const LayoutSafetyMargin = 2

// UsableWidth reduces a raw terminal width by the safety margin, never going
// below 1.
func UsableWidth(raw int) int {
	if w := raw - LayoutSafetyMargin; w >= 1 {
		return w
	}
	return 1
}

// columnsWidth returns the live terminal width from the COLUMNS environment
// variable, which Claude Code sets before invoking the status line (v2.1.153+).
// It falls back to 80 when COLUMNS is unset or unparseable — e.g. older Claude
// Code, or when the binary is run outside Claude Code.
func columnsWidth() int {
	if v := strings.TrimSpace(os.Getenv("COLUMNS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 80
}
