package tournament

import (
	"context"
	"errors"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type mockTeamCreator struct {
	createCoachFunc      func(ctx context.Context, params teams.CreateCoachParams) (teams.Coach, error)
	listAICoachesFunc    func(ctx context.Context) ([]teams.Coach, error)
	getTeamByCoachIDFunc func(ctx context.Context, id int64) (teams.Team, error)
	createTeamFunc       func(ctx context.Context, name string, coachID int64, isDefault bool) (teams.Team, error)
}

func (m *mockTeamCreator) CreateCoach(ctx context.Context, params teams.CreateCoachParams) (teams.Coach, error) {
	return m.createCoachFunc(ctx, params)
}

func (m *mockTeamCreator) ListAICoaches(ctx context.Context) ([]teams.Coach, error) {
	return m.listAICoachesFunc(ctx)
}

func (m *mockTeamCreator) GetTeamByCoachID(ctx context.Context, id int64) (teams.Team, error) {
	return m.getTeamByCoachIDFunc(ctx, id)
}

func (m *mockTeamCreator) CreateTeam(ctx context.Context, name string, coachID int64, isDefault bool) (teams.Team, error) {
	return m.createTeamFunc(ctx, name, coachID, isDefault)
}

type mockPlaybookCreator struct {
	createPlaybookFunc   func(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error)
	getTeamPlaybooksFunc func(ctx context.Context, teamID int64) ([]playbooks.Playbook, error)
}

func (m *mockPlaybookCreator) CreatePlaybook(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error) {
	return m.createPlaybookFunc(ctx, params)
}

func (m *mockPlaybookCreator) GetTeamPlaybooks(ctx context.Context, teamID int64) ([]playbooks.Playbook, error) {
	return m.getTeamPlaybooksFunc(ctx, teamID)
}

func TestGenerateAITeams(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should generate exactly 12 teams with personas and formations", func(t *testing.T) {
			_teams, err := generateAITeams()
			odize.AssertNoError(t, err)
			odize.AssertTrue(t, len(_teams) == 12)

			for i, team := range _teams {
				odize.AssertTrue(t, team.Persona != "")
				odize.AssertTrue(t, len(team.Formations) == 10)

				// Verify personas are assigned in order
				expectedPersona := personas[i%len(personas)].Name
				odize.AssertTrue(t, team.Persona == expectedPersona)
			}
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestAITeamService_GenerateTeams(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var mockTeams *mockTeamCreator
	var mockPlaybooks *mockPlaybookCreator
	var svc *AITeamService

	group.BeforeEach(func() {
		mockTeams = &mockTeamCreator{}
		mockPlaybooks = &mockPlaybookCreator{}
		svc = &AITeamService{
			teamsSvc:     mockTeams,
			playbooksSvc: mockPlaybooks,
		}
	})

	err := group.
		Test("should successfully generate all teams", func(t *testing.T) {
			var coachCount, teamCount, playbookCount int

			mockTeams.createCoachFunc = func(ctx context.Context, params teams.CreateCoachParams) (teams.Coach, error) {
				coachCount++
				return teams.Coach{ID: int64(coachCount), Name: params.Name}, nil
			}
			mockTeams.createTeamFunc = func(ctx context.Context, name string, coachID int64, isDefault bool) (teams.Team, error) {
				teamCount++
				return teams.Team{ID: int64(teamCount), Name: name}, nil
			}
			mockPlaybooks.createPlaybookFunc = func(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error) {
				playbookCount++
				return playbooks.Playbook{ID: int64(playbookCount), Name: params.Name}, nil
			}

			err := svc.GenerateTeams(t.Context())
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, totalTeams, coachCount)
			odize.AssertEqual(t, totalTeams, teamCount)
			odize.AssertEqual(t, totalTeams, playbookCount)
		}).
		Test("should return error when CreateCoach fails", func(t *testing.T) {
			expectedErr := errors.New("coach error")
			mockTeams.createCoachFunc = func(ctx context.Context, params teams.CreateCoachParams) (teams.Coach, error) {
				return teams.Coach{}, expectedErr
			}

			err := svc.GenerateTeams(t.Context())
			odize.AssertError(t, err)
			odize.AssertTrue(t, errors.Is(err, expectedErr))
		}).
		Test("should return error when CreateTeam fails", func(t *testing.T) {
			expectedErr := errors.New("team error")
			mockTeams.createCoachFunc = func(ctx context.Context, params teams.CreateCoachParams) (teams.Coach, error) {
				return teams.Coach{ID: 1}, nil
			}
			mockTeams.createTeamFunc = func(ctx context.Context, name string, coachID int64, isDefault bool) (teams.Team, error) {
				return teams.Team{}, expectedErr
			}

			err := svc.GenerateTeams(t.Context())
			odize.AssertError(t, err)
			odize.AssertTrue(t, errors.Is(err, expectedErr))
		}).
		Test("should return error when CreatePlaybook fails", func(t *testing.T) {
			expectedErr := errors.New("playbook error")
			mockTeams.createCoachFunc = func(ctx context.Context, params teams.CreateCoachParams) (teams.Coach, error) {
				return teams.Coach{ID: 1}, nil
			}
			mockTeams.createTeamFunc = func(ctx context.Context, name string, coachID int64, isDefault bool) (teams.Team, error) {
				return teams.Team{ID: 1, Name: name}, nil
			}
			mockPlaybooks.createPlaybookFunc = func(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error) {
				return playbooks.Playbook{}, expectedErr
			}

			err := svc.GenerateTeams(t.Context())
			odize.AssertError(t, err)
			odize.AssertTrue(t, errors.Is(err, expectedErr))
		}).
		Run()

	odize.AssertNoError(t, err)
}
