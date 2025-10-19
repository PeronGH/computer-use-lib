package computeruse

import (
	"fmt"

	"github.com/go-rod/rod/lib/proto"
)

// Scroll scrolls the page in the specified direction by the given amount
// direction: "up", "down", "left", "right"
// amount: scroll distance (in pixels if not normalized, or 0-999 if normalized)
func (s *Session) Scroll(direction string, amount int) error {
	// Normalize the amount if needed
	actualAmount := amount
	if s.config.NormalizeCoordinates {
		// For vertical scroll, normalize against height; for horizontal, against width
		if direction == "up" || direction == "down" {
			actualAmount = (amount * s.config.ScreenHeight) / 1000
		} else {
			actualAmount = (amount * s.config.ScreenWidth) / 1000
		}
	}

	// Determine scroll direction
	var deltaX, deltaY float64
	switch direction {
	case "up":
		deltaY = -float64(actualAmount)
	case "down":
		deltaY = float64(actualAmount)
	case "left":
		deltaX = -float64(actualAmount)
	case "right":
		deltaX = float64(actualAmount)
	default:
		return fmt.Errorf("invalid scroll direction: %s (must be up, down, left, or right)", direction)
	}

	return s.page.Mouse.Scroll(deltaX, deltaY, 1)
}

// ScrollAt scrolls at a specific location on the page
// x, y: coordinates to scroll at
// direction: "up", "down", "left", "right"
// magnitude: scroll amount (0-999 if normalized, pixels otherwise)
func (s *Session) ScrollAt(x, y int, direction string, magnitude int) error {
	actualX, actualY := s.normalizeCoords(x, y)

	// Normalize magnitude if needed
	actualMagnitude := magnitude
	if s.config.NormalizeCoordinates {
		if direction == "up" || direction == "down" {
			actualMagnitude = (magnitude * s.config.ScreenHeight) / 1000
		} else {
			actualMagnitude = (magnitude * s.config.ScreenWidth) / 1000
		}
	}

	// Move mouse to the position
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)}); err != nil {
		return err
	}

	// Determine scroll deltas
	var deltaX, deltaY float64
	switch direction {
	case "up":
		deltaY = -float64(actualMagnitude)
	case "down":
		deltaY = float64(actualMagnitude)
	case "left":
		deltaX = -float64(actualMagnitude)
	case "right":
		deltaX = float64(actualMagnitude)
	default:
		return fmt.Errorf("invalid scroll direction: %s (must be up, down, left, or right)", direction)
	}

	// Perform mouse wheel scroll
	return s.page.Mouse.Scroll(deltaX, deltaY, 0)
}
