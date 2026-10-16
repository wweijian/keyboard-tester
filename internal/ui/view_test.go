package ui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"keyboard-tester/internal/keyboard"
)

func TestViewsFitTerminal(t *testing.T) {
	m, _ := startTestModel(t)
	for _, size := range [][2]int{{64, 18}, {80, 24}, {120, 35}, {160, 40}} {
		m.width, m.height = size[0], size[1]
		for _, phase := range []int{selecting, capturing} {
			m.phase = phase
			assertViewFits(t, m)
		}
	}
}

func assertViewFits(t *testing.T, m *Model) {
	t.Helper()
	view := m.View()
	if lipgloss.Width(view.Content) > m.width || lipgloss.Height(view.Content) > m.height {
		t.Errorf("phase %d: view %dx%d exceeds terminal %dx%d", m.phase, lipgloss.Width(view.Content), lipgloss.Height(view.Content), m.width, m.height)
	}
	if !view.AltScreen || !strings.Contains(view.Content, "Quit") {
		t.Error("view must expose mouse controls in the alternate screen")
	}
}

func TestUnknownKeysAndCompletionAreVisible(t *testing.T) {
	m, _ := startTestModel(t)
	press(m, 0x3a)
	content := m.unknownView(28)
	if !strings.Contains(content, "0007:003A") || !strings.Contains(content, "F1") {
		t.Fatal("out-of-preset HID ID and label must be visible")
	}
	for id := range m.session.Layout.Keys {
		m.session.Tested[id] = true
	}
	if !strings.Contains(m.summary(), "100%") || !strings.Contains(m.keyTile("a"), "✓") {
		t.Fatal("completion must have both text and a tested marker")
	}
}

func TestScrollClampsAndResizeResets(t *testing.T) {
	m, _ := startTestModel(t)
	m.width, m.height = 64, 18
	m.scrollBy(10000)
	if m.scroll == 0 {
		t.Fatal("long key maps must be scrollable")
	}
	m.scrollBy(-10000)
	if m.scroll != 0 {
		t.Fatal("scroll moved before the beginning")
	}
	if got := visibleLines("a\nb\nc", 100, 2); !strings.Contains(got, "b\nc") {
		t.Fatal("viewport did not clamp to the final lines")
	}
	m.scrollBy(10000)
	m.Update(tea.WindowSizeMsg{Width: 120, Height: 35})
	if m.scroll != 0 {
		t.Fatal("resize must reset scrolling")
	}
}

func TestAllPresetViewsFitTerminal(t *testing.T) {
	for _, size := range keyboard.Sizes {
		m, _ := startTestModel(t)
		layout, _ := keyboard.Lookup(size)
		m.session = keyboard.NewSession(layout)
		for _, width := range []int{64, 80, 120, 160} {
			m.width, m.height = width, 24
			assertViewFits(t, m)
		}
	}
}

func TestKeypadMultiplyReachesUI(t *testing.T) {
	m, _ := startTestModel(t)
	layout, _ := keyboard.Lookup("100")
	m.session = keyboard.NewSession(layout)
	press(m, 0x55)
	if !m.session.Tested["kpmultiply"] || !strings.Contains(m.keyTile("kpmultiply"), "✓") {
		t.Fatal("standard keypad multiply event must mark its tile tested")
	}
	if !strings.Contains(m.lastKey(), "KP *") || !strings.Contains(m.lastKey(), "0007:0055") {
		t.Fatal("last-key feedback must identify keypad multiply and its HID code")
	}
	if !strings.Contains(m.testingView(116), m.lastKey()) {
		t.Fatal("last-key feedback must be above the keyboard map")
	}
}
