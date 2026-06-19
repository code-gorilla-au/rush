package ui

import (
	"strings"

	"github.com/code-gorilla-au/rush/internal/ui/components"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type titleKeyMap struct {
	components.CommonKeys
	CreateCoach key.Binding
	LockerRoom  key.Binding
}

func (k titleKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.CreateCoach, k.LockerRoom, k.Quit}
}

func (k titleKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CreateCoach, k.LockerRoom, k.Quit},
	}
}

func newTitleKeyMap() titleKeyMap {
	return titleKeyMap{
		CommonKeys: components.NewCommonKeys(),
		CreateCoach: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create coach"),
		),
		LockerRoom: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "locker room"),
		),
	}
}

type ModelTitle struct {
	width       int
	height      int
	theme       IceTheme
	globalState *GlobalState
	keys        titleKeyMap
	footer      components.Footer
}

func NewModelTitle(globalState *GlobalState) *ModelTitle {
	keys := newTitleKeyMap()
	return &ModelTitle{
		globalState: globalState,
		keys:        keys,
		footer:      components.NewFooter(keys),
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.CreateCoach):
			if m.globalState.Coach == nil {
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageCreateCoach}
				}
			}
		case key.Matches(msg, m.keys.LockerRoom):
			if m.globalState.Coach != nil {
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageLockerRoom}
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

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styledLogo,
		"",
		navigation,
		"",
		m.footer.View(m.theme.Footer),
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
