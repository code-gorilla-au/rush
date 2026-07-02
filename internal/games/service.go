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
		ResultsLog:   []byte(`[]`),
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

func (s *Service) GetGame(ctx context.Context, id int64) (Game, error) {
	model, err := s.Store.GetGameByID(ctx, id)
	if err != nil {
		return Game{}, fmt.Errorf("getting game: %w", err)
	}

	return fromGameModel(model)
}

func (s *Service) CompleteGame(ctx context.Context, game Game) (Game, error) {
	game.status = StatusComplete

	winner, err := game.CalculateWinner()
	if err != nil {
		return Game{}, fmt.Errorf("calculating winner: %w", err)
	}

	game.winner = &winner

	updatedGame, err := s.UpdateGame(ctx, game)
	if err != nil {
		return Game{}, fmt.Errorf("completing game: %w", err)
	}

	return updatedGame, nil
}

func (s *Service) GetTeamStatistics(ctx context.Context, teamID int64) (TeamStatistics, error) {
	teamIDArg := sql.NullInt64{Int64: teamID, Valid: true}

	models, err := s.Store.ListCompletedGamesByTeam(ctx, database.ListCompletedGamesByTeamParams{
		TeamA: teamIDArg,
		TeamB: teamIDArg,
	})
	if err != nil {
		return TeamStatistics{}, fmt.Errorf("listing completed games: %w", err)
	}

	stats := TeamStatistics{}
	for _, model := range models {
		stats.GamesPlayed++

		switch {
		case !model.Winner.Valid || model.Winner.Int64 == 0:
			stats.Draws++
		case model.Winner.Int64 == teamID:
			stats.Wins++
		default:
			stats.Losses++
		}

		var results []Result
		if err := json.Unmarshal(model.ResultsLog, &results); err != nil {
			return TeamStatistics{}, fmt.Errorf("parsing game %d results: %w", model.ID, err)
		}

		teamRoundsWon := len(filterResultsByTeam(ResultTeamA, results))
		teamRoundsLost := len(filterResultsByTeam(ResultTeamB, results))
		if model.TeamB.Valid && model.TeamB.Int64 == teamID {
			teamRoundsWon, teamRoundsLost = teamRoundsLost, teamRoundsWon
		}

		stats.RoundsWon += teamRoundsWon
		stats.RoundsLost += teamRoundsLost
	}

	if stats.GamesPlayed == 0 {
		return stats, nil
	}

	stats.WinRate = (float64(stats.Wins) / float64(stats.GamesPlayed)) * 100
	stats.RoundDifferential = stats.RoundsWon - stats.RoundsLost
	stats.AverageRoundsWon = float64(stats.RoundsWon) / float64(stats.GamesPlayed)
	stats.AverageRoundsLost = float64(stats.RoundsLost) / float64(stats.GamesPlayed)

	return stats, nil
}
