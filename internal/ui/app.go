package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	width  int
	height int
	theme  IceTheme
}

// New returns a new UI model.
func New() model {
	return model{
		theme: NewIceTheme(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

const logo = `
  ____  _   _ ____  _   _ 
 |  _ \| | | / ___|| | | |
 | |_) | | | \___ \| |_| |
 |  _ <| |_| |___) |  _  |
 |_| \_\\___/|____/|_| |_|
`

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	styledLogo := m.theme.Logo.Render(strings.Trim(logo, "\n"))
	footer := "Press 'q' to quit"
	styledFooter := m.theme.Footer.Render(footer)

	content := styledLogo + "\n\n" + styledFooter

	// Center the content in the terminal
	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)

	// Apply base theme (black background) to the whole view
	finalView := m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	view := tea.NewView(finalView)
	view.AltScreen = true
	return view
}
