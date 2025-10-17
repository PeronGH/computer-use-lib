# Computer Use Library

A Go library for browser-based computer use automation, designed for LLM agents (Claude Computer Use, Google Gemini, etc.). Built on [go-rod](https://github.com/go-rod/rod) for robust browser control.

## Features

- **Unified API**: Single set of commands that work for both Claude and Gemini with minimal adaptation
- **Flexible Coordinate System**: Choose between normalized (0-999 grid) or pixel-based coordinates
- **Idiomatic Go**: Proper error handling and clean interface design
- **Comprehensive Actions**: Supports clicking, typing, scrolling, dragging, keyboard shortcuts, and more
- **Screenshot Capability**: Capture browser state for visual feedback to LLMs
- **Session Management**: Easy browser lifecycle management with context support

## Installation

```bash
go get github.com/PeronGH/computer-use-lib
```

## Quick Start

```go
package main

import (
    "context"

    computeruse "github.com/PeronGH/computer-use-lib"
)

func main() {
    // Create a new browser session
    session, err := computeruse.NewSession(context.Background(), computeruse.SessionConfig{
        ScreenWidth:          1440,
        ScreenHeight:         900,
        NormalizeCoordinates: true, // Use 0-999 grid
        InitialURL:           "https://www.google.com",
    })
    if err != nil {
        panic(err)
    }
    defer session.Close()

    // Use the session
    session.Navigate("https://example.com")
    session.ClickAt(500, 500)
    session.TypeText("Hello, World!")
    screenshot, _ := session.Screenshot()
    _ = screenshot
}
```

## API Reference

### Session Configuration

```go
type SessionConfig struct {
    ScreenWidth          int    // Browser viewport width
    ScreenHeight         int    // Browser viewport height
    NormalizeCoordinates bool   // If true, use 0-999 grid; if false, use pixels
    InitialURL           string // Starting URL (default: "https://www.google.com")
    SearchEngineURL      string // URL for Search() action (default: "https://www.google.com")
    Headless             bool   // Run browser in headless mode
}
```

### Available Commands

All methods return `error` for proper error handling.

| Method | Signature | Claude Mapping | Gemini Mapping |
|--------|-----------|----------------|----------------|
| `Screenshot` | `Screenshot() ([]byte, error)` | `screenshot` | N/A (call separately) |
| `ClickAt` | `ClickAt(x, y int) error` | `left_click` | `click_at` |
| `RightClickAt` | `RightClickAt(x, y int) error` | `right_click` | N/A |
| `MiddleClickAt` | `MiddleClickAt(x, y int) error` | `middle_click` | N/A |
| `DoubleClickAt` | `DoubleClickAt(x, y int) error` | `double_click` | N/A |
| `TripleClickAt` | `TripleClickAt(x, y int) error` | `triple_click` | N/A |
| `MouseDown` | `MouseDown(x, y int) error` | `left_mouse_down` | N/A |
| `MouseUp` | `MouseUp(x, y int) error` | `left_mouse_up` | N/A |
| `MouseMove` | `MouseMove(x, y int) error` | `mouse_move` | N/A |
| `HoverAt` | `HoverAt(x, y int) error` | `mouse_move` | `hover_at` |
| `ClickDrag` | `ClickDrag(fromX, fromY, toX, toY int) error` | `left_click_drag` | `drag_and_drop` |
| `TypeText` | `TypeText(text string) error` | `type` | N/A |
| `TypeTextAt` | `TypeTextAt(x, y int, text string, clearBefore, pressEnter bool) error` | `left_click` + `type` + `key` | `type_text_at` |
| `Key` | `Key(keys ...string) error` | `key` | `key_combination` |
| `Scroll` | `Scroll(direction string, amount int) error` | `scroll` | `scroll_document` |
| `ScrollAt` | `ScrollAt(x, y int, direction string, magnitude int) error` | `mouse_move` + `scroll` | `scroll_at` |
| `Navigate` | `Navigate(url string) error` | N/A | `navigate` |
| `GoBack` | `GoBack() error` | `key` ("Alt+Left") | `go_back` |
| `GoForward` | `GoForward() error` | `key` ("Alt+Right") | `go_forward` |
| `Search` | `Search() error` | N/A | `search` |
| `GetURL` | `GetURL() (string, error)` | N/A | N/A |
| `Close` | `Close() error` | N/A | N/A |

## Architecture

The library provides a unified API layer that translates high-level actions into go-rod browser commands:

```
LLM Agent (Claude/Gemini)
         ↓
Computer Use Library API
         ↓
go-rod (Browser Control)
         ↓
Chrome/Chromium Browser
```

## Key Design Decisions

1. **Coordinate Normalization**: Optional 0-999 grid system allows LLMs to work with consistent coordinates across different screen sizes
2. **Unified API**: Single set of commands that map cleanly to both Claude and Gemini interfaces
3. **Error Handling**: Each method returns an error for proper Go error handling
4. **Minimal Adaptation**: Developers need minimal effort to support both LLM platforms
