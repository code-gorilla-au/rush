package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

func TestModelRules(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("should route to title page when q is pressed", func(t *testing.T) {
			theme := styles.NewIceTheme()
			m := NewModelRules(&uistate.GlobalState{}, theme)
			m.width = 100
			m.height = 50

			// We need to trigger WindowSizeMsg first to initialize viewport
			m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

			_, cmd := m.Update(tea.KeyPressMsg{Text: "esc"})
			odize.AssertTrue(t, cmd != nil)

			msg := cmd()
			switch switchMsg := msg.(type) {
			case uistate.MsgSwitchPage:
				odize.AssertEqual(t, uistate.PageTitle, switchMsg.NewPage)
			default:
				t.Fatalf("expected MsgSwitchPage, got %T", msg)
			}
		}).
		Run()

	odize.AssertNoError(t, err)
}
