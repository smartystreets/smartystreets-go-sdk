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
	if err := json.Unmarshal(response, &candidates); err != nil {
		return err
	}
	for _, candidate := range candidates {
		batch.attach(candidate)
	}
	return nil
}

func buildRequest(batch *Batch) (request *http.Request, err error) {
	if batch == nil || batch.Length() == 0 {
		return nil, emptyBatchError
	}

	request, err = buildPostRequest(batch)
	setHeaders(batch, request)
	return request, err
}

func setHeaders(batch *Batch, request *http.Request) {
	if batch.includeInvalid {
		request.Header.Set(xIncludeInvalidHeader, "true")
	} else if batch.standardizeOnly {
		request.Header.Set(xStandardizeOnlyHeader, "true")
	}
}

func buildPostRequest(batch *Batch) (*http.Request, error) {
	payload, _ := json.Marshal(batch.lookups) // err ignored because since we control the types being serialized it is safe.
	return http.NewRequest("POST", defaultAPIURL, bytes.NewReader(payload))
}

// defaultAPIURL may be overwritten later by a Sender depending on wireup.
const (
	defaultAPIURL          = "https://api.smartystreets.com/street-address"
	xStandardizeOnlyHeader = "X-Standardize-Only"
	xIncludeInvalidHeader  = "X-Include-Invalid"
)

var emptyBatchError = errors.New("The batch was nil or had no records.")
