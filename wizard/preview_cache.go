package wizard

import tea "github.com/charmbracelet/bubbletea"

// previewCache caches a rendered preview string tied to a WizardState.
// Models embed this to avoid recomputing the preview on every View() call.
// The cache is refreshed by calling Refresh with an optional state mutator.
type previewCache struct {
	state  *WizardState
	cached string
}

// newPreviewCache creates a cache and computes the initial preview.
func newPreviewCache(state *WizardState) previewCache {
	return previewCache{
		state:  state,
		cached: previewBlock(state),
	}
}

// Refresh invalidates the layout and recomputes the preview. If mutate is
// non-nil it is applied to a shallow copy of the state before rendering,
// leaving the original state untouched (useful for "hovered" previews).
func (p *previewCache) Refresh(mutate func(*WizardState)) {
	if mutate != nil {
		hovered := *p.state
		mutate(&hovered)
		hovered.InvalidateLayout()
		p.cached = previewBlock(&hovered)
	} else {
		p.state.InvalidateLayout()
		p.cached = previewBlock(p.state)
	}
}

// HandleResize should be called from Update when a tea.WindowSizeMsg is
// received. It invalidates the layout cache and recomputes the preview.
// Returns the tea.Cmd to propagate (nil).
func (p *previewCache) HandleResize(tea.WindowSizeMsg) tea.Cmd {
	p.state.InvalidateLayout()
	p.cached = previewBlock(p.state)
	return nil
}

// String returns the cached preview string.
func (p *previewCache) String() string {
	return p.cached
}
