package ui

import (
	_ "embed"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
)

//go:embed rules.md
var rulesContent string

type rulesKeyMap struct {
	uistate.KeyMap
}

func newRulesKeyMap() rulesKeyMap {
	return rulesKeyMap{
		KeyMap: uistate.NewKeyMap(),
	}
}

type ModelRules struct {
	width       int
	height      int
	theme       styles.IceTheme
	globalState *uistate.GlobalState
	viewport    viewport.Model
	keys        rulesKeyMap
	ready       bool
}

func NewModelRules(globalState *uistate.GlobalState, theme styles.IceTheme) *ModelRules {
	return &ModelRules{
		globalState: globalState,
		theme:       theme,
		keys:        newRulesKeyMap(),
	}
}

func (m *ModelRules) Init() tea.Cmd {
	return nil
}

func (m *ModelRules) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Back) {
			return m, func() tea.Msg {
				return uistate.MsgSwitchPage{NewPage: uistate.PageTitle}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if !m.ready {
			m.viewport = viewport.New(viewport.WithWidth(m.width), viewport.WithHeight(m.height))
			m.ready = true
		} else {
			m.viewport = viewport.New(viewport.WithWidth(m.width), viewport.WithHeight(m.height))
		}

	}

	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *ModelRules) headerView() string {
	title := m.theme.Title.Render("Rules")
	return title
}

func (m *ModelRules) footerView() string {
	info := m.theme.Muted.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	return info
}

func (m *ModelRules) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}

	rendered, err := glamour.Render(rulesContent, "dark")
	if err != nil {
		return tea.NewView("Err: " + err.Error())
	}

	m.viewport.SetContent(rendered)
	view := tea.NewView(m.viewport.View())
	view.AltScreen = true
	return view
}
