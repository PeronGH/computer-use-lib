package computeruse

import (
	"context"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// SessionConfig holds configuration for a browser session
type SessionConfig struct {
	ScreenWidth          int    // Browser viewport width
	ScreenHeight         int    // Browser viewport height
	NormalizeCoordinates bool   // If true, use 0-999 grid; if false, use pixels
	InitialURL           string // Starting URL (default: "https://www.google.com")
	SearchEngineURL      string // URL for Search() action (default: "https://www.google.com")
	Headless             bool   // Run browser in headless mode
}

// Session represents a browser automation session
type Session struct {
	config  SessionConfig
	browser *rod.Browser
	page    *rod.Page
	ctx     context.Context
}

// NewSession creates a new browser session with the given configuration
func NewSession(ctx context.Context, config SessionConfig) (*Session, error) {
	// Set defaults
	if config.ScreenWidth == 0 {
		config.ScreenWidth = 1440
	}
	if config.ScreenHeight == 0 {
		config.ScreenHeight = 900
	}
	if config.InitialURL == "" {
		config.InitialURL = "https://www.google.com"
	}
	if config.SearchEngineURL == "" {
		config.SearchEngineURL = "https://www.google.com"
	}

	// Launch browser
	l := launcher.New().
		Headless(config.Headless).
		Set("disable-extensions").
		Set("disable-file-system").
		Set("disable-plugins").
		Set("disable-dev-shm-usage").
		Set("disable-background-networking").
		Set("disable-default-apps").
		Set("disable-sync")

	url, err := l.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(url).Context(ctx)
	if err := browser.Connect(); err != nil {
		return nil, err
	}

	// Create page with viewport
	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, err
	}

	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:  config.ScreenWidth,
		Height: config.ScreenHeight,
		DeviceScaleFactor: 1,
		Mobile: false,
	}); err != nil {
		return nil, err
	}

	// Navigate to initial URL
	if err := page.Navigate(config.InitialURL); err != nil {
		return nil, err
	}

	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	session := &Session{
		config:  config,
		browser: browser,
		page:    page,
		ctx:     ctx,
	}

	return session, nil
}

// Close closes the browser session
func (s *Session) Close() error {
	if s.browser != nil {
		return s.browser.Close()
	}
	return nil
}

// normalizeCoords converts coordinates based on config
func (s *Session) normalizeCoords(x, y int) (int, int) {
	if s.config.NormalizeCoordinates {
		// Convert from 0-999 grid to actual pixels
		actualX := (x * s.config.ScreenWidth) / 1000
		actualY := (y * s.config.ScreenHeight) / 1000
		return actualX, actualY
	}
	return x, y
}

// GetURL returns the current page URL
func (s *Session) GetURL() (string, error) {
	info, err := s.page.Info()
	if err != nil {
		return "", err
	}
	return info.URL, nil
}

// Wait pauses execution for the specified duration
func (s *Session) Wait(duration time.Duration) error {
	time.Sleep(duration)
	return nil
}
