package teams

import (
	"context"
	"database/sql"

	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type Store interface {
	PlayerStore
	TeamStore
	CoachStore
}

type CoachStore interface {
	GetDefaultCoach(ctx context.Context) (database.Coach, error)
	ClearDefaultCoach(ctx context.Context) error
	CreateCoach(ctx context.Context, arg database.CreateCoachParams) (database.Coach, error)
	GetCoaches(ctx context.Context) ([]database.Coach, error)
	GetAICoaches(ctx context.Context) ([]database.Coach, error)
	SetDefaultCoach(ctx context.Context, id int64) error
	SetDefaultTeam(ctx context.Context, id int64) error
}

type PlayerStore interface {
	CreatePlayer(ctx context.Context, arg database.CreatePlayerParams) (database.Player, error)
	GetTeamMembers(ctx context.Context, teamID sql.NullInt64) ([]database.Player, error)
	UpdatePlayer(ctx context.Context, arg database.UpdatePlayerParams) error
}

type TeamStore interface {
	CreateTeam(ctx context.Context, arg database.CreateTeamParams) (database.Team, error)
	GetTeamByCoachID(ctx context.Context, coachID sql.NullInt64) (database.Team, error)
	DeleteTeam(ctx context.Context, id int64) error
	SetDefaultTeam(ctx context.Context, id int64) error
	ClearDefaultTeam(ctx context.Context) error
}

var _ Store = (*database.Queries)(nil)

type PlaybookCreatGetter interface {
	CreatePlaybook(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error)
	GetTeamPlaybooks(ctx context.Context, teamID int64) ([]playbooks.Playbook, error)
}
