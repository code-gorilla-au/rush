package components

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type PlaybookItem struct {
	playbook playbooks.Playbook
}

func (i PlaybookItem) Title() string       { return i.playbook.Name }
func (i PlaybookItem) Description() string { return i.playbook.Description }
func (i PlaybookItem) FilterValue() string { return i.playbook.Name }

type PlaybookList struct {
	list   list.Model
	active bool
}

func NewPlaybookList(items []playbooks.Playbook) PlaybookList {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = PlaybookItem{playbook: item}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#A5F2F3")).BorderForeground(lipgloss.Color("#A5F2F3"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#87CEEB")).BorderForeground(lipgloss.Color("#A5F2F3"))

	l := list.New(listItems, delegate, 0, 0)
	l.Title = "Playbooks"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#A5F2F3")).Bold(true)

	return PlaybookList{
		list:   l,
		active: true,
	}
}

func (l *PlaybookList) Update(msg tea.Msg) (PlaybookList, tea.Cmd) {
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return *l, cmd
}

func (l *PlaybookList) View() string {
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

func (l *PlaybookList) SetActive(active bool) {
	l.active = active
}

func (l *PlaybookList) SelectedItem() *playbooks.Playbook {
	if item, ok := l.list.SelectedItem().(PlaybookItem); ok {
		return &item.playbook
	}
	return nil
}

func (l *PlaybookList) SetSize(width, height int) {
	l.list.SetSize(width, height)
}

func (l *PlaybookList) SetItems(items []playbooks.Playbook) tea.Cmd {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = PlaybookItem{playbook: item}
	}
	return l.list.SetItems(listItems)
}

func (l *PlaybookList) SetTitle(title string) {
	l.list.Title = title
}

func (l *PlaybookList) Len() int {
	return len(l.list.Items())
}

func (l *PlaybookList) Reset() {
	l.list.Select(0)
	l.list.FilterInput.Reset()
}

func (l *PlaybookList) IsFiltering() bool {
	return l.list.FilterState() == list.Filtering
}
