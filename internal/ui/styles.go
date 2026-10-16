package ui

import "charm.land/lipgloss/v2"

var (
	ink           = lipgloss.Color("#DCE4F2")
	muted         = lipgloss.Color("#8B9BB4")
	accent        = lipgloss.Color("#87DCCD")
	background    = lipgloss.Color("#202B3D")
	titleStyle    = lipgloss.NewStyle().Foreground(accent).Bold(true)
	mutedStyle    = lipgloss.NewStyle().Foreground(muted)
	warningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2C879"))
	buttonStyle   = lipgloss.NewStyle().Foreground(ink).Background(background).Padding(0, 1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#101B27")).Background(accent).Bold(true)
	tileStyle     = lipgloss.NewStyle().Foreground(muted).Background(background).Width(tileWidth - 1).Align(lipgloss.Center)
	seenStyle     = tileStyle.Foreground(lipgloss.Color("#101B27")).Background(accent).Bold(true)
	panelStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#35465E")).Padding(1, 1)
)
