package teams

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/playbooks"
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

func TestService(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var db *sql.DB
	var queries *database.Queries
	var s *Service

	group.BeforeEach(func() {
		db = setupTestDB(t)
		queries = database.New(db)
		s = NewTeamsService(queries, playbooks.NewPlaybooksService(queries))
	})

	group.AfterEach(func() {
		if db != nil {
			db.Close()
		}
	})

	err := group.
		Test("CreateCoach should create a coach and return it", func(t *testing.T) {
			ctx := t.Context()
			name := "Coach Carter"

			coach, err := s.CreateCoach(ctx, CreateCoachParams{
				Name:    name,
				IsHuman: false,
			})
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, coach.ID > 0)
			odize.AssertEqual(t, name, coach.Name)
			odize.AssertFalse(t, coach.CreatedAt.IsZero())
			odize.AssertFalse(t, coach.UpdatedAt.IsZero())
		}).
		Test("SetDefaultCoach should set the default coach", func(t *testing.T) {
			ctx := t.Context()
			name := "Coach Carter"

			coach, err := s.CreateCoach(ctx, CreateCoachParams{
				Name:    name,
				IsHuman: false,
			})
			odize.AssertNoError(t, err)

			err = s.SetDefaultCoach(ctx, coach.ID)
			odize.AssertNoError(t, err)

			// Verify it's set
			queries := database.New(db)
			model, err := queries.GetDefaultCoach(ctx)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, coach.ID, model.ID)
			odize.AssertTrue(t, model.IsDefault.Bool)
		}).
		Test("ClearDefaultCoach should clear the default coach", func(t *testing.T) {
			ctx := t.Context()
			name := "Coach Carter"

			_, err := s.CreateCoach(ctx, CreateCoachParams{
				Name:      name,
				IsDefault: true,
			})
			odize.AssertNoError(t, err)

			err = s.ClearDefaultCoach(ctx)
			odize.AssertNoError(t, err)

			// Verify it's cleared
			queries := database.New(db)
			_, err = queries.GetDefaultCoach(ctx)
			odize.AssertError(t, err)
			odize.AssertTrue(t, err == sql.ErrNoRows)
		}).
		Test("SetDefaultTeam should set the default team", func(t *testing.T) {
			ctx := t.Context()
			queries := database.New(db)
			_, err := queries.CreateTeam(ctx, database.CreateTeamParams{
				Name: "The Bulls",
			})
			odize.AssertNoError(t, err)

			var teamID int64
			err = db.QueryRowContext(ctx, "SELECT id FROM teams WHERE name = ?", "The Bulls").Scan(&teamID)
			odize.AssertNoError(t, err)

			err = s.SetDefaultTeam(ctx, teamID)
			odize.AssertNoError(t, err)

			// Verify it's set
			var isDefault bool
			err = db.QueryRowContext(ctx, "SELECT is_default FROM teams WHERE id = ?", teamID).Scan(&isDefault)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, isDefault)
		}).
		Test("ClearDefaultTeam should clear the default team", func(t *testing.T) {
			_, err := queries.CreateTeam(t.Context(), database.CreateTeamParams{
				Name:      "The Bulls",
				IsDefault: sql.NullBool{Bool: true, Valid: true},
			})
			odize.AssertNoError(t, err)

			var teamID int64
			err = db.QueryRowContext(t.Context(), "SELECT id FROM teams WHERE name = ?", "The Bulls").Scan(&teamID)
			odize.AssertNoError(t, err)

			err = s.ClearDefaultTeam(t.Context())
			odize.AssertNoError(t, err)

			// Verify it's cleared
			var isDefault bool
			err = db.QueryRowContext(t.Context(), "SELECT is_default FROM teams WHERE id = ?", teamID).Scan(&isDefault)
			odize.AssertNoError(t, err)
			odize.AssertFalse(t, isDefault)
		}).
		Test("GetDefaultCoach should return the default coach", func(t *testing.T) {
			ctx := t.Context()
			name := "Coach Carter"
			_, err := s.CreateCoach(ctx, CreateCoachParams{
				Name:      name,
				IsDefault: true,
			})
			odize.AssertNoError(t, err)

			coach, err := s.GetDefaultCoach(ctx)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, name, coach.Name)
		}).
		Test("GetDefaultCoach should return error if no default coach", func(t *testing.T) {
			_, err := s.GetDefaultCoach(t.Context())
			odize.AssertError(t, err)
			odize.AssertTrue(t, errors.Is(err, ErrCoachNotFound))
		}).
		Test("UpdatePlayer should update player name", func(t *testing.T) {
			ctx := t.Context()
			coach, err := s.CreateCoach(ctx, CreateCoachParams{
				Name:      "Coach",
				IsDefault: true,
			})
			odize.AssertNoError(t, err)

			team, err := s.CreateTeam(ctx, "Team", coach.ID, true)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(team.Players) > 0)

			player := team.Players[0]
			newName := "Updated Name"

			err = s.UpdatePlayer(ctx, player.ID, newName)
			odize.AssertNoError(t, err)

			// Verify
			members, err := queries.GetTeamMembers(ctx, sql.NullInt64{Int64: team.ID, Valid: true})
			odize.AssertNoError(t, err)

			found := false
			for _, m := range members {
				if m.ID == player.ID {
					odize.AssertEqual(t, newName, m.Name)
					found = true
					break
				}
			}
			odize.AssertTrue(t, found)
		}).
		Test("GetTeamByCoachID should return team and players", func(t *testing.T) {
			ctx := t.Context()
			coach, err := s.CreateCoach(ctx, CreateCoachParams{
				Name:      "Coach",
				IsDefault: true,
			})
			odize.AssertNoError(t, err)

			_, err = s.CreateTeam(ctx, "The Bulls", coach.ID, true)
			odize.AssertNoError(t, err)

			team, err := s.GetTeamByCoachID(ctx, coach.ID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "The Bulls", team.Name)
			odize.AssertEqual(t, 5, len(team.Players))
		}).
		Test("GetTeamByCoachID should return error if team not found", func(t *testing.T) {
			_, err := s.GetTeamByCoachID(t.Context(), 999)
			odize.AssertError(t, err)
			odize.AssertTrue(t, errors.Is(err, ErrTeamNotFound))
		}).
		Test("CreateTeam should create a team with default players", func(t *testing.T) {
			ctx := t.Context()
			coach, err := s.CreateCoach(ctx, CreateCoachParams{
				Name:    "Coach",
				IsHuman: false,
			})
			odize.AssertNoError(t, err)

			team, err := s.CreateTeam(ctx, "Lakers", coach.ID, true)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, "Lakers", team.Name)
			odize.AssertEqual(t, 5, len(team.Players))
			odize.AssertEqual(t, "Player 1", team.Players[0].Name)
		}).
		Run()

	odize.AssertNoError(t, err)
}
