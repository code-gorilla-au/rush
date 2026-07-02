package games

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	_ "modernc.org/sqlite"
)

func TestService_GetTeamStatistics(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("should return zero values when team has no completed games", func(t *testing.T) {
		teamSvc, gameSvc := setupStatsTestServices(t)
		ctx := t.Context()

		coach, err := teamSvc.CreateCoach(ctx, teams.CreateCoachParams{Name: "Coach Zero", IsHuman: true})
		odize.AssertNoError(t, err)

		team, err := teamSvc.CreateTeam(ctx, "No Games Team", coach.ID, true)
		odize.AssertNoError(t, err)

		stats, err := gameSvc.GetTeamStatistics(ctx, team.ID)
		odize.AssertNoError(t, err)

		odize.AssertEqual(t, 0, stats.GamesPlayed)
		odize.AssertEqual(t, 0, stats.Wins)
		odize.AssertEqual(t, 0, stats.Draws)
		odize.AssertEqual(t, 0, stats.Losses)
		odize.AssertEqual(t, 0, stats.RoundsWon)
		odize.AssertEqual(t, 0, stats.RoundsLost)
		odize.AssertEqual(t, 0, stats.RoundDifferential)
	})

	group.Test("should calculate mixed win draw loss and round metrics", func(t *testing.T) {
		teamSvc, gameSvc := setupStatsTestServices(t)
		ctx := t.Context()

		targetCoach, err := teamSvc.CreateCoach(ctx, teams.CreateCoachParams{Name: "Target Coach", IsHuman: true})
		odize.AssertNoError(t, err)

		opponentCoach, err := teamSvc.CreateCoach(ctx, teams.CreateCoachParams{Name: "Opponent Coach", IsHuman: true})
		odize.AssertNoError(t, err)

		targetTeam, err := teamSvc.CreateTeam(ctx, "Target Team", targetCoach.ID, true)
		odize.AssertNoError(t, err)

		opponentTeam, err := teamSvc.CreateTeam(ctx, "Opponent Team", opponentCoach.ID, false)
		odize.AssertNoError(t, err)

		targetCfg := TeamConfig{TeamID: targetTeam.ID, TeamName: targetTeam.Name, Formations: make([]playbooks.Formation, 10)}
		opponentCfg := TeamConfig{TeamID: opponentTeam.ID, TeamName: opponentTeam.Name, Formations: make([]playbooks.Formation, 10)}

		persistCompletedGame(t, gameSvc, ctx, targetCfg, opponentCfg, targetTeam.ID, []Result{
			{Outcome: ResultTeamA},
			{Outcome: ResultTeamA},
			{Outcome: ResultTeamB},
		})

		persistCompletedGame(t, gameSvc, ctx, opponentCfg, targetCfg, opponentTeam.ID, []Result{
			{Outcome: ResultTeamA},
			{Outcome: ResultTeamB},
			{Outcome: ResultTeamA},
		})

		persistCompletedGame(t, gameSvc, ctx, targetCfg, opponentCfg, 0, []Result{
			{Outcome: ResultTeamA},
			{Outcome: ResultTeamB},
		})

		stats, err := gameSvc.GetTeamStatistics(ctx, targetTeam.ID)
		odize.AssertNoError(t, err)

		odize.AssertEqual(t, 3, stats.GamesPlayed)
		odize.AssertEqual(t, 1, stats.Wins)
		odize.AssertEqual(t, 1, stats.Draws)
		odize.AssertEqual(t, 1, stats.Losses)
		odize.AssertEqual(t, 4, stats.RoundsWon)
		odize.AssertEqual(t, 4, stats.RoundsLost)
		odize.AssertEqual(t, 0, stats.RoundDifferential)
		odize.AssertEqual(t, "33.3", fmt.Sprintf("%.1f", stats.WinRate))
		odize.AssertEqual(t, "1.33", fmt.Sprintf("%.2f", stats.AverageRoundsWon))
		odize.AssertEqual(t, "1.33", fmt.Sprintf("%.2f", stats.AverageRoundsLost))
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}

func setupStatsTestServices(t *testing.T) (*teams.Service, *Service) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	migrator := database.NewMigrator(db, database.SchemaFS)
	if err := migrator.Migrate(t.Context()); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	queries := database.New(db)
	playbookSvc := playbooks.NewPlaybooksService(queries)
	teamSvc := teams.NewTeamsService(queries, playbookSvc)
	gameSvc := NewService(queries)

	return teamSvc, gameSvc
}

func persistCompletedGame(t *testing.T, gameSvc *Service, ctx context.Context, teamA TeamConfig, teamB TeamConfig, winner int64, results []Result) {
	game, err := gameSvc.NewGame(ctx, NewGameParams{TeamA: teamA, TeamB: teamB})
	odize.AssertNoError(t, err)

	game.status = StatusComplete
	game.currentRound = int64(len(game.rounds))
	game.results = results
	game.winner = new(winner)

	_, err = gameSvc.UpdateGame(ctx, game)
	odize.AssertNoError(t, err)
}
