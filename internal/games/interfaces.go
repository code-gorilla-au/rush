package games

import (
	"context"

	"github.com/code-gorilla-au/rush/internal/database"
)

type Store interface {
	CreateGame(ctx context.Context, arg database.CreateGameParams) (database.Game, error)
	GetGameByID(ctx context.Context, id int64) (database.Game, error)
	ListCompletedGamesByTeam(ctx context.Context, arg database.ListCompletedGamesByTeamParams) ([]database.Game, error)
	UpdateGame(ctx context.Context, arg database.UpdateGameParams) (database.Game, error)
	StartGame(ctx context.Context, id int64) error
}

type RollStrategy interface {
	Run(input DecisionInput) DuelResult
}
