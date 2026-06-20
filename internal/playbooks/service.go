package playbooks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/database"
)

type Service struct {
	store Store
}

func NewPlaybooksService(store Store) *Service {
	return &Service{store: store}
}

type PlaybookParams struct {
	TeamID      int64
	Name        string
	Description string
	Formations  []Formation
}

func (s *Service) CreatePlaybook(ctx context.Context, params PlaybookParams) (Playbook, error) {

	data, err := json.Marshal(params.Formations)
	if err != nil {
		return Playbook{}, fmt.Errorf("failed to marshal formations: %w", err)
	}

	model, err := s.store.CreatePlaybook(ctx, database.CreatePlaybookParams{
		Name: params.Name,
		Description: sql.NullString{
			String: params.Description,
			Valid:  true,
		},
		Formations: data,
		TeamID: sql.NullInt64{
			Int64: params.TeamID,
			Valid: true,
		},
	})
	if err != nil {
		return Playbook{}, err
	}

	return fromPlaybookModel(model)
}

func (s *Service) GetTeamPlaybooks(ctx context.Context, teamID int64) ([]Playbook, error) {
	models, err := s.store.GetPlaybooksByTeam(ctx, sql.NullInt64{
		Valid: true,
		Int64: teamID,
	})
	if err != nil {
		return nil, err
	}

	return fromPlaybookModels(models)
}
