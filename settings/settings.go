package settings

import (
	"encoding/json"
	"errors"
	"os"
)

// StatusLine represents the statusLine key in ~/.claude/settings.json.
type StatusLine struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

// Settings is a typed view of ~/.claude/settings.json.
// Unknown fields are preserved via rawFields so round-tripping doesn't lose data.
type Settings struct {
	StatusLine *StatusLine
	rawFields  map[string]json.RawMessage
}

// Read parses the settings file at path. Returns an empty Settings (no error)
// if the file does not exist. Returns an error for malformed JSON.
func Read(path string) (*Settings, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Settings{rawFields: map[string]json.RawMessage{}}, nil
	}
	if err != nil {
		return nil, err
	}

	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	s := &Settings{rawFields: raw}
	if slRaw, ok := raw["statusLine"]; ok {
		var sl StatusLine
		if err := json.Unmarshal(slRaw, &sl); err != nil {
			return nil, err
		}
		s.StatusLine = &sl
	}

	return s, nil
}

// WriteStatusLine surgically inserts or replaces the statusLine key in the
// settings file at path. If the file does not exist it is created with only
// the statusLine key. If the file contains malformed JSON an error is returned
// and the file is left unchanged.
func WriteStatusLine(path string, command string) error {
	s, err := Read(path)
	if err != nil {
		return err
	}

	s.rawFields["statusLine"], err = json.Marshal(StatusLine{
		Type:    "command",
		Command: command,
	})
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.rawFields, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
