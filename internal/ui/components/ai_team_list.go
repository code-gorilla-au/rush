package components

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
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

func NewAITeamList(items []teams.AITeam) AITeamList {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = AITeamItem{team: item}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#A5F2F3")).BorderForeground(lipgloss.Color("#A5F2F3"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#87CEEB")).BorderForeground(lipgloss.Color("#A5F2F3"))

	l := list.New(listItems, delegate, 0, 0)
	l.Title = "AI Teams"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#A5F2F3")).Bold(true)

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

func (l *AITeamList) View() string {
	activeStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#A5F2F3")).
		Padding(1)

	inactiveStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#333333")).
		Padding(1).
		Foreground(lipgloss.Color("#666666"))

	if l.active {
		return activeStyle.Render(l.list.View())
	}
	return inactiveStyle.Render(l.list.View())
}

func (l *AITeamList) SetActive(active bool) {
	l.active = active
}

func (l *AITeamList) SelectedAITeam() *teams.AITeam {
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

func (l *AITeamList) IsFiltering() bool {
	return l.list.FilterState() == list.Filtering
}
