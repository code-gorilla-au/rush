package components

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

// Footer is a reusable component for displaying help information.
type Footer struct {
	Help   help.Model
	KeyMap help.KeyMap
	Width  int
}

// NewFooter creates a new Footer component.
func NewFooter(keyMap help.KeyMap) Footer {
	return Footer{
		Help:   help.New(),
		KeyMap: keyMap,
	}
}

// Update updates the footer component.
func (f *Footer) Update(msg tea.Msg) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		f.Width = msg.Width
	}
}

// View renders the footer component.
func (f Footer) View(theme styles.IceTheme) string {
	return theme.Footer.Render(f.Help.View(f.KeyMap))
}
