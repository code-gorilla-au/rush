package tournaments

import (
	"database/sql"

	"github.com/code-gorilla-au/rush/internal/games"
)

type NumberOfTeams int64

const (
	Four  NumberOfTeams = 4
	Eight NumberOfTeams = 8
)

type Tournament struct {
	ID     int64         `json:"id"`
	Name   string        `json:"name"`
	Number NumberOfTeams `json:"number_of_teams"`
	Stages []Stage       `json:"stages"`
}

type StageStatus string

const (
	StageStatusActive   StageStatus = "active"
	StageStatusPending  StageStatus = "pending"
	StageStatusComplete StageStatus = "complete"
)

const (
	StageNameGroup  = "Group"
	StageNameKnock  = "Knockout"
	StageNameFinals = "finals"
)

type Stage struct {
	ID     int64        `json:"id"`
	Name   string       `json:"name"`
	Status StageStatus  `json:"status"`
	Games  []games.Game `json:"games"`
}

type Service struct {
	teamsSvc TeamsService
	gamesSvc *games.Service
	store    Store
	DB       *sql.DB
	txnFunc  func(db *sql.Tx) Store
}
