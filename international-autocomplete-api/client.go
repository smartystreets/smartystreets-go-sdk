package international_autocomplete_api

import (
	"context"
	"net/http"

	sdk "github.com/smartystreets/smartystreets-go-sdk"
	"github.com/smartystreets/smartystreets-go-sdk/internal/json"
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

	request, err := buildRequest(ctx, lookup)
	if err != nil {
		return err
	}
	response, err := c.sender.Send(request)
	if err != nil {
		return err
	} else {
		return deserializeResponse(response, lookup)
	}
}

func deserializeResponse(response []byte, lookup *Lookup) error {
	var result *Result
	if err := json.Unmarshal(response, &result); err != nil {
		return err
	}
	lookup.Result = result
	return nil
}

func buildRequest(ctx context.Context, lookup *Lookup) (*http.Request, error) {
	var addressID = ""
	if len(lookup.AddressID) > 0 {
		addressID = "/" + lookup.AddressID
	}
	request, err := http.NewRequestWithContext(ctx, "GET", suggestURL+addressID, nil)
	if err != nil {
		return nil, err
	}
	query := request.URL.Query()
	lookup.populate(query)
	request.URL.RawQuery = query.Encode()
	return request, nil
}

const suggestURL = "/v2/lookup" // Remaining parts will be completed later by the sdk.BaseURLClient.
