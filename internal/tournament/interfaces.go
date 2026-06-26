package tournament

import (
	"context"

	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type TeamCreatLister interface {
	CreateCoach(ctx context.Context, params teams.CreateCoachParams) (teams.Coach, error)
	ListAICoaches(ctx context.Context) ([]teams.Coach, error)
	GetTeamByCoachID(ctx context.Context, id int64) (teams.Team, error)
	CreateTeam(ctx context.Context, name string, coachID int64, isDefault bool) (teams.Team, error)
}

type PlaybookCreator interface {
	CreatePlaybook(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error)
	GetTeamPlaybooks(ctx context.Context, teamID int64) ([]playbooks.Playbook, error)
}
