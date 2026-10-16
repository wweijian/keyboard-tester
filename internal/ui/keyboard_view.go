package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"keyboard-tester/internal/keyboard"
)

func (m *Model) testingView(width int) string {
	intro := titleStyle.Width(width).Render(m.device.Name) + "\n" + mutedStyle.Render("· pending   ✓ tested   /   results stay until reset")
	intro += "\n" + mutedStyle.Width(width).Render(m.lastKey())
	if width < sidebarBreakpoint {
		return intro + "\n\n" + m.keyboardView(width) + "\n\n" + m.unknownView(width)
	}
	keys := m.keyboardView(width - sidebarWidth - 3)
	unknown := m.unknownView(sidebarWidth)
	return intro + "\n\n" + lipgloss.JoinHorizontal(lipgloss.Top, keys, "   ", unknown)
}

func (m *Model) lastKey() string {
	if m.session.Presses == 0 {
		return "Last key: waiting for input from the selected keyboard"
	}
	return fmt.Sprintf("Last key: %s · HID %s", keyboard.Describe(m.session.Last), m.session.Last)
}

func (m *Model) keyboardView(width int) string {
	columns := max(1, (width-4)/tileWidth)
	var lines []string
	for _, row := range m.session.Layout.Rows {
		for start := 0; start < len(row); start += columns {
			var tiles []string
			for _, id := range row[start:min(len(row), start+columns)] {
				tiles = append(tiles, m.keyTile(id))
			}
			lines = append(lines, strings.Join(tiles, " "))
		}
		lines = append(lines, "")
	}
	return panelStyle.Width(width - 2).Render(strings.TrimRight(strings.Join(lines, "\n"), "\n"))
}

func (m *Model) keyTile(id string) string {
	style, prefix := tileStyle, "·"
	if m.session.Tested[id] {
		style, prefix = seenStyle, "✓"
	}
	return style.Render(prefix + " " + keyboard.Label(id))
}

func (m *Model) unknownView(width int) string {
	lines := []string{titleStyle.Render("UNKNOWN / OUTSIDE PRESET"), mutedStyle.Render("HID page:usage · presses"), ""}
	if len(m.session.Unknown) == 0 {
		lines = append(lines, mutedStyle.Render("None recorded"))
	}
	for _, key := range m.session.Unknown {
		lines = append(lines, fmt.Sprintf("%s  ×%d", key.Code, key.Count))
		lines = append(lines, mutedStyle.Render(keyboard.Describe(key.Code)))
	}
	if m.session.Presses > 0 {
		lines = append(lines, "", "Last: "+keyboard.Describe(m.session.Last), m.session.Last.String())
	}
	return lipgloss.NewStyle().Width(width).Render(strings.Join(lines, "\n"))
}
