package teams

import (
	"context"
	"errors"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type mockTeamCreator struct {
	createCoachFunc      func(ctx context.Context, params CreateCoachParams) (Coach, error)
	listAICoachesFunc    func(ctx context.Context) ([]Coach, error)
	getTeamByCoachIDFunc func(ctx context.Context, id int64) (Team, error)
	createTeamFunc       func(ctx context.Context, name string, coachID int64, isDefault bool) (Team, error)
}

func (m *mockTeamCreator) CreateCoach(ctx context.Context, params CreateCoachParams) (Coach, error) {
	return m.createCoachFunc(ctx, params)
}

func (m *mockTeamCreator) ListAICoaches(ctx context.Context) ([]Coach, error) {
	return m.listAICoachesFunc(ctx)
}

func (m *mockTeamCreator) GetTeamByCoachID(ctx context.Context, id int64) (Team, error) {
	return m.getTeamByCoachIDFunc(ctx, id)
}

func (m *mockTeamCreator) CreateTeam(ctx context.Context, name string, coachID int64, isDefault bool) (Team, error) {
	return m.createTeamFunc(ctx, name, coachID, isDefault)
}

type mockStore struct {
	Store
	createCoachFunc  func(ctx context.Context, arg database.CreateCoachParams) (database.Coach, error)
	createTeamFunc   func(ctx context.Context, arg database.CreateTeamParams) (database.Team, error)
	createPlayerFunc func(ctx context.Context, arg database.CreatePlayerParams) (database.Player, error)
	getAICoachesFunc func(ctx context.Context) ([]database.Coach, error)
}

func (m *mockStore) CreateCoach(ctx context.Context, arg database.CreateCoachParams) (database.Coach, error) {
	return m.createCoachFunc(ctx, arg)
}

func (m *mockStore) CreateTeam(ctx context.Context, arg database.CreateTeamParams) (database.Team, error) {
	return m.createTeamFunc(ctx, arg)
}

func (m *mockStore) CreatePlayer(ctx context.Context, arg database.CreatePlayerParams) (database.Player, error) {
	return m.createPlayerFunc(ctx, arg)
}

func (m *mockStore) GetAICoaches(ctx context.Context) ([]database.Coach, error) {
	return m.getAICoachesFunc(ctx)
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

func TestTeamsService_GenerateTeams(t *testing.T) {
	group := odize.NewGroup(t, nil)

	var mStore *mockStore
	var mockPlaybooks *mockPlaybookCreator
	var svc *Service

	group.BeforeEach(func() {
		mStore = &mockStore{}
		mockPlaybooks = &mockPlaybookCreator{}
		svc = &Service{
			store:       mStore,
			playbookSvc: mockPlaybooks,
		}
	})

	err := group.
		Test("should successfully generate all teams", func(t *testing.T) {
			var coachCount, teamCount, playbookCount, playerCount int

			mStore.createCoachFunc = func(ctx context.Context, arg database.CreateCoachParams) (database.Coach, error) {
				coachCount++
				return database.Coach{ID: int64(coachCount), Name: arg.Name}, nil
			}
			mStore.createTeamFunc = func(ctx context.Context, arg database.CreateTeamParams) (database.Team, error) {
				teamCount++
				return database.Team{ID: int64(teamCount), Name: arg.Name}, nil
			}
			mStore.createPlayerFunc = func(ctx context.Context, arg database.CreatePlayerParams) (database.Player, error) {
				playerCount++
				return database.Player{ID: int64(playerCount), Name: arg.Name}, nil
			}
			mockPlaybooks.createPlaybookFunc = func(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error) {
				playbookCount++
				return playbooks.Playbook{ID: int64(playbookCount), Name: params.Name}, nil
			}

			err := svc.GenerateAITeams(t.Context())
			odize.AssertNoError(t, err)
			odize.AssertEqual(t, totalTeams, coachCount)
			odize.AssertEqual(t, totalTeams, teamCount)
			odize.AssertEqual(t, totalTeams, playbookCount)
			odize.AssertEqual(t, totalTeams*5, playerCount)
		}).
		Test("should return error when CreateCoach fails", func(t *testing.T) {
			expectedErr := errors.New("coach error")
			mStore.createCoachFunc = func(ctx context.Context, arg database.CreateCoachParams) (database.Coach, error) {
				return database.Coach{}, expectedErr
			}

			err := svc.GenerateAITeams(t.Context())
			odize.AssertError(t, err)
			odize.AssertTrue(t, errors.Is(err, expectedErr))
		}).
		Test("should return error when CreateTeam fails", func(t *testing.T) {
			expectedErr := errors.New("team error")
			mStore.createCoachFunc = func(ctx context.Context, arg database.CreateCoachParams) (database.Coach, error) {
				return database.Coach{ID: 1}, nil
			}
			mStore.createTeamFunc = func(ctx context.Context, arg database.CreateTeamParams) (database.Team, error) {
				return database.Team{}, expectedErr
			}

			err := svc.GenerateAITeams(t.Context())
			odize.AssertError(t, err)
			odize.AssertTrue(t, errors.Is(err, expectedErr))
		}).
		Test("should return error when CreatePlaybook fails", func(t *testing.T) {
			expectedErr := errors.New("playbook error")
			mStore.createCoachFunc = func(ctx context.Context, arg database.CreateCoachParams) (database.Coach, error) {
				return database.Coach{ID: 1}, nil
			}
			mStore.createTeamFunc = func(ctx context.Context, arg database.CreateTeamParams) (database.Team, error) {
				return database.Team{ID: 1, Name: arg.Name}, nil
			}
			mStore.createPlayerFunc = func(ctx context.Context, arg database.CreatePlayerParams) (database.Player, error) {
				return database.Player{}, nil
			}
			mockPlaybooks.createPlaybookFunc = func(ctx context.Context, params playbooks.PlaybookParams) (playbooks.Playbook, error) {
				return playbooks.Playbook{}, expectedErr
			}

			err := svc.GenerateAITeams(t.Context())
			odize.AssertError(t, err)
			odize.AssertTrue(t, errors.Is(err, expectedErr))
		}).
		Run()

	odize.AssertNoError(t, err)
}
