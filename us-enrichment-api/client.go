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
	if len(lookup.SmartyKey) == 0 && len(lookup.DataSet) == 0 && len(lookup.DataSubSet) == 0 {
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
	var err error

	switch strings.ToLower(lookup.DataSubSet) {
	case "financial":
		err = json.Unmarshal(body, &lookup.FinancialResponse)
		break
	case "principal":
		err = json.Unmarshal(body, &lookup.PrincipalResponse)
	}

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
	newLookupURL = strings.Replace(newLookupURL, lookupURLDataSet, lookup.DataSet, 1)
	return strings.Replace(newLookupURL, lookupURLDataSubSet, lookup.DataSubSet, 1)
}

const (
	lookupURLSmartyKey  = ":smartykey"
	lookupURLDataSet    = ":dataset"
	lookupURLDataSubSet = ":datasubset"
	lookupURL           = "/lookup/" + lookupURLSmartyKey + "/" + lookupURLDataSet + "/" + lookupURLDataSubSet // Remaining parts will be completed later by the sdk.BaseURLClient.
)
