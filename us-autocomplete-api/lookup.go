package autocomplete

import (
	"net/url"
	"strconv"
	"strings"
)

type (
	// Lookup represents all input fields documented here:
	// https://smartystreets.com/docs/us-autocomplete-api#http-request-input-fields
	Lookup struct {
		Prefix         string
		MaxSuggestions int
		CityFilter     []string
		StateFilter    []string
		Preferences    []string
		Geolocation    geolocation

		Results []*Suggestion
	}

	geolocation int
)

const (
	GeolocateCity geolocation = iota
	GeolocateState
	GeolocateNone
)

/**************************************************************************/

func (l *Lookup) populate(query url.Values) {
	l.populatePrefix(query)
	l.populateSuggestions(query)
	l.populateStateFilter(query)
	l.populateCityFilter(query)
	l.populatePreferences(query)
	l.populateGeolocation(query)
}

func (l *Lookup) populatePrefix(query url.Values) {
	if len(l.Prefix) > 0 {
		query.Set("prefix", l.Prefix)
	}
}
func (l *Lookup) populateStateFilter(query url.Values) {
	if len(l.StateFilter) > 0 {
		query.Set("state_filter", strings.Join(l.StateFilter, ","))
	}
}
func (l *Lookup) populateCityFilter(query url.Values) {
	if len(l.CityFilter) > 0 {
		query.Set("city_filter", strings.Join(l.CityFilter, ","))
	}
}
func (l *Lookup) populateSuggestions(query url.Values) {
	if l.MaxSuggestions > 0 {
		query.Set("suggestions", strconv.Itoa(l.MaxSuggestions))
	}
}
func (l *Lookup) populatePreferences(query url.Values) {
	if len(l.Preferences) > 0 {
		query.Set("prefer", strings.Join(l.Preferences, ";"))
	}
}
func (l *Lookup) populateGeolocation(query url.Values) {
	switch l.Geolocation {
	case GeolocateNone:
		query.Set("geolocate", "false")
	case GeolocateState:
		query.Set("geolocate_precision", "state")
	}
}
