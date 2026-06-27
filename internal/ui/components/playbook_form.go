package components

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type PlaybookForm struct {
	NameInput        textinput.Model
	DescriptionInput textinput.Model
	focusIndex       int
}

func NewPlaybookForm() PlaybookForm {
	name := textinput.New()
	name.Placeholder = "Playbook Name"
	name.Focus()

	description := textinput.New()
	description.Placeholder = "Description"

	return PlaybookForm{
		NameInput:        name,
		DescriptionInput: description,
		focusIndex:       0,
	}
}

func (f *PlaybookForm) Update(msg tea.Msg) (PlaybookForm, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				f.focusIndex--
			} else {
				f.focusIndex++
			}

			if f.focusIndex > 1 {
				f.focusIndex = 0
			} else if f.focusIndex < 0 {
				f.focusIndex = 1
			}

			var cmd tea.Cmd
			if f.focusIndex == 0 {
				cmd = f.NameInput.Focus()
				f.DescriptionInput.Blur()
			} else {
				f.NameInput.Blur()
				cmd = f.DescriptionInput.Focus()
			}
			cmds = append(cmds, cmd)
		}
	}

	var cmd tea.Cmd
	f.NameInput, cmd = f.NameInput.Update(msg)
	cmds = append(cmds, cmd)

	f.DescriptionInput, cmd = f.DescriptionInput.Update(msg)
	cmds = append(cmds, cmd)

	return *f, tea.Batch(cmds...)
}

func (f *PlaybookForm) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Name:",
		f.NameInput.View(),
		"",
		"Description:",
		f.DescriptionInput.View(),
	)
}

func (f *PlaybookForm) SetValues(name, description string) {
	f.NameInput.SetValue(name)
	f.DescriptionInput.SetValue(description)
}

func (f *PlaybookForm) Values() (string, string) {
	return f.NameInput.Value(), f.DescriptionInput.Value()
}

func (f *PlaybookForm) Reset() {
	f.NameInput.Reset()
	f.DescriptionInput.Reset()
	f.focusIndex = 0
	f.NameInput.Focus()
	f.DescriptionInput.Blur()
}
