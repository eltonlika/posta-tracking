package tracker

import (
	"sort"
	"time"

	q "github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
)

const (
	// DefaultServiceURL url of tracking service page
	DefaultServiceURL = "https://gjurmo.postashqiptare.al/tracking.aspx"
	// DefaultTimeout is 8 seconds of waiting for a tracking request
	DefaultTimeout = time.Second * 8
)

// Event holds tracking event information
type Event struct {
	TrackingNumber string
	Date           time.Time
	Description    string
	Location       string
	Destination    string
}

// Events array type
type Events []Event

func (events Events) Len() int           { return len(events) }
func (events Events) Less(i, j int) bool { return events[i].Date.Before(events[j].Date) }
func (events Events) Swap(i, j int)      { events[i], events[j] = events[j], events[i] }

// Sort events by specified direction
func (events *Events) Sort(descending bool) {
	if descending {
		sort.Stable(sort.Reverse(events))
	} else {
		sort.Stable(events)
	}
}

// Tracker struct mantains tracker instance & configuration
type Tracker struct {
	ServiceURL           string
	SortEventsDescending bool
	browser              *browser.Browser
}

// NewTracker creates new tracker instance
func NewTracker() *Tracker {
	bow := surf.NewBrowser()
	bow.SetAttribute(browser.FollowRedirects, true)
	bow.SetUserAgent(agent.Chrome())
	bow.SetTimeout(DefaultTimeout)

	return &Tracker{
		ServiceURL:           DefaultServiceURL,
		SortEventsDescending: false,
		browser:              bow,
	}
}

// Track returns tracking events for given tracking number
func (tracker *Tracker) Track(trackingNumber string) *Events {
	events := tracker.getTrackingEventsFromPage(trackingNumber)
	events.Sort(tracker.SortEventsDescending)
	return events
}

// SetRequestTimeout set time to wait for response from tracking service
func (tracker *Tracker) SetRequestTimeout(timeout time.Duration) {
	tracker.browser.SetTimeout(timeout)
}

func (tracker *Tracker) getTrackingEventsFromPage(trackingNumber string) *Events {
	bow := tracker.browser

	err := bow.Open(tracker.ServiceURL)
	if err != nil {
		panic(err)
	}

	fm, err := bow.Form("#form1")
	if err != nil {
		panic(err)
	}

	fm.Input("txt_barcode", trackingNumber)
	fm.Input("hBarCodes", trackingNumber)
	err = fm.Submit()
	if err != nil {
		panic(err)
	}

	// get only rows with data cells, exclude header row
	tableDataRows := bow.Dom().Find("table tr").FilterFunction(isTableDataRow)
	events := make(Events, tableDataRows.Length())

	tableDataRows.Each(func(i int, s *q.Selection) {
		event := tableDataRowToEvent(s)
		event.TrackingNumber = trackingNumber
		events[i] = event
	})

	return &events
}

func isTableDataRow(_ int, s *q.Selection) bool {
	return s.ChildrenFiltered("td").Length() > 0
}

func tableDataRowToEvent(s *q.Selection) Event {
	rowValues := s.ChildrenFiltered("td").Map(textExtracter)
	if len(rowValues) < 4 {
		panic("Row has less values than expected: ")
	}

	eventDate, err := time.Parse("02-01-2006 15:04 PM", rowValues[0])
	if err != nil {
		panic(err)
	}

	return Event{
		Date:        eventDate,
		Description: rowValues[1],
		Location:    rowValues[2],
		Destination: rowValues[3],
	}
}

func textExtracter(_ int, s *q.Selection) string { return s.Text() }
