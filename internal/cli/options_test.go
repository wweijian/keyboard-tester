package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpWithoutCapture(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}, {"-h"}} {
		var output, errors bytes.Buffer
		if err := Run(args, &output, &errors); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(output.String(), "Usage:") || errors.Len() != 0 {
			t.Fatal("help should succeed without opening native capture")
		}
	}
}

func TestInvalidArguments(t *testing.T) {
	for _, args := range [][]string{{"--keyboard", "64"}, {"--keyboard"}, {"--bad"}, {"60"}} {
		var output, errors bytes.Buffer
		if err := Run(args, &output, &errors); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
}

func TestKeyboardFlag(t *testing.T) {
	var output bytes.Buffer
	opts, err := parse([]string{"--keyboard", "65%"}, &output)
	if err != nil || opts.keyboard != "65%" || opts.list {
		t.Fatalf("unexpected options %+v: %v", opts, err)
	}
}

func TestCaptureRequiresAdminBeforeEnteringUI(t *testing.T) {
	if err := requireAdmin("65", 0); err != nil {
		t.Fatal(err)
	}
	err := requireAdmin("65", 501)
	if err == nil || !strings.Contains(err.Error(), "sudo ./bin/keyboard-tester --keyboard 65") {
		t.Fatalf("non-admin startup needs an actionable command, got %v", err)
	}
}
