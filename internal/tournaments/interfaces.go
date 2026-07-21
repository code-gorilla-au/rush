package tournaments

import (
	"context"

	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type Store interface {
	AllocateGameToStage(ctx context.Context, arg database.AllocateGameToStageParams) (database.StageGame, error)
	CreateStage(ctx context.Context, arg database.CreateStageParams) (database.Stage, error)
	CreateTournament(ctx context.Context, arg database.CreateTournamentParams) (database.Tournament, error)
	UpdateStage(ctx context.Context, arg database.UpdateStageParams) (database.Stage, error)
	SetStageStatus(ctx context.Context, arg database.SetStageStatusParams) error
}

type TeamsService interface {
	ListAITeams(ctx context.Context) ([]teams.AITeam, error)
	GetCoachByID(ctx context.Context, coachID int64) (teams.Coach, error)
	GetTeamAndPlaybooksByCoachID(ctx context.Context, coachID int64) (teams.TeamWithPlaybooks, error)
}
