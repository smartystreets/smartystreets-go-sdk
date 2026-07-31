package international_postal_code

import (
	"net/url"
	"strings"
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

// normalized lets callers supply any letter case (eg. "Latin"); the API accepts
// either, so this is about sending one canonical form.
func (l Language) normalized() Language {
	return Language(strings.ToLower(string(l)))
}

/**************************************************************************/

func (l *Lookup) populate(query url.Values) {
	populate(query, "input_id", l.InputID)
	populate(query, "country", l.Country)
	//populate(query, "language", string(l.Language.normalized())) // future
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
