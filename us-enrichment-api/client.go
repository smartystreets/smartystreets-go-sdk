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

// Deprecated: SendPropertyFinancialLookup is deprecated. Use SendPropertyFinancial
func (c *Client) SendPropertyFinancialLookup(smartyKey string) (error, []*FinancialResponse) {
	return c.SendPropertyFinancial(&Lookup{SmartyKey: smartyKey})
}
func (c *Client) SendPropertyFinancial(lookup *Lookup) (error, []*FinancialResponse) {
	propertyLookup := &financialLookup{Lookup: lookup}
	err := c.sendLookup(propertyLookup)
	return err, propertyLookup.Response
}

// Deprecated: SendPropertyFinancialLookup is deprecated. Use SendPropertyPrincipal
func (c *Client) SendPropertyPrincipalLookup(smartyKey string) (error, []*PrincipalResponse) {
	return c.SendPropertyPrincipal(&Lookup{SmartyKey: smartyKey})
}
func (c *Client) SendPropertyPrincipal(lookup *Lookup) (error, []*PrincipalResponse) {
	propertyLookup := &principalLookup{Lookup: lookup}
	err := c.sendLookup(propertyLookup)
	return err, propertyLookup.Response
}

func (c *Client) SendGeoReference(lookup *Lookup) (error, []*GeoReferenceResponse) {
	geoRefLookup := &geoReferenceLookup{Lookup: lookup}
	err := c.sendLookup(geoRefLookup)
	return err, geoRefLookup.Response
}

func (c *Client) SendGeoReferenceWithVersion(lookup *Lookup, censusVersion string) (error, []*GeoReferenceResponse) {
	geoRefLookup := &geoReferenceLookup{Lookup: lookup, CensusVersion: censusVersion}
	err := c.sendLookup(geoRefLookup)
	return err, geoRefLookup.Response
}

func (c *Client) SendRiskLookup(lookup *Lookup) (error, []*RiskResponse) {
	rLookup := &riskLookup{Lookup: lookup}
	err := c.sendLookup(rLookup)
	return err, rLookup.Response
}

func (c *Client) SendSecondaryLookup(lookup *Lookup) (error, []*SecondaryResponse) {
	sLookup := &secondaryLookup{Lookup: lookup}
	err := c.sendLookup(sLookup)
	return err, sLookup.Response
}

func (c *Client) SendSecondaryCountLookup(lookup *Lookup) (error, []*SecondaryCountResponse) {
	scLookup := &secondaryCountLookup{Lookup: lookup}
	err := c.sendLookup(scLookup)
	return err, scLookup.Response
}

func (c *Client) SendUniversalLookup(lookup *Lookup, dataSet, dataSubset string) (error, []byte) {
	g := &universalLookup{
		Lookup:     lookup,
		DataSet:    dataSet,
		DataSubset: dataSubset,
	}

	err := c.sendLookup(g)
	return err, g.Response
}

func (c *Client) sendLookup(lookup enrichmentLookup) error {
	return c.sendLookupWithContext(context.Background(), lookup)
}

func (c *Client) sendLookupWithContext(ctx context.Context, lookup enrichmentLookup) error {
	if lookup == nil || lookup.getLookup() == nil {
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

func (c *Client) IsHTTPErrorCode(err error, code int) bool {
	if serr, ok := err.(*sdk.HTTPStatusError); ok && serr.StatusCode() == code {
		return true
	}
	return false
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
	var newLookupURL string
	if len(lookup.getDataSubset()) == 0 {
		newLookupURL = strings.Replace(lookupURLWithoutSubSet, lookupURLSmartyKey, getLookupURLSmartyKeyReplacement(lookup), 1)
	} else {
		newLookupURL = strings.Replace(lookupURLWithSubSet, lookupURLSmartyKey, getLookupURLSmartyKeyReplacement(lookup), 1)
	}

	newLookupURL = strings.Replace(newLookupURL, lookupURLDataSet, lookup.getDataSet(), 1)
	return strings.TrimSuffix(strings.Replace(newLookupURL, lookupURLDataSubSet, lookup.getDataSubset(), 1), "/")
}

func getLookupURLSmartyKeyReplacement(lookup enrichmentLookup) string {
	if len(lookup.getSmartyKey()) > 0 {
		return lookup.getSmartyKey()
	} else {
		return addressSearch
	}
}

const (
	lookupURLSmartyKey     = ":smartykey"
	lookupURLDataSet       = ":dataset"
	lookupURLDataSubSet    = ":datasubset"
	lookupURLWithSubSet    = "/lookup/" + lookupURLSmartyKey + "/" + lookupURLDataSet + "/" + lookupURLDataSubSet // Remaining parts will be completed later by the sdk.BaseURLClient.
	lookupURLWithoutSubSet = "/lookup/" + lookupURLSmartyKey + "/" + lookupURLDataSet                             // Remaining parts will be completed later by the sdk.BaseURLClient.
	lookupETagHeader       = "Etag"
	addressSearch          = "search"
)
