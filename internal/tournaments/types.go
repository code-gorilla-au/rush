package tournaments

import "github.com/code-gorilla-au/rush/internal/games"

type Tournament struct {
	ID     int64    `json:"id"`
	Name   string   `json:"name"`
	Stages []Stages `json:"stages"`
}

type Stages struct {
	ID    int64        `json:"id"`
	Name  string       `json:"name"`
	Games []games.Game `json:"games"`
}

type Service struct {
}
