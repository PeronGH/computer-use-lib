package main

import (
	"context"
	"fmt"
	"log"
	"time"

	computeruse "github.com/PeronGH/computer-use-lib"
)

func main() {
	var err error

	throwIfErr := func() {
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create a new browser session
	session, err := computeruse.NewSession(context.Background(), computeruse.SessionConfig{
		NormalizeCoordinates: true, // Required for Gemini
	})
	throwIfErr()
	defer session.Close()

	fmt.Println("Browser opened successfully!")

	// Accept all cookies
	err = session.ClickAt(571, 819)
	throwIfErr()
	fmt.Println("Cookies accepted")
	time.Sleep(2 * time.Second)

	// Search current time
	err = session.TypeTextAt(494, 444, "current time", true, true)
	throwIfErr()
	fmt.Println("Searching current time")
	time.Sleep(2 * time.Second)

	// Take a screenshot
	screenshot, err := session.Screenshot()
	throwIfErr()
	fmt.Printf("Screenshot taken: %d bytes\n", len(screenshot))
	time.Sleep(2 * time.Second)

	fmt.Println("Example completed successfully!")
}
