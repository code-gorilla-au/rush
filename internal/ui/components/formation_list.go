package components

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/code-gorilla-au/rush/internal/playbooks"
	"github.com/code-gorilla-au/rush/internal/ui/styles"
)

type FormationList struct {
	List[playbooks.Formation]
	showDescription bool
}

type FormationListConfig struct {
	Title           string
	Items           []playbooks.Formation
	EnableFiltering bool
	ShowDescription bool
}

func NewFormationList(config FormationListConfig, theme styles.IceTheme) FormationList {
	l := FormationList{
		List: NewList(ListConfig[playbooks.Formation]{
			Title:             config.Title,
			Items:             config.Items,
			EnableFiltering:   config.EnableFiltering,
			DisableAutoResize: true,
			ItemMapper: func(f playbooks.Formation) ListItem[playbooks.Formation] {
				desc := ""
				if config.ShowDescription {
					desc = f.Description
				}
				return ListItem[playbooks.Formation]{
					Data:      f,
					TitleVal:  fmt.Sprintf("%s (%d-%d-%d)", f.Name, f.Lane1, f.Lane2, f.Lane3),
					DescVal:   desc,
					FilterVal: f.Name,
				}
			},
		}, theme),
		showDescription: config.ShowDescription,
	}
	l.Model.Title = ""
	return l
}

func (l *FormationList) Update(msg tea.Msg) (FormationList, tea.Cmd) {
	var cmd tea.Cmd
	l.List, cmd = l.List.Update(msg)
	return *l, cmd
}

func (l *FormationList) View(theme styles.IceTheme) string {
	return l.List.View()
}

func (l *FormationList) SelectedItem() playbooks.Formation {
	if item, ok := l.List.SelectedItem(); ok {
		return item
	}
	return playbooks.Formation{}
}

func (l *FormationList) SetItems(formations []playbooks.Formation) tea.Cmd {
	return l.List.SetItems(formations)
}

func (l *FormationList) SelectedIndex() int {
	return l.Model.Index()
}
