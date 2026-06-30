package iugame

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type SubPageGame int

const (
	SubPageGameRoot SubPageGame = iota
	SubPageGameComplete
)

type MsgSwitchGamePage struct {
	NewPage SubPageGame
	GameID  int64
}

// GameModel handles all game related pages.
type GameModel struct {
	currentPage         SubPageGame
	subPageGame         tea.Model
	subPageGameRoot     tea.Model
	subPageGameComplete tea.Model
}

func NewGameModel(state *uistate.GlobalState, teamsSvc *teams.Service, gameSvc *games.Service, theme styles.IceTheme) *GameModel {
	return &GameModel{
		subPageGameRoot:     NewModelGame(state, gameSvc, theme),
		subPageGameComplete: NewPageGameComplete(state, teamsSvc, gameSvc, theme),
		currentPage:         SubPageGameRoot,
	}
}

func (m *GameModel) SetGameID(id int64) {
	if p, ok := m.subPageGameRoot.(*PageGameModel); ok {
		p.SetGameID(id)
	}
	if p, ok := m.subPageGameComplete.(*PageGameCompleteModel); ok {
		p.SetGameID(id)
	}
}

func (m *GameModel) SetPage(page SubPageGame) {
	m.currentPage = page
}

func (m *GameModel) View() tea.View {
	switch m.currentPage {
	case SubPageGameRoot:
		return m.subPageGameRoot.View()
	case SubPageGameComplete:
		return m.subPageGameComplete.View()
	default:
		return tea.NewView("unknown game page")
	}
}

func (m *GameModel) Init() tea.Cmd {
	switch m.currentPage {
	case SubPageGameRoot:
		return m.subPageGameRoot.Init()
	case SubPageGameComplete:
		return m.subPageGameComplete.Init()
	default:
		return nil
	}
}

func (m *GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case MsgSwitchGamePage:
		switch msg.NewPage {
		case SubPageGameRoot, SubPageGameComplete:
			m.currentPage = msg.NewPage
		}
		if msg.GameID != 0 {
			m.SetGameID(msg.GameID)
		}

	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.subPageGameRoot, cmd = m.subPageGameRoot.Update(msg)
		cmds = append(cmds, cmd)
		m.subPageGameComplete, cmd = m.subPageGameComplete.Update(msg)
		cmds = append(cmds, cmd)

	}

	var cmd tea.Cmd
	switch m.currentPage {
	case SubPageGameRoot:
		m.subPageGameRoot, cmd = m.subPageGameRoot.Update(msg)
	case SubPageGameComplete:
		m.subPageGameComplete, cmd = m.subPageGameComplete.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
