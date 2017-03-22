package international

import (
	"encoding/json"
	"net/http"

	"github.com/smartystreets/smartystreets-go-sdk"
)

// Client is responsible for sending batches of addresses to the international-street-api.
type Client struct {
	sender sdk.RequestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

// SendBatch sends the batch of inputs, populating the output for each input if the batch was successful.
func (c *Client) SendLookup(lookup *Lookup) error {
	if lookup == nil { // TODO any other conditions?
		return nil
	} else if response, err := c.sender.Send(buildRequest(lookup)); err != nil {
		return err
	} else {
		return deserializeResponse(response, lookup)
	}
}

func deserializeResponse(response []byte, lookup *Lookup) error {
	var results []*Result
	err := json.Unmarshal(response, &results)
	if err != nil {
		return err
	}
	lookup.Results = results
	return nil
}

func buildRequest(lookup *Lookup) *http.Request {
	request, _ := http.NewRequest("GET", verifyURL, nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	lookup.populate(query)
	request.URL.RawQuery = query.Encode()
	return request
}

const verifyURL = "/verify" // Remaining parts will be completed later by the sdk.BaseURLClient.
