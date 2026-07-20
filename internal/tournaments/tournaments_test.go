package tournaments

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

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	odize.AssertNoError(t, err)
	migrator := database.NewMigrator(db, database.SchemaFS)
	err = migrator.Migrate(context.Background())
	odize.AssertNoError(t, err)
	return db
}

func TestGenerateGameParamsFromTeams(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.Test("generates correct number of unique games and ensures no duplicate pairs for 4 teams", func(t *testing.T) {
		n := 4
		var totalTeams []teams.AITeam
		for i := range n {
			totalTeams = append(totalTeams, teams.AITeam{
				Team:  teams.Team{ID: int64(i + 1), Name: fmt.Sprintf("Team%d", i+1)},
				Coach: teams.Coach{},
			})
		}

		tournamentGames := generateGameParamsFromTeams(totalTeams, nil)

		expectedGames := (n * (n - 1)) / 2
		odize.AssertEqual(t, expectedGames, len(tournamentGames))

		seenPairs := make(map[string]bool)
		for _, game := range tournamentGames {
			pairKey := fmt.Sprintf("%d-%d", game.TeamA.TeamID, game.TeamB.TeamID)
			reverseKey := fmt.Sprintf("%d-%d", game.TeamB.TeamID, game.TeamA.TeamID)

			odize.AssertFalse(t, seenPairs[pairKey])
			odize.AssertFalse(t, seenPairs[reverseKey])

			seenPairs[pairKey] = true
		}
	}).
		Test("handles 0 teams", func(t *testing.T) {
			tournamentGames := generateGameParamsFromTeams([]teams.AITeam{}, nil)
			odize.AssertEqual(t, 0, len(tournamentGames))
		}).
		Test("handles 1 team", func(t *testing.T) {
			totalTeams := []teams.AITeam{
				{Team: teams.Team{ID: 1, Name: "Team1"}, Coach: teams.Coach{}},
			}
			tournamentGames := generateGameParamsFromTeams(totalTeams, nil)
			odize.AssertEqual(t, 0, len(tournamentGames))
		}).
		Test("handles 2 teams", func(t *testing.T) {
			totalTeams := []teams.AITeam{
				{Team: teams.Team{ID: 1, Name: "Team1"}, Coach: teams.Coach{}, Playbook: playbooks.Playbook{Formations: []playbooks.Formation{{Name: "F1"}}}},
				{Team: teams.Team{ID: 2, Name: "Team2"}, Coach: teams.Coach{}, Playbook: playbooks.Playbook{Formations: []playbooks.Formation{{Name: "F2"}}}},
			}
			tournamentGames := generateGameParamsFromTeams(totalTeams, nil)
			odize.AssertEqual(t, 1, len(tournamentGames))
			odize.AssertEqual(t, []playbooks.Formation{{Name: "F1"}}, tournamentGames[0].TeamA.Formations)
			odize.AssertEqual(t, []playbooks.Formation{{Name: "F2"}}, tournamentGames[0].TeamB.Formations)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestService_CreateTournament_EdgeCases(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var db *sql.DB
	var s *Service

	group.BeforeEach(func() {
		db = setupTestDB(t)
		queries := database.New(db)
		teamsSvc := teams.NewTeamsService(queries, playbooks.NewPlaybooksService(queries))
		s = &Service{
			teamsSvc: teamsSvc,
		}
	})

	group.AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	err := group.
		Test("returns error when not enough AI teams", func(t *testing.T) {
			ctx := context.Background()
			queries := database.New(db)

			// Create a coach
			coach, err := queries.CreateCoach(ctx, database.CreateCoachParams{Name: "Coach"})
			odize.AssertNoError(t, err)

			// Create a team for the coach
			_, err = queries.CreateTeam(ctx, database.CreateTeamParams{Name: "Team", CoachID: sql.NullInt64{Int64: coach.ID, Valid: true}})
			odize.AssertNoError(t, err)

			// Not enough AI teams (trying to request 2 teams, but only 1 AI team exists)
			err = s.CreateTournament(ctx, coach.ID, 2)
			odize.AssertTrue(t, err != nil)
		}).
		Test("returns error when coach does not exist", func(t *testing.T) {
			err := s.CreateTournament(context.Background(), 999, 2)
			odize.AssertTrue(t, err != nil)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestService_GetHumanTeam_NoPlaybooks(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var db *sql.DB
	var s *Service

	group.BeforeEach(func() {
		db = setupTestDB(t)
		queries := database.New(db)
		teamsSvc := teams.NewTeamsService(queries, playbooks.NewPlaybooksService(queries))
		s = &Service{
			teamsSvc: teamsSvc,
		}
	})

	group.AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	err := group.Test("returns error when no playbooks found", func(t *testing.T) {
		ctx := context.Background()
		queries := database.New(db)
		// Create a coach
		coach, err := queries.CreateCoach(ctx, database.CreateCoachParams{
			Name: "Coach",
		})
		odize.AssertNoError(t, err)

		// Create a team for the coach
		_, err = queries.CreateTeam(ctx, database.CreateTeamParams{
			Name:    "Team",
			CoachID: sql.NullInt64{Int64: coach.ID, Valid: true},
		})
		odize.AssertNoError(t, err)

		// getHumanTeam should return error
		_, err = s.getHumanTeam(ctx, coach.ID)
		odize.AssertTrue(t, err != nil)
		odize.AssertEqual(t, "no playbooks found for team", err.Error())
	}).Run()

	odize.AssertNoError(t, err)
}
