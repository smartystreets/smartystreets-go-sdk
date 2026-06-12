package us_reverse_geo

// Lookup fields defined here: https://smartystreets.com/docs/cloud/us-reverse-geo-api#http-request-input-fields
type Lookup struct {
	Latitude  float64
	Longitude float64
	Source    Source
	Response  Response
}

type Source string

const (
	SourceAll    Source = "all"
	SourcePostal Source = "postal"
)
