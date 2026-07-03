package teams

import (
	"errors"
	"time"

	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type TeamWithPlaybooks struct {
	Playbooks []playbooks.Playbook `json:"playbooks,omitzero"`
	Team      Team                 `json:"team,omitzero"`
}

type Team struct {
	ID        int64     `json:"id,omitzero"`
	Name      string    `json:"name,omitzero"`
	CoachID   int       `json:"coach_id,omitzero"`
	Players   []Player  `json:"players,omitzero"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

type Coach struct {
	ID        int64        `json:"id,omitzero"`
	Name      string       `json:"name,omitzero"`
	Persona   CoachPersona `json:"persona,omitzero"`
	IsDefault bool         `json:"is_default,omitzero"`
	IsHuman   bool         `json:"is_human,omitzero"`
	CreatedAt time.Time    `json:"created_at,omitzero"`
	UpdatedAt time.Time    `json:"updated_at,omitzero"`
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

type TokenName string

const (
	TokenTwistOfFate     TokenName = "Twist of Fate"
	TokenSecondChance    TokenName = "Second Chance"
	TokenPowerPlay       TokenName = "Power Play"
	TokenBrace           TokenName = "Brace"
	TokenPrecisionStrike TokenName = "Precision Strike"
	TokenJammingSignal   TokenName = "Jamming Signal"
	TokenLastStand       TokenName = "Last Stand"
	TokenMomentumSurge   TokenName = "Momentum Surge"
	TokenIceInVeins      TokenName = "Ice in Veins"
	TokenSmokeScreen     TokenName = "Smoke Screen"
)

type CoachPersona string

const (
	CoachPersonaVanguard  CoachPersona = "Vanguard Coach"
	CoachPersonaBastion   CoachPersona = "Bastion Coach"
	CoachPersonaTrickster CoachPersona = "Trickster Coach"
	CoachPersonaWildcard  CoachPersona = "Wildcard Coach"
)
