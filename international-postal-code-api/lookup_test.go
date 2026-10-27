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

func (f *LookupSerializationFixture) TestLanguageBlank() {
	f.lookup.Language = Language("")
	f.populate()
	f.So(f.query, should.BeEmpty)
}

func (f *LookupSerializationFixture) TestLanguageLatin() {
	f.lookup.Language = Latin

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("language"), should.Equal, "latin")
}

func (f *LookupSerializationFixture) TestLanguageNative() {
	f.lookup.Language = Native

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("language"), should.Equal, "native")
}

func (f *LookupSerializationFixture) TestLanguageMixedCase() {
	f.lookup.Language = Language("NATIVE")

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("language"), should.Equal, "native")
}

func (f *LookupSerializationFixture) TestMaxResultsUnset() {
	f.lookup.MaxResults = 0
	f.populate()
	f.So(f.query, should.BeEmpty)
}

func (f *LookupSerializationFixture) TestMaxResults() {
	f.lookup.MaxResults = 1000

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("max_results"), should.Equal, "1000")
}

func (f *LookupSerializationFixture) TestFullLookup() {
	f.lookup.InputID = "Hello, World!"
	f.lookup.Country = "CAN"
	f.lookup.Language = Native
	f.lookup.Locality = "Toronto"
	f.lookup.AdministrativeArea = "ON"
	f.lookup.PostalCode = "ABC DEF"
	f.lookup.MaxResults = 25
	f.populate()

	f.So(f.query, should.HaveLength, 7)
	f.So(f.query.Get("input_id"), should.Equal, "Hello, World!")
	f.So(f.query.Get("country"), should.Equal, "CAN")
	f.So(f.query.Get("language"), should.Equal, "native")
	f.So(f.query.Get("locality"), should.Equal, "Toronto")
	f.So(f.query.Get("administrative_area"), should.Equal, "ON")
	f.So(f.query.Get("postal_code"), should.Equal, "ABC DEF")
	f.So(f.query.Get("max_results"), should.Equal, "25")
}
