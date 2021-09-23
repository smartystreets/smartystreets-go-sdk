package autocomplete_pro

import (
	"net/url"
	"strconv"
	"strings"
)

// Lookup represents all input fields documented here:
// https://smartystreets.com/docs/cloud/us-autocomplete-pro-api#http-request-input-fields
type (
	Lookup struct {
		Search        string
		Source        string
		MaxResults    int
		CityFilter    []string
		StateFilter   []string
		ZIPFilter     []string
		ExcludeStates []string
		PreferCity    []string
		PreferState   []string
		PreferZIP     []string
		PreferRatio   int
		Geolocation   Geolocation

		Results []*Suggestion
	}
	Geolocation string
)

const (
	GeolocateCity Geolocation = "city"
	GeolocateNone Geolocation = "none"
)

func (l Lookup) populate(query url.Values) {
	l.populateSearch(query)
	l.populateMaxResults(query)
	l.populateCityFilter(query)
	l.populateStateFilter(query)
	l.populateZIPFilter(query)
	l.populateExcludeStates(query)
	l.populatePreferCity(query)
	l.populatePreferState(query)
	l.populatePreferZIP(query)
	l.populatePreferRatio(query)
	l.populateGeolocation(query)
	l.populateSource(query)
}

func (l Lookup) populateSearch(query url.Values) {
	if len(l.Search) > 0 {
		query.Set("search", l.Search)
	}
}
func (l Lookup) populateMaxResults(query url.Values) {
	if l.MaxResults > 0 {
		query.Set("max_results", strconv.Itoa(l.MaxResults))
	}
}
func (l Lookup) populateCityFilter(query url.Values) {
	if len(l.CityFilter) > 0 {
		query.Set("include_only_cities", strings.Join(l.CityFilter, ","))
	}
}
func (l Lookup) populateStateFilter(query url.Values) {
	if len(l.StateFilter) > 0 {
		query.Set("include_only_states", strings.Join(l.StateFilter, ","))
	}
}
func (l Lookup) populateZIPFilter(query url.Values) {
	if len(l.ZIPFilter) > 0 {
		query.Set("include_only_zip_codes", strings.Join(l.ZIPFilter, ","))
	}
}
func (l Lookup) populateExcludeStates(query url.Values) {
	if len(l.ExcludeStates) > 0 {
		query.Set("exclude_states", strings.Join(l.ExcludeStates, ","))
	}
}
func (l Lookup) populatePreferCity(query url.Values) {
	if len(l.PreferCity) > 0 {
		query.Set("prefer_cities", strings.Join(l.PreferCity, ","))
	}
}
func (l Lookup) populatePreferState(query url.Values) {
	if len(l.PreferState) > 0 {
		query.Set("prefer_states", strings.Join(l.PreferState, ","))
	}
}
func (l Lookup) populatePreferZIP(query url.Values) {
	if len(l.PreferZIP) > 0 {
		query.Set("prefer_zip_codes", strings.Join(l.PreferZIP, ","))
	}
}
func (l Lookup) populatePreferRatio(query url.Values) {
	if l.PreferRatio > 0 {
		query.Set("prefer_ratio", strconv.Itoa(l.PreferRatio))
	}
}
func (l Lookup) populateGeolocation(query url.Values) {
	if len(l.ZIPFilter) > 0 || len(l.PreferZIP) > 0 {
		query.Set("prefer_geolocation", "none")
		return
	}

	switch l.Geolocation {
	case GeolocateCity:
		query.Set("prefer_geolocation", "city")
	case GeolocateNone:
		query.Set("prefer_geolocation", "none")
	}
}
func (l Lookup) populateSource(query url.Values) {
	if len(l.Source) > 0 {
		query.Set("source", l.Source)
	}
}
