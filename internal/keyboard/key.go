package keyboard

import "fmt"

// Code identifies a physical HID usage independently of case and input language.
type Code struct {
	Page  uint32
	Usage uint32
}

type Definition struct {
	ID    string
	Label string
}

func (c Code) String() string {
	return fmt.Sprintf("%04X:%04X", c.Page, c.Usage)
}

func Describe(c Code) string {
	if key, ok := Mapping[c]; ok {
		return key.Label
	}
	return "Unknown"
}

func Label(id string) string {
	for _, key := range Mapping {
		if key.ID == id {
			return key.Label
		}
	}
	return id
}
