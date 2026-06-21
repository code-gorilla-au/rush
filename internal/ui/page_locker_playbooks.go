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
	width                 int
	height                int
	theme                 IceTheme
	globalState           *GlobalState
	playbookSvc           *playbooks.Service
	keys                  lockerPlaybooksKeyMap
	footer                components.Footer
	playbookList          components.PlaybookList
	formationList         components.FormationList
	selectedFormationList components.SelectedFormationList
	mode                  playbooksViewMode
	activeList            int // 0 for formationList, 1 for selectedFormationList
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
		theme:                 NewIceTheme(),
		globalState:           state,
		playbookSvc:           playbookSvc,
		keys:                  newLockerPlaybooksKeyMap(),
		footer:                components.NewFooter(newLockerPlaybooksKeyMap()),
		newPlaybookName:       ti,
		formationList:         components.NewFormationList(),
		selectedFormationList: components.NewSelectedFormationList(),
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
	case MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case MsgPlaybooksLoaded:
		m.playbooksLoaded = true
		m.playbookList = components.NewPlaybookList(msg.Playbooks)
		m.playbookList.SetSize(m.width, m.height-10) // Set initial size
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
				if m.playbookList.IsFiltering() {
					break
				}
				return m, func() tea.Msg {
					return MsgSwitchPage{NewPage: PageLockerRoom}
				}
			}
			if m.mode == modeAddFormations {
				if m.formationList.IsFiltering() {
					break
				}
			}
			m.mode = modeList
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.playbooksLoaded {
			m.playbookList.SetSize(msg.Width, msg.Height-10)
		}
		m.formationList.SetSize(msg.Width/2-4, msg.Height-15)
		m.selectedFormationList.SetSize(msg.Width/2-4, msg.Height-15)
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
		if m.playbookList.IsFiltering() {
			break
		}
		switch msg.String() {
		case "n": // New playbook
			m.mode = modeCreateName
			m.newPlaybookName.Reset()
			m.newPlaybookName.Focus()
			m.newFormations = nil
			return nil
		}
	}

	var cmd tea.Cmd
	m.playbookList, cmd = m.playbookList.Update(msg)
	return cmd
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.formationList.IsFiltering() {
			break
		}
		switch msg.String() {
		case "tab":
			m.activeList = (m.activeList + 1) % 2
			m.formationList.SetActive(m.activeList == 0)
			m.selectedFormationList.SetActive(m.activeList == 1)
			return nil
		case "enter":
			if m.activeList == 0 {
				if len(m.newFormations) < 10 {
					f := m.formationList.SelectedItem()
					if f.Name != "" {
						m.newFormations = append(m.newFormations, f)
						cmds = append(cmds, m.selectedFormationList.SetItems(m.newFormations))
					}
				}
			} else {
				// Remove from selected
				if len(m.newFormations) > 0 {
					idx := m.selectedFormationList.SelectedIndex()
					if idx >= 0 && idx < len(m.newFormations) {
						m.newFormations = append(m.newFormations[:idx], m.newFormations[idx+1:]...)
						cmds = append(cmds, m.selectedFormationList.SetItems(m.newFormations))
					}
				}
			}
			return tea.Batch(cmds...)
		case "s": // Save
			if len(m.newFormations) > 0 {
				return m.savePlaybook
			}
		}
	}

	m.formationList.SetActive(m.activeList == 0)
	m.selectedFormationList.SetActive(m.activeList == 1)

	var cmd tea.Cmd
	m.formationList, cmd = m.formationList.Update(msg)
	cmds = append(cmds, cmd)

	m.selectedFormationList, cmd = m.selectedFormationList.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
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
				content = m.playbookList.View()
				if !m.playbookList.IsFiltering() {
					content += "\n\nPress 'n' to create new playbook"
				}
			}
		case modeCreateName:
			title = "CREATE PLAYBOOK"
			content = "Enter playbook name:\n\n" + m.newPlaybookName.View()
		case modeAddFormations:
			title = "ADD FORMATIONS"
			// The SetSize will be handled by WindowSizeMsg, but we ensure it here too for the split
			m.formationList.SetSize(m.width/2-4, m.height-15)
			m.selectedFormationList.SetSize(m.width/2-4, m.height-15)

			content = lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.formationList.View(),
				lipgloss.NewStyle().Width(2).Render(""),
				m.selectedFormationList.View(),
			)
			content += "\n\n" + m.theme.Footer.Render(fmt.Sprintf("%d/10 formations • Tab: switch • Enter: add/remove • 's': save", len(m.newFormations)))
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
