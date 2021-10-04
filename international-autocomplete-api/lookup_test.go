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

func (f *LookupFixture) TestNothingToSerialize() {
	f.populate()

	f.So(f.query, should.BeEmpty)
}
func (f *LookupFixture) TestCountry() {
	f.lookup.Country = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("country"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestSearch() {
	f.lookup.Search = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("search"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestAdministrativeArea() {
	f.lookup.AdministrativeArea = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("include_only_administrative_area"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestLocality() {
	f.lookup.Locality = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("include_only_locality"), should.Equal, "Hello, World!")
}
func (f *LookupFixture) TestPostalCode() {
	f.lookup.PostalCode = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("include_only_postal_code"), should.Equal, "Hello, World!")
}
