package teams

import (
	"time"
)

type Team struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CoachID   int       `json:"coach_id,omitempty"`
	Players   []Player  `json:"players,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Coach struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Player struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	TeamID    int       `json:"team_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
