package us_zipcode

// Lookup contains all input fields defined here:
// https://smartystreets.com/docs/us-street-api#input-fields
type Lookup struct {
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	ZIPCode string `json:"zipcode,omitempty"`
	InputID string `json:"input_id,omitempty"`

	Result *Result `json:"-"`
}
