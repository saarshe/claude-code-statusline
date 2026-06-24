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
		return FlowComponents(input, cfg, th, flat, columnsWidth())
	}

	rows := make([][]string, len(cfg.Lines))
	for i, line := range cfg.Lines {
		rows[i] = line.Components
	}
	return rows
}

// FlowComponents arranges comps across lines so that, when rendered, no line is
// wider than width. It starts with everything on one line and repeatedly splits
// the widest overflowing line in half until every line fits (or is reduced to a
// single component, which cannot be split further). Width is measured by
// actually rendering each candidate row, so the result matches real output.
func FlowComponents(input *schema.Input, cfg *config.Config, th *theme.Theme, comps []string, width int) [][]string {
	if len(comps) == 0 {
		return nil
	}

	sep := separator(cfg, th)
	layout := [][]string{comps}

	for {
		widest, widestW := -1, 0
		exceeds := false
		for i, row := range layout {
			w := lipgloss.Width(renderRow(input, cfg, th, sep, row))
			if w > width {
				exceeds = true
			}
			if w > widestW {
				widest, widestW = i, w
			}
		}
		if !exceeds || widest < 0 || len(layout[widest]) <= 1 {
			break
		}

		mid := len(layout[widest]) / 2
		left := layout[widest][:mid]
		right := layout[widest][mid:]
		newLayout := make([][]string, 0, len(layout)+1)
		newLayout = append(newLayout, layout[:widest]...)
		newLayout = append(newLayout, left, right)
		newLayout = append(newLayout, layout[widest+1:]...)
		layout = newLayout
	}

	return layout
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
