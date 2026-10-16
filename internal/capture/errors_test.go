package capture

import (
	"strings"
	"testing"
)

func TestCaptureErrorsGiveSpecificRecovery(t *testing.T) {
	cases := []struct {
		code uint32
		want string
	}{
		{returnNotPrivileged, "sudo"},
		{returnNotPermitted, "Input Monitoring"},
		{returnExclusive, "Another application"},
		{returnNoDevice, "Reconnect"},
	}
	for _, test := range cases {
		message := openError(test.code, "native failure").Error()
		if !strings.Contains(message, test.want) || !strings.Contains(message, "native failure") {
			t.Errorf("error 0x%08X lacks specific recovery: %s", test.code, message)
		}
		if test.code == returnNotPrivileged && strings.Contains(message, "Input Monitoring") {
			t.Fatal("privilege errors must not be misidentified as Input Monitoring denial")
		}
	}
}
