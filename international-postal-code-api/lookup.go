package international_postal_code

import (
	"net/url"
)

// Lookup contains all input fields defined here:
// https://smartystreets.com/docs/cloud/international-street-api#http-input-fields
type Lookup struct {
	InputID string
	Country string
	//Language           Language // future
	//Features           string   // future
	Locality           string
	AdministrativeArea string
	PostalCode         string

	Results []*Candidate
}

/**************************************************************************/

type Language string

const (
	Native = Language("native")
	Latin  = Language("latin")
)

/**************************************************************************/

func (l *Lookup) populate(query url.Values) {
	populate(query, "input_id", l.InputID)
	populate(query, "country", l.Country)
	//populate(query, "language", string(l.Language)) // future
	//populate(query, "features", l.Features)         // future
	populate(query, "locality", l.Locality)
	populate(query, "administrative_area", l.AdministrativeArea)
	populate(query, "postal_code", l.PostalCode)
}

func populate(query url.Values, key, value string) {
	if len(value) > 0 {
		query.Set(key, value)
	}
}
