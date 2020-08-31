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
	return c.SendLookupWithContext(context.Background(), lookup)
}

func (c *Client) SendLookupWithContext(ctx context.Context, lookup *Lookup) error {
	if lookup == nil {
		return nil
	}
	if lookup.Latitude == 0 && lookup.Longitude == 0 {
		return nil
	}

	request := buildRequest(lookup)
	request = request.WithContext(ctx)

	response, err := c.sender.Send(request)
	if err != nil {
		return err
	}

	return deserializeResponse(response, lookup)
}

func deserializeResponse(response []byte, lookup *Lookup) error {
	var results resultListing
	err := json.Unmarshal(response, &results)
	if err != nil {
		return err
	}
	lookup.Results = results.Listing
	return nil
}

func buildRequest(lookup *Lookup) *http.Request {
	request, _ := http.NewRequest("GET", lookupURL, nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	query.Set("latitude", strconv.FormatFloat(lookup.Latitude, 'f', 8, 32))
	query.Set("longitude", strconv.FormatFloat(lookup.Longitude, 'f', 8, 32))
	request.URL.RawQuery = query.Encode()
	return request
}

const lookupURL = "/lookup" // Remaining parts will be completed later by the sdk.BaseURLClient.
