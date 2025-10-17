package computeruse

import (
	"github.com/go-rod/rod/lib/proto"
)

// ClickAt performs a left click at the specified coordinates
func (s *Session) ClickAt(x, y int) error {
	actualX, actualY := s.normalizeCoords(x, y)
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)}); err != nil {
		return err
	}
	return s.page.Mouse.Click(proto.InputMouseButtonLeft, 1)
}

// RightClickAt performs a right click at the specified coordinates
func (s *Session) RightClickAt(x, y int) error {
	actualX, actualY := s.normalizeCoords(x, y)
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)}); err != nil {
		return err
	}
	return s.page.Mouse.Click(proto.InputMouseButtonRight, 1)
}

// MiddleClickAt performs a middle click at the specified coordinates
func (s *Session) MiddleClickAt(x, y int) error {
	actualX, actualY := s.normalizeCoords(x, y)
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)}); err != nil {
		return err
	}
	return s.page.Mouse.Click(proto.InputMouseButtonMiddle, 1)
}

// DoubleClickAt performs a double click at the specified coordinates
func (s *Session) DoubleClickAt(x, y int) error {
	actualX, actualY := s.normalizeCoords(x, y)
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)}); err != nil {
		return err
	}
	return s.page.Mouse.Click(proto.InputMouseButtonLeft, 2)
}

// TripleClickAt performs a triple click at the specified coordinates
func (s *Session) TripleClickAt(x, y int) error {
	actualX, actualY := s.normalizeCoords(x, y)
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)}); err != nil {
		return err
	}
	return s.page.Mouse.Click(proto.InputMouseButtonLeft, 3)
}

// MouseDown presses the left mouse button at the specified coordinates
func (s *Session) MouseDown(x, y int) error {
	actualX, actualY := s.normalizeCoords(x, y)
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)}); err != nil {
		return err
	}
	return s.page.Mouse.Down(proto.InputMouseButtonLeft, 1)
}

// MouseUp releases the left mouse button at the specified coordinates
func (s *Session) MouseUp(x, y int) error {
	actualX, actualY := s.normalizeCoords(x, y)
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)}); err != nil {
		return err
	}
	return s.page.Mouse.Up(proto.InputMouseButtonLeft, 1)
}

// MouseMove moves the cursor to the specified coordinates
func (s *Session) MouseMove(x, y int) error {
	actualX, actualY := s.normalizeCoords(x, y)
	return s.page.Mouse.MoveTo(proto.Point{X: float64(actualX), Y: float64(actualY)})
}

// HoverAt is an alias for MouseMove, hovers at the specified coordinates
func (s *Session) HoverAt(x, y int) error {
	return s.MouseMove(x, y)
}

// ClickDrag performs a click and drag operation from one coordinate to another
func (s *Session) ClickDrag(fromX, fromY, toX, toY int) error {
	actualFromX, actualFromY := s.normalizeCoords(fromX, fromY)
	actualToX, actualToY := s.normalizeCoords(toX, toY)

	// Move to start position
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualFromX), Y: float64(actualFromY)}); err != nil {
		return err
	}

	// Press mouse button
	if err := s.page.Mouse.Down(proto.InputMouseButtonLeft, 1); err != nil {
		return err
	}

	// Move to end position
	if err := s.page.Mouse.MoveTo(proto.Point{X: float64(actualToX), Y: float64(actualToY)}); err != nil {
		return err
	}

	// Release mouse button
	return s.page.Mouse.Up(proto.InputMouseButtonLeft, 1)
}
