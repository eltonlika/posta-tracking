package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/eltonlika/posta-tracking/formatter"
	"github.com/eltonlika/posta-tracking/tracker"
)

// Options cli options holder struct
type Options struct {
	SortDescending bool
	NoHeader       bool
	Timeout        uint
	Delimiter      string
	TrackingNumber string
}

// ParseOptions try to parse options from command line arguments
func ParseOptions() (Options, error) {
	descendingPtr := flag.Bool("descending", false, "sort events in descending order")
	noHeaderPtr := flag.Bool("no-header", false, "do not print header line")
	timeoutPtr := flag.Uint("timeout", 8, "number of seconds to wait for response from tracking service")
	delimiterPtr := flag.String("delimiter", "  ", "string to use as column delimiter (separator)")
	flag.Usage = PrintUsage
	flag.Parse()

	if flag.NArg() == 0 {
		return Options{}, errors.New("No tracking number given")
	}

	return Options{
		SortDescending: *descendingPtr,
		NoHeader:       *noHeaderPtr,
		Timeout:        *timeoutPtr,
		Delimiter:      *delimiterPtr,
		TrackingNumber: flag.Args()[0],
	}, nil
}

// PrintUsage prints cli usage help text
func PrintUsage() {
	fmt.Printf("Usage: %s [OPTIONS] trackingNumber", os.Args[0])
	fmt.Println("")
	flag.PrintDefaults()
}

// Run cli client
func Run() {
	opts, err := ParseOptions()
	if err != nil {
		fmt.Println(err.Error())
		PrintUsage()
		os.Exit(1)
	}

	t := tracker.NewTracker()
	t.SortEventsDescending = opts.SortDescending
	t.SetRequestTimeout(opts.Timeout)

	events, err := t.Track(opts.TrackingNumber)
	if err != nil {
		panic(err)
	}

	formatter := formatter.NewEventsFormatter()
	formatter.NoHeader = opts.NoHeader
	formatter.Delimiter = opts.Delimiter
	err = formatter.Print(events, os.Stdout)
	if err != nil {
		panic(err)
	}
}
