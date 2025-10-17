package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	computeruse "github.com/PeronGH/computer-use-lib"
)

func main() {
	// Create a new browser session
	session, err := computeruse.NewSession(context.Background(), computeruse.SessionConfig{
		ScreenWidth:          1440,
		ScreenHeight:         900,
		NormalizeCoordinates: true, // Use 0-999 grid
		InitialURL:           "https://www.google.com",
		Headless:             false, // Set to true for headless mode
	})
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	fmt.Println("Browser opened successfully!")

	// Example: Take a screenshot
	screenshot, err := session.Screenshot()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Screenshot taken: %d bytes\n", len(screenshot))

	// Save screenshot to file
	if err := os.WriteFile("screenshot.png", screenshot, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Screenshot saved to screenshot.png")

	// Example: Navigate to a different page
	if err := session.Navigate("https://example.com"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Navigated to example.com")

	// Get current URL
	url, err := session.GetURL()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current URL: %s\n", url)

	// Example: Click at center of screen (500, 500 on 0-999 grid)
	if err := session.ClickAt(500, 500); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Clicked at center")

	// Example: Type text
	if err := session.TypeText("Hello from computer-use-lib!"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Typed text")

	// Wait a bit before closing (using standard time.Sleep)
	time.Sleep(2 * time.Second)

	fmt.Println("Example completed successfully!")
}
