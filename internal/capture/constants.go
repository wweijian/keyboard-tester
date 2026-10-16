package capture

import "time"

const (
	commandList = iota
	commandStart
	commandStop
	eventRemoved  = 1
	eventOverflow = 2
	eventBuffer   = 4096
	pollInterval  = 10 * time.Millisecond
)

const (
	returnNotPrivileged uint32 = 0xe00002c1
	returnNoDevice      uint32 = 0xe00002c0
	returnExclusive     uint32 = 0xe00002c5
	returnNotPermitted  uint32 = 0xe00002e2
)
