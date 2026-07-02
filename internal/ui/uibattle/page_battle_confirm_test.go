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

func TestPageBattleConfirmModel(t *testing.T) {
	t.Parallel()
	group := odize.NewGroup(t, nil)

	group.Test("should render confirmation details", func(t *testing.T) {
		state := &uistate.GlobalState{
			Coach: &teams.Coach{Name: "Coach Me"},
			Team: &teams.Team{ID: 1, Name: "My Team", Players: []teams.Player{
				{Name: "Alice"},
				{Name: "Bob"},
			}},
		}
		theme := styles.NewIceTheme()
		m := NewPageBattleConfirm(state, nil, theme)
		m.width = 100
		m.height = 40

		m.SetData(
			&playbooks.Playbook{Name: "Playbook 1"},
			&teams.AITeam{
				Coach:    teams.Coach{Name: "Coach Rival"},
				Playbook: playbooks.Playbook{Name: "Counter Playbook"},
				Team: teams.Team{Name: "Opponent A", Players: []teams.Player{
					{Name: "Eve"},
					{Name: "Mallory"},
				}},
			},
		)

		view := m.View()
		content := view.Content

		odize.AssertTrue(t, strings.Contains(content, "CONFIRM BATTLE"))
		odize.AssertTrue(t, strings.Contains(content, "VS"))
		odize.AssertTrue(t, strings.Contains(content, "My Team"))
		odize.AssertTrue(t, strings.Contains(content, "Coach Me"))
		odize.AssertTrue(t, strings.Contains(content, "Playbook 1"))
		odize.AssertTrue(t, strings.Contains(content, "Opponent A"))
		odize.AssertTrue(t, strings.Contains(content, "Coach Rival"))
		odize.AssertTrue(t, strings.Contains(content, "Counter Playbook"))
		odize.AssertTrue(t, strings.Contains(content, "Alice"))
		odize.AssertTrue(t, strings.Contains(content, "Eve"))
	})

	group.Test("should render fallback values when data is missing", func(t *testing.T) {
		state := &uistate.GlobalState{}
		theme := styles.NewIceTheme()
		m := NewPageBattleConfirm(state, nil, theme)
		m.width = 100
		m.height = 40

		view := m.View()
		content := view.Content

		odize.AssertTrue(t, strings.Contains(content, "Unknown Team"))
		odize.AssertTrue(t, strings.Contains(content, "No players"))
	})

	group.Test("should handle back navigation", func(t *testing.T) {
		state := &uistate.GlobalState{}
		theme := styles.NewIceTheme()
		m := NewPageBattleConfirm(state, nil, theme)

		_, cmd := m.Update(tea.KeyPressMsg{Text: "esc"})

		odize.AssertTrue(t, cmd != nil)
		msg := cmd()
		switch v := msg.(type) {
		case MsgSwitchBattlePage:
			odize.AssertEqual(t, SubPageBattleSelection, v.NewPage)
		default:
			t.Fatalf("expected MsgSwitchBattlePage, got %T", msg)
		}
	})

	err := group.Run()
	odize.AssertNoError(t, err)
}
