package components

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type FormationItem struct {
	formation       playbooks.Formation
	showDescription bool
}

func (i FormationItem) Title() string {
	return fmt.Sprintf("%s (%d-%d-%d)", i.formation.Name, i.formation.Lane1, i.formation.Lane2, i.formation.Lane3)
}
func (i FormationItem) Description() string {
	if i.showDescription {
		return i.formation.Description
	}
	return ""
}
func (i FormationItem) FilterValue() string { return i.formation.Name }

type FormationList struct {
	list   list.Model
	active bool
}

type FormationListConfig struct {
	Title           string
	Items           []playbooks.Formation
	EnableFiltering bool
	ShowDescription bool
}

func NewFormationList(config FormationListConfig) FormationList {
	items := make([]list.Item, len(config.Items))
	for i, f := range config.Items {
		items[i] = FormationItem{
			formation:       f,
			showDescription: config.ShowDescription,
		}
	}

	delegate := list.NewDefaultDelegate()
	// Match PlaybookList styling
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#A5F2F3")).BorderForeground(lipgloss.Color("#A5F2F3"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#87CEEB")).BorderForeground(lipgloss.Color("#A5F2F3"))

	l := list.New(items, delegate, 0, 0)
	l.Title = config.Title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(config.EnableFiltering)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#A5F2F3")).Bold(true)

	return FormationList{
		list: l,
	}
}

func (l *FormationList) Update(msg tea.Msg) (FormationList, tea.Cmd) {
	if !l.active {
		return *l, nil
	}
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return *l, cmd
}

func (l *FormationList) View() string {
	style := lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder())
	if l.active {
		style = style.BorderForeground(lipgloss.Color("#A5F2F3"))
	}
	return style.Render(l.list.View())
}

func (l *FormationList) SelectedItem() playbooks.Formation {
	if item, ok := l.list.SelectedItem().(FormationItem); ok {
		return item.formation
	}
	return playbooks.Formation{}
}

func (l *FormationList) SetSize(width, height int) {
	l.list.SetSize(width, height)
}

func (l *FormationList) SetActive(active bool) {
	l.active = active
}

func (l *FormationList) SetItems(formations []playbooks.Formation) tea.Cmd {
	items := make([]list.Item, len(formations))

	showDescription := false
	if len(l.list.Items()) > 0 {
		if first, ok := l.list.Items()[0].(FormationItem); ok {
			showDescription = first.showDescription
		}
	}

	for i, f := range formations {
		items[i] = FormationItem{
			formation:       f,
			showDescription: showDescription,
		}
	}
	return l.list.SetItems(items)
}

func (l *FormationList) SelectedIndex() int {
	return l.list.Index()
}

func (l *FormationList) Len() int {
	return len(l.list.Items())
}

func (l *FormationList) IsFiltering() bool {
	return l.list.FilterState() == list.Filtering
}
