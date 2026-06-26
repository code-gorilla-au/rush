package games

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/database"
)

func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

type NewGameParams struct {
	TeamA        TeamConfig
	TeamB        TeamConfig
	TournamentID *int64
}

func (s *Service) NewGame(ctx context.Context, params NewGameParams) (Game, error) {
	resolvedTournamentID := sql.NullInt64{}
	if params.TournamentID != nil {
		resolvedTournamentID = sql.NullInt64{
			Int64: *params.TournamentID,
			Valid: true,
		}
	}

	roundsJsonData, rErr := json.Marshal(generateRounds(params.TeamA, params.TeamB))
	if rErr != nil {
		return Game{}, fmt.Errorf("failed to marshal rounds json data: %w", rErr)
	}

	model, err := s.Store.CreateGame(ctx, database.CreateGameParams{
		Name: fmt.Sprintf("%s VS %s", params.TeamA.TeamName, params.TeamB.TeamName),
		TeamA: sql.NullInt64{
			Int64: params.TeamA.TeamID,
			Valid: true,
		},
		TeamB: sql.NullInt64{
			Int64: params.TeamB.TeamID,
			Valid: true,
		},
		TournamentID: resolvedTournamentID,
		ResultsLog:   nil,
		Rounds:       roundsJsonData,
		CurrentRound: 0,
	})
	if err != nil {
		return Game{}, fmt.Errorf("creating game: %w", err)
	}

	return fromGameModel(model)
}

func (s *Service) UpdateGame(ctx context.Context, game Game) (Game, error) {
	model, err := toGameModel(game)
	if err != nil {
		return Game{}, fmt.Errorf("failed to convert game to model: %w", err)
	}

	updated, err := s.Store.UpdateGame(ctx, database.UpdateGameParams{
		Name:         model.Name,
		TeamA:        model.TeamA,
		TeamB:        model.TeamB,
		Winner:       model.Winner,
		Status:       model.Status,
		ResultsLog:   model.ResultsLog,
		Rounds:       model.Rounds,
		CurrentRound: model.CurrentRound,
		TournamentID: model.TournamentID,
		ID:           model.ID,
	})
	if err != nil {
		return Game{}, fmt.Errorf("updating game: %w", err)
	}

	return fromGameModel(updated)
}

func (s *Service) CompleteGame(ctx context.Context, game Game) (Game, error) {
	game.status = StatusComplete

	updatedGame, err := s.UpdateGame(ctx, game)
	if err != nil {
		return Game{}, fmt.Errorf("completing game: %w", err)
	}

	return updatedGame, nil
}
