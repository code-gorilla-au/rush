package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
	"github.com/code-gorilla-au/rush/internal/ui/uitest"
)

func TestTheme(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("NewIceTheme should return a theme with correct colors", func(t *testing.T) {
			theme := styles.NewIceTheme()
			// We can't easily check the color values from the Style object in Lipgloss v2
			// without deep inspection, but we can check if they are not empty.
			odize.AssertTrue(t, theme.Logo.GetForeground() != nil)
			odize.AssertTrue(t, theme.Footer.GetForeground() != nil)
			odize.AssertTrue(t, theme.Base.GetBackground() != nil)
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestNew(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("New should initialize model with IceTheme", func(t *testing.T) {
			s, ps, gs := uitest.SetupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
			})
			odize.AssertTrue(t, m.theme.Logo.GetForeground() != nil)
		}).
		Test("Init should return a command", func(t *testing.T) {
			s, ps, gs := uitest.SetupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
			})
			cmd := m.Init()
			odize.AssertTrue(t, cmd != nil)
		}).
		Test("Update should handle Quit keys", func(t *testing.T) {
			s, ps, gs := uitest.SetupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
			})
			_, cmd := m.Update(tea.KeyPressMsg{Text: "q"})
			odize.AssertTrue(t, cmd != nil)

			_, cmd = m.Update(tea.KeyPressMsg{Text: "ctrl+c"})
			odize.AssertTrue(t, cmd != nil)
		}).
		Test("Update should handle WindowSizeMsg", func(t *testing.T) {
			s, ps, gs := uitest.SetupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
			})
			newModel, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			updatedModel, ok := newModel.(*RootModel)
			odize.AssertTrue(t, ok)
			odize.AssertTrue(t, updatedModel.width == 100)
			odize.AssertTrue(t, updatedModel.height == 50)
		}).
		Test("Update should route from create coach to locker room when state is updated", func(t *testing.T) {
			s, ps, gs := uitest.SetupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
			})

			_, _ = m.Update(uistate.MsgSwitchPage{NewPage: uistate.PageCreateCoach})

			coach := teams.Coach{ID: 1}
			team := teams.Team{ID: 1}

			_, cmd := m.Update(uistate.MsgStateUpdated{Coach: &coach, Team: &team})
			odize.AssertTrue(t, cmd != nil)

			routedMsg, ok := cmd().(uistate.MsgSwitchPage)
			odize.AssertTrue(t, ok)
			odize.AssertEqual(t, uistate.PageLockerRoom, routedMsg.NewPage)
		}).
		Run()

	odize.AssertNoError(t, err)
}
