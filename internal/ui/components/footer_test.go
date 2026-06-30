package components

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
	"github.com/code-gorilla-au/rush/internal/ui/uistate"
)

func TestFooter(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("NewFooter should initialize with given KeyMap", func(t *testing.T) {
			keys := uistate.NewKeyMap()
			footer := NewFooter(keys)
			odize.AssertEqual(t, keys, footer.KeyMap)
		}).
		Test("Update should set width from WindowSizeMsg", func(t *testing.T) {
			footer := NewFooter(uistate.NewKeyMap())
			msg := tea.WindowSizeMsg{Width: 100, Height: 50}
			footer.Update(msg)
			odize.AssertEqual(t, 100, footer.Width)
		}).
		Test("View should render help text", func(t *testing.T) {
			footer := NewFooter(uistate.NewKeyMap())
			theme := styles.NewIceTheme()

			rendered := footer.View(theme)

			odize.AssertTrue(t, len(rendered) > 0)
			odize.AssertTrue(t, strings.Contains(rendered, "q"))
			odize.AssertTrue(t, strings.Contains(rendered, "quit"))
		}).
		Run()

	odize.AssertNoError(t, err)
}

func TestKeyMap(t *testing.T) {
	group := odize.NewGroup(t, nil)

	err := group.
		Test("NewKeyMap should have Quit binding", func(t *testing.T) {
			keys := uistate.NewKeyMap()
			odize.AssertTrue(t, keys.Quit.Enabled())
			odize.AssertEqual(t, "q", keys.Quit.Keys()[0])
		}).
		Test("ShortHelp should contain Quit", func(t *testing.T) {
			keys := uistate.NewKeyMap()
			shortHelp := keys.ShortHelp()
			odize.AssertEqual(t, 1, len(shortHelp))
			odize.AssertEqual(t, keys.Quit, shortHelp[0])
		}).
		Test("FullHelp should contain Quit in the first row", func(t *testing.T) {
			keys := uistate.NewKeyMap()
			fullHelp := keys.FullHelp()
			odize.AssertEqual(t, 1, len(fullHelp))
			odize.AssertEqual(t, 1, len(fullHelp[0]))
			odize.AssertEqual(t, keys.Quit, fullHelp[0][0])
		}).
		Run()

	odize.AssertNoError(t, err)
}
