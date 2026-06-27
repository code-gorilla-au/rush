package games

import (
	"context"

	"github.com/code-gorilla-au/rush/internal/database"
)

type Store interface {
	CreateGame(ctx context.Context, arg database.CreateGameParams) (database.Game, error)
	GetGameByID(ctx context.Context, id int64) (database.Game, error)
	UpdateGame(ctx context.Context, arg database.UpdateGameParams) (database.Game, error)
}
