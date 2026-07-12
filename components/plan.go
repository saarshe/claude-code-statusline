package components

import (
	"regexp"
	"strings"

	"github.com/saarshe/claude-code-statusline/config"
	"github.com/saarshe/claude-code-statusline/schema"
	"github.com/saarshe/claude-code-statusline/theme"
)

type planComponent struct{}

func init() { Register(&planComponent{}) }

func (p *planComponent) Key() ComponentKey { return "plan" }

func (p *planComponent) Render(data *schema.Input, cfg *config.Config, th *theme.Theme) string {
	label := planLabel(data.Plan.OrgType, data.Plan.RateLimitTier)
	if label == "" {
		return ""
	}
	return th.Accent.Render(GetMeta(p.Key()).Prefix(cfg) + label)
}

// tierSuffixRe pulls a trailing usage multiplier (e.g. "20x", "5x") off the
// rate-limit tier string, which looks like "default_claude_max_20x".
var tierSuffixRe = regexp.MustCompile(`(\d+x)$`)

// planLabel turns Claude's raw plan fields into a display label like "Max 20x"
// or "Pro". Returns "" when the plan family is absent so the component renders
// nothing (e.g. when ~/.claude.json is missing or its fields were renamed).
func planLabel(orgType, rateLimitTier string) string {
	family := planFamily(orgType)
	if family == "" {
		return ""
	}
	if m := tierSuffixRe.FindStringSubmatch(rateLimitTier); m != nil {
		return family + " " + m[1]
	}
	return family
}

// planFamily strips the "claude_" prefix and title-cases the plan family:
// "claude_max" -> "Max", "claude_enterprise" -> "Enterprise".
func planFamily(orgType string) string {
	s := strings.TrimPrefix(orgType, "claude_")
	if s == "" {
		return ""
	}
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}
