package schema

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParsePlan_Full(t *testing.T) {
	json := `{
		"userID": "abc",
		"oauthAccount": {
			"emailAddress": "x@y.z",
			"organizationType": "claude_max",
			"organizationRateLimitTier": "default_claude_max_20x"
		}
	}`
	plan := parsePlan(strings.NewReader(json))
	if plan.OrgType != "claude_max" {
		t.Errorf("OrgType = %q, want claude_max", plan.OrgType)
	}
	if plan.RateLimitTier != "default_claude_max_20x" {
		t.Errorf("RateLimitTier = %q, want default_claude_max_20x", plan.RateLimitTier)
	}
}

func TestParsePlan_MissingAccount(t *testing.T) {
	plan := parsePlan(strings.NewReader(`{"userID": "abc"}`))
	if plan.OrgType != "" || plan.RateLimitTier != "" {
		t.Errorf("expected zero Plan, got %+v", plan)
	}
}

func TestParsePlan_Malformed(t *testing.T) {
	plan := parsePlan(strings.NewReader(`{not json`))
	if plan.OrgType != "" || plan.RateLimitTier != "" {
		t.Errorf("expected zero Plan on malformed JSON, got %+v", plan)
	}
}

func TestPopulatePlan_ReadsFromConfigDir(t *testing.T) {
	dir := t.TempDir()
	json := `{"oauthAccount": {"organizationType": "claude_pro"}}`
	if err := os.WriteFile(filepath.Join(dir, ".claude.json"), []byte(json), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_CONFIG_DIR", dir)

	input := &Input{}
	input.PopulatePlan()
	if input.Plan.OrgType != "claude_pro" {
		t.Errorf("Plan.OrgType = %q, want claude_pro", input.Plan.OrgType)
	}
}

func TestPopulatePlan_MissingFileLeavesZero(t *testing.T) {
	t.Setenv("CLAUDE_CONFIG_DIR", t.TempDir()) // empty dir, no .claude.json
	input := &Input{}
	input.PopulatePlan()
	if input.Plan.OrgType != "" {
		t.Errorf("expected zero Plan when file missing, got %+v", input.Plan)
	}
}
