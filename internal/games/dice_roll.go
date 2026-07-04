package games

import (
	"math/rand/v2"
)

func DiceRoll() int {
	return rand.IntN(6) + 1
}
