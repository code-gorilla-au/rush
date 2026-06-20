package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
)

type PlaybookList struct {
	cursor int
	items  []playbooks.Playbook
}

func NewPlaybookList(items []playbooks.Playbook) PlaybookList {
	return PlaybookList{
		items: items,
	}
}

func (l *PlaybookList) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if l.cursor > 0 {
				l.cursor--
			}
		case "down", "j":
			if l.cursor < len(l.items)-1 {
				l.cursor++
			}
		}
	}
}

func (l *PlaybookList) View(itemStyle lipgloss.Style, selectedStyle lipgloss.Style) string {
	if len(l.items) == 0 {
		return "No playbooks found."
	}

	var s string
	for i, item := range l.items {
		content := item.Name
		if i == l.cursor {
			s += selectedStyle.Render("> " + content)
		} else {
			s += itemStyle.Render("  " + content)
		}
		if i < len(l.items)-1 {
			s += "\n"
		}
	}
	return s
}

func (l *PlaybookList) SelectedItem() *playbooks.Playbook {
	if len(l.items) == 0 {
		return nil
	}
	return &l.items[l.cursor]
}

func (l *PlaybookList) Len() int {
	return len(l.items)
}
