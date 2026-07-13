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
	uistate.KeyMap
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
		KeyMap: uistate.NewKeyMap(),
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
	switch msg := msg.(type) {
	case uistate.MsgStateUpdated:
		m.handleStateUpdated(msg)
	case uistate.MsgSwitchPage:
		m.handleSwitchPage(msg)
	case tea.KeyMsg:
		model, cmd, done := m.handleKey(msg)
		if done {
			return model, cmd
		}
	case tea.WindowSizeMsg:
		m.handleWindowSize(msg)
	}

	return m, nil
}

func (m *ModelTitle) handleStateUpdated(msg uistate.MsgStateUpdated) {
	m.globalState.Coach = msg.Coach
	m.globalState.Team = msg.Team
	m.menu.SetHasCoach(m.globalState.Coach != nil)
}

func (m *ModelTitle) handleSwitchPage(msg uistate.MsgSwitchPage) {
	if msg.NewPage == uistate.PageTitle {
		m.menu.SetHasCoach(m.globalState.Coach != nil)
	}
}

func (m *ModelTitle) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit, true
	case key.Matches(msg, m.keys.Up):
		m.menu.MoveUp()
	case key.Matches(msg, m.keys.Down):
		m.menu.MoveDown()
	case key.Matches(msg, m.keys.Enter):
		return m.handleMenuSelect()
	}
	return nil, nil, false
}

func (m *ModelTitle) handleWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
	m.footer.Update(msg)
}

func (m *ModelTitle) handleMenuSelect() (tea.Model, tea.Cmd, bool) {
	selected := m.menu.SelectedItem()
	switch selected {
	case components.TitleItemCreateCoach:
		return m, func() tea.Msg {
			return uistate.MsgSwitchPage{NewPage: uistate.PageCreateCoach}
		}, true
	case components.TitleItemLockerRoom:
		return m, func() tea.Msg {
			return uistate.MsgSwitchPage{NewPage: uistate.PageLockerRoom}
		}, true
	case components.TitleItemNewTournament:
		return m, func() tea.Msg {
			return uistate.MsgSwitchPage{NewPage: uistate.PageNewTournament}
		}, true
	case components.TitleItemNewBattleSelection:
		return m, func() tea.Msg {
			return uistate.MsgSwitchPage{NewPage: uistate.PageNewBattle}
		}, true
	case components.TitleItemRules:
		return m, func() tea.Msg {
			return uistate.MsgSwitchPage{NewPage: uistate.PageRules}
		}, true
	}

	return nil, nil, false
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
