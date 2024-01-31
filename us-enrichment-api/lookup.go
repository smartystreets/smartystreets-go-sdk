package us_enrichment

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Lookup struct {
	SmartyKey string
	Include   string
	Exclude   string
	ETag      string
}

type enrichmentLookup interface {
	getSmartyKey() string
	getDataSet() string
	getDataSubset() string
	getLookup() *Lookup
	getResponse() interface{}
	unmarshalResponse([]byte, http.Header) error
	populate(query url.Values)
}

////////////////////////////////////////////////////////////////////////////////////////

type financialLookup struct {
	Lookup   *Lookup
	Response []*FinancialResponse
}

func (f *financialLookup) getSmartyKey() string {
	return f.Lookup.SmartyKey
}

func (f *financialLookup) getDataSet() string {
	return propertyDataSet
}

func (f *financialLookup) getDataSubset() string {
	return financialDataSubset
}

func (f *financialLookup) getLookup() *Lookup {
	return f.Lookup
}

func (f *financialLookup) getResponse() interface{} {
	return f.Response
}

func (f *financialLookup) unmarshalResponse(bytes []byte, headers http.Header) error {
	if err := json.Unmarshal(bytes, &f.Response); err != nil {
		return err
	}

	if headers != nil {
		if etag, found := headers[lookupETagHeader]; found {
			if len(etag) > 0 && len(f.Response) > 0 {
				f.Response[0].Etag = etag[0]
			}
		}
	}

	return nil
}

func (e *financialLookup) populate(query url.Values) {
	e.Lookup.populateInclude(query)
	e.Lookup.populateExclude(query)
}

////////////////////////////////////////////////////////////////////////////////////////

type principalLookup struct {
	Lookup   *Lookup
	Response []*PrincipalResponse
}

func (p *principalLookup) getSmartyKey() string {
	return p.Lookup.SmartyKey
}

func (p *principalLookup) getDataSet() string {
	return propertyDataSet
}

func (p *principalLookup) getDataSubset() string {
	return principalDataSubset
}

func (p *principalLookup) getLookup() *Lookup {
	return p.Lookup
}

func (f *principalLookup) getResponse() interface{} {
	return f.Response
}

func (p *principalLookup) unmarshalResponse(bytes []byte, headers http.Header) error {
	if err := json.Unmarshal(bytes, &p.Response); err != nil {
		return err
	}

	if headers != nil {
		if etag, found := headers[lookupETagHeader]; found {
			if len(etag) > 0 && len(p.Response) > 0 {
				p.Response[0].Etag = etag[0]
			}
		}
	}

	return nil
}

func (e *principalLookup) populate(query url.Values) {
	e.Lookup.populateInclude(query)
	e.Lookup.populateExclude(query)
}

////////////////////////////////////////////////////////////////////////////////////////

const (
	financialDataSubset = "financial"
	principalDataSubset = "principal"
	propertyDataSet     = "property"
)

func (l Lookup) populateInclude(query url.Values) {
	if len(l.Include) > 0 {
		query.Set("include", l.Include)
	}
}

func (l Lookup) populateExclude(query url.Values) {
	if len(l.Include) > 0 {
		query.Set("exclude", l.Exclude)
	}
}
