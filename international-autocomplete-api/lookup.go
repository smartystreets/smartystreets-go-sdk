package international_autocomplete_api

import (
	"net/url"
	"strconv"
)

const (
	maxResultsDefault      = 5
	maxGroupResultsDefault = 100
	distanceDefault        = 5
)

type Lookup struct {
	Country         string
	Search          string
	AddressID       string
	MaxResults      int
	MaxGroupResults int
	Geolocation     string
	Locality        string
	PostalCode      string
	Result          *Result
}

func (l Lookup) populate(query url.Values) {
	l.populateCountry(query)
	l.populateSearch(query)
	l.populateMaxResults(query)
	l.populateMaxGroupResults(query)
	l.populateGeolocation(query)
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
func (l Lookup) populateMaxGroupResults(query url.Values) {
	maxGroupResults := l.MaxGroupResults
	if maxGroupResults < 1 {
		maxGroupResults = maxGroupResultsDefault
	}
	query.Set("max_group_results", strconv.Itoa(maxGroupResults))
}
func (l Lookup) populateGeolocation(query url.Values) {
	if len(l.Geolocation) > 0 {
		query.Set("geolocation", l.Geolocation)
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
