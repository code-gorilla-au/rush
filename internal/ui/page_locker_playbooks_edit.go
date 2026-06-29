package ui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type lockerPlaybooksEditKeyMap struct {
	components.CommonKeys
	Back   key.Binding
	Enter  key.Binding
	Select key.Binding
}

func (k lockerPlaybooksEditKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit}
}

func (k lockerPlaybooksEditKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Select, k.Back, k.Quit},
	}
}

func newLockerPlaybooksEditKeyMap() lockerPlaybooksEditKeyMap {
	return lockerPlaybooksEditKeyMap{
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

type ModelLockerPlaybooksEdit struct {
	width                 int
	height                int
	theme                 styles.IceTheme
	globalState           *GlobalState
	playbookSvc           *playbooks.Service
	keys                  lockerPlaybooksEditKeyMap
	footer                components.Footer
	formationList         components.FormationList
	selectedFormationList components.FormationList
	activeList            int // 0 for formationList, 1 for selectedFormationList
	playbookID            int64
	playbookName          string
	playbookDescription   string
	newFormations         []playbooks.Formation
	err                   error
}

func NewModelLockerPlaybooksEdit(state *GlobalState, playbookSvc *playbooks.Service) *ModelLockerPlaybooksEdit {
	return &ModelLockerPlaybooksEdit{
		theme:       styles.NewIceTheme(),
		globalState: state,
		playbookSvc: playbookSvc,
		keys:        newLockerPlaybooksEditKeyMap(),
		footer:      components.NewFooter(newLockerPlaybooksEditKeyMap()),
		formationList: components.NewFormationList(components.FormationListConfig{
			Title:           "Available Formations",
			Items:           playbooks.Formations(),
			EnableFiltering: true,
			ShowDescription: true,
		}),
		selectedFormationList: components.NewFormationList(components.FormationListConfig{
			Title:           "Selected Formations (Max 10)",
			Items:           []playbooks.Formation{},
			EnableFiltering: false,
			ShowDescription: false,
		}),
	}
}

func (m *ModelLockerPlaybooksEdit) Init() tea.Cmd {
	return nil
}

func (m *ModelLockerPlaybooksEdit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case MsgSwitchPage:
		if msg.NewPage == PageLockerPlaybooksEdit {
			if msg.Playbook != nil {
				m.load(msg.Playbook)
			} else {
				m.reset()
			}
		}
	case error:
		m.err = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			if m.formationList.IsFiltering() {
				break
			}
			return m, func() tea.Msg {
				return MsgSwitchPage{
					NewPage: PageLockerPlaybooksCreate,
					Playbook: &playbooks.Playbook{
						ID:          m.playbookID,
						Name:        m.playbookName,
						Description: m.playbookDescription,
						Formations:  m.newFormations,
					},
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.formationList.SetSize(msg.Width/2-4, msg.Height-15)
		m.selectedFormationList.SetSize(msg.Width/2-4, msg.Height-15)
		m.footer.Update(msg)
	}

	cmds = append(cmds, m.updateAddFormations(msg))

	return m, tea.Batch(cmds...)
}

func (m *ModelLockerPlaybooksEdit) reset() {
	m.newFormations = nil
	m.playbookID = 0
	m.playbookName = ""
	m.playbookDescription = ""
	m.selectedFormationList.SetItems(nil)
	m.err = nil
}

func (m *ModelLockerPlaybooksEdit) load(p *playbooks.Playbook) {
	m.newFormations = p.Formations
	m.playbookID = p.ID
	m.playbookName = p.Name
	m.playbookDescription = p.Description
	m.selectedFormationList.SetItems(m.newFormations)
	m.err = nil
}

func (m *ModelLockerPlaybooksEdit) updateAddFormations(msg tea.Msg) tea.Cmd {
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

func (m *ModelLockerPlaybooksEdit) savePlaybook() tea.Msg {
	var err error
	if m.playbookID != 0 {
		_, err = m.playbookSvc.UpdatePlaybook(m.globalState.Context(), m.playbookID, playbooks.PlaybookParams{
			TeamID:      m.globalState.Team.ID,
			Name:        m.playbookName,
			Description: m.playbookDescription,
			Formations:  m.newFormations,
		})
	} else {
		_, err = m.playbookSvc.CreatePlaybook(m.globalState.Context(), playbooks.PlaybookParams{
			TeamID:      m.globalState.Team.ID,
			Name:        m.playbookName,
			Description: m.playbookDescription,
			Formations:  m.newFormations,
		})
	}
	if err != nil {
		return err
	}
	return MsgSwitchPage{NewPage: PageLockerPlaybooksList}
}

func (m *ModelLockerPlaybooksEdit) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	var content string
	title := "ALLOCATE FORMATIONS"

	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
	} else {
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
