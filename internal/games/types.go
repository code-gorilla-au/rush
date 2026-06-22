package games

import "errors"

const (
	MaxRounds = 10
)

type Service struct {
	Store Store
}

type Game struct {
	id           int64
	rounds       [10]Round
	currentRound int
	results      []Result
}

type Round struct {
	TeamA Squad
	TeamB Squad
}

type Squad struct {
	TeamID int64
	Lanes  [3][]int
}

type SquadConfig struct {
	TeamID int64
	Lane1  int
	Lane2  int
	Lane3  int
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
