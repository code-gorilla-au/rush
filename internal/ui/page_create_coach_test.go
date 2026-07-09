package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
	"github.com/code-gorilla-au/rush/internal/ui/uitest"
)

func TestModelCreateCoach(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var state *uistate.GlobalState
	var theme styles.IceTheme
	var teamsSvc *teams.Service

	group.BeforeEach(func() {
		state = &uistate.GlobalState{}
		theme = styles.NewIceTheme()

		db := uitest.SetupTestDB(t)
		t.Cleanup(func() { _ = db.Close() })

		queries := database.New(db)
		playbookSvc := playbooks.NewPlaybooksService(queries)
		teamsSvc = teams.NewTeamsService(queries, playbookSvc)
	})

	err := group.
		Test("enter should move focus from coach to team input", func(t *testing.T) {
			m := NewModelCreateCoach(state, teamsSvc, theme)

			_, _ = m.Update(tea.KeyPressMsg{Text: "enter"})

			odize.AssertEqual(t, 1, m.focusIndex)
			odize.AssertFalse(t, m.coachInput.Focused())
			odize.AssertTrue(t, m.teamInput.Focused())
		}).
		Test("submit error from enter should be stored on model", func(t *testing.T) {
			db := uitest.SetupTestDB(t)
			queries := database.New(db)
			playbookSvc := playbooks.NewPlaybooksService(queries)
			brokenTeamsSvc := teams.NewTeamsService(queries, playbookSvc)
			odize.AssertNoError(t, db.Close())

			m := NewModelCreateCoach(state, brokenTeamsSvc, theme)
			m.coachInput.SetValue("Coach")
			m.teamInput.SetValue("Team")
			m.focusIndex = 1
			m.updateFocus()

			_, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})
			odize.AssertTrue(t, cmd != nil)

			msg := cmd()
			_, _ = m.Update(msg)

			odize.AssertTrue(t, m.err != nil)
		}).
		Run()

	odize.AssertNoError(t, err)
}
