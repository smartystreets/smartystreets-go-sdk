package autocomplete_pro

import (
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestLookupSerializationFixture(t *testing.T) {
	gunit.Run(new(LookupSerializationFixture), t)
}

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
func (f *LookupSerializationFixture) TestSearch() {
	f.lookup.Search = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("search"), should.Equal, "Hello, World!")
}
func (f *LookupSerializationFixture) TestSource() {
	f.lookup.Source = "all"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("source"), should.Equal, "all")
}

func (f *LookupSerializationFixture) TestSuggestions() {
	f.lookup.MaxResults = 7

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("max_results"), should.Equal, "7")
}

func (f *LookupSerializationFixture) TestStateFilters() {
	f.lookup.StateFilter = append(f.lookup.StateFilter, "UT", "CA")

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("include_only_states"), should.Equal, "UT,CA")
}

func (f *LookupSerializationFixture) TestCityFilters() {
	f.lookup.CityFilter = append(f.lookup.CityFilter, "Salt Lake City", "Provo")

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("include_only_cities"), should.Equal, "Salt Lake City,Provo")
}

func (f *LookupSerializationFixture) TestZIPFilters() {
	f.lookup.ZIPFilter = append(f.lookup.ZIPFilter, "84660", "84058")

	f.populate()

	f.So(f.query, should.HaveLength, 2)
	f.So(f.query.Get("include_only_zip_codes"), should.Equal, "84660,84058")
}

func (f *LookupSerializationFixture) TestExcludeStates() {
	f.lookup.ExcludeStates = append(f.lookup.ExcludeStates, "UT", "CA")

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("exclude_states"), should.Equal, "UT,CA")
}

func (f *LookupSerializationFixture) TestPreferState() {
	f.lookup.PreferState = append(f.lookup.PreferState, "UT", "CA")

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("prefer_states"), should.Equal, "UT,CA")
}

func (f *LookupSerializationFixture) TestPreferCity() {
	f.lookup.PreferCity = append(f.lookup.PreferCity, "Salt Lake City", "Provo")

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("prefer_cities"), should.Equal, "Salt Lake City,Provo")
}

func (f *LookupSerializationFixture) TestPreferZIP() {
	f.lookup.PreferZIP = append(f.lookup.PreferZIP, "84660", "84058")

	f.populate()

	f.So(f.query, should.HaveLength, 2)
	f.So(f.query.Get("prefer_zip_codes"), should.Equal, "84660,84058")
}

func (f *LookupSerializationFixture) TestPreferRatio() {
	f.lookup.PreferRatio = 6

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("prefer_ratio"), should.Equal, "6")
}

func (f *LookupSerializationFixture) TestGeolocateNone() {
	f.lookup.Geolocation = GeolocateNone

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("prefer_geolocation"), should.Equal, "none")
}

func (f *LookupSerializationFixture) TestGeolocateCity_DefaultValue() {
	f.lookup.Geolocation = GeolocateCity

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("prefer_geolocation"), should.Equal, "city")
}
