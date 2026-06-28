package games

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/database"
)

func generateRounds(teamA TeamConfig, teamB TeamConfig) [10]Round {
	rounds := [10]Round{}

	for i := range rounds {
		r := NewRound()

		r.FillSquad(
			LanesConfig{
				TeamID: teamA.TeamID,
				Lane1:  teamA.Formations[i].Lane1,
				Lane2:  teamA.Formations[i].Lane2,
				Lane3:  teamA.Formations[i].Lane3,
			},
			LanesConfig{
				TeamID: teamB.TeamID,
				Lane1:  teamB.Formations[i].Lane1,
				Lane2:  teamB.Formations[i].Lane2,
				Lane3:  teamB.Formations[i].Lane3,
			},
		)

		rounds[i] = r
	}

	return rounds
}

func (g *Game) ResolveRound(roll RollFn) (Result, error) {
	if g.currentRound < 0 || g.currentRound >= int64(len(g.rounds)) {
		return Result{}, ErrNoRounds
	}

	round := &g.rounds[int(g.currentRound)]
	result := round.ResolveLanes(roll)

	g.results = append(g.results, result)
	g.currentRound++

	return result, nil
}

func (g *Game) IsGameComplete() bool {
	return g.currentRound >= int64(len(g.rounds))
}

func (g *Game) ID() int64 {
	return g.id
}

func (g *Game) Rounds() [10]Round {
	return g.rounds
}

func (g *Game) CurrentRound() int64 {
	return g.currentRound
}

func (g *Game) Name() string {
	return g.name
}

func fromGameModel(m database.Game) (Game, error) {

	var rounds [10]Round
	if err := json.Unmarshal(m.Rounds, &rounds); err != nil {
		return Game{}, fmt.Errorf("failed to unmarshal game model: %w", err)
	}

	var results []Result
	if err := json.Unmarshal(m.ResultsLog, &results); err != nil {
		return Game{}, fmt.Errorf("failed to unmarshal results log: %w", err)
	}

	return Game{
		id:           m.ID,
		name:         m.Name,
		tournamentID: &m.TournamentID.Int64,
		teamA:        m.TeamA.Int64,
		teamB:        m.TeamB.Int64,
		winner:       &m.Winner.Int64,
		status:       GameStatus(m.Status),
		rounds:       rounds,
		currentRound: m.CurrentRound,
		results:      results,
		createdAt:    m.CreatedAt.Time,
		updatedAt:    m.UpdatedAt.Time,
	}, nil
}

func toGameModel(g Game) (database.Game, error) {
	resolvedTournamentID := sql.NullInt64{}
	if g.tournamentID != nil {
		resolvedTournamentID = sql.NullInt64{
			Int64: *g.tournamentID,
			Valid: true,
		}
	}

	resolvedWinner := sql.NullInt64{}
	if g.winner != nil {
		resolvedWinner = sql.NullInt64{
			Int64: *g.winner,
			Valid: true,
		}
	}

	roundData, err := json.Marshal(g.rounds)
	if err != nil {
		return database.Game{}, fmt.Errorf("failed to marshal game model: %w", err)
	}

	resultsData, err := json.Marshal(g.results)
	if err != nil {
		return database.Game{}, fmt.Errorf("failed to marshal results model: %w", err)
	}

	return database.Game{
		ID:           g.id,
		Name:         g.name,
		TournamentID: resolvedTournamentID,
		TeamA: sql.NullInt64{
			Int64: g.teamA,
			Valid: true,
		},
		TeamB: sql.NullInt64{
			Int64: g.teamB,
			Valid: true,
		},
		Winner:       resolvedWinner,
		Status:       string(g.status),
		Rounds:       roundData,
		CurrentRound: g.currentRound,
		ResultsLog:   resultsData,
		CreatedAt: sql.NullTime{
			Time:  g.createdAt,
			Valid: true,
		},
		UpdatedAt: sql.NullTime{
			Time:  g.updatedAt,
			Valid: true,
		},
	}, nil

}
