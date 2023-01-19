package international_autocomplete_api

import (
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestLookupFixture(t *testing.T) {
	gunit.Run(new(LookupFixture), t)
}

type LookupFixture struct {
	*gunit.Fixture

	lookup *Lookup
	query  url.Values
}

func (f *LookupFixture) Setup() {
	f.lookup = new(Lookup)
	f.query = make(url.Values)
}
func (f *LookupFixture) populate() {
	f.lookup.populate(f.query)
}

func (f *LookupFixture) TestDefaults() {
	f.populate()

	f.So(f.query, should.HaveLength, 2)
	f.So(f.query.Get("max_results"), should.Equal, "5")
	f.So(f.query.Get("distance"), should.Equal, "5")
}
func (f *LookupFixture) TestCountry() {
	f.lookup.Country = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 3)
	f.So(f.query.Get("country"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestSearch() {
	f.lookup.Search = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 3)
	f.So(f.query.Get("search"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestMaxResults() {
	f.lookup.MaxResults = 7
	f.populate()

	f.So(f.query.Get("max_results"), should.Equal, "7")
}
func (f *LookupFixture) TestDistance() {
	f.lookup.Distance = 3
	f.populate()

	f.So(f.query.Get("distance"), should.Equal, "3")
}
func (f *LookupFixture) TestGeolocation() {
	typeList := []InternationalGeolocateType{AdminArea, Locality, PostalCode, Geocodes, None}
	for _, geoLocateType := range typeList {
		f.lookup.Geolocation = geoLocateType
		f.populate()
		f.So(f.query.Get("geolocation"), should.Equal, string(geoLocateType))
	}
}
func (f *LookupFixture) TestAdministrativeArea() {
	f.lookup.AdministrativeArea = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 3)
	f.So(f.query.Get("include_only_administrative_area"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestLocality() {
	f.lookup.Locality = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 3)
	f.So(f.query.Get("include_only_locality"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestPostalCode() {
	f.lookup.PostalCode = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 3)
	f.So(f.query.Get("include_only_postal_code"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestLatitude() {
	f.lookup.Latitude = 123.458757987986

	f.populate()

	f.So(f.query.Get("latitude"), should.Equal, "123.45875799") // Here we only care about 8 digits of accuracy
}
func (f *LookupFixture) TestLongitude() {
	f.lookup.Longitude = -134.877532234

	f.populate()

	f.So(f.query.Get("longitude"), should.Equal, "-134.87753223") // Here we only care about 8 digits of accuracy
}
