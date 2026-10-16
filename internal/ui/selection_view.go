package ui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m *Model) selectionView(width int) string {
	heading := "SELECT A PHYSICAL KEYBOARD"
	if m.phase == starting {
		return titleStyle.Render("Starting exclusive capture…")
	}
	if m.phase == stopping {
		return titleStyle.Render("Releasing keyboard…")
	}
	lines := []string{titleStyle.Render(heading), ""}
	if m.loading {
		return strings.Join(append(lines, "Looking for keyboards…"), "\n")
	}
	if len(m.devices) == 0 {
		return strings.Join(append(lines, "No keyboards found. Connect a keyboard, then press R to refresh."), "\n")
	}
	lines = append(lines, m.deviceRows(width)...)
	lines = append(lines, "", mutedStyle.Width(width).Render("Choose your keyboard with ↑/↓, then press Enter or click Start. Capture requires sudo and Input Monitoring permission. The selected keyboard's shortcuts are suppressed until capture is released."))
	return strings.Join(lines, "\n")
}

func (m *Model) deviceRows(width int) []string {
	visible := max(1, (m.bodyHeight()-6)/3)
	first := max(0, m.cursor-visible+1)
	var lines []string
	for i := first; i < min(len(m.devices), first+visible); i++ {
		device := m.devices[i]
		name := fmt.Sprintf("  %d  %s", i+1, device.Name)
		style := lipgloss.NewStyle().Foreground(ink)
		if i == m.cursor {
			name, style = "›"+name[1:], selectedStyle
		}
		lines = append(lines, style.Width(width).MaxWidth(width).Render(name))
		details := fmt.Sprintf("     %s · %04X:%04X · device %X", device.Transport, device.Vendor, device.Product, device.ID)
		lines = append(lines, mutedStyle.MaxWidth(width).Render(details), "")
	}
	return lines
}
