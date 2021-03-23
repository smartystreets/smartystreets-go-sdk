package autocomplete_pro

type suggestionListing struct {
	Listing []*Suggestion `json:"suggestions"`
}

// Suggestion is the primary element in the response array.
// Online documentation: https://smartystreets.com/docs/us-autocomplete-pro-api#http-response
type Suggestion struct {
	StreetLine string `json:"street_line"`
	Secondary  string `json:"secondary"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZIPCode    string `json:"zipcode"`
	Entries    int    `json:"entries"`
}
