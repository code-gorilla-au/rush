package components

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type SelectedFormationItem struct {
	formation playbooks.Formation
}

func (i SelectedFormationItem) Title() string {
	return fmt.Sprintf("%s (%d-%d-%d)", i.formation.Name, i.formation.Lane1, i.formation.Lane2, i.formation.Lane3)
}
func (i SelectedFormationItem) Description() string { return "" }
func (i SelectedFormationItem) FilterValue() string { return i.formation.Name }

type SelectedFormationList struct {
	list   list.Model
	active bool
}

func NewSelectedFormationList() SelectedFormationList {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#A5F2F3")).BorderForeground(lipgloss.Color("#A5F2F3"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#87CEEB")).BorderForeground(lipgloss.Color("#A5F2F3"))

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Selected Formations (Max 10)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#A5F2F3")).Bold(true)

	return SelectedFormationList{
		list: l,
	}
}

func (l *SelectedFormationList) Update(msg tea.Msg) (SelectedFormationList, tea.Cmd) {
	if !l.active {
		return *l, nil
	}
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return *l, cmd
}

func (l *SelectedFormationList) View() string {
	style := lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder())
	if l.active {
		style = style.BorderForeground(lipgloss.Color("#A5F2F3"))
	}
	return style.Render(l.list.View())
}

func (l *SelectedFormationList) SetItems(formations []playbooks.Formation) tea.Cmd {
	items := make([]list.Item, len(formations))
	for i, f := range formations {
		items[i] = SelectedFormationItem{formation: f}
	}
	return l.list.SetItems(items)
}

func (l *SelectedFormationList) SelectedIndex() int {
	return l.list.Index()
}

func (l *SelectedFormationList) SetSize(width, height int) {
	l.list.SetSize(width, height)
}

func (l *SelectedFormationList) SetActive(active bool) {
	l.active = active
}

func (l *SelectedFormationList) Len() int {
	return len(l.list.Items())
}
