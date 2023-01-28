package tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/ryanuber/columnize"
)

// ServiceURL url of tracking service page
const serviceUrl = "https://www.postashqiptare.al/api/"
const originHeader = "https://www.postashqiptare.al"
const requestTimeout = 8 * time.Second

// Event holds tracking event information
type Event struct {
	Num            uint
	Date           time.Time
	TrackingNumber string
	Description    string
	Location       string
	Destination    string
}

type Events []Event

// Sort interface implementations
func (events Events) Len() int           { return len(events) }
func (events Events) Less(i, j int) bool { return events[i].Date.Before(events[j].Date) }
func (events Events) Swap(i, j int)      { events[i], events[j] = events[j], events[i] }

// Event fmt interface implementation
func (e Event) String() string {
	return fmt.Sprintf("%d|%s|%s|%s|%s|%s",
		e.Num,
		e.TrackingNumber,
		e.Date.Format("2006-01-02 15:04 PM"),
		e.Description,
		e.Location,
		e.Destination)
}

// Tracker struct mantains tracker instance & configuration
type Tracker struct {
	client http.Client
}

// NewTracker creates new tracker instance
func NewTracker() Tracker {
	return Tracker{http.Client{Timeout: requestTimeout}}
}

// Track returns tracking events for given tracking number
func (tracker Tracker) Track(trackingNumber string) (Events, error) {
	postBody := "kodi=" + trackingNumber
	request, err := http.NewRequest(http.MethodPost, serviceUrl, bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Origin", originHeader)

	response, err := tracker.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	events, err := parseEvents(response.Body)
	if err != nil {
		return nil, err
	}

	sort.Stable(events)

	// set each event's sequence num and tracking number
	for i := 0; i < events.Len(); i++ {
		events[i].Num = uint(i + 1)
		events[i].TrackingNumber = trackingNumber
	}

	return events, nil
}

// Events fmt interface implementation
func (events Events) String() string {
	if events.Len() == 0 {
		return ""
	}

	table := make([]string, events.Len()+1)

	table[0] = "#|Kodi|Data|Ngjarja|Zyra|Destinacioni"

	for i, e := range events {
		table[i+1] = e.String()
	}

	config := columnize.DefaultConfig()
	config.Delim = "|"
	config.Glue = "  "
	config.Prefix = ""
	config.Empty = ""
	config.NoTrim = false
	return columnize.Format(table, config)
}

func parseEvents(reader io.Reader) (Events, error) {
	var parsed []struct {
		Date        string `json:"Data"`
		Description string `json:"Ngjarja"`
		Location    string `json:"Zyra"`
		Destination string `json:"Destinacioni"`
	}

	if err := json.NewDecoder(reader).Decode(&parsed); err != nil {
		return nil, err
	}

	events := make(Events, len(parsed))

	for i, row := range parsed {
		date, err := time.Parse("02-01-2006 15:04 PM", row.Date)
		if err != nil {
			return nil, err
		}
		events[i] = Event{
			Date:        date,
			Description: row.Description,
			Location:    row.Location,
			Destination: row.Destination,
		}
	}

	return events, nil
}
