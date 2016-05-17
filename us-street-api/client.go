package us_street

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Client struct {
	sender requestSender
}

func NewClient(sender requestSender) *Client {
	return &Client{sender: sender}
}

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

func buildRequest(batch *Batch) (*http.Request, error) {
	if batch == nil || batch.Length() == 0 {
		return nil, emptyBatchError
	} else if length := batch.Length(); length == 1 {
		return buildGetRequest(batch)
	} else {
		return buildPostRequest(batch)
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
const defaultAPIURL = "https://api.smartystreets.com/street-address"
