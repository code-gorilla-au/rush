package components

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type List[T any] struct {
	Model         list.Model
	Active        bool
	AutoResize    bool
	ItemMapper    func(T) ListItem[T]
	width         int
	height        int
	topPadding    int
	bottomPadding int
	leftPadding   int
	rightPadding  int
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
	Title             string
	Items             []T
	ItemMapper        func(T) ListItem[T]
	EnableFiltering   bool
	DisableAutoResize bool
	TopPadding        int
	BottomPadding     int
	LeftPadding       int
	RightPadding      int
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
	l.SetShowTitle(config.Title != "")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(config.EnableFiltering)
	l.SetShowHelp(false)
	l.Styles.Title = theme.Title

	return List[T]{
		Model:         l,
		Active:        true,
		AutoResize:    !config.DisableAutoResize,
		ItemMapper:    config.ItemMapper,
		topPadding:    config.TopPadding,
		bottomPadding: config.BottomPadding,
		leftPadding:   config.LeftPadding,
		rightPadding:  config.RightPadding,
	}
}

func (l *List[T]) Update(msg tea.Msg) (List[T], tea.Cmd) {
	if !l.Active {
		return *l, nil
	}
	var cmd tea.Cmd

	l.Model, cmd = l.Model.Update(msg)

	if wm, ok := msg.(tea.WindowSizeMsg); ok {
		if l.AutoResize {
			l.SetSize(wm.Width, wm.Height)
		} else {
			l.SetSize(l.width, l.height)
		}
	}

	return *l, cmd
}

func (l *List[T]) View() string {
	style := lipgloss.NewStyle().
		Padding(l.topPadding, l.rightPadding, l.bottomPadding, l.leftPadding)

	if l.width > 0 {
		style = style.Width(l.width)
	}
	if l.height > 0 {
		style = style.Height(l.height)
	}

	return style.Render(l.Model.View())
}

func (l *List[T]) SetSize(width, height int) {
	l.width = width
	l.height = height
	l.Model.SetSize(
		max(0, width-l.leftPadding-l.rightPadding),
		max(0, height-l.topPadding-l.bottomPadding),
	)
}

func (l *List[T]) SetPadding(top, right, bottom, left int) {
	l.topPadding = top
	l.rightPadding = right
	l.bottomPadding = bottom
	l.leftPadding = left
	l.Model.SetSize(
		max(0, l.width-l.leftPadding-l.rightPadding),
		max(0, l.height-l.topPadding-l.bottomPadding),
	)
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
	l.Model.SetShowTitle(title != "")
}
