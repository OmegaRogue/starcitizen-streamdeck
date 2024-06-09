package util

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func DiscardErrorOnly(err error) {
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

//goland:noinspection GoUnusedExportedFunction
func DiscardError[T any](first T, err error) T {
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return first
}

var KeyLoc = map[string]map[string]string{
	"de": {
		"minus":        "ß",
		"equals":       "´",
		"grave":        "^",
		"y":            "z",
		"z":            "y",
		"bracketleft":  "Ü",
		"bracketright": "+",
		"backslash":    "#",
		"comma":        ",",
		"period":       ".",
		"slash":        "-",
		"semicolon":    "Ö",
		"apostrophe":   "Ä",
	},
	"en": {
		"grave":      "`",
		"minus":      "-",
		"equals":     "=",
		"bracketl":   "[",
		"bracketr":   "]",
		"backslash":  "\\",
		"semicolon":  ":",
		"apostrophe": "'",
		"comma":      ",",
		"period":     ".",
		"slash":      "/",
	},
}

var ScKeyNames = map[string]string{
	"lalt":        "AltL",
	"ralt":        "AltR",
	"lshift":      "ShiftL",
	"rshift":      "ShiftR",
	"lctrl":       "CtrlL",
	"rctrl":       "CtrlR",
	"f1":          "F1",
	"f2":          "F2",
	"f3":          "F3",
	"f4":          "F4",
	"f5":          "F5",
	"f6":          "F6",
	"f7":          "F7",
	"f8":          "F8",
	"f9":          "F9",
	"f10":         "F10",
	"f11":         "F11",
	"f12":         "F12",
	"f13":         "F13",
	"f14":         "F14",
	"f15":         "F15",
	"numlock":     "Num_Lock",
	"np_divide":   "KP_Divide",
	"np_multiply": "KP_Multiply",
	"np_subtract": "KP_Subtract",
	"np_add":      "KP_Add",
	"np_period":   "KP_Period",
	"np_enter":    "KP_Enter",
	"np_0":        "KP_0",
	"np_1":        "KP_1",
	"np_2":        "KP_2",
	"np_3":        "KP_3",
	"np_4":        "KP_4",
	"np_5":        "KP_5",
	"np_6":        "KP_6",
	"np_7":        "KP_7",
	"np_8":        "KP_8",
	"np_9":        "KP_9",
	"0":           "zero",
	"1":           "one",
	"2":           "two",
	"3":           "three",
	"4":           "four",
	"5":           "five",
	"6":           "six",
	"7":           "seven",
	"8":           "eight",
	"9":           "nine",
	"insert":      "Insert",
	"home":        "Home",
	"delete":      "Delete",
	"end":         "End",
	"pgup":        "PageUp",
	"pgdown":      "PageDown",
	"pgdn":        "PageDown",
	"print":       "Print",
	"scrolllock":  "Scroll_Lock",
	"pause":       "Pause",
	"up":          "Up",
	"down":        "Down",
	"left":        "Left",
	"right":       "Right",
	"escape":      "escape",
	"minus":       "minus",
	"equals":      "equals",
	"grave":       "grave",
	"underline":   "underscore",
	"backspace":   "BackSpace",
	"tab":         "Tab",
	"lbracket":    "bracketleft",
	"rbracket":    "bracketright",
	"enter":       "Return",
	"capslock":    "Caps_lock",
	"colon":       "colon",
	"backslash":   "backslash",
	"comma":       "comma",
	"period":      "period",
	"slash":       "slash",
	"space":       "space",
	"semicolon":   "semicolon",
	"apostrophe":  "apostrophe",
}

func LocalizeKeyString(keyboard, language string) string {
	return strings.Join(lo.Map(strings.Split(keyboard, "+"), func(item string, index int) string {
		k, ok := KeyLoc[language][keyboard]
		if ok {
			return k
		}
		return item
	}), "+")
}

func FromScKey(keyboard string) string {
	return strings.Join(lo.Map(strings.Split(keyboard, "+"), func(item string, index int) string {
		k, ok := ScKeyNames[keyboard]
		if ok {
			return k
		}
		return item
	}), "+")
}
