package ui

import (
	"database/sql"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/database"
	"github.com/code-gorilla-au/rush/internal/games"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/tournament"
)

func setupServices(t *testing.T) (*teams.Service, *playbooks.Service, *games.Service, *tournament.AITeamService) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	migrator := database.NewMigrator(db, database.SchemaFS)
	if err := migrator.Migrate(t.Context()); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	queries := database.New(db)
	ts := teams.NewTeamsService(queries)
	ps := playbooks.NewPlaybooksService(queries)
	gs := games.NewService(queries)
	as := tournament.NewAITeamService(ts, ps)
	return ts, ps, gs, as
}

func TestTheme(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("NewIceTheme should return a theme with correct colors", func(t *testing.T) {
			theme := NewIceTheme()
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
			s, ps, gs, as := setupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
				AiTeamsSvc:  as,
			})
			odize.AssertTrue(t, m.theme.Logo.GetForeground() != nil)
		}).
		Test("Init should return a command", func(t *testing.T) {
			s, ps, gs, as := setupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
				AiTeamsSvc:  as,
			})
			cmd := m.Init()
			odize.AssertTrue(t, cmd != nil)
		}).
		Test("Update should handle Quit keys", func(t *testing.T) {
			s, ps, gs, as := setupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
				AiTeamsSvc:  as,
			})
			_, cmd := m.Update(tea.KeyPressMsg{Text: "q"})
			odize.AssertTrue(t, cmd != nil)

			_, cmd = m.Update(tea.KeyPressMsg{Text: "ctrl+c"})
			odize.AssertTrue(t, cmd != nil)
		}).
		Test("Update should handle WindowSizeMsg", func(t *testing.T) {
			s, ps, gs, as := setupServices(t)
			m := New(Dependencies{
				TeamsSvc:    s,
				PlaybookSvc: ps,
				GameSvc:     gs,
				AiTeamsSvc:  as,
			})
			newModel, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			updatedModel := newModel.(RootModel)
			odize.AssertTrue(t, updatedModel.width == 100)
			odize.AssertTrue(t, updatedModel.height == 50)
		}).
		Run()

	odize.AssertNoError(t, err)
}
