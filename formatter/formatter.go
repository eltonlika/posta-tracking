package formatter

import (
	"fmt"
	"io"
	"strconv"

	"github.com/eltonlika/posta-tracking/tracker"
	"github.com/ryanuber/columnize"
)

// EventsFormatter struct holds events printer configuration
type EventsFormatter struct {
	NoHeader  bool
	Delimiter string
}

// NewEventsFormatter return new instance of EventsFormatter
func NewEventsFormatter() *EventsFormatter {
	return &EventsFormatter{
		NoHeader:  false,
		Delimiter: "  ",
	}
}

// Format return string of formatted events
func (p *EventsFormatter) Format(events tracker.Events) string {
	var table []string
	var rows []string

	if p.NoHeader {
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
	config.Glue = p.Delimiter
	config.Prefix = ""
	config.Empty = ""
	config.NoTrim = false
	return columnize.Format(table, config)
}

// Print events
func (p *EventsFormatter) Print(events tracker.Events, w io.Writer) error {
	formattedEvents := p.Format(events)
	_, err := fmt.Fprintln(w, formattedEvents)
	return err
}
