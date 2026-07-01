package iugame

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
	"github.com/code-gorilla-au/rush/internal/ui/uitest"
)

func TestGameModel(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var queries *database.Queries
	var teamsSvc *teams.Service
	var gameSvc *games.Service
	var state *uistate.GlobalState

	group.BeforeEach(func() {
		db := uitest.SetupTestDB(t)
		t.Cleanup(func() { db.Close() })
		queries = database.New(db)
		teamsSvc = teams.NewTeamsService(queries, playbooks.NewPlaybooksService(queries))
		gameSvc = games.NewService(queries)
		state = &uistate.GlobalState{}
	})

	err := group.
		Test("MsgSwitchGamePage should trigger Init of the new page", func(t *testing.T) {
			theme := styles.NewIceTheme()
			m := NewGameModel(state, teamsSvc, gameSvc, theme)

			// Initial page should be SubPageGameRoot
			odize.AssertEqual(t, SubPageGameRoot, m.currentPage)

			// Send MsgSwitchGamePage to switch to SubPageGameComplete
			msg := MsgSwitchGamePage{
				NewPage: SubPageGameComplete,
				GameID:  123,
			}

			_, cmd := m.Update(msg)

			odize.AssertEqual(t, SubPageGameComplete, m.currentPage)
			odize.AssertTrue(t, cmd != nil)

			// In bubbletea v2, we can't easily inspect the contents of a Batch without executing it
			// or using some internal knowledge. But we can check if it's nil.
			// More importantly, we can check if it returns the expected message when executed.
			// PageGameCompleteModel.Init() returns a func() tea.Msg.

			// If we execute the command, it should try to load game 123.
			// Since game 123 doesn't exist, it should return MsgGameError.

			var resMsg tea.Msg
			if cmd != nil {
				resMsg = cmd()
			}

			// We expect a MsgGameError because game 123 doesn't exist
			// But currently, it probably returns nil because Init() was never called
			odize.AssertTrue(t, resMsg != nil)
			_, ok := resMsg.(MsgGameError)
			odize.AssertTrue(t, ok)
		}).
		Test("uistate.MsgSwitchPage to PageGame should reset to SubPageGameRoot", func(t *testing.T) {
			theme := styles.NewIceTheme()
			m := NewGameModel(state, teamsSvc, gameSvc, theme)

			// Manually set to SubPageGameComplete
			m.currentPage = SubPageGameComplete

			// Send MsgSwitchPage to PageGame
			msg := uistate.MsgSwitchPage{
				NewPage: uistate.PageGame,
				GameID:  123,
			}

			m.Update(msg)

			odize.AssertEqual(t, SubPageGameRoot, m.currentPage)
		}).
		Run()

	odize.AssertNoError(t, err)
}
