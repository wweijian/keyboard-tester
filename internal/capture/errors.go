package capture

import "fmt"

func openError(code uint32, detail string) error {
	message := "The device may be unavailable or protected by macOS."
	switch code {
	case returnNotPrivileged:
		message = "Exclusive keyboard capture requires administrator privileges. Quit and relaunch the built binary with sudo."
	case returnNotPermitted:
		message = "Grant Input Monitoring to your terminal or binary in System Settings > Privacy & Security, then relaunch."
	case returnExclusive:
		message = "Another application owns the keyboard. Close other keyboard capture or remapping applications and retry."
	case returnNoDevice:
		message = "The keyboard is no longer available. Reconnect it and refresh the device list."
	}
	return fmt.Errorf("cannot capture keyboard: %s (0x%08X). %s", detail, code, message)
}
