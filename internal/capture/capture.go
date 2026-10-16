package capture

import "keyboard-tester/internal/keyboard"

type Device struct {
	ID        uint64
	Name      string
	Transport string
	Vendor    uint32
	Product   uint32
}

type Event struct {
	Code       keyboard.Code
	Down       bool
	Generation uint64
	Err        error
}

type Driver interface {
	Devices() ([]Device, error)
	Start(id uint64) (uint64, error)
	Stop() error
	Events() <-chan Event
	Close()
}
