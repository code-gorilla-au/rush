package ui

import (
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type lockerPlayersKeyMap struct {
	components.CommonKeys
	Back key.Binding
}

func (k lockerPlayersKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

func (k lockerPlayersKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
	}
}

func newLockerPlayersKeyMap() lockerPlayersKeyMap {
	return lockerPlayersKeyMap{
		CommonKeys: components.NewCommonKeys(),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to locker room"),
		),
	}
}

type ModelLockerPlayers struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
	teamsSvc    *teams.Service
	keys        lockerPlayersKeyMap
	footer      components.Footer
}

func NewModelLockerPlayers(state *GlobalState, teamsSvc *teams.Service) *ModelLockerPlayers {
	keys := newLockerPlayersKeyMap()
	return &ModelLockerPlayers{
		theme:       NewIceTheme(),
		globalState: state,
		teamsSvc:    teamsSvc,
		keys:        keys,
		footer:      components.NewFooter(keys),
	}
}

func (m *ModelLockerPlayers) Init() tea.Cmd {
	return nil
}

func (m *ModelLockerPlayers) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return MsgSwitchPage{NewPage: PageLockerRoom}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)
	}

	return m, nil
}

func (m *ModelLockerPlayers) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.theme.Logo.Render("PLAYERS"),
		"",
		"Locker Room Players",
		"",
		m.footer.View(m.theme.Footer),
	)

	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)

	view.Content = m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	return view
}
