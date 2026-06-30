package uistate

import (
	"context"

	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type MsgStateUpdated struct {
	Coach *teams.Coach
	Team  *teams.Team
}

type MsgSwitchPage struct {
	NewPage  Page
	Playbook *playbooks.Playbook
	GameID   int64
}

type Page int

const (
	PageTitle Page = iota + 1
	PageCreateCoach
	PageLockerRoom
	PageNewTournament
	PageNewBattleSelection
	PageTitleSettings
	PageGame
)

type GlobalState struct {
	Coach *teams.Coach
	Team  *teams.Team
}

func (m *GlobalState) Context() context.Context {
	return context.Background()
}
