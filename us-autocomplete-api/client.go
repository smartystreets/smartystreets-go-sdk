package autocomplete

import (
	"encoding/json"
	"net/http"

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
		return json.Unmarshal(response, &lookup.Results)
	}
}

func buildRequest(lookup *Lookup) *http.Request {
	request, _ := http.NewRequest("GET", suggestURL, nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	query.Set("prefix", lookup.Prefix)
	// TODO: additional input fields
	request.URL.RawQuery = query.Encode()
	request.Header.Set("Content-Type", "application/json")
	return request
}

const suggestURL = "/suggest" // Remaining parts will be completed later by the sdk.BaseURLClient.
