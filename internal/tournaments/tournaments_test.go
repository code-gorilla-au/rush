package tournaments

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
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

func newTestTournamentService(db *sql.DB) *Service {
	queries := database.New(db)
	teamsSvc := teams.NewTeamsService(queries, playbooks.NewPlaybooksService(queries))
	gamesSvc := games.NewService(queries)
	return NewService(ServiceDependencies{
		GamesSvc: gamesSvc,
		TeamsSvc: teamsSvc,
		Store:    queries,
		DB:       db,
		TxnFunc:  func(tx *sql.Tx) Store { return database.New(tx) },
	})
}

func createTestCoach(ctx context.Context, t *testing.T, db *sql.DB, name string, isHuman bool) database.Coach {
	queries := database.New(db)
	coach, err := queries.CreateCoach(ctx, database.CreateCoachParams{
		Name:    name,
		IsHuman: sql.NullBool{Bool: isHuman, Valid: true},
	})
	odize.AssertNoError(t, err)
	return coach
}

func createTestTeam(ctx context.Context, t *testing.T, db *sql.DB, name string, coachID int64) database.Team {
	queries := database.New(db)
	team, err := queries.CreateTeam(ctx, database.CreateTeamParams{
		Name:    name,
		CoachID: sql.NullInt64{Int64: coachID, Valid: true},
	})
	odize.AssertNoError(t, err)
	return team
}

func createTestPlaybook(ctx context.Context, t *testing.T, db *sql.DB, name string, teamID int64) playbooks.Playbook {
	queries := database.New(db)
	playbooksSvc := playbooks.NewPlaybooksService(queries)
	pb, err := playbooksSvc.CreatePlaybook(ctx, playbooks.PlaybookParams{
		Name:       name,
		TeamID:     teamID,
		Formations: playbooks.Formations(),
	})
	odize.AssertNoError(t, err)
	return pb
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

		tournamentGames := generateGameParamsFromTeams(totalTeams)

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
			tournamentGames := generateGameParamsFromTeams([]teams.AITeam{})
			odize.AssertEqual(t, 0, len(tournamentGames))
		}).
		Test("handles 1 team", func(t *testing.T) {
			totalTeams := []teams.AITeam{
				{Team: teams.Team{ID: 1, Name: "Team1"}, Coach: teams.Coach{}},
			}
			tournamentGames := generateGameParamsFromTeams(totalTeams)
			odize.AssertEqual(t, 0, len(tournamentGames))
		}).
		Test("handles 2 teams", func(t *testing.T) {
			totalTeams := []teams.AITeam{
				{Team: teams.Team{ID: 1, Name: "Team1"}, Coach: teams.Coach{}, Playbook: playbooks.Playbook{Formations: []playbooks.Formation{{Name: "F1"}}}},
				{Team: teams.Team{ID: 2, Name: "Team2"}, Coach: teams.Coach{}, Playbook: playbooks.Playbook{Formations: []playbooks.Formation{{Name: "F2"}}}},
			}
			tournamentGames := generateGameParamsFromTeams(totalTeams)
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
		s = newTestTournamentService(db)
	})

	group.AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	err := group.
		Test("returns error when not enough AI teams", func(t *testing.T) {
			ctx := context.Background()

			// Create a coach
			coach := createTestCoach(ctx, t, db, "Coach", true)

			// Create a team for the coach
			createTestTeam(ctx, t, db, "Team", coach.ID)

			// Not enough AI teams (trying to request 2 teams, but only 1 AI team exists)
			err := s.CreateTournament(ctx, CreateTournamentParams{
				Name:          "Tournament",
				NumberOfTeams: NumberOfTeams(2),
				CoachID:       coach.ID,
			})
			odize.AssertTrue(t, err != nil)
		}).
		Test("returns error when coach does not exist", func(t *testing.T) {
			err := s.CreateTournament(context.Background(), CreateTournamentParams{
				Name:          "Tournament",
				NumberOfTeams: NumberOfTeams(2),
				CoachID:       999,
			})
			odize.AssertTrue(t, err != nil)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestService_CreateTournament_Success(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var db *sql.DB
	var s *Service

	group.BeforeEach(func() {
		db = setupTestDB(t)
		s = newTestTournamentService(db)
	})

	group.AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	err := group.
		Test("successfully creates a tournament", func(t *testing.T) {
			ctx := context.Background()

			// Create a coach
			coach := createTestCoach(ctx, t, db, "Coach", true)

			// Create a team for the coach
			team := createTestTeam(ctx, t, db, "HumanTeam", coach.ID)

			// Create a playbook for the team
			createTestPlaybook(ctx, t, db, "PB", team.ID)

			// Create AI teams
			for i := 0; i < 2; i++ {
				aiCoach := createTestCoach(ctx, t, db, fmt.Sprintf("AICoach%d", i), false)
				aiTeam := createTestTeam(ctx, t, db, fmt.Sprintf("Team%d", i), aiCoach.ID)

				createTestPlaybook(ctx, t, db, "PB", aiTeam.ID)
			}

			// Create Tournament
			err := s.CreateTournament(ctx, CreateTournamentParams{
				Name:          "TestTournament",
				NumberOfTeams: NumberOfTeams(3),
				CoachID:       coach.ID,
			})
			odize.AssertNoError(t, err)
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
		s = newTestTournamentService(db)
	})

	group.AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	err := group.Test("returns error when no playbooks found", func(t *testing.T) {
		ctx := context.Background()
		// Create a coach
		coach := createTestCoach(ctx, t, db, "Coach", true)

		// Create a team for the coach
		createTestTeam(ctx, t, db, "Team", coach.ID)

		// getHumanTeam should return error
		_, err := s.getHumanTeam(ctx, coach.ID)
		odize.AssertTrue(t, err != nil)
		odize.AssertEqual(t, "no playbooks found for team", err.Error())
	}).Run()

	odize.AssertNoError(t, err)
}

func TestService_UpdateStageStatus(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var db *sql.DB
	var s *Service

	group.BeforeEach(func() {
		db = setupTestDB(t)
		s = newTestTournamentService(db)
	})

	group.AfterEach(func() {
		if db != nil {
			_ = db.Close()
		}
	})

	err := group.
		Test("successfully updates stage status", func(t *testing.T) {
			ctx := context.Background()
			queries := database.New(db)

			// Create a tournament
			tournament, err := queries.CreateTournament(ctx, database.CreateTournamentParams{
				Name:          "TestTournament",
				NumberOfTeams: 4,
			})
			odize.AssertNoError(t, err)

			// Create a stage
			stage, err := queries.CreateStage(ctx, database.CreateStageParams{
				Name: "TestStage",
				TournamentID: sql.NullInt64{
					Int64: tournament.ID,
					Valid: true,
				},
				Status: string(StageStatusPending),
			})
			odize.AssertNoError(t, err)

			// Update stage status to active
			err = s.UpdateStageStatus(ctx, tournament.ID, stage.ID, StageStatusActive)
			odize.AssertNoError(t, err)

			// Verify status
			updatedStage, err := queries.GetStageByID(ctx, stage.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, string(StageStatusActive), updatedStage.Status)

			// Update stage status to complete
			err = s.UpdateStageStatus(ctx, tournament.ID, stage.ID, StageStatusComplete)
			odize.AssertNoError(t, err)

			// Verify status
			updatedStage, err = queries.GetStageByID(ctx, stage.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, string(StageStatusComplete), updatedStage.Status)
		}).
		Run()

	odize.AssertNoError(t, err)
}
