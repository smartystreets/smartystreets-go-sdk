package us_enrichment

import (
	"encoding/json"
)

type enrichmentLookup interface {
	GetSmartyKey() string
	GetDataSet() string
	GetDataSubset() string

	UnmarshalResponse([]byte) error
}

////////////////////////////////////////////////////////////////////////////////////////

type financialLookup struct {
	SmartyKey string
	Response  []*FinancialResponse
}

func (f *financialLookup) GetSmartyKey() string {
	return f.SmartyKey
}

func (f *financialLookup) GetDataSet() string {
	return propertyDataSet
}

func (f *financialLookup) GetDataSubset() string {
	return financialDataSubset
}

func (f *financialLookup) UnmarshalResponse(bytes []byte) error {
	return json.Unmarshal(bytes, &f.Response)
}

////////////////////////////////////////////////////////////////////////////////////////

type principalLookup struct {
	SmartyKey string
	Response  []*PrincipalResponse
}

func (p *principalLookup) GetSmartyKey() string {
	return p.SmartyKey
}

func (p *principalLookup) GetDataSet() string {
	return propertyDataSet
}

func (p *principalLookup) GetDataSubset() string {
	return principalDataSubset
}

func (p *principalLookup) UnmarshalResponse(bytes []byte) error {
	return json.Unmarshal(bytes, &p.Response)
}

////////////////////////////////////////////////////////////////////////////////////////

const (
	financialDataSubset = "financial"
	principalDataSubset = "principal"
	propertyDataSet     = "property"
)
