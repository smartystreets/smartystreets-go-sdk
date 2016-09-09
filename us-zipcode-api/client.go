package zipcode

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/smartystreets/smartystreets-go-sdk"
)

// Client is responsible for sending batches of addresses to the us-zipcode-api.
type Client struct {
	sender sdk.RequestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

// SendBatch sends the batch of inputs, populating the output for each input if the batch was successful.
func (c *Client) SendBatch(batch *Batch) error {
	if batch == nil || batch.Length() == 0 {
		return nil
	} else if response, err := c.sender.Send(buildRequest(batch)); err != nil {
		return err
	} else {
		return deserializeResponse(response, batch)
	}
}

func deserializeResponse(response []byte, batch *Batch) error {
	var results []*Result
	err := json.Unmarshal(response, &results)
	if err == nil {
		batch.attach(results)
	}
	return err
}

func buildRequest(batch *Batch) *http.Request {
	payload, _ := json.Marshal(batch.lookups)                                       // We control the types being serialized. This is safe.
	request, _ := http.NewRequest("POST", placeholderURL, bytes.NewReader(payload)) // We control the method and the URL. This is safe.
	request.Header.Set("Content-Type", "application/json")
	return request
}

const placeholderURL = "/lookup" // Remaining parts will be completed later by the sdk.BaseURLClient.
