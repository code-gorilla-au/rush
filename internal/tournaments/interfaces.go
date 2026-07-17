package tournaments

import (
	"context"

	"github.com/code-gorilla-au/rush/internal/database"
)

type Store interface {
	AllocateGameToStage(ctx context.Context, arg database.AllocateGameToStageParams) (database.StageGame, error)
	CreateStage(ctx context.Context, arg database.CreateStageParams) (database.Stage, error)
	CreateTournament(ctx context.Context, arg database.CreateTournamentParams) (database.Tournament, error)
	UpdateStage(ctx context.Context, arg database.UpdateStageParams) (database.Stage, error)
}
