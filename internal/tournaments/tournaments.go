package tournaments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type ServiceDependencies struct {
	GamesSvc *games.Service
	TeamsSvc *teams.Service
	Store    Store
	DB       *sql.DB
	TxnFunc  func(db *sql.Tx) Store
}

func NewService(deps ServiceDependencies) *Service {
	return &Service{
		gamesSvc: deps.GamesSvc,
		teamsSvc: deps.TeamsSvc,
		store:    deps.Store,
		DB:       deps.DB,
		txnFunc:  deps.TxnFunc,
	}
}

func (s *Service) CreateTournament(ctx context.Context, coachId int64, numberOfTeams NumberOfTeams) error {
	return nil
}

func (s *Service) getNonHumanTeams(ctx context.Context, aiTeams int) ([]teams.AITeam, error) {
	var tournamentList []teams.AITeam
	var err error

	tournamentList, err = s.teamsSvc.ListAITeams(ctx)
	if err != nil {
		return tournamentList, fmt.Errorf("failed to list teams: %w", err)
	}

	if len(tournamentList) < aiTeams {
		return tournamentList, errors.New("not enough teams")
	}

	return tournamentList[0:aiTeams], nil

}
