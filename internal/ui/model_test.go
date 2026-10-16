package ui

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"

	"keyboard-tester/internal/capture"
	"keyboard-tester/internal/keyboard"
)

type fakeDriver struct {
	started  uint64
	stopped  bool
	startErr error
	devices  []capture.Device
	events   chan capture.Event
}

func (d *fakeDriver) Devices() ([]capture.Device, error) { return d.devices, nil }
func (d *fakeDriver) Events() <-chan capture.Event       { return d.events }
func (d *fakeDriver) Close()                             {}
func (d *fakeDriver) Stop() error                        { d.stopped = true; return nil }
func (d *fakeDriver) Start(id uint64) (uint64, error) {
	d.started = id
	return 1, d.startErr
}

func testModel(t *testing.T) (*Model, *fakeDriver) {
	t.Helper()
	layout, err := keyboard.Lookup("60")
	if err != nil {
		t.Fatal(err)
	}
	d := &fakeDriver{devices: []capture.Device{{ID: 10, Name: "First"}, {ID: 20, Name: "Second"}}}
	m := New(d, layout)
	m.Update(m.loadDevices())
	return m, d
}

func startTestModel(t *testing.T) (*Model, *fakeDriver) {
	t.Helper()
	m, d := testModel(t)
	m.Update(m.startCapture()())
	return m, d
}

func press(m *Model, usage uint32) tea.Cmd {
	return m.receiveEvent(capture.Event{Code: keyboard.Code{Page: keyboard.KeyboardPage, Usage: usage}, Down: true, Generation: m.generation})
}

func TestStartsSelectedPhysicalDevice(t *testing.T) {
	m, d := testModel(t)
	m.terminalKey("down")
	m.Update(m.terminalKey("enter")())
	if d.started != 20 || m.device.ID != 20 || m.phase != capturing {
		t.Fatal("capture did not use the selected physical device")
	}
}

func TestTerminalAndStaleEventsCannotCount(t *testing.T) {
	m, _ := startTestModel(t)
	m.terminalKey("a")
	m.receiveEvent(capture.Event{Code: keyboard.Code{Page: keyboard.KeyboardPage, Usage: 4}, Down: true, Generation: 999})
	if m.session.Presses != 0 {
		t.Fatal("terminal or stale device input contributed results")
	}
	press(m, 4)
	if !m.session.Tested["a"] {
		t.Fatal("current device input was not recorded")
	}
}

func TestShortcutAndPlainQ(t *testing.T) {
	m, _ := startTestModel(t)
	if press(m, 0x14) != nil || !m.session.Tested["q"] {
		t.Fatal("plain Q must be testable")
	}
	m.session.Release()
	press(m, keyboard.LeftControl)
	press(m, keyboard.RightAlt)
	command := press(m, 0x14)
	if command == nil {
		t.Fatal("quit chord was not recognized")
	}
	if _, ok := command().(tea.QuitMsg); !ok {
		t.Fatal("quit chord did not quit")
	}
}

func TestCaptureFailureAndDisconnect(t *testing.T) {
	m, d := testModel(t)
	d.startErr = errors.New("permission denied")
	m.Update(m.startCapture()())
	if m.phase != selecting || m.notice != "permission denied" {
		t.Fatal("failed capture must remain in device selection")
	}
	d.startErr = nil
	m.Update(m.startCapture()())
	press(m, 4)
	m.receiveEvent(capture.Event{Generation: m.generation, Err: errors.New("disconnected")})
	if m.phase != selecting || !m.session.Tested["a"] || m.notice != "disconnected" {
		t.Fatal("disconnect must preserve results and show a recovery message")
	}
}

func TestDeviceChangeClearsResults(t *testing.T) {
	m, _ := startTestModel(t)
	press(m, 4)
	m.Update(m.stopCapture()())
	m.loading = false
	m.Update(m.startCapture()())
	if !m.session.Tested["a"] {
		t.Fatal("resuming the same device lost its results")
	}
	m.Update(m.stopCapture()())
	m.loading, m.cursor = false, 1
	m.Update(m.startCapture()())
	if m.session.Presses != 0 || len(m.session.Tested) != 0 {
		t.Fatal("changing devices must clear previous results")
	}
}

func TestMouseCanReleaseKeyboard(t *testing.T) {
	m, d := startTestModel(t)
	command := m.click(tea.MouseClickMsg{Button: tea.MouseLeft, X: 3, Y: toolbarY})
	if command == nil {
		t.Fatal("Devices button did not release capture")
	}
	m.Update(command())
	if !d.stopped || m.phase != selecting {
		t.Fatal("keyboard was not released")
	}
}
