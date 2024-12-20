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
	Freeform  string
	Street    string
	City      string
	State     string
	ZIPCode   string
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

func (g *universalLookup) getSmartyKey() string {
	if len(g.Lookup.SmartyKey) == 0 {
		return smartyKeyBypass
	}
	return g.Lookup.SmartyKey
}
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
	g.Lookup.populateFreeform(query)
	g.Lookup.populateStreet(query)
	g.Lookup.populateCity(query)
	g.Lookup.populateState(query)
	g.Lookup.populateZIPCode(query)

}

////////////////////////////////////////////////////////////////////////////////////////

type financialLookup struct {
	Lookup   *Lookup
	Response []*FinancialResponse
}

func (f *financialLookup) getSmartyKey() string {
	if len(f.Lookup.SmartyKey) == 0 {
		return smartyKeyBypass
	}
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
	e.Lookup.populateFreeform(query)
	e.Lookup.populateStreet(query)
	e.Lookup.populateCity(query)
	e.Lookup.populateState(query)
	e.Lookup.populateZIPCode(query)
}

////////////////////////////////////////////////////////////////////////////////////////

type principalLookup struct {
	Lookup   *Lookup
	Response []*PrincipalResponse
}

func (p *principalLookup) getSmartyKey() string {
	if len(p.Lookup.SmartyKey) == 0 {
		return smartyKeyBypass
	}
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
	e.Lookup.populateFreeform(query)
	e.Lookup.populateStreet(query)
	e.Lookup.populateCity(query)
	e.Lookup.populateState(query)
	e.Lookup.populateZIPCode(query)
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
	g.Lookup.populateFreeform(query)
	g.Lookup.populateStreet(query)
	g.Lookup.populateCity(query)
	g.Lookup.populateState(query)
	g.Lookup.populateZIPCode(query)
}

func (g *geoReferenceLookup) getSmartyKey() string {
	if len(g.Lookup.SmartyKey) == 0 {
		return smartyKeyBypass
	}
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

type secondaryLookup struct {
	*Lookup
	Response []*SecondaryResponse
}

func (s *secondaryLookup) getSmartyKey() string {
	if len(s.Lookup.SmartyKey) == 0 {
		return smartyKeyBypass
	}
	return s.SmartyKey
}

func (s *secondaryLookup) getDataSet() string {
	return secondaryData
}

func (s *secondaryLookup) getDataSubset() string {
	return emptyDataSubset
}

func (s *secondaryLookup) getLookup() *Lookup {
	return s.Lookup
}

func (s *secondaryLookup) getResponse() interface{} {
	return s.Response
}

func (s *secondaryLookup) unmarshalResponse(bytes []byte, header http.Header) error {
	if err := json.Unmarshal(bytes, &s.Response); err != nil {
		return err
	}

	if header != nil {
		if etag, found := header[lookupETagHeader]; found {
			if len(etag) > 0 && len(s.Response) > 0 {
				s.Response[0].Etag = etag[0]
			}
		}
	}

	return nil
}

func (s *secondaryLookup) populate(query url.Values) {
	s.Lookup.populateInclude(query)
	s.Lookup.populateExclude(query)
	s.Lookup.populateFreeform(query)
	s.Lookup.populateStreet(query)
	s.Lookup.populateCity(query)
	s.Lookup.populateState(query)
	s.Lookup.populateZIPCode(query)
}

////////////////////////////////////////////////////////////////////////////////////////

type secondaryCountLookup struct {
	*Lookup
	Response []*SecondaryCountResponse
}

func (s *secondaryCountLookup) getSmartyKey() string {
	if len(s.Lookup.SmartyKey) == 0 {
		return smartyKeyBypass
	}
	return s.SmartyKey
}

func (s *secondaryCountLookup) getDataSet() string {
	return secondaryData
}

func (s *secondaryCountLookup) getDataSubset() string {
	return secondaryDataCount
}

func (s *secondaryCountLookup) getLookup() *Lookup {
	return s.Lookup
}

func (s *secondaryCountLookup) getResponse() interface{} {
	return s.Response
}

func (s *secondaryCountLookup) unmarshalResponse(bytes []byte, header http.Header) error {
	if err := json.Unmarshal(bytes, &s.Response); err != nil {
		return err
	}

	if header != nil {
		if etag, found := header[lookupETagHeader]; found {
			if len(etag) > 0 && len(s.Response) > 0 {
				s.Response[0].Etag = etag[0]
			}
		}
	}

	return nil
}

func (s *secondaryCountLookup) populate(query url.Values) {
	s.Lookup.populateInclude(query)
	s.Lookup.populateExclude(query)
	s.Lookup.populateFreeform(query)
	s.Lookup.populateStreet(query)
	s.Lookup.populateCity(query)
	s.Lookup.populateState(query)
	s.Lookup.populateZIPCode(query)
}

const (
	financialDataSubset = "financial"
	principalDataSubset = "principal"
	propertyDataSet     = "property"
	geoReferenceDataSet = "geo-reference"
	secondaryData       = "secondary"
	secondaryDataCount  = "count"
	emptyDataSubset     = ""
	smartyKeyBypass     = "search"
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

func (l Lookup) populateFreeform(query url.Values) {
	if len(l.Freeform) > 0 {
		query.Set("freeform", l.Freeform)
	}
}

func (l Lookup) populateStreet(query url.Values) {
	if len(l.Street) > 0 {
		query.Set("street", l.Street)
	}
}

func (l Lookup) populateCity(query url.Values) {
	if len(l.City) > 0 {
		query.Set("city", l.City)
	}
}

func (l Lookup) populateState(query url.Values) {
	if len(l.State) > 0 {
		query.Set("state", l.State)
	}
}

func (l Lookup) populateZIPCode(query url.Values) {
	if len(l.ZIPCode) > 0 {
		query.Set("zipcode", l.ZIPCode)
	}
}
