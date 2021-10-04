package international_autocomplete_api

import "net/url"

type Lookup struct {
	Country            string
	Search             string
	AdministrativeArea string
	Locality           string
	PostalCode         string
	Results            []*Suggestion
}

func (l Lookup) populate(query url.Values) {
	l.populateCountry(query)
	l.populateSearch(query)
	l.populateAdministrativeArea(query)
	l.populateLocality(query)
	l.populatePostalCode(query)
}
func (l Lookup) populateCountry(query url.Values) {
	if len(l.Country) > 0 {
		query.Set("country", l.Country)
	}
}
func (l Lookup) populateSearch(query url.Values) {
	if len(l.Search) > 0 {
		query.Set("search", l.Search)
	}
}
func (l Lookup) populateAdministrativeArea(query url.Values) {
	if len(l.AdministrativeArea) > 0 {
		query.Set("include_only_administrative_area", l.AdministrativeArea)
	}
}
func (l Lookup) populateLocality(query url.Values) {
	if len(l.Locality) > 0 {
		query.Set("include_only_locality", l.Locality)
	}
}
func (l Lookup) populatePostalCode(query url.Values) {
	if len(l.PostalCode) > 0 {
		query.Set("include_only_postal_code", l.PostalCode)
	}
}
