package tracking

import "time"

// Event holds information about tracking event
type Event struct {
	Date        time.Time
	Description string
	Location    string
	Destination string
}

// SortDirection int alias to specify event sorting
type SortDirection int

const (
	// SortAscending flag to sort events in ascending order
	SortAscending SortDirection = 0
	// SortDescending flag to sort events in descending order
	SortDescending SortDirection = 1
	// DefaultSorting of events is ascending
	DefaultSorting SortDirection = SortAscending
	// DefaultRequestTimeout is 8 seconds of waiting for a tracking request
	DefaultRequestTimeout = time.Second * 8
)

// Tracker struct mantains Tracker configuration
type Tracker struct {
	ServiceURL     string
	RequestTimeout time.Duration
	EventSorting   SortDirection
}

// Track returns tracking events for given tracking number
func (*Tracker) Track(trackingNumber string) []Event {
	return make([]Event, 0)
}
