package components

import (
	"strings"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

func TestTeamTile(t *testing.T) {
	t.Parallel()
	group := odize.NewGroup(t, nil)

	err := group.
		Test("View should render team, coach, playbook, and players", func(t *testing.T) {
			tile := NewTeamTile("Frost Wolves", "Coach Ice", "North Formation", []string{"Alice", "Bob"})

			rendered := tile.View(styles.NewIceTheme(), 40)

			odize.AssertTrue(t, strings.Contains(rendered, "Frost Wolves"))
			odize.AssertTrue(t, strings.Contains(rendered, "Coach Ice"))
			odize.AssertTrue(t, strings.Contains(rendered, "North Formation"))
			odize.AssertTrue(t, strings.Contains(rendered, "Alice"))
			odize.AssertTrue(t, strings.Contains(rendered, "Bob"))
		}).
		Test("View should show fallback values when data is empty", func(t *testing.T) {
			tile := NewTeamTile("", "", "", []string{"", "   "})

			rendered := tile.View(styles.NewIceTheme(), 30)

			odize.AssertTrue(t, strings.Contains(rendered, "Unknown Team"))
			odize.AssertTrue(t, strings.Contains(rendered, "Coach"))
			odize.AssertTrue(t, strings.Contains(rendered, "Playbook"))
			odize.AssertTrue(t, strings.Count(rendered, "Unknown") >= 2)
			odize.AssertTrue(t, strings.Contains(rendered, "No players"))
		}).
		Run()

	odize.AssertNoError(t, err)
}
