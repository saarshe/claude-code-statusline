package components

import (
	"testing"
)

func TestRegistry_GetKnownKey(t *testing.T) {
	c := Get("model")
	if c == nil {
		t.Error("Get('model') returned nil")
	}
}

func TestRegistry_GetUnknownKey(t *testing.T) {
	c := Get("nonexistent")
	if c != nil {
		t.Errorf("Get('nonexistent') = %v, want nil", c)
	}
}

func TestRegistry_KeyMatchesComponentKey(t *testing.T) {
	for _, key := range registeredKeys() {
		c := Get(key)
		if c == nil {
			t.Errorf("Get(%q) returned nil", key)
			continue
		}
		if string(c.Key()) != key {
			t.Errorf("component registered as %q but Key() returns %q", key, c.Key())
		}
	}
}
