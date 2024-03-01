package us_enrichment

import (
	bytesPackage "bytes"
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

type universalLookup struct {
	Lookup     *Lookup
	DataSet    string
	DataSubset string
	Response   []byte
}

func (g *universalLookup) getSmartyKey() string     { return g.Lookup.SmartyKey }
func (g *universalLookup) getDataSet() string       { return g.DataSet }
func (g *universalLookup) getDataSubset() string    { return g.DataSubset }
func (g *universalLookup) getLookup() *Lookup       { return g.Lookup }
func (g *universalLookup) getResponse() interface{} { return g.Response }
func (g *universalLookup) unmarshalResponse(bytes []byte, headers http.Header) error {
	g.Response = bytes
	if headers != nil {
		if etag, found := headers[lookupETagHeader]; found && len(etag) > 0 {

			eTagAttribute := []byte(`"eTag": "` + etag[0] + `",`)
			insertLocation := bytesPackage.IndexByte(bytes, '{') + 1

			if insertLocation > 0 && insertLocation < len(bytes) {
				var modifiedResponse bytesPackage.Buffer
				modifiedResponse.Write(bytes[:insertLocation])
				modifiedResponse.Write(eTagAttribute)
				modifiedResponse.Write(bytes[insertLocation:])
				g.Response = modifiedResponse.Bytes()
			}

		}
	}

	return nil

}
func (g *universalLookup) populate(query url.Values) {
	g.Lookup.populateInclude(query)
	g.Lookup.populateExclude(query)
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

type geoReferenceLookup struct {
	Lookup   *Lookup
	Response []*GeoReferenceResponse
}

func (g *geoReferenceLookup) getDataSubset() string {
	return emptyDataSubset
}

func (g *geoReferenceLookup) populate(query url.Values) {
	g.Lookup.populateInclude(query)
	g.Lookup.populateExclude(query)
}

func (g *geoReferenceLookup) getSmartyKey() string {
	return g.Lookup.SmartyKey
}

func (g *geoReferenceLookup) getDataSet() string {
	return geoReferenceDataSet
}

func (g *geoReferenceLookup) getLookup() *Lookup {
	return g.Lookup
}

func (g *geoReferenceLookup) getResponse() interface{} {
	return g.Response
}

func (g *geoReferenceLookup) unmarshalResponse(bytes []byte, headers http.Header) error {
	if err := json.Unmarshal(bytes, &g.Response); err != nil {
		return err
	}

	if headers != nil {
		if etag, found := headers[lookupETagHeader]; found {
			if len(etag) > 0 && len(g.Response) > 0 {
				g.Response[0].Etag = etag[0]
			}
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////

const (
	financialDataSubset = "financial"
	principalDataSubset = "principal"
	propertyDataSet     = "property"
	geoReferenceDataSet = "geo-reference"
	emptyDataSubset     = ""
)

func (l Lookup) populateInclude(query url.Values) {
	if len(l.Include) > 0 {
		query.Set("include", l.Include)
	}
}

func (l Lookup) populateExclude(query url.Values) {
	if len(l.Exclude) > 0 {
		query.Set("exclude", l.Exclude)
	}
}
