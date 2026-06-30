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

type lockerPlaybooksCreateKeyMap struct {
	components.CommonKeys
	Back  key.Binding
	Enter key.Binding
}

func (k lockerPlaybooksCreateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Back, k.Quit}
}

func (k lockerPlaybooksCreateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Back, k.Quit},
	}
}

func newLockerPlaybooksCreateKeyMap() lockerPlaybooksCreateKeyMap {
	return lockerPlaybooksCreateKeyMap{
		CommonKeys: components.NewCommonKeys(),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
	}
}

type ModelLockerPlaybooksCreate struct {
	width        int
	height       int
	theme        styles.IceTheme
	playbookSvc  *playbooks.Service
	keys         lockerPlaybooksCreateKeyMap
	footer       components.Footer
	playbookID   int64
	playbookForm components.PlaybookForm
	formations   []playbooks.Formation
	err          error
}

func NewModelLockerPlaybooksCreate(state *uistate.GlobalState, playbookSvc *playbooks.Service, theme styles.IceTheme) *ModelLockerPlaybooksCreate {
	return &ModelLockerPlaybooksCreate{
		theme:        theme,
		playbookSvc:  playbookSvc,
		keys:         newLockerPlaybooksCreateKeyMap(),
		footer:       components.NewFooter(newLockerPlaybooksCreateKeyMap()),
		playbookForm: components.NewPlaybookForm(),
	}
}

func (m *ModelLockerPlaybooksCreate) Init() tea.Cmd {
	return nil
}

func (m *ModelLockerPlaybooksCreate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgSwitchLockerPage:
		if msg.NewPage == SubPageLockerPlaybooksCreate {
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
			return m, func() tea.Msg {
				return MsgSwitchLockerPage{NewPage: SubPageLockerPlaybooksList}
			}
		case key.Matches(msg, m.keys.Enter):
			name, description := m.playbookForm.Values()
			if name != "" {
				return m, func() tea.Msg {
					return MsgSwitchLockerPage{
						NewPage: SubPageLockerPlaybooksEdit,
						Playbook: &playbooks.Playbook{
							ID:          m.playbookID,
							Name:        name,
							Description: description,
							Formations:  m.formations,
						},
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.Update(msg)
	}

	var cmd tea.Cmd
	m.playbookForm, cmd = m.playbookForm.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ModelLockerPlaybooksCreate) reset() {
	m.playbookForm.Reset()
	m.formations = nil
	m.playbookID = 0
	m.err = nil
}

func (m *ModelLockerPlaybooksCreate) load(p *playbooks.Playbook) {
	m.playbookForm.SetValues(p.Name, p.Description)
	m.formations = p.Formations
	m.playbookID = p.ID
	m.err = nil
}

func (m *ModelLockerPlaybooksCreate) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true

	var content string
	title := "PLAYBOOK NAME"

	if m.err != nil {
		content = m.theme.Logo.Render(fmt.Sprintf("Error: %v", m.err))
	} else {
		if m.playbookID != 0 {
			title = "EDIT PLAYBOOK"
		} else {
			title = "CREATE PLAYBOOK"
		}
		content = "Enter playbook details:\n\n" + m.playbookForm.View()
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
