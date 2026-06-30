package ui

import (
	"context"
	"encoding/json"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/ui/components"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

type pageGameMockStore struct {
	games.Store
	updateCalled bool
}

func (m *pageGameMockStore) GetGameByID(ctx context.Context, id int64) (database.Game, error) {
	rounds := [10]games.Round{}
	roundsData, _ := json.Marshal(rounds)

	return database.Game{
		ID:           id,
		Name:         "Team A VS Team B",
		Status:       "pending",
		Rounds:       roundsData,
		ResultsLog:   []byte(`[]`),
		CurrentRound: 0,
	}, nil
}

func (m *pageGameMockStore) UpdateGame(ctx context.Context, arg database.UpdateGameParams) (database.Game, error) {
	m.updateCalled = true
	return database.Game{
		ID:           arg.ID,
		Name:         arg.Name,
		Status:       arg.Status,
		Rounds:       arg.Rounds,
		CurrentRound: arg.CurrentRound,
		ResultsLog:   arg.ResultsLog,
	}, nil
}

func TestPageGameModel(t *testing.T) {
	group := odize.NewGroup(t, nil)

	state := &uistate.GlobalState{}
	store := &pageGameMockStore{}
	gameSvc := games.NewService(store)

	group.Test("should handle MsgGameLoaded and initialize gameComp", func(t *testing.T) {
		theme := styles.NewIceTheme()
		m := NewModelGame(state, gameSvc, theme)
		m.SetGameID(1)

		game, _ := gameSvc.GetGame(context.Background(), 1)
		_, cmd := m.Update(MsgGameLoaded{Game: game})

		odize.AssertTrue(t, m.game != nil)
		odize.AssertTrue(t, cmd != nil)

		msg := cmd()
		_, ok := msg.(components.MsgResolveRound)
		odize.AssertTrue(t, ok)
	})

	group.Test("should persist game on MsgResolveRound", func(t *testing.T) {
		theme := styles.NewIceTheme()
		m := NewModelGame(state, gameSvc, theme)
		m.SetGameID(1)
		game, _ := gameSvc.GetGame(context.Background(), 1)
		m.Update(MsgGameLoaded{Game: game})

		store.updateCalled = false
		_, cmd := m.Update(components.MsgResolveRound{})

		odize.AssertTrue(t, cmd != nil)

		// Execute persistence command
		// If it's a batch, we might need to handle it.
		// But in Bubble Tea v2, Batch returns a Cmd that returns a BatchMsg.
		msg := cmd()
		if batchMsg, ok := msg.(tea.BatchMsg); ok {
			for _, c := range batchMsg {
				if c != nil {
					c()
				}
			}
		}

		odize.AssertTrue(t, store.updateCalled)
	})

	group.Test("should handle MsgNextRound and reset gameComp", func(t *testing.T) {
		theme := styles.NewIceTheme()
		m := NewModelGame(state, gameSvc, theme)
		m.SetGameID(1)
		game, _ := gameSvc.GetGame(context.Background(), 1)
		m.Update(MsgGameLoaded{Game: game})

		// Simulate round resolved (currentRound incremented)
		m.game.ResolveRound(func() int { return 1 })

		_, cmd := m.Update(components.MsgNextRound{})

		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		_, ok := msg.(components.MsgResolveRound)
		odize.AssertTrue(t, ok)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
