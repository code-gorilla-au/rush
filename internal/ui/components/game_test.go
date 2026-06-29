package components

import (
	"context"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type mockStore struct {
	games.Store
}

func (m *mockStore) CreateGame(ctx context.Context, arg database.CreateGameParams) (database.Game, error) {
	return database.Game{
		ID:           1,
		Name:         arg.Name,
		TeamA:        arg.TeamA,
		TeamB:        arg.TeamB,
		Status:       "pending",
		Rounds:       arg.Rounds,
		CurrentRound: arg.CurrentRound,
		ResultsLog:   arg.ResultsLog,
	}, nil
}

func TestGameComponent(t *testing.T) {
	group := odize.NewGroup(t, nil)

	teamA := games.TeamConfig{TeamID: 1, TeamName: "Team A", Formations: make([]playbooks.Formation, 10)}
	teamB := games.TeamConfig{TeamID: 2, TeamName: "Team B", Formations: make([]playbooks.Formation, 10)}

	for i := 0; i < 10; i++ {
		teamA.Formations[i] = playbooks.Formation{Lane1: 1, Lane2: 1, Lane3: 1}
		teamB.Formations[i] = playbooks.Formation{Lane1: 1, Lane2: 1, Lane3: 1}
	}

	rolls := []int{6, 1}
	idx := 0
	rollFn := func() int {
		val := rolls[idx%2]
		idx++
		return val
	}

	group.Test("should resolve round after tick", func(t *testing.T) {
		theme := styles.NewIceTheme()
		svc := games.NewService(&mockStore{})
		game, err := svc.NewGame(context.Background(), games.NewGameParams{
			TeamA: teamA,
			TeamB: teamB,
		})
		odize.AssertNoError(t, err)

		gComp := NewGame(&game, "Team A", "Team B", rollFn)

		// Initial state
		odize.AssertFalse(t, gComp.resolved)
		view := gComp.View(theme)
		odize.AssertTrue(t, strings.Contains(view, "ROUND 1"))
		odize.AssertTrue(t, strings.Contains(view, "Resolving..."))

		// Handle MsgResolveRound
		cmd := gComp.Update(MsgResolveRound{})
		odize.AssertNil(t, cmd)
		odize.AssertTrue(t, gComp.resolved)

		// Resolved state
		view = gComp.View(theme)
		odize.AssertTrue(t, strings.Contains(view, "WINNER: Team A"))
		odize.AssertTrue(t, strings.Contains(view, "Press Enter for next round..."))

		// Handle Enter key
		cmd = gComp.Update(tea.KeyPressMsg{Text: "enter"})
		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		_, ok := msg.(MsgNextRound)
		odize.AssertTrue(t, ok)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
