package us_enrichment

import (
	"context"
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

func (c *Client) SendPropertyFinancialLookup(smartyKey string) (error, []*FinancialResponse) {
	return c.SendPropertyFinancialWithLookup(&Lookup{SmartyKey: smartyKey})
}

func (c *Client) SendPropertyFinancialWithLookup(lookup *Lookup) (error, []*FinancialResponse) {
	propertyLookup := &financialLookup{Lookup: lookup}
	err := c.sendLookup(propertyLookup)
	return err, propertyLookup.Response
}

func (c *Client) SendPropertyPrincipalLookup(smartyKey string) (error, []*PrincipalResponse) {
	return c.SendPropertyPrincipalWithLookup(&Lookup{SmartyKey: smartyKey})
}

func (c *Client) SendPropertyPrincipalWithLookup(lookup *Lookup) (error, []*PrincipalResponse) {
	propertyLookup := &principalLookup{Lookup: lookup}
	err := c.sendLookup(propertyLookup)
	return err, propertyLookup.Response
}

func (c *Client) sendLookup(lookup enrichmentLookup) error {
	return c.sendLookupWithContext(context.Background(), lookup)
}

func (c *Client) sendLookupWithContext(ctx context.Context, lookup enrichmentLookup) error {
	if lookup == nil || lookup.getLookup() == nil {
		return nil
	}
	if len(lookup.getSmartyKey()) == 0 {
		return nil
	}

	request := buildRequest(lookup)
	request = request.WithContext(ctx)

	response, err := c.sender.Send(request)
	if err != nil {
		return err
	}

	var headers http.Header
	if request.Response != nil {
		headers = request.Response.Header
	}

	return lookup.unmarshalResponse(response, headers)
}

func buildRequest(lookup enrichmentLookup) *http.Request {
	request, _ := http.NewRequest("GET", buildLookupURL(lookup), nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	lookup.populate(query)
	request.Header.Add(lookupETagHeader, lookup.getLookup().ETag)
	request.URL.RawQuery = query.Encode()
	return request
}

func buildLookupURL(lookup enrichmentLookup) string {
	newLookupURL := strings.Replace(lookupURL, lookupURLSmartyKey, lookup.getSmartyKey(), 1)
	newLookupURL = strings.Replace(newLookupURL, lookupURLDataSet, lookup.getDataSet(), 1)
	return strings.Replace(newLookupURL, lookupURLDataSubSet, lookup.getDataSubset(), 1)
}

const (
	lookupURLSmartyKey  = ":smartykey"
	lookupURLDataSet    = ":dataset"
	lookupURLDataSubSet = ":datasubset"
	lookupURL           = "/lookup/" + lookupURLSmartyKey + "/" + lookupURLDataSet + "/" + lookupURLDataSubSet // Remaining parts will be completed later by the sdk.BaseURLClient.
	lookupETagHeader    = "Etag"
)
