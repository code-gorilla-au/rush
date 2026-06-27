package components

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/odize"
	"github.com/code-gorilla-au/rush/internal/teams"
)

func TestCoachAvatar(t *testing.T) {
	group := odize.NewGroup(t, nil)

	coach := &teams.Coach{Name: "Ted Lasso"}
	team := &teams.Team{Name: "AFC Richmond"}
	teamStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF"))
	coachStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))

	err := group.
		Test("NewCoachAvatar should initialize with given coach and team", func(t *testing.T) {
			avatar := NewCoachAvatar(coach, team)
			odize.AssertEqual(t, coach, avatar.Coach)
			odize.AssertEqual(t, team, avatar.Team)
		}).
		Test("View should render team and coach names when both are present", func(t *testing.T) {
			avatar := NewCoachAvatar(coach, team)
			rendered := avatar.View(teamStyle, coachStyle)

			odize.AssertTrue(t, strings.Contains(rendered, "AFC Richmond"))
			odize.AssertTrue(t, strings.Contains(rendered, "Coach: Ted Lasso"))
		}).
		Test("View should return empty string if coach or team is nil", func(t *testing.T) {
			avatarNoCoach := NewCoachAvatar(nil, team)
			odize.AssertEqual(t, "", avatarNoCoach.View(teamStyle, coachStyle))

			avatarNoTeam := NewCoachAvatar(coach, nil)
			odize.AssertEqual(t, "", avatarNoTeam.View(teamStyle, coachStyle))
		}).
		Run()

	odize.AssertNoError(t, err)
}
