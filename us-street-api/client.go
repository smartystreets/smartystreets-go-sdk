package us_street

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Client is responsible for sending batches of addresses to the us-street-api.
type Client struct {
	sender requestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender requestSender) *Client {
	return &Client{sender: sender}
}

// Ping returns an error if the service is not reachable or not responding.
// The error is of type HTTPStatusError.
func (c *Client) Ping() error {
	_, err := c.sender.Send(buildPingRequest())
	return err
}

// SendBatch sends the batch of inputs, populating the output for each input if the batch was successful.
func (c *Client) SendBatch(batch *Batch) error {
	if request, err := buildRequest(batch); err != nil {
		return err
	} else if response, err := c.sender.Send(request); err != nil {
		return err
	} else {
		return deserializeResponse(response, batch)
	}
}

func deserializeResponse(response []byte, batch *Batch) error {
	var candidates []*Candidate
	err := json.Unmarshal(response, &candidates)
	if err == nil {
		batch.attach(candidates)
	}
	return err
}

func buildRequest(batch *Batch) (*http.Request, error) {
	if batch == nil || batch.Length() == 0 {
		return nil, emptyBatchError
	}
	return buildPostRequest(batch)
}

func buildPostRequest(batch *Batch) (*http.Request, error) {
	payload, _ := json.Marshal(batch.lookups) // err ignored because since we control the types being serialized it is safe.
	return http.NewRequest("POST", placeholderURL, bytes.NewReader(payload))
}

func buildPingRequest() *http.Request {
	request, _ := http.NewRequest("GET", statusURL, nil) // TODO: /status is the url to use
	return request
}

var (
	placeholderURL = "/street-address" // Remaining parts will be completed later by the sdk.BaseURLClient.
	statusURL      = "/status"
)

var emptyBatchError = errors.New("The batch was nil or had no records.")
