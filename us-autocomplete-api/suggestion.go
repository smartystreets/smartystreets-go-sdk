package autocomplete

// Suggestion is the primary element in the response array.
// Online documentation: https://smartystreets.com/docs/us-autocomplete-api#http-response
type Suggestion struct {
	Text       string `json:"text"`
	StreetLine string `json:"street_line"`
	City       string `json:"city"`
	State      string `json:"state"`
}
