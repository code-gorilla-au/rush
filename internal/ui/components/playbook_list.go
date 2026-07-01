package components

import (
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type PlaybookList struct {
	List[playbooks.Playbook]
}

func NewPlaybookList(items []playbooks.Playbook, theme styles.IceTheme) PlaybookList {
	return PlaybookList{
		List: NewList(ListConfig[playbooks.Playbook]{
			Items: items,
			ItemMapper: func(item playbooks.Playbook) ListItem[playbooks.Playbook] {
				return ListItem[playbooks.Playbook]{
					Data:      item,
					TitleVal:  item.Name,
					DescVal:   item.Description,
					FilterVal: item.Name,
				}
			},
			EnableFiltering:   true,
			DisableAutoResize: true,
		}, theme),
	}
}

func (l *PlaybookList) Update(msg tea.Msg) (PlaybookList, tea.Cmd) {
	var cmd tea.Cmd
	l.List, cmd = l.List.Update(msg)
	return *l, cmd
}

func (l *PlaybookList) View(theme styles.IceTheme) string {
	return l.List.View()
}

func (l *PlaybookList) SelectedItem() *playbooks.Playbook {
	if item, ok := l.List.SelectedItem(); ok {
		return &item
	}
	return nil
}

func (l *PlaybookList) SetItems(items []playbooks.Playbook) tea.Cmd {
	return l.List.SetItems(items)
}

func (l *PlaybookList) SetTitle(title string) {
	l.Model.Title = title
}
