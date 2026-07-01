package uibattle

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

func TestPageBattleSelectionModel_Rendering(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("should load data and render selection", func(t *testing.T) {
		state := &uistate.GlobalState{
			Team: &teams.Team{ID: 1, Name: "My Team"},
		}
		theme := styles.NewIceTheme()

		m := NewModelBattleSelection(state, nil, nil, nil, theme)
		m.width = 100
		m.height = 40
		m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

		msg := MsgBattleSelectionDataLoaded{
			Playbooks: []playbooks.Playbook{{ID: 1, Name: "Playbook 1"}},
			AITeams: []teams.AITeam{
				{
					Coach: teams.Coach{Name: "Coach A"},
					Team:  teams.Team{Name: "Team A"},
				},
			},
		}
		m.Update(msg)

		odize.AssertEqual(t, 1, m.playbookList.Len())
		odize.AssertEqual(t, 1, m.aiTeamList.Len())

		view := m.View()
		content := view.Content

		odize.AssertTrue(t, strings.Contains(content, "NEW BATTLE"))
		odize.AssertTrue(t, strings.Contains(content, "Playbook 1"))
		odize.AssertTrue(t, strings.Contains(content, "Team A"))
	})

	group.Test("should handle state transitions", func(t *testing.T) {
		state := &uistate.GlobalState{
			Team: &teams.Team{ID: 1, Name: "My Team"},
		}
		theme := styles.NewIceTheme()
		m := NewModelBattleSelection(state, nil, nil, nil, theme)
		m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
		m.Update(MsgBattleSelectionDataLoaded{
			Playbooks: []playbooks.Playbook{{ID: 1, Name: "Playbook 1"}},
			AITeams: []teams.AITeam{
				{
					Team: teams.Team{ID: 2, Name: "Team A"},
				},
			},
		})

		// 1. Initial state: selecting playbook
		odize.AssertEqual(t, stateSelectingPlaybook, m.state)

		// 2. Select playbook -> selecting opponent
		m.Update(tea.KeyPressMsg{Text: "enter"})
		odize.AssertEqual(t, stateSelectingOpponent, m.state)
		odize.AssertTrue(t, m.selectedPlaybook != nil)

		// 3. Select opponent -> returns MsgSwitchBattlePage
		_, cmd := m.Update(tea.KeyPressMsg{Text: "enter"})
		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		switch v := msg.(type) {
		case MsgSwitchBattlePage:
			odize.AssertEqual(t, SubPageBattleConfirm, v.NewPage)
			odize.AssertEqual(t, m.selectedPlaybook, v.SelectedPlaybook)
			odize.AssertEqual(t, m.selectedAITeam, v.SelectedAITeam)
		default:
			t.Fatalf("expected MsgSwitchBattlePage, got %T", msg)
		}

		// 4. Back from selecting opponent -> selecting playbook
		m.Update(tea.KeyPressMsg{Text: "esc"})
		odize.AssertEqual(t, stateSelectingPlaybook, m.state)
	})

	group.Test("should handle back navigation to title", func(t *testing.T) {
		state := &uistate.GlobalState{}
		theme := styles.NewIceTheme()
		m := NewModelBattleSelection(state, nil, nil, nil, theme)

		_, cmd := m.Update(tea.KeyPressMsg{Text: "esc"})

		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		switch v := msg.(type) {
		case uistate.MsgSwitchPage:
			odize.AssertEqual(t, uistate.PageTitle, v.NewPage)
		default:
			t.Fatalf("expected MsgSwitchPage, got %T", msg)
		}
	})

	group.Test("should reset state on Init", func(t *testing.T) {
		state := &uistate.GlobalState{
			Team: &teams.Team{ID: 1, Name: "My Team"},
		}
		theme := styles.NewIceTheme()
		m := NewModelBattleSelection(state, nil, nil, nil, theme)

		// 1. Set some state
		m.state = stateSelectingOpponent
		m.selectedPlaybook = &playbooks.Playbook{ID: 1, Name: "Playbook 1"}
		m.selectedAITeam = &teams.AITeam{Team: teams.Team{ID: 2, Name: "Team A"}}

		// 2. Call Init
		m.Init()

		// 3. Verify it's reset
		odize.AssertEqual(t, stateSelectingPlaybook, m.state)
		odize.AssertTrue(t, m.selectedPlaybook == nil)
		odize.AssertTrue(t, m.selectedAITeam == nil)
	})

	group.Test("should call Init on MsgSwitchPage", func(t *testing.T) {
		state := &uistate.GlobalState{
			Team: &teams.Team{ID: 1, Name: "My Team"},
		}
		theme := styles.NewIceTheme()
		// We test BattleModel here as it now handles MsgSwitchPage
		m := NewBattleModel(state, nil, nil, nil, theme)

		// 1. Set some state to verify reset
		if p, ok := m.subPageBattleSelection.(*PageBattleSelectionModel); ok {
			p.state = stateSelectingOpponent
		}

		// 2. Send MsgSwitchPage
		_, cmd := m.Update(uistate.MsgSwitchPage{NewPage: uistate.PageNewBattleSelection})

		// 3. Verify it's reset and returns loadData command
		if p, ok := m.subPageBattleSelection.(*PageBattleSelectionModel); ok {
			odize.AssertEqual(t, stateSelectingPlaybook, p.state)
		}
		odize.AssertTrue(t, cmd != nil)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
