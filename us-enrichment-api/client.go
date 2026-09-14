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

// Deprecated: SendPropertyPrincipalLookup is deprecated. Use SendPropertyPrincipal
func (c *Client) SendPropertyPrincipalLookup(smartyKey string) ([]*PrincipalResponse, error) {
	return c.SendPropertyPrincipal(&Lookup{SmartyKey: smartyKey})
}

func (c *Client) SendPropertyPrincipal(lookup *Lookup) ([]*PrincipalResponse, error) {
	propertyLookup := &principalLookup{Lookup: lookup}
	err := c.sendLookup(propertyLookup)
	return propertyLookup.Response, err
}

// SendPropertyPrincipalWithContextAndAuth sends a lookup with the provided context and per-request credentials.
// If credential is non-nil, it will be used to sign this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) SendPropertyPrincipalWithContextAndAuth(ctx context.Context, lookup *Lookup, credential sdk.Credential) ([]*PrincipalResponse, error) {
	propertyLookup := &principalLookup{Lookup: lookup}
	err := c.sendLookupWithContextAndAuth(ctx, propertyLookup, credential)
	return propertyLookup.Response, err
}

func (c *Client) SendGeoReference(lookup *Lookup) ([]*GeoReferenceResponse, error) {
	geoRefLookup := &geoReferenceLookup{Lookup: lookup}
	err := c.sendLookup(geoRefLookup)
	return geoRefLookup.Response, err
}

// SendGeoReferenceWithContextAndAuth sends a lookup with the provided context and per-request credentials.
// If credential is non-nil, it will be used to sign this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) SendGeoReferenceWithContextAndAuth(ctx context.Context, lookup *Lookup, credential sdk.Credential) ([]*GeoReferenceResponse, error) {
	geoRefLookup := &geoReferenceLookup{Lookup: lookup}
	err := c.sendLookupWithContextAndAuth(ctx, geoRefLookup, credential)
	return geoRefLookup.Response, err
}

func (c *Client) SendGeoReferenceWithVersion(lookup *Lookup, censusVersion string) ([]*GeoReferenceResponse, error) {
	geoRefLookup := &geoReferenceLookup{Lookup: lookup, CensusVersion: censusVersion}
	err := c.sendLookup(geoRefLookup)
	return geoRefLookup.Response, err
}

// SendGeoReferenceWithVersionContextAndAuth sends a lookup with the provided context and per-request credentials.
// If credential is non-nil, it will be used to sign this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) SendGeoReferenceWithVersionContextAndAuth(ctx context.Context, lookup *Lookup, censusVersion string, credential sdk.Credential) ([]*GeoReferenceResponse, error) {
	geoRefLookup := &geoReferenceLookup{Lookup: lookup, CensusVersion: censusVersion}
	err := c.sendLookupWithContextAndAuth(ctx, geoRefLookup, credential)
	return geoRefLookup.Response, err
}

func (c *Client) SendSecondary(lookup *Lookup) ([]*SecondaryResponse, error) {
	sLookup := &secondaryLookup{Lookup: lookup}
	err := c.sendLookup(sLookup)
	return sLookup.Response, err
}

// SendSecondaryWithContextAndAuth sends a lookup with the provided context and per-request credentials.
// If credential is non-nil, it will be used to sign this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) SendSecondaryWithContextAndAuth(ctx context.Context, lookup *Lookup, credential sdk.Credential) ([]*SecondaryResponse, error) {
	sLookup := &secondaryLookup{Lookup: lookup}
	err := c.sendLookupWithContextAndAuth(ctx, sLookup, credential)
	return sLookup.Response, err
}

// Deprecated: SendSecondaryLookup is deprecated. Use SendSecondary
func (c *Client) SendSecondaryLookup(lookup *Lookup) ([]*SecondaryResponse, error) {
	return c.SendSecondary(lookup)
}

func (c *Client) SendSecondaryCount(lookup *Lookup) ([]*SecondaryCountResponse, error) {
	scLookup := &secondaryCountLookup{Lookup: lookup}
	err := c.sendLookup(scLookup)
	return scLookup.Response, err
}

// SendSecondaryCountWithContextAndAuth sends a lookup with the provided context and per-request credentials.
// If credential is non-nil, it will be used to sign this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) SendSecondaryCountWithContextAndAuth(ctx context.Context, lookup *Lookup, credential sdk.Credential) ([]*SecondaryCountResponse, error) {
	scLookup := &secondaryCountLookup{Lookup: lookup}
	err := c.sendLookupWithContextAndAuth(ctx, scLookup, credential)
	return scLookup.Response, err
}

// Deprecated: SendSecondaryCountLookup is deprecated. Use SendSecondaryCount
func (c *Client) SendSecondaryCountLookup(lookup *Lookup) ([]*SecondaryCountResponse, error) {
	return c.SendSecondaryCount(lookup)
}

func (c *Client) SendBusinessSummary(lookup *Lookup) ([]*BusinessSummaryResponse, error) {
	bLookup := &businessSummaryLookup{Lookup: lookup}
	err := c.sendLookup(bLookup)
	return bLookup.Response, err
}

func (c *Client) SendBusinessSummaryWithContextAndAuth(ctx context.Context, lookup *Lookup, credential sdk.Credential) ([]*BusinessSummaryResponse, error) {
	bLookup := &businessSummaryLookup{Lookup: lookup}
	err := c.sendLookupWithContextAndAuth(ctx, bLookup, credential)
	return bLookup.Response, err
}

func (c *Client) SendBusinessDetail(lookup *Lookup) ([]*BusinessDetailResponse, error) {
	bLookup := &businessDetailLookup{Lookup: lookup}
	err := c.sendLookup(bLookup)
	return bLookup.Response, err
}

func (c *Client) SendBusinessDetailWithContextAndAuth(ctx context.Context, lookup *Lookup, credential sdk.Credential) ([]*BusinessDetailResponse, error) {
	bLookup := &businessDetailLookup{Lookup: lookup}
	err := c.sendLookupWithContextAndAuth(ctx, bLookup, credential)
	return bLookup.Response, err
}

func (c *Client) SendUniversalLookup(lookup *Lookup, dataSet, dataSubset string) ([]byte, error) {
	return c.SendUniversalLookupWithContext(context.Background(), lookup, dataSet, dataSubset)
}

// SendUniversalLookupWithContext sends a lookup with the provided context.
func (c *Client) SendUniversalLookupWithContext(ctx context.Context, lookup *Lookup, dataSet, dataSubset string) ([]byte, error) {
	return c.SendUniversalLookupWithContextAndAuth(ctx, lookup, dataSet, dataSubset, nil)
}

// SendUniversalLookupWithContextAndAuth sends a lookup with the provided context and per-request credentials.
// If credential is non-nil, it will be used to sign this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) SendUniversalLookupWithContextAndAuth(ctx context.Context, lookup *Lookup, dataSet, dataSubset string, credential sdk.Credential) ([]byte, error) {
	u := &universalLookup{
		Lookup:     lookup,
		DataSet:    dataSet,
		DataSubset: dataSubset,
	}

	err := c.sendLookupWithContextAndAuth(ctx, u, credential)
	return u.Response, err
}

func (c *Client) sendLookup(lookup enrichmentLookup) error {
	return c.sendLookupWithContextAndAuth(context.Background(), lookup, nil)
}

func (c *Client) sendLookupWithContext(ctx context.Context, lookup enrichmentLookup) error {
	return c.sendLookupWithContextAndAuth(ctx, lookup, nil)
}

// sendLookupWithContextAndAuth sends an enrichmentLookup with the provided context and per-request credentials.
// If credential is non-nil, it will be used to sign this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) sendLookupWithContextAndAuth(ctx context.Context, lookup enrichmentLookup, credential sdk.Credential) error {
	if lookup == nil || lookup.getLookup() == nil {
		return nil
	}

	request := buildRequest(ctx, lookup)
	if credential != nil {
		if err := credential.Sign(request); err != nil {
			return err
		}
	}
	response, err := c.sender.Send(request)
	if err != nil {
		return err
	}

	var headers http.Header
	if request.Response != nil {
		headers = request.Response.Header
	}

	if etag := headers.Get(lookupETagHeader); etag != "" {
		lookup.getLookup().ResponseETag = etag
	}

	if request.Response != nil && request.Response.StatusCode == http.StatusNotModified {
		return nil
	}

	return lookup.unmarshalResponse(response, headers)
}

func (c *Client) IsHTTPErrorCode(err error, code int) bool {
	if serr, ok := err.(*sdk.HTTPStatusError); ok && serr.StatusCode() == code {
		return true
	}
	return false
}

func buildRequest(ctx context.Context, lookup enrichmentLookup) *http.Request {
	request, _ := http.NewRequestWithContext(ctx, "GET", buildLookupURL(lookup), nil) // We control the method and the URL. This is safe.
	query := request.URL.Query()
	lookup.populate(query)
	request.Header.Add(lookupETagHeader, lookup.getLookup().ETag)
	request.URL.RawQuery = query.Encode()
	return request
}

func buildLookupURL(lookup enrichmentLookup) string {
	if lookup.getDataSet() == businessDataSet {
		if bl, ok := lookup.(businessIDProvider); ok && len(bl.getBusinessID()) > 0 {
			return "/lookup/" + lookup.getDataSet() + "/" + bl.getBusinessID()
		}
	}

	var newLookupURL string
	if len(lookup.getDataSubset()) == 0 {
		newLookupURL = strings.Replace(lookupURLWithoutSubSet, lookupURLSmartyKey, getLookupURLSmartyKeyReplacement(lookup), 1)
	} else {
		newLookupURL = strings.Replace(lookupURLWithSubSet, lookupURLSmartyKey, getLookupURLSmartyKeyReplacement(lookup), 1)
	}

	newLookupURL = strings.Replace(newLookupURL, lookupURLDataSet, lookup.getDataSet(), 1)
	return strings.TrimSuffix(strings.Replace(newLookupURL, lookupURLDataSubSet, lookup.getDataSubset(), 1), "/")
}

type businessIDProvider interface {
	getBusinessID() string
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
