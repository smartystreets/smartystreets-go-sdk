package international_autocomplete_api

import (
	"math"
	"net/url"
	"strconv"
)

const (
	maxResultsDefault = 5
	distanceDefault   = 5
)

type Lookup struct {
	Country            string
	Search             string
	MaxResults         int
	Distance           int
	Geolocation        InternationalGeolocateType
	AdministrativeArea string
	Locality           string
	PostalCode         string
	Latitude           float64
	Longitude          float64
	Result             *Result
}

func (l Lookup) populate(query url.Values) {
	l.populateCountry(query)
	l.populateSearch(query)
	l.populateMaxResults(query)
	l.populateDistance(query)
	l.populateGeolocation(query)
	l.populateAdministrativeArea(query)
	l.populateLocality(query)
	l.populatePostalCode(query)
	l.populateLatitude(query)
	l.populateLongitude(query)
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
func (l Lookup) populateDistance(query url.Values) {
	distance := l.Distance
	if distance < 1 {
		distance = distanceDefault
	}
	query.Set("distance", strconv.Itoa(distance))
}
func (l Lookup) populateGeolocation(query url.Values) {
	if l.Geolocation != None {
		query.Set("geolocation", string(l.Geolocation))
	} else {
		query.Del("geolocation")
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
func (l Lookup) populateLatitude(query url.Values) {
	if math.Floor(l.Latitude) != 0 {
		query.Set("latitude", strconv.FormatFloat(l.Latitude, 'f', 8, 64))
	}
}
func (l Lookup) populateLongitude(query url.Values) {
	if math.Floor(l.Longitude) != 0 {
		query.Set("longitude", strconv.FormatFloat(l.Longitude, 'f', 8, 64))
	}
}

type InternationalGeolocateType string

const (
	AdminArea  = InternationalGeolocateType("adminarea")
	Locality   = InternationalGeolocateType("locality")
	PostalCode = InternationalGeolocateType("postalcode")
	Geocodes   = InternationalGeolocateType("geocodes")
	None       = InternationalGeolocateType("")
)
