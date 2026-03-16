package components

import "github.com/saarshe/claude-code-statusline/config"

// Meta holds display metadata for a component or feature.
type Meta struct {
	Emoji      string // prefix when emojis enabled (e.g. "📊")
	TextPrefix string // prefix when emojis disabled (e.g. "Cache: ")
	TextSuffix string // suffix when emojis disabled (e.g. "]" for model)
	Name       string // human-readable label (used by wizard)
	Desc       string // short description (used by wizard)
}

// Prefix returns the display prefix for this meta entry based on emoji mode.
func (m Meta) Prefix(cfg *config.Config) string {
	if cfg.Emojis != config.EmojiNone && m.Emoji != "" {
		return m.Emoji + " "
	}
	return m.TextPrefix
}

// Suffix returns the display suffix (only meaningful when emojis are off).
func (m Meta) Suffix(cfg *config.Config) string {
	if cfg.Emojis != config.EmojiNone {
		return ""
	}
	return m.TextSuffix
}

var componentMeta = map[ComponentKey]Meta{
	"model":              {Emoji: "🤖", TextPrefix: "[", TextSuffix: "]", Name: "Model name", Desc: "which Claude model is active"},
	"context_pct":        {Emoji: "📊", Name: "Context %", Desc: "how full the context is"},
	"context_bar":        {Emoji: "📊", Name: "Context bar", Desc: "visual context usage bar"},
	"context_tokens":     {Emoji: "📊", Name: "Context tokens", Desc: "token counts with context size"},
	"context_tokens_bar": {Emoji: "📊", Name: "Context tokens + bar", Desc: "token counts with visual bar"},
	"tokens":             {Emoji: "🎟️", TextPrefix: "Tok: ", Name: "Per-turn totals", Desc: "total input and output this turn"},
	"tokens_cache":       {Emoji: "🎟️", TextPrefix: "Tok: ", Name: "Per-turn + cache", Desc: "per-turn totals with cache hit rate"},
	"tokens_session":     {Emoji: "🎟️", TextPrefix: "Tok: ", Name: "Session output", Desc: "cumulative output tokens for the session"},
	"tokens_full":        {Emoji: "🎟️", TextPrefix: "Tok: ", Name: "Full breakdown", Desc: "per-turn with cache + session output total"},
	"cache":              {Emoji: "💾", TextPrefix: "Cache: ", Name: "Cache counts", Desc: "tokens reused vs stored"},
	"cache_hit":          {Emoji: "⚡", TextPrefix: "Cache: ", Name: "Cache efficiency", Desc: "cache hit percentage"},
	"cost":               {Emoji: "💰", Name: "Cost", Desc: "total session cost in USD"},
	"duration":           {Emoji: "⏱️", TextPrefix: "Time: ", Name: "Duration", Desc: "total session time"},
	"git_branch":         {Emoji: "🌿", Name: "Git branch", Desc: "current branch name"},
	"git_status":         {Emoji: "🌿", Name: "Git status", Desc: "branch name and file change counts"},
	"lines_changed":      {Emoji: "📝", Name: "Lines changed", Desc: "lines added / removed"},
	"lines_summary":      {Emoji: "📝", Name: "Lines summary", Desc: "total lines touched"},
	"directory":          {Emoji: "📁", Name: "Directory", Desc: "current working directory"},
	"agent":              {Emoji: "🤖", Name: "Agent name", Desc: "shown when running as a sub-agent"},
	"worktree":           {Emoji: "🌿", Name: "Worktree", Desc: "shown when working in a git worktree"},
}

// GetMeta returns the display metadata for a component key.
func GetMeta(key ComponentKey) Meta {
	return componentMeta[key]
}

// StyleOption maps a wizard style value to its resolved component key.
type StyleOption struct {
	Value        string       // wizard style value (e.g. "turn", "hit")
	ComponentKey ComponentKey // resolved component key (e.g. "tokens", "cache_hit")
}

// FeatureStyles maps feature keys to their available style options.
// Features not listed here have a 1:1 mapping (feature key == component key).
// Context is excluded because its styles don't map 1:1 to component keys.
var FeatureStyles = map[string][]StyleOption{
	"tokens": {
		{"turn", "tokens"},
		{"turn_cache", "tokens_cache"},
		{"session", "tokens_session"},
		{"full", "tokens_full"},
	},
	"cache": {
		{"hit", "cache_hit"},
		{"counts", "cache"},
	},
	"git": {
		{"branch", "git_branch"},
		{"status", "git_status"},
	},
	"lines_changed": {
		{"summary", "lines_summary"},
		{"detail", "lines_changed"},
	},
}

// FeatureMeta maps wizard feature keys (which may group multiple components)
// to display metadata. The wizard uses these for its selection UI.
var FeatureMeta = []struct {
	Key  string
	Meta Meta
}{
	{"model", componentMeta["model"]},
	{"context", Meta{Emoji: "📊", Name: "Context window", Desc: "how full the context is"}},
	{"tokens", componentMeta["tokens"]},
	{"cache", Meta{Emoji: "💾", Name: "Cache", Desc: "how much context is reused vs freshly processed"}},
	{"cost", componentMeta["cost"]},
	{"duration", componentMeta["duration"]},
	{"git", Meta{Emoji: "🌿", Name: "Git", Desc: "branch name and file change counts"}},
	{"lines_changed", componentMeta["lines_changed"]},
	{"directory", componentMeta["directory"]},
	{"agent", componentMeta["agent"]},
	{"worktree", componentMeta["worktree"]},
}
