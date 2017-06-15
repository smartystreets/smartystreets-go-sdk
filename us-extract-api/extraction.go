package extract

import "github.com/smartystreets/smartystreets-go-sdk/us-street-api"

// Result, Metadata, and ExtractedAddress represent all output fields documented here:
// https://smartystreets.com/docs/cloud/us-extract-api#http-response
type Result struct {
	Metadata  Metadata           `json:"meta"`
	Addresses []ExtractedAddress `json:"addresses"`
}

type Metadata struct {
	Lines                   int  `json:"lines"`
	Characters              int  `json:"character_count"`
	Bytes                   int  `json:"bytes"`
	Addresses               int  `json:"address_count"`
	VerifiedAddresses       int  `json:"verified_count"`
	ContainsNonASCIIUnicode bool `json:"unicode"`
}

type ExtractedAddress struct {
	Text      string             `json:"text"`
	Verified  bool               `json:"verified"`
	Line      int                `json:"line"`
	Start     int                `json:"start"`
	End       int                `json:"end"`
	APIOutput []street.Candidate `json:"api_output"`
}
