package proxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/titusjaka/autocomplete"
)

// Constants for aviasales autocomplete services
const (
	// PlacesURI stands for aviasales autocomplete service for places (airports, countries, cities)
	PlacesURI = "v2/places.json"
)

// Proxy implements AVSProxy interface
// It is used to fetch data from aviasales autocomplete service
type Proxy struct {
	avsURL string
	client *http.Client
	logger *log.Logger
}

// New returns new Proxy structure
func New(url string, timeout time.Duration) *Proxy {
	return &Proxy{
		avsURL: url,
		client: &http.Client{Timeout: timeout},
		logger: log.New(ioutil.Discard, "", 0),
	}
}

// SetLogger sets logger for Proxy structure
func (p *Proxy) SetLogger(logger *log.Logger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}
	p.logger = logger
	p.logger.Printf("[DEBUG] logger set successfully")
	return nil
}

// SearchPlaces processes HTTP query to aviasales autocomplete service
// It fetches data from AVS service and returns data in widget format
func (p *Proxy) SearchPlaces(params url.Values) ([]*autocomplete.WidgetData, error) {
	resp, err := p.makeRequest(PlacesURI, http.MethodGet, params)
	if err != nil {
		p.logger.Printf("[ERROR] Failed to search places in Aviasales autocomplete service: %v", err)
		return nil, errors.Wrap(err, "failed to search places")
	}

	var places []*autocomplete.AVSPlace
	if err = json.Unmarshal(resp, &places); err != nil {
		p.logger.Printf("[ERROR] Failed to unmarshal HTTP response: %v", err)
		return nil, errors.Wrap(err, "failed to unmarshal HTTP response")
	}

	wd := make([]*autocomplete.WidgetData, len(places))
	for key, place := range places {
		wd[key] = autocomplete.CreateWidgetData(place)
	}

	return wd, nil
}

func (p *Proxy) makeRequest(serviceURI string, method string, params url.Values) ([]byte, error) {
	req, err := buildRequest(fmt.Sprintf("%s/%s", p.avsURL, serviceURI), method, params)
	if err != nil {
		p.logger.Printf("[ERROR] Failed to build HTTP request: %v", err)
		return nil, err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Printf("[ERROR] Failed to process HTTP request: %v", err)
		return nil, errors.Wrap(err, "failed to process HTTP request")
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			p.logger.Printf("[ERROR] Failed to close HTTP body: %v", err)
		}
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.logger.Printf("[ERROR] Failed to read HTTP response body: %v", err)
		return nil, errors.Wrap(err, "failed to read HTTP body")
	}
	return respBody, nil
}

func buildRequest(uri string, method string, params url.Values) (*http.Request, error) {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build request")
	}

	query := req.URL.Query()
	for key, value := range params {
		for index := range value {
			query.Add(key, value[index])
		}
	}
	req.URL.RawQuery = query.Encode()
	return req, nil
}
