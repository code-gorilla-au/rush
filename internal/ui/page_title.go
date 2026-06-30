package ui

import (
	"strings"

	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type titleKeyMap struct {
	components.CommonKeys
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
}

func (k titleKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

func (k titleKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Quit},
	}
}

func newTitleKeyMap() titleKeyMap {
	return titleKeyMap{
		CommonKeys: components.NewCommonKeys(),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

type ModelTitle struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
	keys        titleKeyMap
	footer      components.Footer
	menu        components.TitleMenu
}

func NewModelTitle(globalState *uistate.GlobalState, theme styles.IceTheme) *ModelTitle {
	keys := newTitleKeyMap()
	return &ModelTitle{
		globalState: globalState,
		keys:        keys,
		footer:      components.NewFooter(keys),
		menu:        components.NewTitleMenu(globalState.Coach != nil),
		theme:       theme,
	}
}

func (m *ModelTitle) Init() tea.Cmd {
	return nil
}

func (m *ModelTitle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch vMsg := msg.(type) {
	case uistate.MsgStateUpdated:
		m.globalState.Coach = vMsg.Coach
		m.globalState.Team = vMsg.Team
		m.menu.SetHasCoach(m.globalState.Coach != nil)
	case tea.KeyMsg:
		switch {
		case key.Matches(vMsg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(vMsg, m.keys.Up):
			m.menu.MoveUp()
		case key.Matches(vMsg, m.keys.Down):
			m.menu.MoveDown()
		case key.Matches(vMsg, m.keys.Enter):
			selected := m.menu.SelectedItem()
			switch selected {
			case components.TitleItemCreateCoach:
				return m, func() tea.Msg {
					return uistate.MsgSwitchPage{NewPage: uistate.PageCreateCoach}
				}
			case components.TitleItemLockerRoom:
				return m, func() tea.Msg {
					return uistate.MsgSwitchPage{NewPage: uistate.PageLockerRoom}
				}
			case components.TitleItemNewTournament:
				return m, func() tea.Msg {
					return uistate.MsgSwitchPage{NewPage: uistate.PageNewTournament}
				}
			case components.TitleItemNewBattleSelection:
				return m, func() tea.Msg {
					return uistate.MsgSwitchPage{NewPage: uistate.PageNewBattleSelection}
				}
			case components.TitleItemSettings:
				return m, func() tea.Msg {
					return uistate.MsgSwitchPage{NewPage: uistate.PageTitleSettings}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = vMsg.Width
		m.height = vMsg.Height
		m.footer.Update(vMsg)
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

func (m *ModelTitle) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	styledLogo := m.theme.Logo.Render(strings.Trim(logo, "\n"))

	navigation := m.menu.View(m.theme)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styledLogo,
		"",
		navigation,
		"",
		m.footer.View(m.theme),
	)

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
