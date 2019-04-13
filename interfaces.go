package autocomplete

import (
	"net/url"
)

// Database writes and reads widget data
type Database interface {
	Get(request string) ([]*WidgetData, error)
	Write(request string, data []*WidgetData) error
}

// AVSProxy is proxy interface. It is used to handle HTTP requests to aviasales autocomplete service
type AVSProxy interface {
	SearchPlaces(params url.Values) ([]*WidgetData, error)
}
