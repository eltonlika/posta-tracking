package main

import (
	"fmt"
	"os"
	"flag"
	"encoding/json"
	"github.com/eltonlika/posta-tracking/tracker"
)

func main() {
	descendingPtr := flag.Bool("desc", false, "sort events in descending order")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1{
		fmt.Println("No tracking number given")
		os.Exit(1)
	}

	trackingNumber := args[0]

	t := tracker.NewTracker()
	if  descendingPtr != nil && *descendingPtr {
		t.EventSortingDirection = tracker.SortDescending
	}
	events := t.Track(trackingNumber)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("","")
	encoder.SetEscapeHTML(false)
	encoder.Encode(events)
}
