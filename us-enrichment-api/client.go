package us_enrichment

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

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
	if len(lookup.SmartyKey) == 0 && len(lookup.DataSet) == 0 {
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

func deserializeResponse(body []byte, lookup *Lookup) error {
	err := json.Unmarshal(body, &lookup.Response)
	if err != nil {
		return err
	}
	return nil
}

func buildRequest(lookup *Lookup) *http.Request {
	request, _ := http.NewRequest("GET", buildLookupURL(lookup), nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	request.URL.RawQuery = query.Encode()
	return request
}

func buildLookupURL(lookup *Lookup) string {
	newLookupURL := strings.Replace(lookupURL, lookupURLSmartyKey, lookup.SmartyKey, 1)
	return strings.Replace(newLookupURL, lookupURLDataSet, lookup.DataSet, 1)
}

const (
	lookupURLSmartyKey = ":smartykey"
	lookupURLDataSet   = ":dataset"
	lookupURL          = "/lookup/" + lookupURLSmartyKey + "/" + lookupURLDataSet // Remaining parts will be completed later by the sdk.BaseURLClient.
)
