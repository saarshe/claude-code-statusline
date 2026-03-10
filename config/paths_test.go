package config

import (
	"strings"
	"testing"
)

func TestConfigPath(t *testing.T) {
	path := ConfigPath()

	if path == "" {
		t.Error("ConfigPath() returned empty string")
	}
	if !strings.HasSuffix(path, "config.toml") {
		t.Errorf("ConfigPath() = %q, want suffix 'config.toml'", path)
	}
	if !strings.Contains(path, ".claude-code-statusline") {
		t.Errorf("ConfigPath() = %q, want '.claude-code-statusline' in path", path)
	}
}
