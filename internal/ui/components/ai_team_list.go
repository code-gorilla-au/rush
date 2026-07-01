package components

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type AITeamList struct {
	List[teams.AITeam]
}

func NewAITeamList(items []teams.AITeam, theme styles.IceTheme) AITeamList {
	return AITeamList{
		List: NewList(ListConfig[teams.AITeam]{
			Items: items,
			ItemMapper: func(item teams.AITeam) ListItem[teams.AITeam] {
				return ListItem[teams.AITeam]{
					Data:      item,
					TitleVal:  item.Team.Name,
					DescVal:   item.Coach.Name,
					FilterVal: item.Team.Name,
				}
			},
			EnableFiltering:   true,
			DisableAutoResize: true,
		}, theme),
	}
}

func (l *AITeamList) Update(msg tea.Msg) (AITeamList, tea.Cmd) {
	var cmd tea.Cmd
	l.List, cmd = l.List.Update(msg)
	return *l, cmd
}

func (l *AITeamList) View(theme styles.IceTheme) string {
	return l.List.View()
}

func (l *AITeamList) SelectedItem() *teams.AITeam {
	if item, ok := l.List.SelectedItem(); ok {
		return &item
	}
	return nil
}

func (l *AITeamList) SetItems(items []teams.AITeam) tea.Cmd {
	return l.List.SetItems(items)
}

func (l *AITeamList) SetTitle(title string) {
	l.Model.Title = title
}
