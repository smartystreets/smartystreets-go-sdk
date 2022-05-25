package zipcode

import "net/url"

// Lookup contains all input fields defined here:
// https://smartystreets.com/docs/us-zipcode-api#input-fields
type Lookup struct {
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	ZIPCode string `json:"zipcode,omitempty"`
	InputID string `json:"input_id,omitempty"`

	Result *Result `json:"result,omitempty"`
}

func (l *Lookup) encodeQueryString(query url.Values) {
	encode(query, l.City, "city")
	encode(query, l.State, "state")
	encode(query, l.ZIPCode, "zipcode")
	encode(query, l.InputID, "input_id")
}

func encode(query url.Values, source string, target string) {
	if source != "" {
		query.Set(target, source)
	}
}
