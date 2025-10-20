package computeruse

import (
	"context"
	"sync"

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
	config       SessionConfig
	browser      *rod.Browser
	page         *rod.Page
	mainTargetID proto.TargetTargetID
	ctx          context.Context
	navMutex     sync.Mutex // Protects navigation from concurrent popup redirects
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
		config.InitialURL = "about:blank"
	}
	if config.SearchEngineURL == "" {
		config.SearchEngineURL = "https://duckduckgo.com"
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
		Set("disable-sync").
		Set("disable-blink-features", "AutomationControlled")

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
		Width:             config.ScreenWidth,
		Height:            config.ScreenHeight,
		DeviceScaleFactor: 1,
		Mobile:            false,
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
		config:       config,
		browser:      browser,
		page:         page,
		mainTargetID: page.TargetID,
		ctx:          ctx,
	}

	// Handle new tabs/popups by redirecting their content to the main session page.
	// This prevents the session from getting "trapped" when links open in new tabs
	// (e.g., search results with target="_blank").
	//
	// How it works:
	// 1. Listen for new page targets created by our main page
	// 2. Wait for the new page to load and extract its URL
	// 3. Navigate our main session page to that URL
	// 4. Close the popup/new tab
	//
	// This ensures the session always controls a single page, avoiding issues
	// where user actions open new tabs and leave the session orphaned.
	session.startPopupRedirector()

	return session, nil
}

// startPopupRedirector starts a background goroutine that redirects popup/new tab content
// back to the main session page. This prevents the session from becoming orphaned when
// websites open links in new tabs.
func (s *Session) startPopupRedirector() {
	go s.browser.EachEvent(func(e *proto.TargetTargetCreated) {
		// Only handle page-type targets (not workers, service workers, etc.)
		if e.TargetInfo.Type != proto.TargetTargetInfoTypePage {
			return
		}

		// Only handle pages opened by our main page (ignore unrelated tabs)
		if e.TargetInfo.OpenerID != s.mainTargetID {
			return
		}

		// Handle each popup in a separate goroutine to avoid blocking the event listener
		targetID := e.TargetInfo.TargetID
		go s.handlePopup(targetID)
	})()
}

// handlePopup processes a single popup by extracting its URL and navigating the main page to it
func (s *Session) handlePopup(targetID proto.TargetTargetID) {
	// Get the page object for the new tab
	newPage, err := s.browser.PageFromTarget(targetID)
	if err != nil {
		return
	}
	defer func() { _ = newPage.Close() }()

	// Wait for the new page to finish loading so we can capture its final URL
	// (important for redirects or JavaScript-based navigation)
	if err := newPage.WaitLoad(); err != nil {
		return
	}

	// Extract the URL from the new page
	newURL := s.extractURL(newPage)
	if newURL == "" || newURL == "about:blank" {
		return // Ignore blank pages
	}

	// Navigate the main session page to the popup's URL
	// Use mutex to prevent race conditions if multiple popups open simultaneously
	s.navMutex.Lock()
	defer s.navMutex.Unlock()

	if err := s.page.Navigate(newURL); err != nil {
		return
	}
	_ = s.page.WaitLoad()
}

// extractURL attempts to get the URL from a page, trying multiple methods
func (s *Session) extractURL(page *rod.Page) string {
	// First try: Get URL from page info
	info, err := page.Info()
	if err != nil {
		return ""
	}
	url := info.URL

	// Second try: For some pages the info URL may still be blank;
	// fall back to evaluating location.href in the page's JavaScript context
	if url == "" || url == "about:blank" {
		if res, err := page.Eval("() => location.href"); err == nil && res != nil && !res.Value.Nil() {
			if href := res.Value.Str(); href != "" && href != "about:blank" {
				url = href
			}
		}
	}

	return url
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
