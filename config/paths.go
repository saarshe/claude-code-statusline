package config

import (
	"os"
	"path/filepath"
)

func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".claude-code-statusline/config.toml"
	}
	return filepath.Join(home, ".claude-code-statusline", "config.toml")
}
