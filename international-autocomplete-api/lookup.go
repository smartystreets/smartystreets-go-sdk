package international_autocomplete_api

import (
	"net/url"
	"strconv"
)

const (
	maxResultsDefault = 5
	distanceDefault   = 5
)

type Lookup struct {
	Country    string
	Search     string
	AddressID  string
	MaxResults int
	Locality   string
	PostalCode string
	Result     *Result
}

func (l Lookup) populate(query url.Values) {
	l.populateCountry(query)
	l.populateSearch(query)
	l.populateMaxResults(query)
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
func (l Lookup) populateMaxResults(query url.Values) {
	maxResults := l.MaxResults
	if maxResults < 1 {
		maxResults = maxResultsDefault
	}
	query.Set("max_results", strconv.Itoa(maxResults))
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
