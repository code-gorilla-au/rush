package uilocker

import (
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type lockerRoomKeyMap struct {
	uistate.KeyMap
}

func (k lockerRoomKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Back, k.Quit}
}

func (k lockerRoomKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Back, k.Quit},
	}
}

func newLockerRoomKeyMap() lockerRoomKeyMap {
	return lockerRoomKeyMap{
		KeyMap: uistate.NewKeyMap(),
	}
}

type ModelLockerRoom struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
	keys        lockerRoomKeyMap
	footer      components.Footer
	list        components.LockerRoomList
}

func NewModelLockerRoom(globalState *uistate.GlobalState, theme styles.IceTheme) *ModelLockerRoom {
	keys := newLockerRoomKeyMap()
	return &ModelLockerRoom{
		globalState: globalState,
		keys:        keys,
		footer:      components.NewFooter(keys),
		theme:       theme,
		list:        components.NewLockerRoomList(theme),
	}
}

func (m *ModelLockerRoom) Init() tea.Cmd {
	return nil
}

func (m *ModelLockerRoom) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case uistate.MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return uistate.MsgSwitchPage{NewPage: uistate.PageTitle}
			}
		case key.Matches(msg, m.keys.Select):
			switch m.list.SelectedItem() {
			case components.ItemPlayers:
				return m, func() tea.Msg {
					return MsgSwitchLockerPage{NewPage: SubPageLockerPlayers}
				}
			case components.ItemPlaybooks:
				return m, func() tea.Msg {
					return MsgSwitchLockerPage{NewPage: SubPageLockerPlaybooksList}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(40, 20)
		m.footer.Update(msg)
	}

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

func (m *ModelLockerRoom) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.theme.Logo.Render("LOCKER ROOM"),
		"",
		m.list.View(m.theme),
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
