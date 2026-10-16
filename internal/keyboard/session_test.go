package keyboard

import "testing"

func testSession(t *testing.T) *Session {
	t.Helper()
	layout, err := Lookup("60")
	if err != nil {
		t.Fatal(err)
	}
	return NewSession(layout)
}

func TestPhysicalIdentityAndPressEdges(t *testing.T) {
	s := testSession(t)
	a, shift := Code{KeyboardPage, 0x04}, Code{KeyboardPage, 0xe1}
	if !s.Record(a, true) || s.Record(a, true) || s.Record(a, false) {
		t.Fatal("press/release edges were not deduplicated")
	}
	s.Record(shift, true)
	s.Record(a, true)
	if len(s.Tested) != 2 || !s.Tested["a"] || s.Presses != 3 {
		t.Fatalf("shifted/unshifted A should share one physical key: %+v", s)
	}
}

func TestUnknownAndOutOfPresetKeys(t *testing.T) {
	s := testSession(t)
	codes := []Code{{KeyboardPage, 0x3a}, {ConsumerPage, 0x123}, {0xff00, 0x123}}
	for _, code := range codes {
		s.Record(code, true)
		s.Record(code, true)
		s.Record(code, false)
		s.Record(code, true)
	}
	if len(s.Tested) != 0 || len(s.Unknown) != 3 {
		t.Fatalf("unknown HID usages or out-of-preset keys were lost: %+v", s)
	}
	for _, key := range s.Unknown {
		if key.Count != 2 {
			t.Errorf("%v: count = %d; want 2", key.Code, key.Count)
		}
	}
}

func TestModifiersAreDistinct(t *testing.T) {
	s := testSession(t)
	for usage := uint32(0xe0); usage <= 0xe7; usage++ {
		s.Record(Code{KeyboardPage, usage}, true)
	}
	if len(s.Tested) != 8 || !s.ControlAlt() {
		t.Fatal("left/right modifier identities were lost")
	}
	s.Release()
	if s.ControlAlt() || len(s.Tested) != 8 {
		t.Fatal("releasing capture must clear down keys and preserve results")
	}
}

func TestResetDuringShortcut(t *testing.T) {
	s := testSession(t)
	s.Record(Code{KeyboardPage, RightControl}, true)
	s.Record(Code{KeyboardPage, RightAlt}, true)
	s.Record(Code{KeyboardPage, 0x15}, true)
	s.Record(Code{0xff00, 0x01}, true)
	s.Reset()
	if s.Presses != 0 || len(s.Tested) != 0 || len(s.Unknown) != 0 || s.Last != (Code{}) {
		t.Fatal("reset did not clear results")
	}
	if !s.ControlAlt() || s.Record(Code{KeyboardPage, 0x15}, true) {
		t.Fatal("reset must retain physical edges until keys are released")
	}
}

func TestKeypadMultiplyIsDistinctFromShiftEight(t *testing.T) {
	for _, size := range []string{"96", "100"} {
		layout, err := Lookup(size)
		if err != nil {
			t.Fatal(err)
		}
		s := NewSession(layout)
		s.Record(Code{KeyboardPage, 0xe1}, true)
		s.Record(Code{KeyboardPage, 0x25}, true)
		if s.Tested["kpmultiply"] {
			t.Fatal("Shift+8 must not count as keypad multiply")
		}
		s.Record(Code{KeyboardPage, 0x55}, true)
		if !s.Tested["kpmultiply"] || len(s.Unknown) != 0 {
			t.Fatalf("%s%%: standard keypad multiply was not recorded", size)
		}
	}
}
