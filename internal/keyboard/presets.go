package keyboard

// Each preset lists the physical key usages expected in its reference layout.
// Percentages are size categories, not a universal key inventory; see README.
var presets = map[string][][]string{
	"40": {
		{"esc", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "backspace"},
		{"tab", "a", "s", "d", "f", "g", "h", "j", "k", "l", "semicolon", "enter"},
		{"lshift", "z", "x", "c", "v", "b", "n", "m", "comma", "period", "slash", "rshift"},
		{"lctrl", "lgui", "lalt", "space", "ralt", "rgui", "menu", "rctrl"},
	},
	"60": {
		{"esc", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "minus", "equal", "backspace"},
		{"tab", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "lbracket", "rbracket", "backslash"},
		{"caps", "a", "s", "d", "f", "g", "h", "j", "k", "l", "semicolon", "quote", "enter"},
		{"lshift", "z", "x", "c", "v", "b", "n", "m", "comma", "period", "slash", "rshift"},
		{"lctrl", "lgui", "lalt", "space", "ralt", "rgui", "menu", "rctrl"},
	},
	"65": {
		{"esc", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "minus", "equal", "backspace", "delete"},
		{"tab", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "lbracket", "rbracket", "backslash", "pgup"},
		{"caps", "a", "s", "d", "f", "g", "h", "j", "k", "l", "semicolon", "quote", "enter", "pgdn"},
		{"lshift", "z", "x", "c", "v", "b", "n", "m", "comma", "period", "slash", "rshift", "up"},
		{"lctrl", "lgui", "lalt", "space", "ralt", "rgui", "menu", "rctrl", "left", "down", "right"},
	},
	"75": {
		{"esc", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "insert"},
		{"grave", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "minus", "equal", "backspace", "delete"},
		{"tab", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "lbracket", "rbracket", "backslash", "pgup"},
		{"caps", "a", "s", "d", "f", "g", "h", "j", "k", "l", "semicolon", "quote", "enter", "pgdn"},
		{"lshift", "z", "x", "c", "v", "b", "n", "m", "comma", "period", "slash", "rshift", "up", "home"},
		{"lctrl", "lgui", "lalt", "space", "ralt", "rgui", "menu", "rctrl", "left", "down", "right", "end"},
	},
	"80": {
		{"esc", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "print", "scroll", "pause"},
		{"grave", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "minus", "equal", "backspace", "insert", "home", "pgup"},
		{"tab", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "lbracket", "rbracket", "backslash", "delete", "end", "pgdn"},
		{"caps", "a", "s", "d", "f", "g", "h", "j", "k", "l", "semicolon", "quote", "enter"},
		{"lshift", "z", "x", "c", "v", "b", "n", "m", "comma", "period", "slash", "rshift", "up"},
		{"lctrl", "lgui", "lalt", "space", "ralt", "rgui", "menu", "rctrl", "left", "down", "right"},
	},
	"96": {
		{"esc", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "delete"},
		{"grave", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "minus", "equal", "backspace", "numlock", "kpdivide", "kpmultiply", "kpminus"},
		{"tab", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "lbracket", "rbracket", "backslash", "kp7", "kp8", "kp9", "kpplus"},
		{"caps", "a", "s", "d", "f", "g", "h", "j", "k", "l", "semicolon", "quote", "enter", "kp4", "kp5", "kp6", "pgup"},
		{"lshift", "z", "x", "c", "v", "b", "n", "m", "comma", "period", "slash", "rshift", "up", "kp1", "kp2", "kp3", "pgdn"},
		{"lctrl", "lgui", "lalt", "space", "ralt", "rgui", "menu", "rctrl", "left", "down", "right", "home", "end", "kp0", "kpperiod", "kpenter"},
	},
	"100": {
		{"esc", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "print", "scroll", "pause"},
		{"grave", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "minus", "equal", "backspace", "insert", "home", "pgup", "numlock", "kpdivide", "kpmultiply", "kpminus"},
		{"tab", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "lbracket", "rbracket", "backslash", "delete", "end", "pgdn", "kp7", "kp8", "kp9", "kpplus"},
		{"caps", "a", "s", "d", "f", "g", "h", "j", "k", "l", "semicolon", "quote", "enter", "kp4", "kp5", "kp6"},
		{"lshift", "z", "x", "c", "v", "b", "n", "m", "comma", "period", "slash", "rshift", "up", "kp1", "kp2", "kp3", "kpenter"},
		{"lctrl", "lgui", "lalt", "space", "ralt", "rgui", "menu", "rctrl", "left", "down", "right", "kp0", "kpperiod"},
	},
}
