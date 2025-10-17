package computeruse

import (
	"runtime"
	"strings"

	"github.com/go-rod/rod/lib/input"
)

// Key mapping from user-friendly names to rod key codes
var keyMap = map[string]input.Key{
	"backspace":    input.Backspace,
	"tab":          input.Tab,
	"return":       input.Enter,
	"enter":        input.Enter,
	"shift":        input.ShiftLeft,
	"shiftleft":    input.ShiftLeft,
	"shiftright":   input.ShiftRight,
	"control":      input.ControlLeft,
	"ctrl":         input.ControlLeft,
	"controlleft":  input.ControlLeft,
	"controlright": input.ControlRight,
	"alt":          input.AltLeft,
	"altleft":      input.AltLeft,
	"altright":     input.AltRight,
	"escape":       input.Escape,
	"esc":          input.Escape,
	"space":        input.Space,
	"pageup":       input.PageUp,
	"pagedown":     input.PageDown,
	"end":          input.End,
	"home":         input.Home,
	"left":         input.ArrowLeft,
	"arrowleft":    input.ArrowLeft,
	"up":           input.ArrowUp,
	"arrowup":      input.ArrowUp,
	"right":        input.ArrowRight,
	"arrowright":   input.ArrowRight,
	"down":         input.ArrowDown,
	"arrowdown":    input.ArrowDown,
	"insert":       input.Insert,
	"delete":       input.Delete,
	"f1":           input.F1,
	"f2":           input.F2,
	"f3":           input.F3,
	"f4":           input.F4,
	"f5":           input.F5,
	"f6":           input.F6,
	"f7":           input.F7,
	"f8":           input.F8,
	"f9":           input.F9,
	"f10":          input.F10,
	"f11":          input.F11,
	"f12":          input.F12,
	"command":      input.MetaLeft,
	"cmd":          input.MetaLeft,
	"meta":         input.MetaLeft,
	"metaleft":     input.MetaLeft,
	"metaright":    input.MetaRight,
}

// TypeText types the given text string
func (s *Session) TypeText(text string) error {
	return s.page.Keyboard.Type([]input.Key(text)...)
}

// TypeTextAt clicks at the specified coordinates and types text
// clearBefore: if true, selects all and deletes before typing
// pressEnter: if true, presses Enter after typing
func (s *Session) TypeTextAt(x, y int, text string, clearBefore, pressEnter bool) error {
	// Click at the position first
	if err := s.ClickAt(x, y); err != nil {
		return err
	}

	// Wait for click to register
	if err := s.page.WaitLoad(); err != nil {
		return err
	}

	// Clear existing text if requested
	if clearBefore {
		// Select all (Cmd+A on macOS, Ctrl+A elsewhere)
		if runtime.GOOS == "darwin" {
			if err := s.Key("Command", "A"); err != nil {
				return err
			}
		} else {
			if err := s.Key("Control", "A"); err != nil {
				return err
			}
		}
		// Delete selected text
		if err := s.Key("Delete"); err != nil {
			return err
		}
	}

	// Type the text
	if err := s.TypeText(text); err != nil {
		return err
	}

	if err := s.page.WaitLoad(); err != nil {
		return err
	}

	// Press Enter if requested
	if pressEnter {
		if err := s.Key("Enter"); err != nil {
			return err
		}
	}

	return s.page.WaitLoad()
}

// Key presses a key or key combination
// Examples: Key("Enter"), Key("Control", "C"), Key("Alt", "F4")
func (s *Session) Key(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	// Normalize and convert keys
	rodKeys := make([]input.Key, 0, len(keys))
	for _, key := range keys {
		normalizedKey := strings.ToLower(strings.TrimSpace(key))

		// Check if it's a special key
		if rodKey, ok := keyMap[normalizedKey]; ok {
			rodKeys = append(rodKeys, rodKey)
		} else {
			// Treat as literal character(s)
			for _, ch := range key {
				rodKeys = append(rodKeys, input.Key(ch))
			}
		}
	}

	if len(rodKeys) == 0 {
		return nil
	}

	// If only one key, just press and release it
	if len(rodKeys) == 1 {
		return s.page.Keyboard.Type(rodKeys...)
	}

	// For key combinations: press all modifier keys except the last
	for i := 0; i < len(rodKeys)-1; i++ {
		if err := s.page.Keyboard.Press(rodKeys[i]); err != nil {
			return err
		}
	}

	// Press and release the last key (the main key)
	if err := s.page.Keyboard.Type(rodKeys[len(rodKeys)-1]); err != nil {
		return err
	}

	// Release modifier keys in reverse order
	for i := len(rodKeys) - 2; i >= 0; i-- {
		if err := s.page.Keyboard.Release(rodKeys[i]); err != nil {
			return err
		}
	}

	return nil
}
