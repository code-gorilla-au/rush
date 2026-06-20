package playbooks

import "time"

type Playbook struct {
	ID          int64
	Name        string
	TeamID      int64
	Description string
	Formations  []Formation
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Formation struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Lane1       int    `json:"lane_1,omitempty"`
	Lane2       int    `json:"lane_2,omitempty"`
	Lane3       int    `json:"lane_3,omitempty"`
}
