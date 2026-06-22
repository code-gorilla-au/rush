package game

import "errors"

const (
	MaxRounds = 10
)

type Game struct {
	rounds       [10]Round
	currentRound int
	results      []Result
}

type Round struct {
	TeamA Squad
	TeamB Squad
}

type Squad struct {
	Lanes [3][]int
}

type SquadLanes struct {
	Lane1 int
	Lane2 int
	Lane3 int
}

type Result struct {
	TeamA            bool
	TeamB            bool
	RemainingPlayers int
}

type RollFn func() int

var (
	ErrNoPlayer = errors.New("no player left in lane")
	ErrNoRounds = errors.New("no rounds left")
)
