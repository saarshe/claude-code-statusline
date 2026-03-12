package components

import (
	"strings"
	"testing"

	"github.com/saars/claude-code-statusline/config"
	"github.com/saars/claude-code-statusline/schema"
	"github.com/saars/claude-code-statusline/theme"
)

func TestAgent_ShowsNameWhenActive(t *testing.T) {
	data := &schema.Input{Agent: &schema.Agent{Name: "my-agent"}}
	result := Get("agent").Render(data, config.Default(), theme.Get("default"))
	if !strings.Contains(result, "my-agent") {
		t.Errorf("expected agent name in output, got %q", result)
	}
}

func TestAgent_EmptyWhenNil(t *testing.T) {
	data := &schema.Input{}
	result := Get("agent").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when no agent, got %q", result)
	}
}

func TestAgent_EmptyWhenNameEmpty(t *testing.T) {
	data := &schema.Input{Agent: &schema.Agent{Name: ""}}
	result := Get("agent").Render(data, config.Default(), theme.Get("default"))
	if result != "" {
		t.Errorf("expected empty output when agent name is empty, got %q", result)
	}
}
