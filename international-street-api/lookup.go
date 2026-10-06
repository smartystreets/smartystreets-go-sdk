package street

import (
	"net/url"
	"strings"
)

// Lookup contains all input fields defined here:
// https://smartystreets.com/docs/cloud/international-street-api#http-input-fields
type Lookup struct {
	InputID            string
	Country            string
	Geocode            bool
	Language           Language
	Features           []Feature
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

// normalized lets callers supply any letter case (eg. "Latin"); the API accepts
// either, so this is about sending one canonical form.
func (l Language) normalized() Language {
	return Language(strings.ToLower(string(l)))
}

type Feature string

const (
	FeatureOccupantUse              = Feature("occupant-use")
	FeatureGeocodeClassification    = Feature("geocode-classification")
	FeatureGeocodePrecisionEnhanced = Feature("geocode-precision-enhanced")
	FeatureSmartyKey                = Feature("smartykey")
)

// featuresValue renders the requested features as the comma-separated list
// the API expects; an empty slice yields "" so populate omits the parameter.
func (l *Lookup) featuresValue() string {
	s := make([]string, len(l.Features))
	for i, f := range l.Features {
		s[i] = string(f)
	}
	return strings.Join(s, ",")
}

/**************************************************************************/

func (l *Lookup) populate(query url.Values) {
	populate(query, "input_id", l.InputID)
	populate(query, "country", l.Country)
	populate(query, "geocode", boolString(l.Geocode))
	populate(query, "language", string(l.Language.normalized()))
	populate(query, "features", l.featuresValue())
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
