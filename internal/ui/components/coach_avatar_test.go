package components

import (
	"strings"
	"testing"

	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

func TestCoachAvatar(t *testing.T) {
	group := odize.NewGroup(t, nil)

	coach := &teams.Coach{Name: "Ted Lasso"}
	team := &teams.Team{Name: "AFC Richmond"}
	theme := styles.NewIceTheme()

	err := group.
		Test("NewCoachAvatar should initialize with given coach and team", func(t *testing.T) {
			avatar := NewCoachAvatar(coach, team)
			odize.AssertEqual(t, coach, avatar.Coach)
			odize.AssertEqual(t, team, avatar.Team)
		}).
		Test("View should render team and coach names when both are present", func(t *testing.T) {
			avatar := NewCoachAvatar(coach, team)
			rendered := avatar.View(theme)

			odize.AssertTrue(t, strings.Contains(rendered, "AFC Richmond"))
			odize.AssertTrue(t, strings.Contains(rendered, "Coach: Ted Lasso"))
		}).
		Test("View should return empty string if coach or team is nil", func(t *testing.T) {
			avatarNoCoach := NewCoachAvatar(nil, team)
			odize.AssertEqual(t, "", avatarNoCoach.View(theme))

			avatarNoTeam := NewCoachAvatar(coach, nil)
			odize.AssertEqual(t, "", avatarNoTeam.View(theme))
		}).
		Run()

	odize.AssertNoError(t, err)
}
