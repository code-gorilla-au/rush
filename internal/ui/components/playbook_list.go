package components

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
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

func NewPlaybookList(items []playbooks.Playbook, theme styles.IceTheme) PlaybookList {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = PlaybookItem{playbook: item}
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

func (l *PlaybookList) View(theme styles.IceTheme) string {
	return l.list.View()
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
