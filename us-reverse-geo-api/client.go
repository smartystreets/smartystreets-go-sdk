package us_reverse_geo

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/smartystreets/smartystreets-go-sdk"
)

type Client struct {
	sender sdk.RequestSender
}

func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

func (c *Client) SendLookup(lookup *Lookup) error {
	return c.SendLookupWithContextAndAuth(context.Background(), lookup, "", "")
}

func (c *Client) SendLookupWithContext(ctx context.Context, lookup *Lookup) error {
	return c.SendLookupWithContextAndAuth(ctx, lookup, "", "")
}

// SendLookupWithContextAndAuth sends a lookup with the provided context and per-request credentials.
// If authID and authToken are both non-empty, they will be used for this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) SendLookupWithContextAndAuth(ctx context.Context, lookup *Lookup, authID, authToken string) error {
	if lookup == nil {
		return nil
	}
	if lookup.Latitude == 0 && lookup.Longitude == 0 {
		return nil
	}

	request := buildRequest(lookup)
	request = request.WithContext(ctx)
	if len(authID) > 0 && len(authToken) > 0 {
		request.SetBasicAuth(authID, authToken)
	}

	response, err := c.sender.Send(request)
	if err != nil {
		return err
	}

	return deserializeResponse(response, lookup)
}

func deserializeResponse(body []byte, lookup *Lookup) error {
	err := json.Unmarshal(body, &lookup.Response)
	if err != nil {
		return err
	}
	return nil
}

func buildRequest(lookup *Lookup) *http.Request {
	request, _ := http.NewRequest("GET", lookupURL, nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	query.Set("latitude", strconv.FormatFloat(lookup.Latitude, 'f', 8, 64))
	query.Set("longitude", strconv.FormatFloat(lookup.Longitude, 'f', 8, 64))
	query.Set("source", lookup.Source)
	request.URL.RawQuery = query.Encode()
	return request
}

const lookupURL = "/lookup" // Remaining parts will be completed later by the sdk.BaseURLClient.
