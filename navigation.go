package computeruse

import (
	"strings"
)

// Navigate navigates the browser to the specified URL
func (s *Session) Navigate(url string) error {
	// Normalize URL - add https:// if no scheme is present
	normalizedURL := url
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		normalizedURL = "https://" + url
	}

	if err := s.page.Navigate(normalizedURL); err != nil {
		return err
	}

	return s.page.WaitLoad()
}

// GoBack navigates back in browser history
func (s *Session) GoBack() error {
	if err := s.page.NavigateBack(); err != nil {
		return err
	}
	return s.page.WaitLoad()
}

// GoForward navigates forward in browser history
func (s *Session) GoForward() error {
	if err := s.page.NavigateForward(); err != nil {
		return err
	}
	return s.page.WaitLoad()
}

// Search navigates to the configured search engine URL
func (s *Session) Search() error {
	return s.Navigate(s.config.SearchEngineURL)
}
