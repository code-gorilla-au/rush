package ui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/ui/components"
)

type playbooksViewMode int

const (
	modeList playbooksViewMode = iota
	modeCreateName
	modeAddFormations
)

type lockerPlaybooksKeyMap struct {
	components.CommonKeys
	Back   key.Binding
	Enter  key.Binding
	Select key.Binding
}

func (k lockerPlaybooksKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit}
}

func (k lockerPlaybooksKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Select, k.Back, k.Quit},
	}
}

func newLockerPlaybooksKeyMap() lockerPlaybooksKeyMap {
	return lockerPlaybooksKeyMap{
		CommonKeys: components.NewCommonKeys(),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

type ModelLockerPlaybooks struct {
	width         int
	height        int
	theme         IceTheme
	globalState   *GlobalState
	playbookSvc   *playbooks.Service
	keys          lockerPlaybooksKeyMap
	footer        components.Footer
	playbookList  components.PlaybookList
	formationList components.FormationList
	mode          playbooksViewMode
	// Create flow state
	newPlaybookName textinput.Model
	newFormations   []playbooks.Formation
	playbooksLoaded bool
	err             error
}

func NewModelLockerPlaybooks(state *GlobalState, playbookSvc *playbooks.Service) *ModelLockerPlaybooks {
	ti := textinput.New()
	ti.Placeholder = "Playbook Name"
	ti.Focus()

	return &ModelLockerPlaybooks{
		theme:           NewIceTheme(),
		globalState:     state,
		playbookSvc:     playbookSvc,
		keys:            newLockerPlaybooksKeyMap(),
		footer:          components.NewFooter(newLockerPlaybooksKeyMap()),
		newPlaybookName: ti,
		formationList:   components.NewFormationList(),
	}
}

func (m *ModelLockerPlaybooks) Init() tea.Cmd {
	return m.loadPlaybooks
}

func (m *ModelLockerPlaybooks) loadPlaybooks() tea.Msg {
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

func (m *ModelLockerPlaybooks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgPlaybooksLoaded:
		m.playbooksLoaded = true
		m.playbookList = components.NewPlaybookList(msg.Playbooks)
		m.mode = modeList
	case MsgSwitchPage:
		if msg.NewPage == PageLockerPlaybooks {
			cmds = append(cmds, m.loadPlaybooks)
		}
	case error:
		m.err = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			if m.mode == modeList {
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageLockerRoom}
				}
			}
			m.mode = modeList
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)
	}

	switch m.mode {
	case modeList:
		cmds = append(cmds, m.updateList(msg))
	case modeCreateName:
		cmds = append(cmds, m.updateCreateName(msg))
	case modeAddFormations:
		cmds = append(cmds, m.updateAddFormations(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m *ModelLockerPlaybooks) updateList(msg tea.Msg) tea.Cmd {
	if !m.playbooksLoaded {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n": // New playbook
			m.mode = modeCreateName
			m.newPlaybookName.Reset()
			m.newPlaybookName.Focus()
			m.newFormations = nil
			return nil
		}
	}

	m.playbookList.Update(msg)
	return nil
}

func (m *ModelLockerPlaybooks) updateCreateName(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" && m.newPlaybookName.Value() != "" {
			m.mode = modeAddFormations
			return nil
		}
	}
	var cmd tea.Cmd
	m.newPlaybookName, cmd = m.newPlaybookName.Update(msg)
	return cmd
}

func (m *ModelLockerPlaybooks) updateAddFormations(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if len(m.newFormations) < 10 {
				m.newFormations = append(m.newFormations, m.formationList.SelectedItem())
			}
		} else if msg.String() == "s" { // Save
			if len(m.newFormations) > 0 {
				return m.savePlaybook
			}
		}
	}
	m.formationList.Update(msg)
	return nil
}

func (m *ModelLockerPlaybooks) savePlaybook() tea.Msg {
	_, err := m.playbookSvc.CreatePlaybook(m.globalState.Context(), playbooks.PlaybookParams{
		TeamID:     m.globalState.Team.ID,
		Name:       m.newPlaybookName.Value(),
		Formations: m.newFormations,
	})
	if err != nil {
		return err
	}
	return m.loadPlaybooks()
}

func (m *ModelLockerPlaybooks) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	var content string
	title := "PLAYBOOKS"

	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
	} else if !m.playbooksLoaded {
		content = "Loading..."
	} else {
		switch m.mode {
		case modeList:
			if m.playbookList.Len() == 0 {
				content = "No playbooks yet. Press 'n' to create one."
			} else {
				content = m.playbookList.View(lipgloss.NewStyle(), m.theme.ListSelected)
				content += "\n\nPress 'n' to create new playbook"
			}
		case modeCreateName:
			title = "CREATE PLAYBOOK"
			content = "Enter playbook name:\n\n" + m.newPlaybookName.View()
		case modeAddFormations:
			title = "ADD FORMATIONS"
			content = fmt.Sprintf("Formations (%d/10):\n", len(m.newFormations))
			for _, f := range m.newFormations {
				content += m.theme.ListSelected.Render(" + "+f.Name) + "\n"
			}
			content += "\nAvailable Formations:\n"
			content += m.formationList.View(lipgloss.NewStyle(), m.theme.ListSelected)
			content += "\n\nPress 'enter' to add, 's' to save"
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
