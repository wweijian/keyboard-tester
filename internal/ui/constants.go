package ui

const (
	selecting = iota
	starting
	capturing
	stopping
	defaultWidth      = 100
	defaultHeight     = 30
	minimumWidth      = 64
	minimumHeight     = 18
	tileWidth         = 8
	sidebarWidth      = 28
	sidebarBreakpoint = 108
	toolbarY          = 3
	scrollStep        = 3
)

const (
	actionStart = "start"
	actionReset = "reset"
	actionHelp  = "help"
	actionQuit  = "quit"
)

var shortcuts = map[uint32]string{
	0x14: actionQuit,  // Q
	0x15: actionReset, // R
	0x13: actionStart, // P: release capture and choose a device
	0x0b: actionHelp,  // H
}
