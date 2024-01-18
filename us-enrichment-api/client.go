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

func (c *Client) SendPropertyLookup(lookup *Lookup) (error, interface{}) {
	var propertyLookup enrichmentLookup

	switch lookup.DataSubset {
	case FinancialDataSubset:
		propertyLookup = &financialLookup{Lookup: lookup}
	case PrincipalDataSubset:
		propertyLookup = &principalLookup{Lookup: lookup}
	default:
		//TODO: return unknown datasubset
		return nil, nil
	}

	err := c.sendLookup(propertyLookup)
	return err, propertyLookup.GetResponse()
}

func (c *Client) SendPropertyFinancialLookup(smartyKey string) (error, []*FinancialResponse) {
	l := &financialLookup{
		Lookup: &Lookup{SmartyKey: smartyKey, DataSubset: FinancialDataSubset},
	}
	err := c.sendLookup(l)

	return err, l.Response
}

func (c *Client) SendPropertyPrincipalLookup(smartyKey string) (error, []*PrincipalResponse) {
	l := &principalLookup{
		Lookup: &Lookup{SmartyKey: smartyKey, DataSubset: PrincipalDataSubset},
	}
	err := c.sendLookup(l)

	return err, l.Response
}

func (c *Client) sendLookup(lookup enrichmentLookup) error {
	return c.sendLookupWithContext(context.Background(), lookup)
}

func (c *Client) sendLookupWithContext(ctx context.Context, lookup enrichmentLookup) error {
	if lookup == nil || lookup.GetLookup() == nil {
		return nil
	}
	if len(lookup.GetSmartyKey()) == 0 {
		return nil
	}

	request := buildRequest(lookup)
	request = request.WithContext(ctx)

	response, err := c.sender.Send(request)
	if err != nil {
		return err
	}

	return lookup.UnmarshalResponse(response)
}

func buildRequest(lookup enrichmentLookup) *http.Request {
	request, _ := http.NewRequest("GET", buildLookupURL(lookup), nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	request.URL.RawQuery = query.Encode()
	return request
}

func buildLookupURL(lookup enrichmentLookup) string {
	newLookupURL := strings.Replace(lookupURL, lookupURLSmartyKey, lookup.GetSmartyKey(), 1)
	newLookupURL = strings.Replace(newLookupURL, lookupURLDataSet, lookup.GetDataSet(), 1)
	return strings.Replace(newLookupURL, lookupURLDataSubSet, lookup.GetDataSubset(), 1)
}

const (
	lookupURLSmartyKey  = ":smartykey"
	lookupURLDataSet    = ":dataset"
	lookupURLDataSubSet = ":datasubset"
	lookupURL           = "/lookup/" + lookupURLSmartyKey + "/" + lookupURLDataSet + "/" + lookupURLDataSubSet // Remaining parts will be completed later by the sdk.BaseURLClient.
)
