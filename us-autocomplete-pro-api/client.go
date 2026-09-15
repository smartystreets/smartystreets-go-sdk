package autocomplete_pro

import (
	"context"
	"encoding/json/v2"
	"net/http"

	sdk "github.com/smartystreets/smartystreets-go-sdk/v2"
)

// Client is responsible for sending of lookups to the us-autocomplete-pro-api.
type Client struct {
	sender sdk.RequestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

// SendBatch sends the batch of inputs, populating the output for each input if the batch was successful.
func (c *Client) SendLookup(lookup *Lookup) error {
	return c.SendLookupWithContext(context.Background(), lookup)
}

func (c *Client) SendLookupWithContext(ctx context.Context, lookup *Lookup) error {
	if lookup == nil || len(lookup.Search) == 0 {
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
	var suggestions suggestionListing
	err := json.Unmarshal(response, &suggestions)
	if err != nil {
		return err
	}
	lookup.Results = suggestions.Listing
	return nil
}

func buildRequest(ctx context.Context, lookup *Lookup) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", suggestURL, nil)
	if err != nil {
		return nil, err
	}
	query := request.URL.Query()
	lookup.populate(query)
	request.URL.RawQuery = query.Encode()
	return request, nil
}

const suggestURL = "/lookup" // Remaining parts will be completed later by the sdk.BaseURLClient.
