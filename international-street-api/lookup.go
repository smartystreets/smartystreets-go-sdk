package street

import "net/url"

// Lookup contains all input fields defined here:
// https://smartystreets.com/docs/cloud/international-street-api#http-input-fields
type Lookup struct {
	InputID            string
	Country            string
	Geocode            bool
	Language           Language
	Freeform           string
	Address1           string
	Address2           string
	Address3           string
	Address4           string
	Organization       string
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
	populate(query, "geocode", boolString(l.Geocode))
	populate(query, "language", string(l.Language))
	populate(query, "freeform", l.Freeform)
	populate(query, "address1", l.Address1)
	populate(query, "address2", l.Address2)
	populate(query, "address3", l.Address3)
	populate(query, "address4", l.Address4)
	populate(query, "organization", l.Organization)
	populate(query, "locality", l.Locality)
	populate(query, "administrative_area", l.AdministrativeArea)
	populate(query, "postal_code", l.PostalCode)
}

func boolString(value bool) string {
	if value {
		return "true"
	}
	return ""
}

func populate(query url.Values, key, value string) {
	if len(value) > 0 {
		query.Set(key, value)
	}
}
