package international_autocomplete_api

import (
	"context"
	"encoding/json"
	"net/http"

	sdk "github.com/smartystreets/smartystreets-go-sdk"
)

type Client struct {
	sender sdk.RequestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

func (c *Client) SendLookup(lookup *Lookup) error {
	return c.SendLookupWithContext(context.Background(), lookup)
}

func (c *Client) SendLookupWithContext(ctx context.Context, lookup *Lookup) error {
	if lookup == nil || len(lookup.Country) == 0 || (len(lookup.Search) == 0 && len(lookup.AddressID) == 0) {
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
	err := json.Unmarshal(response, &lookup.Result)
	if err != nil {
		return err
	}
	return nil
}

func buildRequest(lookup *Lookup) *http.Request {
	var addressID = ""
	if len(lookup.AddressID) > 0 {
		addressID = "/" + lookup.AddressID
	}
	request, _ := http.NewRequest("GET", suggestURL+addressID, nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	lookup.populate(query)
	request.URL.RawQuery = query.Encode()
	return request
}

const suggestURL = "/v2/lookup" // Remaining parts will be completed later by the sdk.BaseURLClient.
