package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type button struct {
	label  string
	action string
}

func (m *Model) buttons() []button {
	label := "Start"
	if m.phase == capturing {
		label = "Devices"
	}
	return []button{{label, actionStart}, {"Reset", actionReset}, {"Help", actionHelp}, {"Quit", actionQuit}}
}

func (m *Model) action(action string) tea.Cmd {
	switch action {
	case actionQuit:
		return tea.Quit
	case actionStart:
		if m.phase == capturing {
			return m.stopCapture()
		}
		return m.startCapture()
	case actionReset:
		m.session.Reset()
		m.scroll = 0
	case actionHelp:
		m.help, m.scroll = !m.help, 0
	}
	return nil
}

func (m *Model) terminalKey(key string) tea.Cmd {
	if key == "ctrl+c" {
		return tea.Quit
	}
	if m.phase != selecting {
		return nil
	}
	switch key {
	case "up", "k":
		m.cursor = max(0, m.cursor-1)
	case "down", "j":
		m.cursor = min(max(0, len(m.devices)-1), m.cursor+1)
	case "enter":
		return m.startCapture()
	case "r":
		return m.refresh()
	case "?", "h":
		return m.action(actionHelp)
	case "q", "esc":
		return tea.Quit
	}
	return nil
}

func (m *Model) click(msg tea.MouseClickMsg) tea.Cmd {
	if msg.Button != tea.MouseLeft || msg.Y != toolbarY {
		return nil
	}
	x := 2
	for _, button := range m.buttons() {
		width := lipgloss.Width(buttonStyle.Render(button.label))
		if msg.X >= x && msg.X < x+width {
			return m.action(button.action)
		}
		x += width + 1
	}
	return nil
}

func (m *Model) wheel(msg tea.MouseWheelMsg) {
	delta := scrollStep
	if msg.Button == tea.MouseWheelUp {
		delta = -scrollStep
	}
	if m.phase == selecting && !m.help {
		m.cursor = max(0, min(max(0, len(m.devices)-1), m.cursor+delta))
		return
	}
	m.scrollBy(delta)
}
