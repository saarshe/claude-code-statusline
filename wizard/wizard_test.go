package wizard

import (
	"os"
	"strings"
	"testing"

	"github.com/saarshe/claude-code-statusline/config"
)

func TestWizardState_HasContext_True(t *testing.T) {
	state := &WizardState{Features: []string{"model", "context", "cost"}}
	if !state.HasContext() {
		t.Error("expected HasContext() = true")
	}
}

func TestWizardState_HasContext_False(t *testing.T) {
	state := &WizardState{Features: []string{"model", "cost"}}
	if state.HasContext() {
		t.Error("expected HasContext() = false")
	}
}

func TestWizardState_ToConfig_IdentityOnly_SingleRow(t *testing.T) {
	state := &WizardState{
		Features:  []string{"model", "git"},
		GitStyle:  "status",
		Emojis:    "all",
	}
	cfg := state.ToConfig()
	if len(cfg.Lines) != 1 {
		t.Errorf("expected 1 row for identity-only, got %d", len(cfg.Lines))
	}
}

func TestWizardState_ToConfig_StatsOnly_SingleRow(t *testing.T) {
	state := &WizardState{
		Features:     []string{"context", "cost"},
		ContextStyle: "pct",
		Emojis:       "all",
	}
	cfg := state.ToConfig()
	if len(cfg.Lines) != 1 {
		t.Errorf("expected 1 row for stats-only, got %d", len(cfg.Lines))
	}
}

func TestWizardState_ToConfig_FewFeatures_SingleRow(t *testing.T) {
	state := &WizardState{
		Features:     []string{"model", "cost"},
		ContextStyle: "pct",
		Emojis:       "all",
	}
	cfg := state.ToConfig()
	if len(cfg.Lines) != 1 {
		t.Errorf("expected 1 row for few features, got %d", len(cfg.Lines))
	}
	// Components should be in canonical order: identity first, then stats.
	if cfg.Lines[0].Components[0] != "model" {
		t.Errorf("expected 'model' first, got %q", cfg.Lines[0].Components[0])
	}
	if cfg.Lines[0].Components[1] != "cost" {
		t.Errorf("expected 'cost' second, got %q", cfg.Lines[0].Components[1])
	}
}

func TestWizardState_ToConfig_ContextPct(t *testing.T) {
	state := &WizardState{
		Features:     []string{"context"},
		ContextStyle: "pct",
	}
	cfg := state.ToConfig()
	row := cfg.Lines[0].Components
	found := false
	for _, c := range row {
		if c == "context_pct" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'context_pct' component, got %v", row)
	}
}

func TestWizardState_ToConfig_ContextBar(t *testing.T) {
	state := &WizardState{
		Features:     []string{"context"},
		ContextStyle: "block",
	}
	cfg := state.ToConfig()
	row := cfg.Lines[0].Components
	found := false
	for _, c := range row {
		if c == "context_bar" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'context_bar' component, got %v", row)
	}
	if cfg.ContextBar.Style != config.BarBlock {
		t.Errorf("expected BarBlock style, got %q", cfg.ContextBar.Style)
	}
}

func TestWizardState_ToConfig_EmojisNone(t *testing.T) {
	state := &WizardState{
		Features: []string{"model"},
		Emojis:   "none",
	}
	cfg := state.ToConfig()
	if cfg.Emojis != config.EmojiNone {
		t.Errorf("expected EmojiNone, got %q", cfg.Emojis)
	}
}

func TestWizardState_ToTOML_RoundTrip(t *testing.T) {
	state := &WizardState{
		Features:     []string{"model", "git_branch", "context", "cost"},
		ContextStyle: "solid",
		Emojis:       "all",
		BarWidth:     12,
	}

	tomlStr, err := state.ToTOML()
	if err != nil {
		t.Fatalf("ToTOML() error: %v", err)
	}
	if tomlStr == "" {
		t.Fatal("ToTOML() returned empty string")
	}

	f, err := os.CreateTemp(t.TempDir(), "config*.toml")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	f.WriteString(tomlStr)
	f.Close()

	loaded, err := config.LoadFile(f.Name())
	if err != nil {
		t.Fatalf("LoadFile() error: %v\nTOML:\n%s", err, tomlStr)
	}
	if len(loaded.Lines) == 0 {
		t.Errorf("expected at least 1 line, got 0\nTOML:\n%s", tomlStr)
	}
	if loaded.ContextBar.Style != config.BarSolid {
		t.Errorf("expected BarSolid, got %q", loaded.ContextBar.Style)
	}
}

func TestWizardState_ToTOML_ContainsComponents(t *testing.T) {
	state := &WizardState{
		Features: []string{"model", "cost"},
		Emojis:   "none",
	}

	tomlStr, err := state.ToTOML()
	if err != nil {
		t.Fatalf("ToTOML() error: %v", err)
	}
	if !strings.Contains(tomlStr, "model") {
		t.Error("expected TOML to contain 'model'")
	}
	if !strings.Contains(tomlStr, "cost") {
		t.Error("expected TOML to contain 'cost'")
	}
}

func TestWizardState_DefaultState(t *testing.T) {
	state := DefaultState()
	if len(state.Features) == 0 {
		t.Error("DefaultState() should have at least one feature")
	}
	if state.Emojis == "" {
		t.Error("DefaultState() should have non-empty Emojis")
	}
	if state.ContextStyle == "" {
		t.Error("DefaultState() should have non-empty ContextStyle")
	}
	if state.TokenStyle == "" {
		t.Error("DefaultState() should have non-empty TokenStyle")
	}
	if state.CacheStyle == "" {
		t.Error("DefaultState() should have non-empty CacheStyle")
	}
	if state.LinesStyle == "" {
		t.Error("DefaultState() should have non-empty LinesStyle")
	}
	if state.GitStyle == "" {
		t.Error("DefaultState() should have non-empty GitStyle")
	}
	if state.BarWidth <= 0 {
		t.Error("DefaultState() should have positive BarWidth")
	}
}

func TestWizardState_CacheHit_UsesCorrectComponent(t *testing.T) {
	state := &WizardState{
		Features:   []string{"cache"},
		CacheStyle: "hit",
		Emojis:     "all",
	}
	cfg := state.ToConfig()
	row := cfg.Lines[0].Components
	found := false
	for _, c := range row {
		if c == "cache_hit" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'cache_hit' component, got %v", row)
	}
}

func TestWizardState_CacheCounts_UsesCorrectComponent(t *testing.T) {
	state := &WizardState{
		Features:   []string{"cache"},
		CacheStyle: "counts",
		Emojis:     "all",
	}
	cfg := state.ToConfig()
	row := cfg.Lines[0].Components
	found := false
	for _, c := range row {
		if c == "cache" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'cache' component, got %v", row)
	}
}

func TestWizardState_LinesSummary_UsesCorrectComponent(t *testing.T) {
	state := &WizardState{
		Features:   []string{"lines_changed"},
		LinesStyle: "summary",
		Emojis:     "all",
	}
	cfg := state.ToConfig()
	row := cfg.Lines[0].Components
	found := false
	for _, c := range row {
		if c == "lines_summary" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'lines_summary' component, got %v", row)
	}
}

func TestWizardState_HasCache_True(t *testing.T) {
	state := &WizardState{Features: []string{"model", "cache"}}
	if !state.HasCache() {
		t.Error("expected HasCache() = true")
	}
}

func TestWizardState_HasCache_False(t *testing.T) {
	state := &WizardState{Features: []string{"model", "cost"}}
	if state.HasCache() {
		t.Error("expected HasCache() = false")
	}
}

func TestWizardState_HasLines_True(t *testing.T) {
	state := &WizardState{Features: []string{"lines_changed"}}
	if !state.HasLines() {
		t.Error("expected HasLines() = true")
	}
}

func TestWizardState_HasGit_True(t *testing.T) {
	state := &WizardState{Features: []string{"model", "git"}}
	if !state.HasGit() {
		t.Error("expected HasGit() = true")
	}
}

func TestWizardState_HasGit_False(t *testing.T) {
	state := &WizardState{Features: []string{"model", "cost"}}
	if state.HasGit() {
		t.Error("expected HasGit() = false")
	}
}

func TestWizardState_GitStyle_Branch(t *testing.T) {
	state := &WizardState{
		Features: []string{"git"},
		GitStyle: "branch",
	}
	cfg := state.ToConfig()
	row := cfg.Lines[0].Components
	found := false
	for _, c := range row {
		if c == "git_branch" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'git_branch' component, got %v", row)
	}
}

func TestWizardState_GitStyle_Status(t *testing.T) {
	state := &WizardState{
		Features: []string{"git"},
		GitStyle: "status",
	}
	cfg := state.ToConfig()
	row := cfg.Lines[0].Components
	found := false
	for _, c := range row {
		if c == "git_status" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'git_status' component, got %v", row)
	}
}

func TestWizardState_TokenStyle_Turn(t *testing.T) {
	state := &WizardState{Features: []string{"tokens"}, TokenStyle: "turn"}
	cfg := state.ToConfig()
	if cfg.Lines[0].Components[0] != "tokens" {
		t.Errorf("expected 'tokens', got %v", cfg.Lines[0].Components)
	}
}

func TestWizardState_TokenStyle_TurnCache(t *testing.T) {
	state := &WizardState{Features: []string{"tokens"}, TokenStyle: "turn_cache"}
	cfg := state.ToConfig()
	if cfg.Lines[0].Components[0] != "tokens_cache" {
		t.Errorf("expected 'tokens_cache', got %v", cfg.Lines[0].Components)
	}
}

func TestWizardState_TokenStyle_Session(t *testing.T) {
	state := &WizardState{Features: []string{"tokens"}, TokenStyle: "session"}
	cfg := state.ToConfig()
	if cfg.Lines[0].Components[0] != "tokens_session" {
		t.Errorf("expected 'tokens_session', got %v", cfg.Lines[0].Components)
	}
}

func TestWizardState_TokenStyle_Full(t *testing.T) {
	state := &WizardState{Features: []string{"tokens"}, TokenStyle: "full"}
	cfg := state.ToConfig()
	if cfg.Lines[0].Components[0] != "tokens_full" {
		t.Errorf("expected 'tokens_full', got %v", cfg.Lines[0].Components)
	}
}

func TestWizardState_HasTokens(t *testing.T) {
	state := &WizardState{Features: []string{"tokens"}}
	if !state.HasTokens() {
		t.Error("expected HasTokens() = true")
	}
	state2 := &WizardState{Features: []string{"model"}}
	if state2.HasTokens() {
		t.Error("expected HasTokens() = false")
	}
}

func TestWizardState_NoFeaturesSelected(t *testing.T) {
	state := &WizardState{Features: []string{}, Emojis: "all"}
	cfg := state.ToConfig()
	// Should not panic; lines may be empty
	_ = cfg
}

func TestWizardState_Theme_DefaultState(t *testing.T) {
	state := DefaultState()
	if state.Theme != "default" {
		t.Errorf("DefaultState().Theme = %q, want %q", state.Theme, "default")
	}
}

func TestWizardState_Theme_ToConfig(t *testing.T) {
	state := &WizardState{
		Features: []string{"model"},
		Theme:    "dracula",
		Emojis:   "all",
	}
	cfg := state.ToConfig()
	if cfg.Theme != "dracula" {
		t.Errorf("ToConfig().Theme = %q, want %q", cfg.Theme, "dracula")
	}
}

func TestWizardState_Theme_ToTOML(t *testing.T) {
	state := &WizardState{
		Features: []string{"model"},
		Theme:    "nord",
		Emojis:   "all",
	}
	tomlStr, err := state.ToTOML()
	if err != nil {
		t.Fatalf("ToTOML() error: %v", err)
	}
	if !strings.Contains(tomlStr, `theme = "nord"`) {
		t.Errorf("expected TOML to contain theme = \"nord\", got:\n%s", tomlStr)
	}
}

func TestPrevRunnableStep_GoesToPreviousStep(t *testing.T) {
	steps := []Step{
		{Run: func(*WizardState) error { return nil }},
		{Run: func(*WizardState) error { return nil }},
		{Run: func(*WizardState) error { return nil }},
	}
	state := DefaultState()
	got := prevRunnableStep(steps, 2, state)
	if got != 1 {
		t.Errorf("prevRunnableStep(steps, 2) = %d, want 1", got)
	}
}

func TestPrevRunnableStep_SkipsNonRunnable(t *testing.T) {
	never := func(*WizardState) bool { return false }
	steps := []Step{
		{Run: func(*WizardState) error { return nil }},
		{ShouldRun: never, Run: func(*WizardState) error { return nil }},
		{Run: func(*WizardState) error { return nil }},
	}
	state := DefaultState()
	got := prevRunnableStep(steps, 2, state)
	if got != 0 {
		t.Errorf("prevRunnableStep(steps, 2) = %d, want 0 (skipping non-runnable)", got)
	}
}

func TestPrevRunnableStep_StaysOnFirstStep(t *testing.T) {
	steps := []Step{
		{Run: func(*WizardState) error { return nil }},
		{Run: func(*WizardState) error { return nil }},
	}
	state := DefaultState()
	got := prevRunnableStep(steps, 0, state)
	if got != 0 {
		t.Errorf("prevRunnableStep(steps, 0) = %d, want 0", got)
	}
}

func TestPrevRunnableStep_AllPreviousSkipped(t *testing.T) {
	never := func(*WizardState) bool { return false }
	steps := []Step{
		{ShouldRun: never, Run: func(*WizardState) error { return nil }},
		{ShouldRun: never, Run: func(*WizardState) error { return nil }},
		{Run: func(*WizardState) error { return nil }},
	}
	state := DefaultState()
	got := prevRunnableStep(steps, 2, state)
	if got != 2 {
		t.Errorf("prevRunnableStep(steps, 2) = %d, want 2 (no runnable previous step)", got)
	}
}

func TestWizardState_Theme_RoundTrip(t *testing.T) {
	state := &WizardState{
		Features: []string{"model", "cost"},
		Theme:    "catppuccin",
		Emojis:   "all",
	}
	tomlStr, err := state.ToTOML()
	if err != nil {
		t.Fatalf("ToTOML() error: %v", err)
	}

	f, err := os.CreateTemp(t.TempDir(), "config*.toml")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	f.WriteString(tomlStr)
	f.Close()

	loaded, err := config.LoadFile(f.Name())
	if err != nil {
		t.Fatalf("LoadFile() error: %v\nTOML:\n%s", err, tomlStr)
	}
	if loaded.Theme != "catppuccin" {
		t.Errorf("loaded.Theme = %q, want %q", loaded.Theme, "catppuccin")
	}
}
