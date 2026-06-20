package playbooks

import (
	"database/sql"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	odize.AssertNoError(t, err)

	migrator := database.NewMigrator(db, database.SchemaFS)
	err = migrator.Migrate(t.Context())
	odize.AssertNoError(t, err)

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
		s = NewPlaybooksService(queries)
	})

	group.AfterEach(func() {
		if db != nil {
			db.Close()
		}
	})

	err := group.
		Test("CreatePlaybook should create a playbook and return it", func(t *testing.T) {
			ctx := t.Context()
			params := PlaybookParams{
				TeamID:      1,
				Name:        "Test Playbook",
				Description: "Test Description",
				Formations: []Formation{
					{
						Name:  "Formation 1",
						Lane1: 1,
						Lane2: 2,
						Lane3: 3,
					},
				},
			}

			pb, err := s.CreatePlaybook(ctx, params)
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, pb.ID > 0)
			odize.AssertEqual(t, params.Name, pb.Name)
			odize.AssertEqual(t, params.Description, pb.Description)
			odize.AssertEqual(t, 1, len(pb.Formations))
			odize.AssertEqual(t, "Formation 1", pb.Formations[0].Name)
		}).
		Test("GetTeamPlaybooks should return playbooks for a team", func(t *testing.T) {
			ctx := t.Context()

			// Create a team first (using queries directly to avoid dependency on teams package)
			res, err := db.ExecContext(ctx, "INSERT INTO teams (name, is_default) VALUES (?, ?)", "Team 1", true)
			odize.AssertNoError(t, err)
			teamID, err := res.LastInsertId()
			odize.AssertNoError(t, err)

			// Create a playbook for this team
			params := PlaybookParams{
				TeamID:     teamID,
				Name:       "Team Playbook",
				Formations: []Formation{{Name: "F1"}},
			}

			_, err = s.CreatePlaybook(ctx, params)
			odize.AssertNoError(t, err)

			playbooks, err := s.GetTeamPlaybooks(ctx, teamID)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, 1, len(playbooks))
			odize.AssertEqual(t, "Team Playbook", playbooks[0].Name)
		}).
		Test("DeletePlaybook should delete a playbook", func(t *testing.T) {
			ctx := t.Context()
			params := PlaybookParams{
				TeamID:     1,
				Name:       "To Be Deleted",
				Formations: []Formation{{Name: "F1"}},
			}

			pb, err := s.CreatePlaybook(ctx, params)
			odize.AssertNoError(t, err)

			err = s.DeletePlaybook(ctx, pb.ID)
			odize.AssertNoError(t, err)

			// Verify it is deleted
			playbooks, err := s.GetTeamPlaybooks(ctx, 1)
			odize.AssertNoError(t, err)

			found := false
			for _, p := range playbooks {
				if p.ID == pb.ID {
					found = true
					break
				}
			}
			odize.AssertFalse(t, found)
		}).
		Test("UpdatePlaybookFormations should update formations", func(t *testing.T) {
			ctx := t.Context()
			params := PlaybookParams{
				TeamID:     1,
				Name:       "To Be Updated",
				Formations: []Formation{{Name: "F1"}},
			}

			pb, err := s.CreatePlaybook(ctx, params)
			odize.AssertNoError(t, err)

			newFormations := []Formation{
				{
					Name:  "New Formation",
					Lane1: 5,
				},
			}

			updatedPb, err := s.UpdatePlaybookFormations(ctx, pb.ID, newFormations)
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, pb.ID, updatedPb.ID)
			odize.AssertEqual(t, 1, len(updatedPb.Formations))
			odize.AssertEqual(t, "New Formation", updatedPb.Formations[0].Name)
			odize.AssertEqual(t, 5, updatedPb.Formations[0].Lane1)
		}).
		Run()

	odize.AssertNoError(t, err)
}
