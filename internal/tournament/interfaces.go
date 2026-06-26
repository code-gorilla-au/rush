package tournament

import (
	"context"

	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type TeamCreator interface {
	CreateCoach(ctx context.Context, params teams.CreateCoachParams) (teams.Coach, error)
	ListAICoaches(ctx context.Context) ([]teams.Coach, error)
	CreateTeam(ctx context.Context, name string, coachID int64, isDefault bool) (teams.Team, error)
}

type PlaybookCreator interface {
	CreatePlaybook(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error)
}
