package us_enrichment

import (
	"encoding/json"
)

type Lookup struct {
	SmartyKey  string
	DataSubset string
	Include    string
	Exclude    string

	ETag string
}

type enrichmentLookup interface {
	GetSmartyKey() string
	GetDataSet() string
	GetDataSubset() string
	GetLookup() *Lookup
	GetResponse() any
	UnmarshalResponse([]byte) error
}

////////////////////////////////////////////////////////////////////////////////////////

type financialLookup struct {
	Lookup   *Lookup
	Response []*FinancialResponse
}

func (f *financialLookup) GetSmartyKey() string {
	return f.Lookup.SmartyKey
}

func (f *financialLookup) GetDataSet() string {
	return propertyDataSet
}

func (f *financialLookup) GetDataSubset() string {
	return f.Lookup.DataSubset
}

func (f *financialLookup) GetLookup() *Lookup {
	return f.Lookup
}

func (f *financialLookup) GetResponse() any {
	return f.Response
}

func (f *financialLookup) UnmarshalResponse(bytes []byte) error {
	return json.Unmarshal(bytes, &f.Response)
}

////////////////////////////////////////////////////////////////////////////////////////

type principalLookup struct {
	Lookup   *Lookup
	Response []*PrincipalResponse
}

func (p *principalLookup) GetSmartyKey() string {
	return p.Lookup.SmartyKey
}

func (p *principalLookup) GetDataSet() string {
	return propertyDataSet
}

func (p *principalLookup) GetDataSubset() string {
	return p.Lookup.DataSubset
}

func (p *principalLookup) GetLookup() *Lookup {
	return p.Lookup
}

func (f *principalLookup) GetResponse() any {
	return f.Response
}

func (p *principalLookup) UnmarshalResponse(bytes []byte) error {
	return json.Unmarshal(bytes, &p.Response)
}

////////////////////////////////////////////////////////////////////////////////////////

const (
	FinancialDataSubset = "financial"
	PrincipalDataSubset = "principal"
	propertyDataSet     = "property"
)
