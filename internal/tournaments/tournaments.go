package tournaments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type ServiceDependencies struct {
	GamesSvc *games.Service
	TeamsSvc TeamsService
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

type CreateTournamentParams struct {
	Name          string
	NumberOfTeams NumberOfTeams
	CoachID       int64
}

func (s *Service) CreateTournament(ctx context.Context, params CreateTournamentParams) error {
	t, err := s.insertNewTournament(ctx, params.Name, params.NumberOfTeams)
	if err != nil {
		return fmt.Errorf("inserting new tournament: %w", err)
	}

	gameConfigs, err := s.generateGroupStageGames(ctx, params.CoachID, int64(params.NumberOfTeams))
	if err != nil {
		return fmt.Errorf("generating games for tournament: %w", err)
	}

	groupStage := t.Stages[0]

	if err = s.allocateGamesToStage(ctx, gameConfigs, groupStage); err != nil {
		return err
	}

	return nil
}

func (s *Service) allocateGamesToStage(ctx context.Context, gameConfigs []games.NewGameParams, stage Stage) error {
	for _, gameConfig := range gameConfigs {
		g, gErr := s.gamesSvc.NewGame(ctx, gameConfig)

		if gErr != nil {
			return fmt.Errorf("creating new game: %w", gErr)
		}

		if _, err := s.store.AllocateGameToStage(ctx, database.AllocateGameToStageParams{
			StageID: sql.NullInt64{
				Int64: stage.ID,
				Valid: true,
			},
			GameID: sql.NullInt64{
				Int64: g.ID(),
				Valid: true,
			},
		}); err != nil {
			return fmt.Errorf("allocating game to stage: %w", err)
		}
	}
	return nil
}

func (s *Service) UpdateStageStatus(ctx context.Context, tournamentID int64, stageID int64, status StageStatus) error {
	if err := s.store.SetStageStatus(ctx, database.SetStageStatusParams{
		ID:           stageID,
		TournamentID: sql.NullInt64{Int64: tournamentID, Valid: true},
		Status:       string(status),
	}); err != nil {
		return fmt.Errorf("setting stage status: %w", err)
	}
	return nil
}

func (s *Service) insertNewTournament(ctx context.Context, name string, numberOfTeams NumberOfTeams) (Tournament, error) {
	var newTournament database.Tournament
	var stage database.Stage

	err := database.WithTxnCtx(s.DB, func(tx *sql.Tx) error {
		txDb := s.txnFunc(tx)

		var err error

		newTournament, err = txDb.CreateTournament(ctx, database.CreateTournamentParams{
			Name:          name,
			NumberOfTeams: int64(numberOfTeams),
		})
		if err != nil {
			return err
		}

		stage, err = txDb.CreateStage(ctx, database.CreateStageParams{
			Name: StageNameGroup,
			TournamentID: sql.NullInt64{
				Int64: newTournament.ID,
				Valid: true,
			},
			Status: string(games.StatusPending),
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Tournament{}, fmt.Errorf("failed to create tournament: %w", err)
	}

	return toTournament(newTournament, []database.Stage{stage}), nil

}

func (s *Service) generateGroupStageGames(ctx context.Context, coachId int64, numberOfTeams int64) ([]games.NewGameParams, error) {
	var tournamentGames []games.NewGameParams

	var totalTeams []teams.AITeam
	h, err := s.getHumanTeam(ctx, coachId)
	if err != nil {
		return tournamentGames, fmt.Errorf("failed to get human team: %w", err)
	}

	totalTeams = append(totalTeams, h)

	aiTeams, err := s.getNonHumanTeams(ctx, numberOfTeams-1)
	if err != nil {
		return tournamentGames, fmt.Errorf("failed to get non-human teams: %w", err)
	}

	totalTeams = append(totalTeams, aiTeams...)

	tournamentGames = generateGameParamsFromTeams(totalTeams)

	return tournamentGames, nil
}

func (s *Service) getNonHumanTeams(ctx context.Context, aiTeams int64) ([]teams.AITeam, error) {
	var tournamentList []teams.AITeam
	var err error

	tournamentList, err = s.teamsSvc.ListAITeams(ctx)
	if err != nil {
		return tournamentList, fmt.Errorf("failed to list teams: %w", err)
	}

	if int64(len(tournamentList)) < aiTeams {
		return tournamentList, errors.New("not enough teams")
	}

	return tournamentList[0:aiTeams], nil

}

func (s *Service) getHumanTeam(ctx context.Context, coachID int64) (teams.AITeam, error) {
	coach, err := s.teamsSvc.GetCoachByID(ctx, coachID)
	if err != nil {
		return teams.AITeam{}, fmt.Errorf("failed to get coach: %w", err)
	}

	team, err := s.teamsSvc.GetTeamAndPlaybooksByCoachID(ctx, coachID)
	if err != nil {
		return teams.AITeam{}, fmt.Errorf("failed to get team: %w", err)
	}

	if len(team.Playbooks) == 0 {
		return teams.AITeam{}, errors.New("no playbooks found for team")
	}

	return teams.AITeam{
		Coach:    coach,
		Team:     team.Team,
		Playbook: team.Playbooks[0],
	}, nil
}

func generateGameParamsFromTeams(totalTeams []teams.AITeam) []games.NewGameParams {

	var tournamentGames []games.NewGameParams

	for i := 0; i < len(totalTeams); i++ {

		for j := i + 1; j < len(totalTeams); j++ {
			first := totalTeams[i]
			second := totalTeams[j]

			tournamentGames = append(tournamentGames, games.NewGameParams{
				TeamA: games.TeamConfig{
					TeamID:     first.Team.ID,
					TeamName:   first.Team.Name,
					Players:    teams.GetPlayerIDs(first.Team),
					Augments:   first.Coach.AvailableAugments(),
					Formations: first.Playbook.Formations,
				},
				TeamB: games.TeamConfig{
					TeamID:     second.Team.ID,
					TeamName:   second.Team.Name,
					Players:    teams.GetPlayerIDs(second.Team),
					Augments:   second.Coach.AvailableAugments(),
					Formations: second.Playbook.Formations,
				},
			})
		}

	}

	return tournamentGames
}
