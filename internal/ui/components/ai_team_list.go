package components

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type AITeamItem struct {
	team teams.AITeam
}

func (i AITeamItem) Title() string       { return i.team.Team.Name }
func (i AITeamItem) Description() string { return i.team.Coach.Name }
func (i AITeamItem) FilterValue() string { return i.team.Team.Name }

type AITeamList struct {
	list   list.Model
	active bool
}

func NewAITeamList(items []teams.AITeam, theme styles.IceTheme) AITeamList {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = AITeamItem{team: item}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = theme.Base.PaddingLeft(3)
	delegate.Styles.NormalDesc = theme.Muted.PaddingLeft(3)
	delegate.Styles.SelectedTitle = theme.SelectedTitle
	delegate.Styles.SelectedDesc = theme.SelectedDesc
	delegate.Styles.DimmedTitle = theme.Muted.PaddingLeft(3)
	delegate.Styles.DimmedDesc = theme.Muted.PaddingLeft(3)

	l := list.New(listItems, delegate, 0, 0)
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.Styles.Title = theme.Title

	return AITeamList{
		list:   l,
		active: true,
	}
}

func (l *AITeamList) Update(msg tea.Msg) (AITeamList, tea.Cmd) {
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return *l, cmd
}

func (l *AITeamList) View(theme styles.IceTheme) string {
	return l.list.View()
}

func (l *AITeamList) SetActive(active bool) {
	l.active = active
}

func (l *AITeamList) SelectedItem() *teams.AITeam {
	if item, ok := l.list.SelectedItem().(AITeamItem); ok {
		return &item.team
	}
	return nil
}

func (l *AITeamList) SetSize(width, height int) {
	l.list.SetSize(width, height)
}

func (l *AITeamList) SetItems(items []teams.AITeam) tea.Cmd {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = AITeamItem{team: item}
	}
	return l.list.SetItems(listItems)
}

func (l *AITeamList) SetTitle(title string) {
	l.list.Title = title
}

func (l *AITeamList) Len() int {
	return len(l.list.Items())
}

func (l *AITeamList) Reset() {
	l.list.Select(0)
	l.list.FilterInput.Reset()
}

func (l *AITeamList) IsFiltering() bool {
	return l.list.FilterState() == list.Filtering
}
