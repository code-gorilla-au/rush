package components

import (
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type TeamTile struct {
	TeamName     string
	CoachName    string
	PlaybookName string
	PlayerNames  []string
}

func NewTeamTile(teamName, coachName, playbookName string, playerNames []string) TeamTile {
	return TeamTile{
		TeamName:     teamName,
		CoachName:    coachName,
		PlaybookName: playbookName,
		PlayerNames:  slices.Clone(playerNames),
	}
}

func (t TeamTile) View(theme styles.IceTheme, width int) string {
	tileBody := lipgloss.JoinVertical(lipgloss.Left,
		theme.SecondaryHeader.Render(valueOrDefault(t.TeamName, "Unknown Team")),
		renderField(theme, "Coach", t.CoachName),
		renderField(theme, "Playbook", t.PlaybookName),
		"",
		theme.Muted.Render("Players"),
		renderPlayers(theme, t.PlayerNames),
	)

	tileStyle := theme.RoundBorder.Copy().Align(lipgloss.Left)

	if width > 0 {
		tileStyle = tileStyle.Width(width)
	}

	return tileStyle.Render(tileBody)
}

func renderField(theme styles.IceTheme, label, value string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		theme.Muted.Render(label+": "),
		theme.Base.Render(valueOrDefault(value, "Unknown")),
	)
}

func renderPlayers(theme styles.IceTheme, playerNames []string) string {
	players := make([]string, 0, len(playerNames))
	for _, playerName := range playerNames {
		name := strings.TrimSpace(playerName)
		if name == "" {
			continue
		}

		players = append(players, theme.Player.Render("• "+name))
	}

	if len(players) == 0 {
		return theme.Muted.Render("No players")
	}

	return lipgloss.JoinVertical(lipgloss.Left, players...)
}

func valueOrDefault(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}

	return trimmed
}
