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

type RoundOutcome string

const (
	ResultDraw  RoundOutcome = "draw"
	ResultTeamA RoundOutcome = "team_a"
	ResultTeamB RoundOutcome = "team_b"
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
	results      []RoundResult
	createdAt    time.Time
	updatedAt    time.Time
}

type Round struct {
	TeamA TeamFormation
	TeamB TeamFormation
}

type TeamStatistics struct {
	GamesPlayed       int
	Wins              int
	Draws             int
	Losses            int
	WinRate           float64
	RoundsWon         int
	RoundsLost        int
	RoundDifferential int
	AverageRoundsWon  float64
	AverageRoundsLost float64
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

type RoundResult struct {
	Outcome          RoundOutcome
	RemainingPlayers int
}

type RollFn func() int

var (
	ErrNoPlayer        = errors.New("no player left in lane")
	ErrNoRounds        = errors.New("no rounds left")
	ErrGameNotComplete = errors.New("game not complete")
)
