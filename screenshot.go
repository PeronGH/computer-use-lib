package computeruse

import (
	"time"

	"github.com/go-rod/rod/lib/proto"
)

// Screenshot captures the current browser viewport as a PNG image
// Returns the PNG image data as a byte slice
func (s *Session) Screenshot() ([]byte, error) {
	// Wait for the page to be fully loaded
	if err := s.page.WaitLoad(); err != nil {
		return nil, err
	}

	// Add a small delay to ensure rendering is complete
	time.Sleep(500 * time.Millisecond)

	// Capture screenshot of the viewport only (not full page)
	screenshot, err := s.page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatPng,
		Quality: nil, // PNG doesn't use quality parameter
	})

	return screenshot, err
}
