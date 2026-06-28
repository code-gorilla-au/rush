package games

import (
	"errors"
	"time"

	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type Service struct {
	Store Store
}

type GameStatus string

const (
	StatusPending  GameStatus = "pending"
	StatusRunning  GameStatus = "running"
	StatusComplete GameStatus = "complete"
)

type Game struct {
	id           int64
	name         string
	tournamentID *int64
	teamA        int64
	teamB        int64
	winner       *int64
	status       GameStatus
	rounds       [10]Round
	currentRound int64
	results      []Result
	createdAt    time.Time
	updatedAt    time.Time
}

func (g Game) ID() int64 {
	return g.id
}

type Round struct {
	TeamA TeamFormation
	TeamB TeamFormation
}

type TeamConfig struct {
	TeamID     int64
	TeamName   string
	Formations []playbooks.Formation
}

type TeamFormation struct {
	TeamID int64
	Lanes  [3][]int
}

type LanesConfig struct {
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
