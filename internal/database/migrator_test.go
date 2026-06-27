package database

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestMigrator_Migrate(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	migrator := NewMigrator(db, SchemaFS)

	if err := migrator.Migrate(t.Context()); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Verify that tables were created
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatalf("failed to query tables: %v", err)
	}
	defer rows.Close()

	tables := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("failed to scan table name: %v", err)
		}
		tables[name] = true
	}

	expectedTables := []string{"coaches", "teams", "players"}
	for _, table := range expectedTables {
		if !tables[table] {
			t.Errorf("expected table %s not found", table)
		}
	}
}
