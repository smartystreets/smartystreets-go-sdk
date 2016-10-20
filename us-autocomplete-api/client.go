package autocomplete

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/smartystreets/smartystreets-go-sdk"
)

// Client is responsible for sending batches of addresses to the us-street-api.
type Client struct {
	sender sdk.RequestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

// SendBatch sends the batch of inputs, populating the output for each input if the batch was successful.
func (c *Client) SendLookup(lookup *Lookup) error {
	if lookup == nil || len(lookup.Prefix) == 0 {
		return nil
	} else if response, err := c.sender.Send(buildRequest(lookup)); err != nil {
		return err
	} else {
		return deserializeResponse(response, lookup)
	}
}

func deserializeResponse(response []byte, lookup *Lookup) error {
	var suggestions suggestionListing
	err := json.Unmarshal(response, &suggestions)
	if err != nil {
		return err
	}
	lookup.Results = suggestions.Listing
	return nil
}

func buildRequest(lookup *Lookup) *http.Request {
	request, _ := http.NewRequest("GET", suggestURL, nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()

	// Apply lookup values to query
	query.Set("prefix", lookup.Prefix)
	query.Set("suggestions", strconv.Itoa(lookup.MaxSuggestions))
	query.Set("city_filter", strings.Join(lookup.CityFilter, ","))
	query.Set("state_filter", strings.Join(lookup.StateFilter, ","))
	query.Set("prefer", strings.Join(lookup.CityStatePreferences, ";"))

	request.URL.RawQuery = query.Encode()
	return request
}

const suggestURL = "/suggest" // Remaining parts will be completed later by the sdk.BaseURLClient.
