package teams

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/database"
)

type Service struct {
	store       Store
	playbookSvc PlaybookCreator
}

func NewTeamsService(store Store, playbookSvc PlaybookCreator) *Service {
	return &Service{
		store:       store,
		playbookSvc: playbookSvc,
	}
}

type CreateCoachParams struct {
	Name      string
	IsHuman   bool
	IsDefault bool
}

func (s *Service) CreateCoach(ctx context.Context, params CreateCoachParams) (Coach, error) {
	model, err := s.store.CreateCoach(ctx, database.CreateCoachParams{
		Name: params.Name,
		IsHuman: sql.NullBool{
			Bool:  params.IsHuman,
			Valid: true,
		},
		IsDefault: sql.NullBool{
			Bool:  params.IsDefault,
			Valid: true,
		},
	})
	if err != nil {
		return Coach{}, fmt.Errorf("creating coach: %w", err)
	}

	return fromCoachModel(model), nil
}

func (s *Service) ListCoaches(ctx context.Context) ([]Coach, error) {
	models, err := s.store.GetCoaches(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing coaches: %w", err)
	}

	return fromCoachesModel(models), nil
}

func (s *Service) ListAICoaches(ctx context.Context) ([]Coach, error) {
	models, err := s.store.GetAICoaches(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing ai coaches: %w", err)
	}

	return fromCoachesModel(models), nil
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

func (s *Service) CreateTeam(ctx context.Context, name string, coachID int64, isDefault bool) (Team, error) {
	model, err := s.store.CreateTeam(ctx, database.CreateTeamParams{
		Name: name,
		IsDefault: sql.NullBool{
			Bool:  isDefault,
			Valid: true,
		},
		CoachID: sql.NullInt64{
			Int64: coachID,
			Valid: true,
		},
	})
	if err != nil {
		return Team{}, fmt.Errorf("creating team %s: %w", name, err)
	}

	playersModel, err := s.createPlayers(ctx, model.ID)
	if err != nil {
		return Team{}, fmt.Errorf("creating players: %w", err)
	}

	return fromTeamModel(model, playersModel), nil
}

func (s *Service) createPlayers(ctx context.Context, teamID int64) ([]database.Player, error) {
	modelPlayers := make([]database.Player, 5)

	for i := 0; i < 5; i++ {
		model, err := s.store.CreatePlayer(ctx, database.CreatePlayerParams{
			Name: "Player " + fmt.Sprint(i+1),
			TeamID: sql.NullInt64{
				Int64: teamID,
				Valid: true,
			},
		})
		if err != nil {
			return modelPlayers, fmt.Errorf("creating player: %w", err)
		}

		modelPlayers[i] = model
	}

	return modelPlayers, nil
}

func (s *Service) SetDefaultTeam(ctx context.Context, id int64) error {
	return s.store.SetDefaultTeam(ctx, id)
}

func (s *Service) UpdatePlayer(ctx context.Context, id int64, name string) error {
	err := s.store.UpdatePlayer(ctx, database.UpdatePlayerParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("updating player: %w", err)
	}

	return nil
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
