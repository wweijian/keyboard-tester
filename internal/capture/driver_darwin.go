//go:build darwin && cgo

package capture

/*
#cgo CFLAGS: -mmacosx-version-min=11.0
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#include "native_darwin.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"sort"
	"sync"
	"time"

	"keyboard-tester/internal/keyboard"
)

type request struct {
	kind  int
	id    uint64
	reply chan response
}

type response struct {
	devices    []Device
	generation uint64
	err        error
}

type driver struct {
	requests chan request
	events   chan Event
	done     chan struct{}
	closed   chan struct{}
	once     sync.Once
}

func New() (Driver, error) {
	d := &driver{requests: make(chan request), events: make(chan Event, eventBuffer),
		done: make(chan struct{}), closed: make(chan struct{})}
	ready := make(chan error, 1)
	go d.run(ready)
	if err := <-ready; err != nil {
		return nil, err
	}
	return d, nil
}

func (d *driver) Devices() ([]Device, error) {
	result := d.submit(commandList, 0)
	return result.devices, result.err
}

func (d *driver) Start(id uint64) (uint64, error) {
	result := d.submit(commandStart, id)
	return result.generation, result.err
}

func (d *driver) Stop() error {
	return d.submit(commandStop, 0).err
}

func (d *driver) Events() <-chan Event {
	return d.events
}

func (d *driver) Close() {
	d.once.Do(func() { close(d.done) })
	<-d.closed
}

func (d *driver) submit(kind int, id uint64) response {
	req := request{kind: kind, id: id, reply: make(chan response, 1)}
	select {
	case d.requests <- req:
		return <-req.reply
	case <-d.done:
		return response{err: errors.New("keyboard capture is closed")}
	}
}

func (d *driver) run(ready chan<- error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer close(d.closed)
	defer close(d.events)
	ctx := C.kt_create()
	if ctx == nil {
		ready <- errors.New("cannot initialize macOS HID manager")
		return
	}
	defer C.kt_destroy(ctx)
	ready <- nil
	d.loop(ctx)
}

func (d *driver) loop(ctx *C.kt_context) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	var generation uint64
	var pending error
	for {
		select {
		case <-d.done:
			return
		case req := <-d.requests:
			if req.kind != commandList {
				generation++
				pending = nil
			}
			req.reply <- handle(ctx, req, generation)
		case <-ticker.C:
			pending = d.poll(ctx, generation, pending)
		}
	}
}

func handle(ctx *C.kt_context, req request, generation uint64) response {
	result := response{generation: generation}
	switch req.kind {
	case commandList:
		result.devices, result.err = devices(ctx)
	case commandStart:
		if code := C.kt_start(ctx, C.uint64_t(req.id)); code != 0 {
			result.err = openError(uint32(code), C.GoString(C.kt_error(code)))
		}
	case commandStop:
		C.kt_stop(ctx)
	}
	return result
}

func devices(ctx *C.kt_context) ([]Device, error) {
	count := int(C.kt_refresh(ctx))
	if count < 0 {
		return nil, errors.New("cannot enumerate keyboards: out of memory")
	}
	result := make([]Device, 0)
	seen := make(map[uint64]bool)
	for i := 0; i < count; i++ {
		var info C.kt_device
		if !bool(C.kt_get_device(ctx, C.int(i), &info)) || seen[uint64(info.id)] {
			continue
		}
		device := deviceInfo(info)
		seen[device.ID] = true
		result = append(result, device)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func deviceInfo(info C.kt_device) Device {
	name := C.GoString(&info.name[0])
	if name == "" {
		name = "Unnamed keyboard"
	}
	return Device{ID: uint64(info.id), Name: name, Transport: C.GoString(&info.transport[0]),
		Vendor: uint32(info.vendor), Product: uint32(info.product)}
}

func nativeEvent(raw C.kt_event, generation uint64) Event {
	event := Event{Code: keyboard.Code{Page: uint32(raw.page), Usage: uint32(raw.usage)},
		Down: bool(raw.down), Generation: generation}
	switch int(raw.kind) {
	case eventRemoved:
		event.Err = errors.New("keyboard disconnected; capture released. Reconnect it and refresh the device list")
	case eventOverflow:
		event.Err = errors.New("input queue overflowed; capture released. Start a new capture to continue")
	}
	return event
}

func (d *driver) poll(ctx *C.kt_context, generation uint64, pending error) error {
	if pending != nil {
		select {
		case d.events <- Event{Generation: generation, Err: pending}:
			return nil
		default:
			return pending
		}
	}
	C.kt_pump()
	var raw C.kt_event
	for bool(C.kt_next(ctx, &raw)) {
		select {
		case d.events <- nativeEvent(raw, generation):
		default:
			C.kt_stop(ctx)
			return errors.New("UI input queue overflowed; capture released. Restart capture to continue")
		}
	}
	return nil
}
