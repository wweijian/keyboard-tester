package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m *Model) View() tea.View {
	width := max(20, m.width-4)
	header := titleStyle.Render("KEYBOARD LAB") + "  " + mutedStyle.Render(m.session.Layout.Size+"% / US ANSI")
	stats := m.summary()
	body := visibleLines(m.content(width), m.scroll, m.bodyHeight())
	content := strings.Join([]string{header, stats, m.toolbar(), "", body, "", m.footer(width)}, "\n")
	view := tea.NewView(lipgloss.NewStyle().Padding(1, 2).Render(content))
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}

func (m *Model) summary() string {
	tested, expected := len(m.session.Tested), len(m.session.Layout.Keys)
	percent := 100 * tested / expected
	bar := strings.Repeat("━", percent/10) + strings.Repeat("─", 10-percent/10)
	return fmt.Sprintf("%s  %d/%d tested · %d%% · %d presses", titleStyle.Render(bar), tested, expected, percent, m.session.Presses)
}

func (m *Model) toolbar() string {
	var parts []string
	for _, button := range m.buttons() {
		parts = append(parts, buttonStyle.Render(button.label))
	}
	return strings.Join(parts, " ")
}

func (m *Model) content(width int) string {
	if m.width < minimumWidth || m.height < minimumHeight {
		return fmt.Sprintf("Resize to at least %d columns × %d rows.\nMouse controls above remain available.", minimumWidth, minimumHeight)
	}
	if m.help {
		return helpView(width)
	}
	if m.phase != capturing {
		return m.selectionView(width)
	}
	return m.testingView(width)
}

func (m *Model) bodyHeight() int {
	return max(1, m.height-10-lipgloss.Height(m.footer(max(20, m.width-4))))
}

func (m *Model) scrollBy(delta int) {
	maximum := max(0, lipgloss.Height(m.content(max(20, m.width-4)))-m.bodyHeight())
	m.scroll = max(0, min(maximum, m.scroll+delta))
}

func visibleLines(content string, offset, height int) string {
	lines := strings.Split(content, "\n")
	offset = max(0, min(offset, max(0, len(lines)-height)))
	end := min(len(lines), offset+height)
	return lipgloss.NewStyle().Height(height).Render(strings.Join(lines[offset:end], "\n"))
}

func (m *Model) footer(width int) string {
	text := "↑/↓ select · Enter start · R refresh · ? help · Q quit"
	if m.phase == capturing {
		text = "Ctrl+Opt: Q quit · R reset · P devices · H help · ↑/↓ scroll"
	}
	if m.notice != "" {
		return warningStyle.Width(width).Render(m.notice) + "\n" + mutedStyle.Width(width).Render(text)
	}
	return mutedStyle.Width(width).Render(text + "\nMouse: click buttons · wheel to scroll")
}
