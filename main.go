package main

import (
	"fmt"
	"os"

	"github.com/eltonlika/posta-tracking/tracker"
)

func cli() int {
	var log = os.Stderr.WriteString

	if len(os.Args) < 2 {
		log(fmt.Sprintf("Error: No tracking number given\nUsage: %s <tracking number>\n", os.Args[0]))
		return 1
	}

	events, err := tracker.NewTracker().Track(os.Args[1])
	if err != nil {
		log(fmt.Sprintf("Error: %s\n", err.Error()))
		return 1
	}

	if events.Len() > 0 {
		fmt.Println(events)
	}

	return 0
}

func main() {
	os.Exit(cli())
}
