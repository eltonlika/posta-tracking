package tracker

import(
	"strings"
	"time"
	"sort"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
	"github.com/PuerkitoBio/goquery"
)

// Event holds tracking event information
type Event struct {
	Date        time.Time
	Description string
	Location    string
	Destination string
}

// ToString convert Event struct to string representation
func(e Event) ToString() string{
	d := "|"
	return e.Date.Format("2006-01-02 15:04 PM") + d + e.Description + d + e.Location + d + e.Destination;
}

// Events array type
type Events []Event

func (events Events) Len() int           { return len(events) }
func (events Events) Less(i, j int) bool { return events[i].Date.Before(events[j].Date) }
func (events Events) Swap(i, j int)      { events[i], events[j] = events[j], events[i] }

// SortDirection int alias to specify event sorting direction
type SortDirection int

const (
	// SortAscending flag to sort events in ascending order
	SortAscending SortDirection = 0
	// SortDescending flag to sort events in descending order
	SortDescending SortDirection = 1
	// DefaultEventSortingDirection of events is Ascending
	DefaultEventSortingDirection SortDirection = SortAscending
	// DefaultRequestTimeout is 8 seconds of waiting for a tracking request
	DefaultRequestTimeout = time.Second * 8
	// DefaultServiceURL url of tracking service page
	DefaultServiceURL = "https://gjurmo.postashqiptare.al/tracking.aspx"
)

// Tracker struct mantains tracker instance & configuration
type Tracker struct {
	ServiceURL    			string
	EventSortingDirection   SortDirection
	browser 	   			*browser.Browser
}

// SetRequestTimeout set time to wait for response from tracking service
func (tracker *Tracker) SetRequestTimeout(timeout time.Duration){
	tracker.browser.SetTimeout(timeout)
}

// NewTracker creates new tracker instance
func NewTracker() *Tracker{
	bow := surf.NewBrowser()
	bow.SetAttribute(browser.FollowRedirects, true)
	bow.SetUserAgent(agent.Chrome())
	bow.SetTimeout(DefaultRequestTimeout)

	return &Tracker{
		ServiceURL: DefaultServiceURL,
		EventSortingDirection: DefaultEventSortingDirection,
		browser: bow,
	}
}

func (tracker *Tracker) getTrackingEventsFromPage(trackingNumber string) *Events{
	bow := tracker.browser
	
	err := bow.Open(tracker.ServiceURL)
	if err != nil{
		panic(err)
	}

	fm, err := bow.Form("#form1")
	if err != nil{
		panic(err)
	}

	fm.Input("txt_barcode", trackingNumber)
	fm.Input("hBarCodes", trackingNumber)
	err = fm.Submit()
	if err != nil {
		panic(err)
	}

	// get only rows with data cells, exclude header row
	tableDataRows := bow.Dom().Find("table tr").FilterFunction(func (_ int, s *goquery.Selection) bool {
		return s.ChildrenFiltered("td").Length() > 0
	})

	events := make(Events, tableDataRows.Length())

	tableDataRows.Each(func (i int, s *goquery.Selection){
		event := Event{}
		s.ChildrenFiltered("td").Each(func (j int, s2 *goquery.Selection){
			value := strings.TrimSpace(s2.Text())
			if j==0 {
				if dt, err := time.Parse("02-01-2006 15:04 PM", value); err == nil {
					event.Date = dt
				} else {
					panic(err)
				}
			}else if j==1{
				event.Description = value
			}else if j==2{
				event.Location = value
			}else if j==3{
				event.Destination = value
			}
		})
		events[i] = event
	})

	return &events
}

// Track returns tracking events for given tracking number
func (tracker *Tracker) Track(trackingNumber string) *Events {
	events := tracker.getTrackingEventsFromPage(trackingNumber)
	
	if(tracker.EventSortingDirection == SortAscending){
		sort.Stable(events)
	}else if(tracker.EventSortingDirection == SortDescending){
		sort.Stable(sort.Reverse(events))
	}

	return events;
}