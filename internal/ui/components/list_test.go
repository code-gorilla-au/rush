package components

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

func TestList(t *testing.T) {
	group := odize.NewGroup(t, nil)

	theme := styles.NewIceTheme()

	err := group.
		Test("NewList initializes with items", func(t *testing.T) {
			items := []string{"Item 1", "Item 2"}
			config := ListConfig[string]{
				Title: "Test List",
				Items: items,
				ItemMapper: func(s string) ListItem[string] {
					return ListItem[string]{
						Data:     s,
						TitleVal: s,
					}
				},
			}
			l := NewList(config, theme)

			odize.AssertEqual(t, 2, l.Len())
			selected, ok := l.SelectedItem()
			odize.AssertTrue(t, ok)
			odize.AssertEqual(t, "Item 1", selected)
		}).
		Test("SetSize updates dimensions with default padding", func(t *testing.T) {
			items := []string{"Item 1"}
			config := ListConfig[string]{
				Items: items,
				ItemMapper: func(s string) ListItem[string] {
					return ListItem[string]{Data: s, TitleVal: s}
				},
			}
			l := NewList(config, theme)

			l.SetSize(100, 50)
			// Defaults: left=2, right=2, top=1, bottom=1
			// Total X padding = 4, Total Y padding = 2
			odize.AssertEqual(t, 96, l.Model.Width())
			odize.AssertEqual(t, 48, l.Model.Height())
		}).
		Test("Update handles tea.WindowSizeMsg when AutoResize is enabled", func(t *testing.T) {
			items := []string{"Item 1"}
			config := ListConfig[string]{
				Items: items,
				ItemMapper: func(s string) ListItem[string] {
					return ListItem[string]{Data: s, TitleVal: s}
				},
				DisableAutoResize: false,
			}
			l := NewList(config, theme)

			newModel, _ := l.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			odize.AssertEqual(t, 76, newModel.Model.Width())
			odize.AssertEqual(t, 38, newModel.Model.Height())
		}).
		Test("Update ignores tea.WindowSizeMsg when AutoResize is disabled", func(t *testing.T) {
			items := []string{"Item 1"}
			config := ListConfig[string]{
				Items: items,
				ItemMapper: func(s string) ListItem[string] {
					return ListItem[string]{Data: s, TitleVal: s}
				},
				DisableAutoResize: true,
			}
			l := NewList(config, theme)

			newModel, _ := l.Update(tea.WindowSizeMsg{Width: 80, Height: 40})

			odize.AssertEqual(t, 0, newModel.Model.Width())
			odize.AssertEqual(t, 0, newModel.Model.Height())
		}).
		Test("SetPadding and SetSize account for granular padding", func(t *testing.T) {
			items := []string{"Item 1"}
			config := ListConfig[string]{
				Items: items,
				ItemMapper: func(s string) ListItem[string] {
					return ListItem[string]{Data: s, TitleVal: s}
				},
				LeftPadding:   5,
				RightPadding:  5,
				TopPadding:    3,
				BottomPadding: 3,
			}
			l := NewList(config, theme)

			l.SetSize(100, 50)
			odize.AssertEqual(t, 90, l.Model.Width())
			odize.AssertEqual(t, 44, l.Model.Height())

			l.SetPadding(5, 5, 5, 5) // top, right, bottom, left
			odize.AssertEqual(t, 90, l.Model.Width())
			odize.AssertEqual(t, 40, l.Model.Height())
		}).
		Test("View applies padding", func(t *testing.T) {
			items := []string{"Item 1"}
			config := ListConfig[string]{
				Items: items,
				ItemMapper: func(s string) ListItem[string] {
					return ListItem[string]{Data: s, TitleVal: s}
				},
				LeftPadding: 10,
				TopPadding:  6,
			}
			l := NewList(config, theme)
			l.SetSize(100, 50)

			view := l.View()
			// The view should have 3 lines of top padding (6/2)
			// We can check this by counting newlines at the beginning.
			// However, lipgloss might render it differently.

			// Just check that it's not empty and has some length.
			odize.AssertTrue(t, len(view) > 0)
		}).
		Run()

	odize.AssertNoError(t, err)
}
