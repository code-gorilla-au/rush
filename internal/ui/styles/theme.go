package styles

import "charm.land/lipgloss/v2"

// IceTheme defines the "ice" color palette and styles.
type IceTheme struct {
	Logo            lipgloss.Style
	Footer          lipgloss.Style
	Base            lipgloss.Style
	Button          lipgloss.Style
	Hotkey          lipgloss.Style
	ListSelected    lipgloss.Style
	CoachTeam       lipgloss.Style
	CoachName       lipgloss.Style
	Title           lipgloss.Style
	Header          lipgloss.Style
	Winner          lipgloss.Style
	Muted           lipgloss.Style
	ActiveBorder    lipgloss.Style
	InactiveBorder  lipgloss.Style
	SelectedTitle   lipgloss.Style
	SelectedDesc    lipgloss.Style
	TeamA           lipgloss.Style
	TeamB           lipgloss.Style
	Player          lipgloss.Style
	Separator       lipgloss.Style
	Label           lipgloss.Style
	Highlight       lipgloss.Style
	SecondaryHeader lipgloss.Style
	RoundBorder     lipgloss.Style
}

// NewIceTheme returns a new IceTheme with the "ice" palette.
func NewIceTheme() IceTheme {
	// Colors
	iceBlue := lipgloss.Color("#A5F2F3")
	skyBlue := lipgloss.Color("#87CEEB")
	white := lipgloss.Color("#FFFFFF")
	black := lipgloss.Color("#000000")
	darkGrey := lipgloss.Color("#333333")
	grey := lipgloss.Color("#666666")
	mutedGrey := lipgloss.Color("#888888")
	faintGrey := lipgloss.Color("#555555")
	gold := lipgloss.Color("#FFD700")

	return IceTheme{
		Logo: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true),
		Footer: lipgloss.NewStyle().
			Foreground(skyBlue).
			Italic(true),
		Base: lipgloss.NewStyle().
			Foreground(white),
		Button: lipgloss.NewStyle().
			Foreground(black).
			Background(iceBlue).
			Padding(0, 3).
			MarginTop(1).
			Bold(true),
		Hotkey: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true),
		ListSelected: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true),
		CoachTeam: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true),
		CoachName: lipgloss.NewStyle().
			Foreground(faintGrey).
			Italic(true).
			Faint(true),
		Title: lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(iceBlue).
			Bold(true),
		Header: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true).
			MarginBottom(1),
		Winner: lipgloss.NewStyle().
			Foreground(gold).
			Bold(true).
			MarginTop(1),
		Muted: lipgloss.NewStyle().
			Foreground(mutedGrey),
		ActiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(iceBlue).
			Padding(1),
		InactiveBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(darkGrey).
			Padding(1).
			Foreground(grey),
		SelectedTitle: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(iceBlue).
			PaddingLeft(2),
		SelectedDesc: lipgloss.NewStyle().
			Foreground(mutedGrey).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(iceBlue).
			PaddingLeft(2),
		TeamA: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true).
			Width(10).
			Align(lipgloss.Right),
		TeamB: lipgloss.NewStyle().
			Foreground(white).
			Bold(true).
			Width(10).
			Align(lipgloss.Left),
		Player: lipgloss.NewStyle().
			Foreground(skyBlue),
		Separator: lipgloss.NewStyle().
			Foreground(faintGrey),
		Label: lipgloss.NewStyle().
			Foreground(mutedGrey).
			PaddingLeft(2),
		Highlight: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true),
		SecondaryHeader: lipgloss.NewStyle().
			Foreground(white).
			Bold(true).
			MarginBottom(1),
		RoundBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(skyBlue).
			Padding(1),
	}
}
