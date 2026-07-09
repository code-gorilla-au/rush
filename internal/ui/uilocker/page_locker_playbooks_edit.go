package uilocker

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type lockerPlaybooksEditKeyMap struct {
	uistate.KeyMap
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
		KeyMap: uistate.NewKeyMap(),
	}
}

const (
	availableFormationListIndex = iota
	selectedFormationListIndex
	maxFormationsPerPlaybook = 10
)

type ModelLockerPlaybooksEdit struct {
	width                 int
	height                int
	theme                 styles.IceTheme
	globalState           *uistate.GlobalState
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

func NewModelLockerPlaybooksEdit(state *uistate.GlobalState, playbookSvc *playbooks.Service, theme styles.IceTheme) *ModelLockerPlaybooksEdit {
	keys := newLockerPlaybooksEditKeyMap()
	model := &ModelLockerPlaybooksEdit{
		theme:       theme,
		globalState: state,
		playbookSvc: playbookSvc,
		keys:        keys,
		footer:      components.NewFooter(keys),
		formationList: components.NewFormationList(components.FormationListConfig{
			Title:           "Available Formations",
			Items:           playbooks.Formations(),
			EnableFiltering: true,
			ShowDescription: true,
		}, theme),
		selectedFormationList: components.NewFormationList(components.FormationListConfig{
			Title:           "Selected Formations (Max 10)",
			Items:           []playbooks.Formation{},
			EnableFiltering: false,
			ShowDescription: false,
		}, theme),
	}
	model.setActiveList(availableFormationListIndex)
	return model
}

func (m *ModelLockerPlaybooksEdit) Init() tea.Cmd {
	return m.reset()
}

func (m *ModelLockerPlaybooksEdit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.footer.Update(msg)
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case uistate.MsgStateUpdated:
		m.globalState.Coach = msg.Coach
		m.globalState.Team = msg.Team
	case MsgSwitchLockerPage:
		if msg.NewPage == SubPageLockerPlaybooksEdit {
			if msg.Playbook != nil {
				cmds = append(cmds, m.load(msg.Playbook))
			} else {
				cmds = append(cmds, m.reset())
			}
		}
	case error:
		m.err = msg
	case tea.KeyMsg:
		model, cmd, handled := m.handleKey(msg)
		if handled {
			cmds = append(cmds, cmd)
			return model, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		listWidth := (m.width - 10) / 2
		listHeight := m.height - 20
		m.formationList.SetSize(listWidth, listHeight)
		m.selectedFormationList.SetSize(listWidth, listHeight)
	}

	var cmd tea.Cmd
	m.formationList, cmd = m.formationList.Update(msg)
	cmds = append(cmds, cmd)

	m.selectedFormationList, cmd = m.selectedFormationList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ModelLockerPlaybooksEdit) reset() tea.Cmd {
	m.newFormations = nil
	m.playbookID = 0
	m.playbookName = ""
	m.playbookDescription = ""
	m.formationList.Reset()
	m.selectedFormationList.Reset()
	m.setActiveList(availableFormationListIndex)
	m.err = nil
	return m.selectedFormationList.SetItems(nil)
}

func (m *ModelLockerPlaybooksEdit) load(p *playbooks.Playbook) tea.Cmd {
	m.newFormations = p.Formations
	m.playbookID = p.ID
	m.playbookName = p.Name
	m.playbookDescription = p.Description
	m.formationList.Reset()
	m.selectedFormationList.Reset()
	m.setActiveList(availableFormationListIndex)
	m.err = nil
	return m.selectedFormationList.SetItems(m.newFormations)
}

func (m *ModelLockerPlaybooksEdit) handleKey(msg tea.KeyMsg) (*ModelLockerPlaybooksEdit, tea.Cmd, bool) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit, true
	case key.Matches(msg, m.keys.Back):
		return m.handleBack()
	case key.Matches(msg, m.keys.Tab):
		return m.handleTab()
	case key.Matches(msg, m.keys.Enter):
		if m.formationList.IsFiltering() {
			return m, nil, false
		}
		return m, m.toggleFormationSelection(), true
	case key.Matches(msg, m.keys.Save):
		if len(m.newFormations) > 0 {
			return m, m.savePlaybook, true
		}
	}

	return m, nil, false
}

func (m *ModelLockerPlaybooksEdit) handleBack() (*ModelLockerPlaybooksEdit, tea.Cmd, bool) {
	if m.formationList.IsFiltering() {
		return m, nil, false
	}
	return m, func() tea.Msg {
		return MsgSwitchLockerPage{
			NewPage: SubPageLockerPlaybooksCreate,
			Playbook: &playbooks.Playbook{
				ID:          m.playbookID,
				Name:        m.playbookName,
				Description: m.playbookDescription,
				Formations:  m.newFormations,
			},
		}
	}, true
}

func (m *ModelLockerPlaybooksEdit) handleTab() (*ModelLockerPlaybooksEdit, tea.Cmd, bool) {
	if m.formationList.IsFiltering() {
		return m, nil, false
	}
	if m.activeList == availableFormationListIndex {
		m.setActiveList(selectedFormationListIndex)
	} else {
		m.setActiveList(availableFormationListIndex)
	}
	return m, nil, true
}

func (m *ModelLockerPlaybooksEdit) setActiveList(activeList int) {
	m.activeList = activeList
	m.formationList.SetActive(activeList == availableFormationListIndex)
	m.selectedFormationList.SetActive(activeList == selectedFormationListIndex)
}

func (m *ModelLockerPlaybooksEdit) toggleFormationSelection() tea.Cmd {
	if m.activeList == availableFormationListIndex {
		if len(m.newFormations) >= maxFormationsPerPlaybook {
			return nil
		}

		formation := m.formationList.SelectedItem()
		if formation.Name == "" {
			return nil
		}

		m.newFormations = append(m.newFormations, formation)
		return m.selectedFormationList.SetItems(m.newFormations)
	}

	if len(m.newFormations) == 0 {
		return nil
	}

	idx := m.selectedFormationList.SelectedIndex()
	if idx < 0 || idx >= len(m.newFormations) {
		return nil
	}

	m.newFormations = append(m.newFormations[:idx], m.newFormations[idx+1:]...)
	return m.selectedFormationList.SetItems(m.newFormations)
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
	return MsgSwitchLockerPage{NewPage: SubPageLockerPlaybooksList}
}

func (m *ModelLockerPlaybooksEdit) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing...")
	}

	view := tea.NewView("")
	view.AltScreen = true

	var content string
	title := "ALLOCATE FORMATIONS"

	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
	} else {
		availableView := m.formationList.View(m.theme)
		selectedView := m.selectedFormationList.View(m.theme)

		if m.activeList == availableFormationListIndex {
			availableView = m.theme.ActiveBorder.Render(availableView)
			selectedView = m.theme.InactiveBorder.Render(selectedView)
		} else {
			availableView = m.theme.InactiveBorder.Render(availableView)
			selectedView = m.theme.ActiveBorder.Render(selectedView)
		}

		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Left,
				m.theme.SecondaryHeader.Render("AVAILABLE FORMATIONS"),
				availableView,
			),
			lipgloss.NewStyle().Width(2).Render(""),
			lipgloss.JoinVertical(lipgloss.Left,
				m.theme.SecondaryHeader.Render("SELECTED FORMATIONS (MAX 10)"),
				selectedView,
			),
		)
		content += "\n\n" + m.theme.Muted.Render(fmt.Sprintf("%d/%d formations • Tab: switch • Enter: add/remove • 's': save", len(m.newFormations), maxFormationsPerPlaybook))
	}

	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		m.theme.Logo.Render(title),
		"",
		content,
		"",
		m.footer.View(m.theme),
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
