package components

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/code-gorilla-au/rush/internal/teams"
)

type PlayerItem struct {
	Player    teams.Player
	Input     textinput.Model
	IsEditing bool
}

type PlayerList struct {
	cursor int
	Items  []PlayerItem
}

func NewPlayerList(players []teams.Player) PlayerList {
	items := make([]PlayerItem, len(players))
	for i, p := range players {
		ti := textinput.New()
		ti.SetValue(p.Name)
		ti.CharLimit = 50
		ti.SetWidth(20)
		items[i] = PlayerItem{
			Player: p,
			Input:  ti,
		}
	}
	return PlayerList{
		Items: items,
	}
}

func (l *PlayerList) Update(msg tea.Msg) tea.Cmd {
	if len(l.Items) == 0 {
		return nil
	}

	item := &l.Items[l.cursor]

	if item.IsEditing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				item.IsEditing = false
				item.Input.Blur()
				item.Player.Name = item.Input.Value()
				return func() tea.Msg {
					return MsgPlayerUpdated{Player: item.Player}
				}
			case "esc":
				item.IsEditing = false
				item.Input.Blur()
				item.Input.SetValue(item.Player.Name)
				return nil
			}
		}
		var cmd tea.Cmd
		item.Input, cmd = item.Input.Update(msg)
		return cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if l.cursor > 0 {
				l.cursor--
			}
		case "down", "j":
			if l.cursor < len(l.Items)-1 {
				l.cursor++
			}
		case "enter":
			item.IsEditing = true
			return item.Input.Focus()
		}
	}
	return nil
}

func (l *PlayerList) View(itemStyle lipgloss.Style, selectedStyle lipgloss.Style) string {
	var s string
	for i, item := range l.Items {
		var content string
		if item.IsEditing {
			content = item.Input.View()
		} else {
			content = item.Player.Name
		}

		if i == l.cursor {
			s += selectedStyle.Render("> " + content)
		} else {
			s += itemStyle.Render("  " + content)
		}
		if i < len(l.Items)-1 {
			s += "\n"
		}
	}
	return s
}

type MsgPlayerUpdated struct {
	Player teams.Player
}
