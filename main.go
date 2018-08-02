package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/eltonlika/posta-tracking/tracker"
	"github.com/ryanuber/columnize"
)

var tableHeadear = "Data|Ngjarja|Zyra|Destinacioni"

func main() {
	descendingPtr := flag.Bool("desc", false, "sort events in descending order")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("No tracking number given")
		os.Exit(1)
	}

	trackingNumber := args[0]

	t := tracker.NewTracker()
	if descendingPtr != nil && *descendingPtr {
		t.EventSortingDirection = tracker.SortDescending
	}
	events := *t.Track(trackingNumber)

	lines := make([]string, len(events)+1)
	lines[0] = tableHeadear
	for i, e := range events {
		lines[i+1] = e.ToString()
	}

	config := columnize.DefaultConfig()
	config.Delim = "|"
	config.Glue = "\t"
	config.Prefix = ""
	config.Empty = ""
	config.NoTrim = false
	result := columnize.Format(lines, config)
	fmt.Println(result)
}
