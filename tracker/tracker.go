package tracker

import (
	"errors"
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
	Date           time.Time
	TrackingNumber string
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
func (events Events) Sort(descending bool) {
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

	browser *browser.Browser
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
func (tracker *Tracker) Track(trackingNumber string) (Events, error) {
	events, err := tracker.findTrackingEvents(trackingNumber)
	if err != nil {
		return nil, err
	}

	events.Sort(tracker.SortEventsDescending)
	return events, err
}

// SetRequestTimeout set time to wait for response from tracking service
func (tracker *Tracker) SetRequestTimeout(seconds uint) {
	tracker.browser.SetTimeout(time.Second * time.Duration(seconds))
}

func (tracker *Tracker) findTrackingEvents(trackingNumber string) (Events, error) {
	bow := tracker.browser

	err := bow.Open(tracker.ServiceURL)
	if err != nil {
		return nil, err
	}

	fm, err := bow.Form("#form1")
	if err != nil {
		return nil, err
	}

	fm.Input("txt_barcode", trackingNumber)
	fm.Input("hBarCodes", trackingNumber)
	err = fm.Submit()
	if err != nil {
		return nil, err
	}

	events, err := extractEvents(bow.Dom())
	if err != nil {
		return nil, err
	}

	// set each event's tracking number field
	for i := range events {
		events[i].TrackingNumber = trackingNumber
	}

	return events, nil
}

func extractText(_ int, s *q.Selection) string { return s.Text() }

func extractEvent(s *q.Selection) (Event, error) {
	event := Event{}

	rowValues := s.ChildrenFiltered("td").Map(extractText)
	if len(rowValues) < 4 {
		return event, errors.New("Row has less values than expected")
	}

	eventDate, err := time.Parse("02-01-2006 15:04 PM", rowValues[0])
	if err != nil {
		return event, err
	}

	event.Date = eventDate
	event.Description = rowValues[1]
	event.Location = rowValues[2]
	event.Destination = rowValues[3]

	return event, nil
}

func isTableDataRow(_ int, tr *q.Selection) bool { return tr.ChildrenFiltered("td").Length() > 0 }

func extractEvents(dom *q.Selection) (Events, error) {
	// get only rows with data cells, exclude header row
	tableDataRows := dom.Find("table#gvTraking tr").FilterFunction(isTableDataRow)
	events := make(Events, tableDataRows.Length())

	var err error
	tableDataRows.EachWithBreak(func(i int, tr *q.Selection) bool {
		events[i], err = extractEvent(tr)
		return err == nil // if has error then stop iteration
	})

	return events, err
}
