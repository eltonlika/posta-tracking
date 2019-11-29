package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/eltonlika/posta-tracking/formatter"
	"github.com/eltonlika/posta-tracking/tracker"
)

const (
	// NoError status returned when program compeltes successfully
	NoError = iota

	// ErrorInvalidArgs status returned when program has no valid arguments
	ErrorInvalidArgs

	// ErrorTrackerFailed status returned when tracker execution failed
	ErrorTrackerFailed

	// ErrorNoEventsFound status returned when no events found for given tracking number
	ErrorNoEventsFound

	// ErrorOther status returned for any other error
	ErrorOther
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

func print(err error) {
	fmt.Println("Error: " + err.Error())
}

// Run cli client
func Run() int {
	opts, err := ParseOptions()
	if err != nil {
		print(err)
		PrintUsage()
		return ErrorInvalidArgs
	}

	t := tracker.NewTracker()
	t.SortEventsDescending = opts.SortDescending
	t.SetRequestTimeout(opts.Timeout)

	events, err := t.Track(opts.TrackingNumber)
	if err != nil {
		print(err)
		return ErrorTrackerFailed
	}

	if len(events) == 0 {
		return ErrorNoEventsFound
	}

	formatter := formatter.NewEventsFormatter()
	formatter.NoHeader = opts.NoHeader
	formatter.Delimiter = opts.Delimiter
	err = formatter.Print(events, os.Stdout)
	if err != nil {
		print(err)
		return ErrorOther
	}

	return NoError
}
