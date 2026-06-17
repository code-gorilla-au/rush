package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/code-gorilla-au/env"
	"github.com/code-gorilla-au/rush/internal/database"
)

func main() {
	ctx := context.Background()
	env.LoadEnvFile(".env.local")
	config := NewConfig()

	db, err := database.NewSqLiteProvider(config.DatabaseUrl)
	if err != nil {
		slog.Error("Failed to create database provider", "error", err)
		os.Exit(1)

	}

	defer func() {
		dbErr := db.Close()
		if dbErr != nil {
			slog.Error("Failed to close database", "error", err)
		}
	}()

	migrator := database.NewMigrator(db, database.SchemaFS)
	if mErr := migrator.Migrate(ctx); mErr != nil {
		slog.Error("Failed to migrate database", "error", mErr)
		os.Exit(1)
	}

	slog.Info("Hello, world!")
}
