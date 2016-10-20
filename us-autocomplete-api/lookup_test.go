package autocomplete

import (
	"net/url"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type LookupSerializationFixture struct {
	*gunit.Fixture

	lookup *Lookup
	query  url.Values
}

func (f *LookupSerializationFixture) Setup() {
	f.lookup = new(Lookup)
	f.query = make(url.Values)
}
func (f *LookupSerializationFixture) populate() {
	f.lookup.populate(f.query)
}

func (f *LookupSerializationFixture) TestNothingToSerialize() {
	f.populate()

	f.So(f.query, should.BeEmpty)
}
func (f *LookupSerializationFixture) TestPrefix() {
	f.lookup.Prefix = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("prefix"), should.Equal, "Hello, World!")
}

func (f *LookupSerializationFixture) TestSuggestions() {
	f.lookup.MaxSuggestions = 7

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("suggestions"), should.Equal, "7")
}

func (f *LookupSerializationFixture) TestStateFilters() {
	f.lookup.StateFilter = append(f.lookup.StateFilter, "UT", "CA")

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("state_filter"), should.Equal, "UT,CA")
}

func (f *LookupSerializationFixture) TestCityFilters() {
	f.lookup.CityFilter = append(f.lookup.CityFilter, "Salt Lake City", "Provo")

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("city_filter"), should.Equal, "Salt Lake City,Provo")
}

func (f *LookupSerializationFixture) TestCityStatePreferences() {
	f.lookup.Preferences = append(f.lookup.Preferences, "Provo,UT", "Dallas,TX", "NV", "Sacramento")
	f.populate()
	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("prefer"), should.Equal, "Provo,UT;Dallas,TX;NV;Sacramento")
}

func (f *LookupSerializationFixture) TestGeolocateNone() {
	f.lookup.Geolocation = GeolocateNone

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("geolocate"), should.Equal, "false")
}
func (f *LookupSerializationFixture) TestGeolocateState() {
	f.lookup.Geolocation = GeolocateState

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("geolocate_precision"), should.Equal, "state")
}
func (f *LookupSerializationFixture) TestGeolocateCity_DefaultValue() {
	f.lookup.Geolocation = GeolocateCity

	f.populate()

	f.So(f.query, should.BeEmpty)
}
