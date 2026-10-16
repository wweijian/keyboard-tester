package keyboard

type Unknown struct {
	Code  Code
	Count int
}

type Session struct {
	Layout  Layout
	Tested  map[string]bool
	Unknown []Unknown
	Last    Code
	Presses int
	down    map[Code]bool
}

func NewSession(layout Layout) *Session {
	return &Session{Layout: layout, Tested: make(map[string]bool), down: make(map[Code]bool)}
}

// Press/release bookkeeping prevents duplicate downs; no hold state is tested.
func (s *Session) Record(code Code, down bool) bool {
	previous := s.down[code]
	if !down {
		delete(s.down, code)
		return false
	}
	if previous {
		return false
	}
	s.down[code] = true
	s.Last, s.Presses = code, s.Presses+1
	key, known := Mapping[code]
	if known && s.Layout.Keys[key.ID] {
		s.Tested[key.ID] = true
		return true
	}
	s.recordUnknown(code)
	return true
}

func (s *Session) recordUnknown(code Code) {
	for i := range s.Unknown {
		if s.Unknown[i].Code == code {
			s.Unknown[i].Count++
			return
		}
	}
	s.Unknown = append(s.Unknown, Unknown{Code: code, Count: 1})
}

func (s *Session) ControlAlt() bool {
	ctrl := s.down[Code{KeyboardPage, LeftControl}] || s.down[Code{KeyboardPage, RightControl}]
	alt := s.down[Code{KeyboardPage, LeftAlt}] || s.down[Code{KeyboardPage, RightAlt}]
	return ctrl && alt
}

func (s *Session) Reset() {
	s.Tested = make(map[string]bool)
	s.Unknown = nil
	s.Last, s.Presses = Code{}, 0
}

func (s *Session) Release() {
	clear(s.down)
}
