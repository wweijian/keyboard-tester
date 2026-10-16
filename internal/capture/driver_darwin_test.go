//go:build darwin && cgo

package capture

import "testing"

func TestDriverClosesWithoutStartingCapture(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	if _, err := d.Devices(); err != nil {
		t.Fatal(err)
	}
	if _, err := d.Start(0); err == nil {
		t.Fatal("invalid device ID must fail before requesting permission or seizing input")
	}
	d.Close()
	if _, err := d.Devices(); err == nil {
		t.Fatal("commands after close must return an error")
	}
	if _, open := <-d.Events(); open {
		t.Fatal("event stream should close on shutdown")
	}
}
