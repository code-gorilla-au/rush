package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

func TestModelNewBattleSelection_Rendering(t *testing.T) {
	group := odize.NewGroup(t, nil)

	group.Test("should load data and render selection", func(t *testing.T) {
		state := &GlobalState{
			Team: &teams.Team{ID: 1, Name: "My Team"},
		}
		theme := styles.NewIceTheme()

		m := NewModelNewBattleSelection(state, nil, nil, nil, theme)
		m.width = 100
		m.height = 40
		m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

		msg := msgDataLoaded{
			playbooks: []playbooks.Playbook{{ID: 1, Name: "Playbook 1"}},
			aiTeams: []teams.AITeam{
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
		state := &GlobalState{
			Team: &teams.Team{ID: 1, Name: "My Team"},
		}
		theme := styles.NewIceTheme()
		m := NewModelNewBattleSelection(state, nil, nil, nil, theme)
		m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
		m.Update(msgDataLoaded{
			playbooks: []playbooks.Playbook{{ID: 1, Name: "Playbook 1"}},
			aiTeams: []teams.AITeam{
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

		// 3. Select opponent -> confirming
		m.Update(tea.KeyPressMsg{Text: "enter"})
		odize.AssertEqual(t, stateConfirming, m.state)
		odize.AssertTrue(t, m.selectedAITeam != nil)

		// 4. Back from confirming -> selecting opponent
		m.Update(tea.KeyPressMsg{Text: "esc"})
		odize.AssertEqual(t, stateSelectingOpponent, m.state)

		// 5. Back from selecting opponent -> selecting playbook
		m.Update(tea.KeyPressMsg{Text: "esc"})
		odize.AssertEqual(t, stateSelectingPlaybook, m.state)
	})

	group.Test("should handle back navigation to title", func(t *testing.T) {
		state := &GlobalState{}
		theme := styles.NewIceTheme()
		m := NewModelNewBattleSelection(state, nil, nil, nil, theme)

		_, cmd := m.Update(tea.KeyPressMsg{Text: "esc"})

		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		switch v := msg.(type) {
		case MsgSwitchPage:
			odize.AssertEqual(t, PageTitle, v.NewPage)
		default:
			t.Fatalf("expected MsgSwitchPage, got %T", msg)
		}
	})

	group.Test("should reset state on Init", func(t *testing.T) {
		state := &GlobalState{
			Team: &teams.Team{ID: 1, Name: "My Team"},
		}
		theme := styles.NewIceTheme()
		m := NewModelNewBattleSelection(state, nil, nil, nil, theme)

		// 1. Set some state
		m.state = stateConfirming
		m.selectedPlaybook = &playbooks.Playbook{ID: 1, Name: "Playbook 1"}
		m.selectedAITeam = &teams.AITeam{Team: teams.Team{ID: 2, Name: "Team A"}}

		// 2. Call Init
		m.Init()

		// 3. Verify it's reset
		odize.AssertEqual(t, stateSelectingPlaybook, m.state)
		odize.AssertTrue(t, m.selectedPlaybook == nil)
		odize.AssertTrue(t, m.selectedAITeam == nil)
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
