package computeruse

import "github.com/go-rod/rod/lib/proto"

// Screenshot captures the current browser viewport as a PNG image
// Returns the PNG image data as a byte slice
func (s *Session) Screenshot() ([]byte, error) {
	if err := s.page.WaitLoad(); err != nil {
		return nil, err
	}

	// Capture screenshot of the viewport only (not full page)
	screenshot, err := s.page.Screenshot(false, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatPng,
		Quality: nil, // PNG doesn't use quality parameter
	})

	return screenshot, err
}
