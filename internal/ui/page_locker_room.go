package ui

import (
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type lockerRoomKeyMap struct {
	components.CommonKeys
	Back   key.Binding
	Select key.Binding
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
		CommonKeys: components.NewCommonKeys(),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back to title"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

type ModelLockerRoom struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *GlobalState
	keys        lockerRoomKeyMap
	footer      components.Footer
	list        components.LockerRoomList
}

func NewModelLockerRoom(globalState *GlobalState, theme styles.IceTheme) *ModelLockerRoom {
	keys := newLockerRoomKeyMap()
	return &ModelLockerRoom{
		globalState: globalState,
		keys:        keys,
		footer:      components.NewFooter(keys),
		theme:       theme,
		list:        components.NewLockerRoomList(),
	}
}

func (m *ModelLockerRoom) Init() tea.Cmd {
	return nil
}

func (m *ModelLockerRoom) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.list.Update(msg)

	switch msg := msg.(type) {
	case MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return MsgSwitchPage{NewPage: PageTitle}
			}
		case key.Matches(msg, m.keys.Select):
			switch m.list.SelectedItem() {
			case components.ItemPlayers:
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageLockerPlayers}
				}
			case components.ItemPlaybooks:
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageLockerPlaybooksList}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)
	}

	return m, nil
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
