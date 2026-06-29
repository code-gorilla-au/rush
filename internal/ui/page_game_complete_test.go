package ui

import (
	"database/sql"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	migrator := database.NewMigrator(db, database.SchemaFS)
	if err := migrator.Migrate(t.Context()); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	return db
}

func TestPageGameCompleteModel(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var db *sql.DB
	var queries *database.Queries
	var teamsSvc *teams.Service
	var gameSvc *games.Service
	var state *GlobalState

	group.BeforeEach(func() {
		db = setupTestDB(t)
		queries = database.New(db)
		teamsSvc = teams.NewTeamsService(queries, playbooks.NewPlaybooksService(queries))
		gameSvc = games.NewService(queries)
		state = &GlobalState{}
	})

	group.AfterEach(func() {
		if db != nil {
			db.Close()
		}
	})

	err := group.
		Test("should load winner details successfully", func(t *testing.T) {
			ctx := t.Context()

			// Setup winner coach and team
			coach, err := teamsSvc.CreateCoach(ctx, teams.CreateCoachParams{
				Name:    "Winning Coach",
				IsHuman: true,
			})
			odize.AssertNoError(t, err)

			team, err := teamsSvc.CreateTeam(ctx, "Winning Team", coach.ID, false)
			odize.AssertNoError(t, err)

			// Setup loser coach and team
			coachB, err := teamsSvc.CreateCoach(ctx, teams.CreateCoachParams{
				Name:    "Loser Coach",
				IsHuman: false,
			})
			odize.AssertNoError(t, err)

			teamB, err := teamsSvc.CreateTeam(ctx, "Loser Team", coachB.ID, false)
			odize.AssertNoError(t, err)

			// Create a game
			formationsA := make([]playbooks.Formation, 10)
			for i := range formationsA {
				formationsA[i] = playbooks.Formation{Lane1: 1}
			}
			formationsB := make([]playbooks.Formation, 10)
			for i := range formationsB {
				formationsB[i] = playbooks.Formation{Lane1: 1}
			}

			game, err := gameSvc.NewGame(ctx, games.NewGameParams{
				TeamA: games.TeamConfig{
					TeamID:     team.ID,
					TeamName:   team.Name,
					Formations: formationsA,
				},
				TeamB: games.TeamConfig{
					TeamID:     teamB.ID,
					TeamName:   teamB.Name,
					Formations: formationsB,
				},
			})
			odize.AssertNoError(t, err)

			// Simulate game completion with Team A as winner
			rollIdx := 0
			rollFn := func() int {
				rolls := []int{10, 1}
				r := rolls[rollIdx%2]
				rollIdx++
				return r
			}

			// Finish all rounds
			for i := 0; i < 10; i++ {
				game.ResolveRound(rollFn)
			}

			_, err = gameSvc.CompleteGame(ctx, game)
			odize.AssertNoError(t, err)

			theme := styles.NewIceTheme()
			m := NewPageGameComplete(state, teamsSvc, gameSvc, theme)
			m.SetGameID(game.ID())

			// Execute Init
			cmd := m.Init()
			msg := cmd()

			winnerMsg, ok := msg.(MsgWinnerLoaded)
			odize.AssertTrue(t, ok)
			odize.AssertEqual(t, team.ID, winnerMsg.Team.ID)
			odize.AssertEqual(t, coach.ID, winnerMsg.Coach.ID)
			odize.AssertFalse(t, winnerMsg.IsDraw)
		}).
		Test("should handle enter key and switch to title page", func(t *testing.T) {
			theme := styles.NewIceTheme()
			m := NewPageGameComplete(state, teamsSvc, gameSvc, theme)
			_, cmd := m.Update(tea.KeyPressMsg{
				Text: "enter",
			})

			odize.AssertTrue(t, cmd != nil)
			msg := cmd()
			switchMsg, ok := msg.(MsgSwitchPage)
			odize.AssertTrue(t, ok)
			odize.AssertEqual(t, PageTitle, switchMsg.NewPage)
		}).
		Run()

	odize.AssertNoError(t, err)
}
