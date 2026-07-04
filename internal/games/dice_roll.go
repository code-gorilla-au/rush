package games

import (
	"math/rand/v2"

	"github.com/code-gorilla-au/rush/internal/augments"
)

func DiceRoll() int {
	return rand.IntN(6) + 1
}

type TeamDecisionInput struct {
	augment augments.Name
	roll    int
}

type DecisionInput struct {
	lastRound    RoundResult
	teamAAugment augments.Name
	teamBAugment augments.Name
}
