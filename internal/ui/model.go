package ui

import (
	tea "charm.land/bubbletea/v2"

	"keyboard-tester/internal/capture"
	"keyboard-tester/internal/keyboard"
)

type Model struct {
	driver        capture.Driver
	session       *keyboard.Session
	devices       []capture.Device
	device        capture.Device
	cursor        int
	phase         int
	generation    uint64
	width, height int
	scroll        int
	help          bool
	loading       bool
	notice        string
}

func New(driver capture.Driver, layout keyboard.Layout) *Model {
	return &Model{driver: driver, session: keyboard.NewSession(layout), loading: true,
		width: defaultWidth, height: defaultHeight}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.loadDevices, m.waitEvent)
}

func (m *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.scroll = 0
	case devicesMsg:
		m.receiveDevices(msg)
	case startedMsg:
		m.receiveStart(msg)
	case stoppedMsg:
		return m, m.receiveStop(msg)
	case capture.Event:
		return m, tea.Batch(m.receiveEvent(msg), m.waitEvent)
	case captureClosedMsg:
		m.phase, m.notice = selecting, "Keyboard capture service closed. Quit and relaunch."
	case tea.KeyPressMsg:
		return m, m.terminalKey(msg.String())
	case tea.MouseClickMsg:
		return m, m.click(msg)
	case tea.MouseWheelMsg:
		m.wheel(msg)
	}
	return m, nil
}
