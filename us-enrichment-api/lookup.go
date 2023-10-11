package us_enrichment

type Lookup struct {
	SmartyKey  string
	DataSet    string
	DataSubSet string

	FinancialResponse []FinancialResponse
	PrincipalResponse []PrincipalResponse
}
