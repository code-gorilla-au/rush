package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	width  int
	height int
}

// New returns a new UI model.
func New() model {
	return model{}
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

	logoLines := strings.Split(strings.Trim(logo, "\n"), "\n")
	logoWidth := 0
	for _, line := range logoLines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}

	leftPadding := (m.width - logoWidth) / 2
	topPadding := (m.height - len(logoLines) - 2) / 2 // -2 for spacing and footer

	if leftPadding < 0 {
		leftPadding = 0
	}
	if topPadding < 0 {
		topPadding = 0
	}

	var b strings.Builder
	for range topPadding {
		b.WriteString("\n")
	}

	pad := strings.Repeat(" ", leftPadding)
	for _, line := range logoLines {
		b.WriteString(pad)
		b.WriteString(line)
		b.WriteString("\n")
	}

	footer := "Press 'q' to quit"
	footerPadding := (m.width - len(footer)) / 2
	if footerPadding < 0 {
		footerPadding = 0
	}
	b.WriteString("\n")
	b.WriteString(strings.Repeat(" ", footerPadding))
	b.WriteString(footer)

	// Fill the rest of the height to prevent jumping if terminal is small
	remainingLines := m.height - topPadding - len(logoLines) - 2
	for range remainingLines {
		b.WriteString("\n")
	}

	view := tea.NewView(b.String())
	view.AltScreen = true
	return view
}
