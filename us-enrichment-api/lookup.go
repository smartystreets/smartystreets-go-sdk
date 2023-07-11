package us_enrichment

// Lookup fields defined here: https://smarty.com/docs/cloud/us-enrichment-api#http-request-input-fields
type Lookup struct {
	SmartyKey string
	DataSet   string

	Response []*Response
}
