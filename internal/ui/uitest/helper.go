package uitest

import (
	"database/sql"
	"testing"

	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	_ "modernc.org/sqlite"
)

func SetupTestDB(t *testing.T) *sql.DB {
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

func SetupServices(t *testing.T) (*teams.Service, *playbooks.Service, *games.Service) {
	db := SetupTestDB(t)
	queries := database.New(db)
	ps := playbooks.NewPlaybooksService(queries)
	ts := teams.NewTeamsService(queries, ps)
	gs := games.NewService(queries)
	return ts, ps, gs
}
