package ui

import (
	tea "charm.land/bubbletea/v2"

	"keyboard-tester/internal/capture"
	"keyboard-tester/internal/keyboard"
)

func (m *Model) receiveDevices(msg devicesMsg) {
	m.devices, m.loading = msg.devices, false
	m.cursor = max(0, min(m.cursor, len(m.devices)-1))
	if msg.err != nil {
		m.notice = msg.err.Error()
	}
}

func (m *Model) receiveStart(msg startedMsg) {
	if msg.err != nil {
		m.phase, m.notice = selecting, msg.err.Error()
		return
	}
	if m.device.ID != msg.device.ID {
		m.session = keyboard.NewSession(m.session.Layout)
	}
	m.session.Release()
	m.device, m.generation, m.phase = msg.device, msg.generation, capturing
}

func (m *Model) receiveStop(msg stoppedMsg) tea.Cmd {
	m.phase = selecting
	m.session.Release()
	if msg.err != nil {
		m.notice = msg.err.Error()
		return nil
	}
	return m.refresh()
}

func (m *Model) receiveEvent(event capture.Event) tea.Cmd {
	if m.phase != capturing || event.Generation != m.generation {
		return nil
	}
	if event.Err != nil {
		m.phase, m.notice = selecting, event.Err.Error()
		m.session.Release()
		return nil
	}
	if !m.session.Record(event.Code, event.Down) || !m.session.ControlAlt() {
		return nil
	}
	if event.Code.Page != keyboard.KeyboardPage {
		return nil
	}
	return m.physicalShortcut(event.Code.Usage)
}

func (m *Model) physicalShortcut(usage uint32) tea.Cmd {
	switch usage {
	case 0x52:
		m.scrollBy(-scrollStep)
	case 0x51:
		m.scrollBy(scrollStep)
	default:
		return m.action(shortcuts[usage])
	}
	return nil
}
