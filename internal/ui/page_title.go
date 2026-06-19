package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ModelTitle struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
}

func NewModelTitle(globalState *GlobalState) *ModelTitle {
	return &ModelTitle{
		globalState: globalState,
	}
}

func (m ModelTitle) Init() tea.Cmd {
	return nil
}

func (m ModelTitle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			if m.globalState.Coach == nil {
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageCreateCoach}
				}
			}
		case "l":
			if m.globalState.Coach != nil {
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageLockerRoom}
				}
			}
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

func (m ModelTitle) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	styledLogo := m.theme.Logo.Render(strings.Trim(logo, "\n"))

	var navigation string
	if m.globalState.Coach == nil {
		navigation = "Press " + m.theme.Hotkey.Render("c") + " to create a coach"
	} else {
		navigation = m.theme.Button.Render("Locker Room (l)")
	}

	footer := "Press 'q' to quit"
	styledFooter := m.theme.Footer.Render(footer)

	content := styledLogo + "\n\n" + navigation + "\n\n" + styledFooter

	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)

	finalView := m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	view := tea.NewView(finalView)
	view.AltScreen = true
	return view
}
