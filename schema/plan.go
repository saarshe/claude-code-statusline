package schema

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// Plan holds the user's Claude subscription plan. Unlike most fields it does
// not come from the statusline JSON on stdin — Claude Code does not send the
// plan there. It is read at runtime from Claude Code's global config file
// (~/.claude.json) via PopulatePlan, mirroring how Git is populated.
type Plan struct {
	OrgType       string // oauthAccount.organizationType, e.g. "claude_max"
	RateLimitTier string // oauthAccount.organizationRateLimitTier, e.g. "default_claude_max_20x"
}

// claudeConfigPath returns the path to Claude Code's global config file.
// It honors CLAUDE_CONFIG_DIR when set, otherwise falls back to ~/.claude.json.
// Returns "" if the home directory cannot be determined.
func claudeConfigPath() string {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, ".claude.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude.json")
}

// PopulatePlan reads the subscription plan from Claude Code's config file and
// stores it in the Plan field. Any failure (missing file, malformed JSON,
// renamed fields) leaves Plan zero-valued so the plan component renders nothing
// rather than breaking the status line.
func (i *Input) PopulatePlan() {
	path := claudeConfigPath()
	if path == "" {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	i.Plan = parsePlan(f)
}

// parsePlan extracts just the plan fields from a Claude config JSON stream.
// It deliberately reads only organizationType and organizationRateLimitTier —
// no email, name, or organization name — even though the file contains them.
func parsePlan(r io.Reader) Plan {
	var doc struct {
		OAuthAccount struct {
			OrganizationType          string `json:"organizationType"`
			OrganizationRateLimitTier string `json:"organizationRateLimitTier"`
		} `json:"oauthAccount"`
	}
	if err := json.NewDecoder(r).Decode(&doc); err != nil {
		return Plan{}
	}
	return Plan{
		OrgType:       doc.OAuthAccount.OrganizationType,
		RateLimitTier: doc.OAuthAccount.OrganizationRateLimitTier,
	}
}
