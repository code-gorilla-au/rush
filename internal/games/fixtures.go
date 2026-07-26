package games

import (
	"database/sql"
	"encoding/json"

	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

// NewTestTeamConfig creates a TeamConfig for testing purposes
func NewTestTeamConfig(id int64, name string) TeamConfig {
	return TeamConfig{
		TeamID:     id,
		TeamName:   name,
		Players:    []int64{101, 102, 103},
		Formations: make([]playbooks.Formation, 10),
	}
}

// CreateTestGame creates a Game for testing purposes
func CreateTestGame(teamA, teamB int64, winner int64, rounds [10]Round, results []RoundResult, status GameStatus) (Game, error) {
	winnerID := sql.NullInt64{Int64: winner, Valid: winner != 0}
	if winner == 0 && (len(results) > 0) {
		winnerID = sql.NullInt64{Int64: 0, Valid: true}
	}

	roundsData, _ := json.Marshal(rounds)
	resultsData, _ := json.Marshal(results)

	model := database.Game{
		ID:           1,
		TeamA:        sql.NullInt64{Int64: teamA, Valid: true},
		TeamB:        sql.NullInt64{Int64: teamB, Valid: true},
		Winner:       winnerID,
		Status:       string(status),
		Rounds:       roundsData,
		ResultsLog:   resultsData,
		CurrentRound: int64(len(results)),
	}

	return fromGameModel(model)
}
