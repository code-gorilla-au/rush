package main

import (
	"context"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/env"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui"
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

	queries := database.New(db)
	playbooksSvc := playbooks.NewPlaybooksService(queries)
	teamsSvc := teams.NewTeamsService(queries, playbooksSvc)
	gameSvc := games.NewService(queries)

	go func() {
		hasAICoaches, tErr := teamsSvc.HasAICoaches(ctx)
		if tErr != nil {
			slog.Error("Failed to check for AI coaches", "error", tErr)
			return
		}

		if hasAICoaches {
			return
		}

		if tErr = teamsSvc.GenerateAITeams(ctx); tErr != nil {
			slog.Error("Failed to generate teams", "error", tErr)
		}
	}()

	p := tea.NewProgram(ui.New(ui.Dependencies{
		TeamsSvc:    teamsSvc,
		PlaybookSvc: playbooksSvc,
		GameSvc:     gameSvc,
	}))
	if _, err = p.Run(); err != nil {
		slog.Error("Failed to run program", "error", err)
		os.Exit(1)
	}
}
