// +build integration

package proxy

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/titusjaka/autocomplete"
)

func TestProxy_SearchPlaces(t *testing.T) {
	proxy := New("https://places.aviasales.ru", time.Second*3)

	queryParams := url.Values{}

	queryParams.Add("term", "mow")
	queryParams.Add("locale", "us")
	queryParams.Add("types[]", "city")
	queryParams.Add("types[]", "airport")

	resp, err := proxy.SearchPlaces(queryParams)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedResponse, resp)
}

var expectedResponse = []*autocomplete.WidgetData{
	{Slug: "MOW", Title: "Moscow", Subtitle: "Russia"},
	{Slug: "DME", Title: "Moscow Domodedovo Airport", Subtitle: "Moscow"},
	{Slug: "SVO", Title: "Sheremetyevo International Airport", Subtitle: "Moscow"},
	{Slug: "VKO", Title: "Vnukovo Airport", Subtitle: "Moscow"},
	{Slug: "ZIA", Title: "Zhukovsky International Airport", Subtitle: "Moscow"},
	{Slug: "MBA", Title: "Mombasa", Subtitle: "Kenya"},
	{Slug: "STL", Title: "Lambert-St. Louis International Airport", Subtitle: "Saint Louis"},
	{Slug: "MKC", Title: "Charles B. Wheeler Downtown Airport", Subtitle: "Kansas City"},
	{Slug: "MCI", Title: "Kansas City International Airport", Subtitle: "Kansas City"},
	{Slug: "JLN", Title: "Joplin", Subtitle: "United States"},
	{Slug: "SOW", Title: "Show Low", Subtitle: "United States"},
	{Slug: "NSH", Title: "Now Shahr", Subtitle: "Iran"},
	{Slug: "MQN", Title: "Mo i Rana", Subtitle: "Norway"},
	{Slug: "YYY", Title: "Mont Joli", Subtitle: "Canada"},
	{Slug: "SUS", Title: "Spirit of St. Louis Airport", Subtitle: "Saint Louis"},
	{Slug: "BFM", Title: "Mobile Downtown Airport", Subtitle: "Mobile"},
}
