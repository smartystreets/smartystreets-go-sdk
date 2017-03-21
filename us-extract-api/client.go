package extract

import (
	"encoding/json"
	"net/http"

	"github.com/smartystreets/smartystreets-go-sdk"
)

// Client is responsible for sending requests to the us-extract-api.
type Client struct {
	sender sdk.RequestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

// SendBatch sends the batch of inputs, populating the output for each input if the batch was successful.
func (c *Client) SendLookup(lookup *Lookup) error {
	if lookup == nil || len(lookup.Text) == 0 {
		return nil
	} else if response, err := c.sender.Send(buildRequest(lookup)); err != nil {
		return err
	} else {
		return deserializeResponse(response, lookup)
	}
}

func deserializeResponse(response []byte, lookup *Lookup) error {
	var extraction Result
	err := json.Unmarshal(response, &extraction)
	if err != nil {
		return err
	}
	lookup.Result = &extraction
	return nil
}

func buildRequest(lookup *Lookup) *http.Request {
	request, _ := http.NewRequest("POST", extractURL, nil) // We control the method and the URL. This is safe.
	lookup.populate(request)
	return request
}

const extractURL = "/" // Remaining parts will be completed later by the sdk.BaseURLClient.
