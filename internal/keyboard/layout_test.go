package keyboard

import (
	"testing"
)

func TestPresetInventories(t *testing.T) {
	counts := map[string]int{"40": 44, "60": 61, "65": 68, "75": 84, "80": 87, "96": 100, "100": 104}
	for size, expected := range counts {
		layout, err := Lookup(size + "%")
		if err != nil || len(layout.Keys) != expected {
			t.Fatalf("%s: got %d keys, err %v; want %d", size, len(layout.Keys), err, expected)
		}
		assertInventory(t, layout)
	}
}

func assertInventory(t *testing.T, layout Layout) {
	t.Helper()
	seen := make(map[string]bool)
	mapped := make(map[string]bool)
	for _, definition := range Mapping {
		mapped[definition.ID] = true
	}
	for _, row := range layout.Rows {
		for _, id := range row {
			if seen[id] {
				t.Errorf("%s: duplicate key %s", layout.Size, id)
			}
			seen[id] = true
			if !mapped[id] {
				t.Errorf("%s: key %s has no mapped label", layout.Size, id)
			}
		}
	}
}

func TestInvalidSizes(t *testing.T) {
	for _, size := range []string{"", "64", "ansi", "60%%", "-1"} {
		if _, err := Lookup(size); err == nil {
			t.Errorf("accepted invalid size %q", size)
		}
	}
}
