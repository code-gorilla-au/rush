package games

import (
	"math/rand/v2"

	"github.com/code-gorilla-au/rush/internal/teams"
)

type RollEngine struct {
	TeamA teams.TokenName
	TeamB teams.TokenName
}

type RollResult struct {
	Roll int
	Team RoundOutcome
}

func DiceRoll() int {
	return rand.IntN(6) + 1
}
