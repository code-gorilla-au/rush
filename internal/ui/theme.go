package ui

import "charm.land/lipgloss/v2"

// IceTheme defines the "ice" color palette and styles.
type IceTheme struct {
	Logo   lipgloss.Style
	Footer lipgloss.Style
	Base   lipgloss.Style
}

// NewIceTheme returns a new IceTheme with the "ice" palette.
func NewIceTheme() IceTheme {
	// Colors
	iceBlue := lipgloss.Color("#A5F2F3")
	skyBlue := lipgloss.Color("#87CEEB")
	white := lipgloss.Color("#FFFFFF")
	black := lipgloss.Color("#000000")

	return IceTheme{
		Logo: lipgloss.NewStyle().
			Foreground(iceBlue).
			Bold(true),
		Footer: lipgloss.NewStyle().
			Foreground(skyBlue).
			Italic(true),
		Base: lipgloss.NewStyle().
			Background(black).
			Foreground(white),
	}
}
