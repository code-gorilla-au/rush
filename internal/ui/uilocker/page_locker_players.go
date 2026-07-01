package uilocker

import (
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type lockerPlayersKeyMap struct {
	uistate.KeyMap
}

func (k lockerPlayersKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit}
}

func (k lockerPlayersKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Back, k.Quit},
		{k.Up, k.Down},
	}
}

func newLockerPlayersKeyMap() lockerPlayersKeyMap {
	return lockerPlayersKeyMap{
		KeyMap: uistate.NewKeyMap(),
	}
}

type ModelLockerPlayers struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
	teamsSvc    *teams.Service
	keys        lockerPlayersKeyMap
	footer      components.Footer
	playerList  components.PlayerList
}

func NewModelLockerPlayers(state *uistate.GlobalState, teamsSvc *teams.Service, theme styles.IceTheme) *ModelLockerPlayers {
	keys := newLockerPlayersKeyMap()
	return &ModelLockerPlayers{
		theme:       theme,
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case uistate.MsgStateUpdated:
		if m.globalState.Team != nil {
			m.playerList = components.NewPlayerList(m.globalState.Team.Players)
		}
	case MsgSwitchLockerPage:
		if msg.NewPage == SubPageLockerPlayers && m.globalState.Team != nil {
			m.playerList = components.NewPlayerList(m.globalState.Team.Players)
		}
	case components.MsgPlayerUpdated:
		cmds = append(cmds, m.handlePlayerUpdated(msg))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return MsgSwitchLockerPage{NewPage: SubPageLockerRoom}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)
	}

	cmd := m.playerList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ModelLockerPlayers) handlePlayerUpdated(msg components.MsgPlayerUpdated) tea.Cmd {
	return func() tea.Msg {
		err := m.teamsSvc.UpdatePlayer(m.globalState.Context(), msg.Player.ID, msg.Player.Name)
		if err != nil {
			return err // Should probably handle this better
		}

		// Update global state
		if m.globalState.Team != nil {
			for i, p := range m.globalState.Team.Players {
				if p.ID == msg.Player.ID {
					m.globalState.Team.Players[i] = msg.Player
					break
				}
			}
		}
		return nil
	}
}

func (m *ModelLockerPlayers) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	playersView := "No players found"
	if len(m.playerList.Items) > 0 {
		playersView = m.playerList.View(m.theme)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.theme.Logo.Render(m.globalState.Team.Name),
		"",
		playersView,
		"",
		m.footer.View(m.theme),
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
