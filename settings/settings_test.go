package settings

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRead_MissingFile(t *testing.T) {
	s, err := Read("/nonexistent/path/settings.json")
	if err != nil {
		t.Errorf("missing file should not error, got: %v", err)
	}
	if s == nil {
		t.Error("expected non-nil Settings for missing file")
	}
	if s.StatusLine != nil {
		t.Error("expected nil StatusLine for missing file")
	}
}

func TestRead_ValidWithoutStatusLine(t *testing.T) {
	path := writeTempFile(t, `{"someOtherKey": "value", "nested": {"a": 1}}`)

	s, err := Read(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.StatusLine != nil {
		t.Errorf("expected nil StatusLine, got %+v", s.StatusLine)
	}
}

func TestRead_ValidWithStatusLine(t *testing.T) {
	path := writeTempFile(t, `{
		"other": "stuff",
		"statusLine": {
			"type": "command",
			"command": "/usr/local/bin/my-statusline"
		}
	}`)

	s, err := Read(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.StatusLine == nil {
		t.Fatal("expected non-nil StatusLine")
	}
	if s.StatusLine.Type != "command" {
		t.Errorf("expected type 'command', got %q", s.StatusLine.Type)
	}
	if s.StatusLine.Command != "/usr/local/bin/my-statusline" {
		t.Errorf("expected command '/usr/local/bin/my-statusline', got %q", s.StatusLine.Command)
	}
}

func TestRead_MalformedJSON(t *testing.T) {
	path := writeTempFile(t, `{not valid json}`)

	_, err := Read(path)
	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestWriteStatusLine_MissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	err := WriteStatusLine(path, "/path/to/binary")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	s, err := Read(path)
	if err != nil {
		t.Fatalf("could not read back file: %v\ncontent: %s", err, data)
	}
	if s.StatusLine == nil {
		t.Fatal("expected statusLine in created file")
	}
	if s.StatusLine.Command != "/path/to/binary" {
		t.Errorf("expected command '/path/to/binary', got %q", s.StatusLine.Command)
	}
}

func TestWriteStatusLine_ExistingFileWithoutStatusLine(t *testing.T) {
	path := writeTempFile(t, `{"someKey": "someValue", "another": 42}`)

	err := WriteStatusLine(path, "/new/binary")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s, err := Read(path)
	if err != nil {
		t.Fatalf("could not read back file: %v", err)
	}
	if s.StatusLine == nil {
		t.Fatal("expected statusLine to be added")
	}
	if s.StatusLine.Command != "/new/binary" {
		t.Errorf("expected command '/new/binary', got %q", s.StatusLine.Command)
	}

	// Verify other keys are preserved
	raw, _ := os.ReadFile(path)
	content := string(raw)
	if !strings.Contains(content, "someKey") {
		t.Error("expected 'someKey' to be preserved in file")
	}
	if !strings.Contains(content, "someValue") {
		t.Error("expected 'someValue' to be preserved in file")
	}
}

func TestWriteStatusLine_ExistingFileWithStatusLine(t *testing.T) {
	path := writeTempFile(t, `{
		"other": "data",
		"statusLine": {"type": "command", "command": "/old/path"}
	}`)

	err := WriteStatusLine(path, "/new/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s, err := Read(path)
	if err != nil {
		t.Fatalf("could not read back file: %v", err)
	}
	if s.StatusLine.Command != "/new/path" {
		t.Errorf("expected command '/new/path', got %q", s.StatusLine.Command)
	}

	// Other data still present
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), "other") {
		t.Error("expected 'other' key to be preserved")
	}
}

func TestWriteStatusLine_MalformedFile(t *testing.T) {
	path := writeTempFile(t, `{not valid json}`)
	original, _ := os.ReadFile(path)

	err := WriteStatusLine(path, "/some/binary")
	if err == nil {
		t.Error("expected error for malformed file")
	}

	// File should be unchanged
	after, _ := os.ReadFile(path)
	if string(original) != string(after) {
		t.Error("malformed file should not be modified")
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "settings*.json")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}
