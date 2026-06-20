package playbooks

import (
	"context"
	"database/sql"

	"github.com/code-gorilla-au/rush/internal/database"
)

type Store interface {
	CreatePlaybook(ctx context.Context, arg database.CreatePlaybookParams) (database.Playbook, error)
	DeletePlaybook(ctx context.Context, id int64) error
	GetPlaybooksByTeam(ctx context.Context, teamID sql.NullInt64) ([]database.Playbook, error)
	UpdatePlaybookFormations(ctx context.Context, arg database.UpdatePlaybookFormationsParams) (database.Playbook, error)
}
