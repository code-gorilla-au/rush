package game

import "errors"

type Game struct {
	Rounds [10]Round
}

type Round struct {
	TeamA Squad
	TeamB Squad
}

type Squad struct {
	Lanes [3][]int
}

type Result struct {
	TeamA            bool
	TeamB            bool
	RemainingPlayers int
}

var (
	ErrNoPlayer = errors.New("no player left in lane")
)
