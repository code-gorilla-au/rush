package games

import (
	"errors"
	"time"

	"github.com/code-gorilla-au/rush/internal/augments"
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

type Outcome string

const (
	Draw  Outcome = "draw"
	TeamA Outcome = "team_a"
	TeamB Outcome = "team_b"
)

type Game struct {
	id           int64
	name         string
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
	TeamA       TeamFormation `json:"team_a"`
	TeamB       TeamFormation `json:"team_b"`
	DuelResults []DuelResult  `json:"duel_results"`
}

type DuelResult struct {
	Player    int64      `json:"player"`
	Outcome   Outcome    `json:"outcome"`
	Roll      int        `json:"roll"`
	RollDelta int        `json:"roll_delta"`
	Lane      int        `json:"lane"`
	TeamA     PlayerRoll `json:"team_a"`
	TeamB     PlayerRoll `json:"team_b"`
}

type PlayerRoll struct {
	Player           int64         `json:"player"`
	Roll             int           `json:"roll"`
	TriggeredAugment augments.Name `json:"triggered_augment"`
}

type TeamStatistics struct {
	GamesPlayed       int     `json:"games_played"`
	Wins              int     `json:"wins"`
	Draws             int     `json:"draws"`
	Losses            int     `json:"losses"`
	WinRate           float64 `json:"win_rate"`
	RoundsWon         int     `json:"rounds_won"`
	RoundsLost        int     `json:"rounds_lost"`
	RoundDifferential int     `json:"round_differential"`
	AverageRoundsWon  float64 `json:"average_rounds_won"`
	AverageRoundsLost float64 `json:"average_rounds_lost"`
}

type TeamConfig struct {
	TeamID     int64                 `json:"team_id"`
	TeamName   string                `json:"team_name"`
	Players    []int64               `json:"players"`
	Augments   []augments.Effect     `json:"augments"`
	Formations []playbooks.Formation `json:"formations"`
}

type TeamFormation struct {
	TeamID   int64             `json:"team_id"`
	Augments []augments.Effect `json:"augments"`
	Lanes    [3][]int64        `json:"lanes"`
}

type LanesConfig struct {
	TeamID   int64             `json:"team_id"`
	Players  []int64           `json:"players"`
	Augments []augments.Effect `json:"augments"`
	Lane1    int               `json:"lane_1"`
	Lane2    int               `json:"lane_2"`
	Lane3    int               `json:"lane_3"`
}

type RoundResult struct {
	Outcome          Outcome `json:"outcome"`
	RemainingPlayers int     `json:"remaining_players"`
}

type RollFn func() int

var (
	ErrNoPlayer        = errors.New("no player left in lane")
	ErrNoRounds        = errors.New("no rounds left")
	ErrGameNotComplete = errors.New("game not complete")
)

type TeamDecisionInput struct {
	triggeredAugment augments.Name
	passivesAugments []augments.Effect
	player           int64
	roll             int
}

type DecisionInput struct {
	lastRound *DuelResult
	teamA     TeamDecisionInput
	teamB     TeamDecisionInput
}

type DecisionEngineFunc func(input DecisionInput) DecisionInput

type DecisionEngine struct {
	beforeRoll    []DecisionEngineFunc
	afterRoll     []DecisionEngineFunc
	afterAugments []DecisionEngineFunc
	rollFn        RollFn
}
