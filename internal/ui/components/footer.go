package components

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
func (f Footer) View(style lipgloss.Style) string {
	return style.Render(f.Help.View(f.KeyMap))
}

// CommonKeys defines keys that are shared across many pages.
type CommonKeys struct {
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k CommonKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k CommonKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit},
	}
}

// NewCommonKeys returns a default set of common keys.
func NewCommonKeys() CommonKeys {
	return CommonKeys{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}
