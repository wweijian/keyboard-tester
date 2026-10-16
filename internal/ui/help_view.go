package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

var helpLines = []string{
	"HOW TO TEST",
	"",
	"Select a keyboard with ↑/↓ and press Enter, or click Start.",
	"Press each physical key once. ✓ means the key has been tested.",
	"Uppercase and lowercase count as the same physical key.",
	"Left/right modifiers and keypad keys are tested separately.",
	"Keys outside the preset appear in UNKNOWN with their HID IDs.",
	"There is no hold test. All results live only in RAM.",
	"",
	"DURING CAPTURE — hold Control + Option, then:",
	"Q quit · R reset · P release capture / select a device · H help",
	"↑/↓ scroll the key map and unknown list. Mouse wheel also scrolls.",
	"Mouse buttons above work even if keyboard keys are broken.",
	"Testing continues while this help is open. Click Help to close it.",
	"",
	"Capture suppresses the selected keyboard's events system-wide.",
	"Use Devices to release it before switching to another application.",
	"Other keyboards do not contribute results. Ctrl+C from one can quit.",
	"Selecting a different device clears the previous device's results.",
	"",
	"If macOS denies access, grant Input Monitoring to the terminal or",
	"binary in System Settings > Privacy & Security, then relaunch.",
	"Some built-in keyboards are protected against exclusive access.",
	"Firmware-only Fn/layer keys cannot be detected unless they emit HID input.",
	"Presets are reference layouts; keyboard size does not guarantee exact keys.",
}

func helpView(width int) string {
	return lipgloss.NewStyle().Width(width).Render(strings.Join(helpLines, "\n"))
}
