package teams

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/database"
)

type Service struct {
	store Store
}

func NewTeamsService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateCoach(ctx context.Context, name string, isDefault bool) (Coach, error) {
	model, err := s.store.CreateCoach(ctx, database.CreateCoachParams{
		Name: name,
		IsDefault: sql.NullBool{
			Bool:  isDefault,
			Valid: true,
		},
	})
	if err != nil {
		return Coach{}, fmt.Errorf("creating coach: %w", err)
	}

	return fromCoachModel(model), nil
}

func (s *Service) GetDefaultCoach(ctx context.Context) (Coach, error) {
	model, err := s.store.GetDefaultCoach(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Coach{}, ErrCoachNotFound
		}

		return Coach{}, fmt.Errorf("getting default coach: %w", err)
	}

	return fromCoachModel(model), nil
}

func (s *Service) GetTeamByCoachID(ctx context.Context, id int64) (Team, error) {
	model, err := s.store.GetTeamByCoachID(ctx, sql.NullInt64{
		Valid: true,
		Int64: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Team{}, ErrTeamNotFound
		}

		return Team{}, fmt.Errorf("getting team: %w", err)
	}

	pModel, err := s.store.GetTeamMembers(ctx, sql.NullInt64{
		Int64: model.ID,
		Valid: true,
	})
	if err != nil {
		return Team{}, fmt.Errorf("getting team members: %w", err)
	}

	return fromTeamModel(model, pModel), nil
}

func (s *Service) SetDefaultTeam(ctx context.Context, id int64) error {
	return s.store.SetDefaultTeam(ctx, id)
}

func (s *Service) ClearDefaultTeam(ctx context.Context) error {
	return s.store.ClearDefaultTeam(ctx)
}

func (s *Service) SetDefaultCoach(ctx context.Context, id int64) error {
	return s.store.SetDefaultCoach(ctx, id)
}

func (s *Service) ClearDefaultCoach(ctx context.Context) error {
	return s.store.ClearDefaultCoach(ctx)
}
