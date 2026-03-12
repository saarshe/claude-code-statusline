package main

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	out, code := run([]string{"--version"}, strings.NewReader(""))
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out, "claude-code-statusline") {
		t.Errorf("expected version output to contain binary name, got %q", out)
	}
}

func TestHelp(t *testing.T) {
	out, code := run([]string{"--help"}, strings.NewReader(""))
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(strings.ToLower(out), "usage") {
		t.Errorf("expected help output to contain 'usage', got %q", out)
	}
}

func TestGarbageStdin(t *testing.T) {
	out, code := run([]string{}, strings.NewReader("not valid json!!!"))
	if code != 0 {
		t.Fatalf("expected exit 0 for garbage stdin, got %d", code)
	}
	if out != "" {
		t.Errorf("expected empty output for garbage stdin, got %q", out)
	}
}

func TestEmptyStdin(t *testing.T) {
	out, code := run([]string{}, strings.NewReader(""))
	if code != 0 {
		t.Fatalf("expected exit 0 for empty stdin, got %d", code)
	}
	if out != "" {
		t.Errorf("expected empty output for empty stdin, got %q", out)
	}
}
