package components

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type List[T any] struct {
	Model      list.Model
	Active     bool
	ItemMapper func(T) ListItem[T]
}

type ListItem[T any] struct {
	Data      T
	TitleVal  string
	DescVal   string
	FilterVal string
}

func (i ListItem[T]) Title() string       { return i.TitleVal }
func (i ListItem[T]) Description() string { return i.DescVal }
func (i ListItem[T]) FilterValue() string { return i.FilterVal }

type ListConfig[T any] struct {
	Title           string
	Items           []T
	ItemMapper      func(T) ListItem[T]
	EnableFiltering bool
}

func NewList[T any](config ListConfig[T], theme styles.IceTheme) List[T] {
	items := make([]list.Item, len(config.Items))
	for i, item := range config.Items {
		items[i] = config.ItemMapper(item)
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = theme.Base.PaddingLeft(3)
	delegate.Styles.NormalDesc = theme.Muted.PaddingLeft(3)
	delegate.Styles.SelectedTitle = theme.SelectedTitle
	delegate.Styles.SelectedDesc = theme.SelectedDesc
	delegate.Styles.DimmedTitle = theme.Muted.PaddingLeft(3)
	delegate.Styles.DimmedDesc = theme.Muted.PaddingLeft(3)

	l := list.New(items, delegate, 0, 0)
	l.Title = config.Title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(config.EnableFiltering)
	l.SetShowHelp(false)
	l.Styles.Title = theme.Title

	return List[T]{
		Model:      l,
		Active:     true,
		ItemMapper: config.ItemMapper,
	}
}

func (l *List[T]) Update(msg tea.Msg) (List[T], tea.Cmd) {
	if !l.Active {
		return *l, nil
	}
	var cmd tea.Cmd
	l.Model, cmd = l.Model.Update(msg)
	return *l, cmd
}

func (l *List[T]) View() string {
	return l.Model.View()
}

func (l *List[T]) SetSize(width, height int) {
	l.Model.SetSize(width, height)
}

func (l *List[T]) SetItems(items []T) tea.Cmd {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = l.ItemMapper(item)
	}
	return l.Model.SetItems(listItems)
}

func (l *List[T]) SelectedItem() (T, bool) {
	item := l.Model.SelectedItem()
	if item == nil {
		var zero T
		return zero, false
	}
	if li, ok := item.(ListItem[T]); ok {
		return li.Data, true
	}
	var zero T
	return zero, false
}

func (l *List[T]) Len() int {
	return len(l.Model.Items())
}

func (l *List[T]) IsFiltering() bool {
	return l.Model.FilterState() == list.Filtering
}

func (l *List[T]) Reset() {
	l.Model.Select(0)
	l.Model.FilterInput.Reset()
}

func (l *List[T]) SetActive(active bool) {
	l.Active = active
}

func (l *List[T]) SetTitle(title string) {
	l.Model.Title = title
}
