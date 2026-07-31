package international_postal_code

import (
	"net/url"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
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

// The Language field is not yet serialized (see Lookup), so this covers the
// helper directly to keep it honest until the field goes live.
func (f *LookupSerializationFixture) TestLanguageNormalized() {
	f.So(Language("Latin").normalized(), should.Equal, Latin)
	f.So(Language("NATIVE").normalized(), should.Equal, Native)
	f.So(Latin.normalized(), should.Equal, Latin)
	f.So(Language("").normalized(), should.Equal, Language(""))
}

func (f *LookupSerializationFixture) TestFullLookup() {
	f.lookup.InputID = "Hello, World!"
	f.lookup.Country = "CAN"
	f.lookup.Locality = "Toronto"
	f.lookup.AdministrativeArea = "ON"
	f.lookup.PostalCode = "ABC DEF"
	f.populate()

	f.So(f.query, should.HaveLength, 5)
	f.So(f.query.Get("input_id"), should.Equal, "Hello, World!")
	f.So(f.query.Get("country"), should.Equal, "CAN")
	f.So(f.query.Get("locality"), should.Equal, "Toronto")
	f.So(f.query.Get("administrative_area"), should.Equal, "ON")
	f.So(f.query.Get("postal_code"), should.Equal, "ABC DEF")
}
