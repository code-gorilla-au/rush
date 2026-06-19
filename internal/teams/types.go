package teams

import (
	"errors"
	"time"
)

type Team struct {
	ID        int64     `json:"id,omitzero"`
	Name      string    `json:"name,omitzero"`
	CoachID   int       `json:"coach_id,omitzero"`
	Players   []Player  `json:"players,omitzero"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

type Coach struct {
	ID        int64     `json:"id,omitzero"`
	Name      string    `json:"name,omitzero"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

type Player struct {
	ID        int64     `json:"id,omitzero"`
	Name      string    `json:"name,omitzero"`
	TeamID    int       `json:"team_id,omitzero"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

var (
	ErrCoachNotFound = errors.New("coach not found")
	ErrTeamNotFound  = errors.New("team not found")
)
