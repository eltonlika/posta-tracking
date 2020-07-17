package tracker

import (
	"fmt"
	"sort"
	"time"

	q "github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"github.com/ryanuber/columnize"
	"gopkg.in/headzoo/surf.v1"
)

// ServiceURL url of tracking service page
const ServiceURL = "https://gjurmo.postashqiptare.al/tracking.aspx"

// Event holds tracking event information
type Event struct {
	Num            uint
	Date           time.Time
	TrackingNumber string
	Description    string
	Location       string
	Destination    string
}

// Events array type
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

// Tracker struct mantains tracker instance & configuration
type Tracker struct {
	browser *browser.Browser
}

// NewTracker creates new tracker instance
func NewTracker() Tracker {
	bow := surf.NewBrowser()
	bow.SetAttribute(browser.FollowRedirects, true)
	bow.SetUserAgent(agent.Chrome())
	return Tracker{browser: bow}
}

// Track returns tracking events for given tracking number
func (tracker Tracker) Track(trackingNumber string) (Events, error) {
	if err := tracker.browser.Open(ServiceURL); err != nil {
		return nil, err
	}

	form, err := tracker.browser.Form("#form1")
	if err != nil {
		return nil, err
	}

	form.Input("txt_barcode", trackingNumber)
	form.Input("hBarCodes", trackingNumber)
	if err = form.Submit(); err != nil {
		return nil, err
	}

	events, err := parseEvents(trackingTableValues(tracker.browser.Dom()))
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

func cellValue(_ int, s *q.Selection) string {
	return s.Text()
}

func tdRow(_ int, s *q.Selection) bool {
	return s.ChildrenFiltered("td").Length() > 0
}

func trackingTableValues(dom *q.Selection) [][]string {
	rows := dom.Find("table#gvTraking tr").FilterFunction(tdRow)

	values := make([][]string, rows.Length())

	rows.Each(func(i int, s *q.Selection) {
		values[i] = s.ChildrenFiltered("td").Map(cellValue)
	})

	return values
}

func parseEvents(values [][]string) (Events, error) {
	events := make(Events, len(values))

	for i, row := range values {
		if len(row) < 4 {
			return nil, fmt.Errorf("Row %v has less values than expected", row)
		}

		date, err := time.Parse("02-01-2006 15:04 PM", row[0])
		if err != nil {
			return nil, err
		}

		events[i] = Event{
			Date:        date,
			Description: row[1],
			Location:    row[2],
			Destination: row[3]}
	}

	return events, nil
}
