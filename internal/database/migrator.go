package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"slices"
	"strings"
)

import "embed"

//go:embed schema/*.sql
var SchemaFS embed.FS

// Service defines the interface for database migrations.
type Service interface {
	Migrate(ctx context.Context) error
}

// Migrator handles database schema migrations.
type Migrator struct {
	db *sql.DB
	fs fs.FS
}

// NewMigrator creates a new Migrator instance that implements the Service interface.
func NewMigrator(db *sql.DB, fsys fs.FS) Service {
	return &Migrator{
		db: db,
		fs: fsys,
	}
}

// Migrate executes all SQL files found in the schema directory.
func (m *Migrator) Migrate(ctx context.Context) error {
	entries, err := fs.ReadDir(m.fs, "schema")
	if err != nil {
		return fmt.Errorf("reading schema directory: %w", err)
	}

	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}

	slices.Sort(sqlFiles)

	for _, fileName := range sqlFiles {
		content, err := fs.ReadFile(m.fs, "schema/"+fileName)
		if err != nil {
			return fmt.Errorf("reading schema file %s: %w", fileName, err)
		}

		if _, err := m.db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("executing schema file %s: %w", fileName, err)
		}
	}

	return nil
}
