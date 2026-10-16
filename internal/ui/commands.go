package ui

import (
	tea "charm.land/bubbletea/v2"

	"keyboard-tester/internal/capture"
)

type devicesMsg struct {
	devices []capture.Device
	err     error
}

type startedMsg struct {
	device     capture.Device
	generation uint64
	err        error
}

type stoppedMsg struct{ err error }
type captureClosedMsg struct{}

func (m *Model) loadDevices() tea.Msg {
	devices, err := m.driver.Devices()
	return devicesMsg{devices, err}
}

func (m *Model) waitEvent() tea.Msg {
	event, open := <-m.driver.Events()
	if !open {
		return captureClosedMsg{}
	}
	return event
}

func (m *Model) startCapture() tea.Cmd {
	if m.phase != selecting || m.loading || len(m.devices) == 0 {
		return nil
	}
	device := m.devices[m.cursor]
	m.phase, m.notice, m.scroll = starting, "", 0
	return func() tea.Msg {
		generation, err := m.driver.Start(device.ID)
		return startedMsg{device, generation, err}
	}
}

func (m *Model) stopCapture() tea.Cmd {
	if m.phase != capturing {
		return nil
	}
	m.phase, m.help, m.scroll = stopping, false, 0
	return func() tea.Msg { return stoppedMsg{m.driver.Stop()} }
}

func (m *Model) refresh() tea.Cmd {
	if m.phase != selecting || m.loading {
		return nil
	}
	m.loading, m.notice = true, ""
	return m.loadDevices
}
