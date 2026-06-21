package ui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/ui/components"
)

type lockerPlaybooksListKeyMap struct {
	components.CommonKeys
	Back  key.Binding
	Enter key.Binding
	New   key.Binding
}

func (k lockerPlaybooksListKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.New, k.Back, k.Quit}
}

func (k lockerPlaybooksListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.New, k.Back, k.Quit},
	}
}

func newLockerPlaybooksListKeyMap() lockerPlaybooksListKeyMap {
	return lockerPlaybooksListKeyMap{
		CommonKeys: components.NewCommonKeys(),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		New: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new playbook"),
		),
	}
}

type ModelLockerPlaybooksList struct {
	width           int
	height          int
	theme           IceTheme
	globalState     *GlobalState
	playbookSvc     *playbooks.Service
	keys            lockerPlaybooksListKeyMap
	footer          components.Footer
	playbookList    components.PlaybookList
	playbooksLoaded bool
	err             error
}

func NewModelLockerPlaybooksList(state *GlobalState, playbookSvc *playbooks.Service) *ModelLockerPlaybooksList {
	return &ModelLockerPlaybooksList{
		theme:       NewIceTheme(),
		globalState: state,
		playbookSvc: playbookSvc,
		keys:        newLockerPlaybooksListKeyMap(),
		footer:      components.NewFooter(newLockerPlaybooksListKeyMap()),
	}
}

func (m *ModelLockerPlaybooksList) Init() tea.Cmd {
	return m.loadPlaybooks
}

func (m *ModelLockerPlaybooksList) loadPlaybooks() tea.Msg {
	if m.globalState.Team == nil {
		return nil
	}
	items, err := m.playbookSvc.GetTeamPlaybooks(m.globalState.Context(), m.globalState.Team.ID)
	if err != nil {
		return err
	}
	return MsgPlaybooksLoaded{Playbooks: items}
}

type MsgPlaybooksLoaded struct {
	Playbooks []playbooks.Playbook
}

func (m *ModelLockerPlaybooksList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case MsgPlaybooksLoaded:
		m.playbooksLoaded = true
		m.playbookList = components.NewPlaybookList(msg.Playbooks)
		m.playbookList.SetSize(m.width, m.height-10)
	case MsgSwitchPage:
		if msg.NewPage == PageLockerPlaybooksList {
			cmds = append(cmds, m.loadPlaybooks)
		}
	case error:
		m.err = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			if m.playbookList.IsFiltering() {
				break
			}
			return m, func() tea.Msg {
				return MsgSwitchPage{NewPage: PageLockerRoom}
			}
		case key.Matches(msg, m.keys.New):
			if !m.playbookList.IsFiltering() {
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageLockerPlaybooksEdit}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.playbooksLoaded {
			m.playbookList.SetSize(msg.Width, msg.Height-10)
		}
		m.footer.Update(msg)
	}

	if m.playbooksLoaded {
		var cmd tea.Cmd
		m.playbookList, cmd = m.playbookList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *ModelLockerPlaybooksList) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	var content string
	title := "PLAYBOOKS"

	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
	} else if !m.playbooksLoaded {
		content = "Loading..."
	} else {
		if m.playbookList.Len() == 0 {
			content = "No playbooks yet. Press 'n' to create one."
		} else {
			content = m.playbookList.View()
			if !m.playbookList.IsFiltering() {
				content += "\n\nPress 'n' to create new playbook"
			}
		}
	}

	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		m.theme.Logo.Render(title),
		"",
		content,
		"",
		m.footer.View(m.theme.Footer),
	)

	centeredContent := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		mainContent,
	)

	view.Content = m.theme.Base.
		Width(m.width).
		Height(m.height).
		Render(centeredContent)

	return view
}
