package teams

import "context"

type Service struct {
	store Store
}

func NewTeamsService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateCoach(ctx context.Context, name string) error {
	return s.store.CreateCoach(ctx, name)
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
