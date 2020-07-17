package main

import (
	"fmt"
	"os"

	"github.com/eltonlika/posta-tracking/tracker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No tracking number given")
		fmt.Printf("Usage: %s <tracking number>\n", os.Args[0])
		return
	}

	events, err := tracker.NewTracker().Track(os.Args[1])
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	if events.Len() > 0 {
		fmt.Println(events)
	}
}
