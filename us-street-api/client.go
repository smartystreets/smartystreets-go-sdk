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

// Ping:
// 	Returns (true,  nil) if the us-street-api service is reachable and responding nominally.
//  Returns (false, nil) if the us-street-api service is reachable but not responding.
//  Returns (false, err) if the status request could not be completed.
func (c *Client) Ping() (bool, error) {
	result, err := c.sender.Send(buildPingRequest())
	if err != nil {
		return false, err
	}
	return string(result) == "OK", nil
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
	request, _ := http.NewRequest("GET", placeholderURL, nil)
	return request
}

var (
	placeholderURL = "/street-address" // Remaining parts will be completed later by the sdk.BaseURLClient.
)

var emptyBatchError = errors.New("The batch was nil or had no records.")
