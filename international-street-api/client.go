package street

import (
	"context"
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

// SendLookup sends the lookup, populating the output if the request was successful.
func (c *Client) SendLookup(lookup *Lookup) error {
	if lookup == nil || (len(lookup.Address1) == 0 && len(lookup.Freeform) == 0) {
		return nil
	} else if response, err := c.sender.Send(buildRequest(lookup)); err != nil {
		return err
	} else {
		return deserializeResponse(response, lookup)
	}
}

func (c *Client) SendLookupWithContext(ctx context.Context, lookup *Lookup) error {
	if lookup == nil || (len(lookup.Address1) == 0 && len(lookup.Freeform) == 0) {
		return nil
	}

	request := buildRequest(lookup)
	request = request.WithContext(ctx)
	response, err := c.sender.Send(request)
	if err != nil {
		return err
	} else {
		return deserializeResponse(response, lookup)
	}
}

func deserializeResponse(response []byte, lookup *Lookup) error {
	var candidates []*Candidate
	err := json.Unmarshal(response, &candidates)
	if err != nil {
		return err
	}
	lookup.Results = candidates
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
