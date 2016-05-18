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

// Send sends the batch of inputs, populating the output for each input if the batch was successful.
func (c *Client) Send(batch *Batch) error {
	if request, err := buildRequest(batch); err != nil {
		return err
	} else if response, err := c.sender.Send(request); err != nil {
		return err
	} else {
		return deserializeResponse(response, batch)
	}
}

func deserializeResponse(response []byte, batch *Batch) error {
	var candidates []Candidate
	err := json.Unmarshal(response, &candidates)
	if err == nil {
		for _, candidate := range candidates {
			batch.attach(candidate) // TODO: what about inputs that don't produce any outputs?
		}
	}
	return err
}

func buildRequest(batch *Batch) (request *http.Request, err error) {
	if batch == nil || batch.Length() == 0 {
		return nil, emptyBatchError
	}

	if length := batch.Length(); length == 1 {
		request, err = buildGetRequest(batch)
	} else {
		request, err = buildPostRequest(batch)
	}

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

func buildGetRequest(batch *Batch) (*http.Request, error) {
	url := defaultAPIURL + "?" + batch.marshalQueryString()
	return http.NewRequest("GET", url, nil)
}

func buildPostRequest(batch *Batch) (*http.Request, error) {
	payload, _ := batch.marshalJSON() // err ignored because since we control the types being serialized it is safe.
	return http.NewRequest("POST", defaultAPIURL, bytes.NewReader(payload))
}

// defaultAPIURL may be overwritten later by a Sender depending on wireup.
const (
	defaultAPIURL          = "https://api.smartystreets.com/street-address"
	xStandardizeOnlyHeader = "X-Standardize-Only"
	xIncludeInvalidHeader  = "X-Include-Invalid"
)

var emptyBatchError = errors.New("The batch was nil or had no records.")
