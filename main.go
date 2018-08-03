package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/eltonlika/posta-tracking/tracker"
	"github.com/ryanuber/columnize"
)

func main() {
	descendingPtr := flag.Bool("descending", false, "sort events in descending order")
	noHeaderPtr := flag.Bool("no-header", false, "do not print header line")
	timeoutPtr := flag.Uint("timeout", 8, "number of seconds to wait for response from tracking service")
	delimiterPtr := flag.String("delimiter", "  ", "string to use as column delimiter (separator)")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("No tracking number given")
		os.Exit(1)
	}

	t := tracker.NewTracker()
	t.SortEventsDescending = *descendingPtr
	t.SetRequestTimeout(time.Second * time.Duration(*timeoutPtr))

	trackingNumber := args[0]
	events, err := t.Track(trackingNumber)
	if err != nil {
		panic(err)
	}

	eventsTable := formatTrackingEventsAsTable(*events, *noHeaderPtr, *delimiterPtr)
	fmt.Println(eventsTable)
}

func formatTrackingEventsAsTable(events tracker.Events, noHeader bool, delimiter string) string {
	var table []string
	var rows []string

	if noHeader {
		table = make([]string, len(events))
		rows = table
	} else {
		table = make([]string, len(events)+1)
		table[0] = "#|Kodi|Data|Ngjarja|Zyra|Destinacioni"
		rows = table[1:]
	}

	sp := "|"
	for i, e := range events {
		num := strconv.Itoa(i + 1)
		date := e.Date.Format("2006-01-02 15:04 PM")
		rows[i] = num + sp + e.TrackingNumber + sp + date + sp + e.Description + sp + e.Location + sp + e.Destination
	}

	config := columnize.DefaultConfig()
	config.Delim = sp
	config.Glue = delimiter
	config.Prefix = ""
	config.Empty = ""
	config.NoTrim = false
	return columnize.Format(table, config)
}
