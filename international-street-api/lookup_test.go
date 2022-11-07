package street

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
func (f *LookupSerializationFixture) TestInputID() {
	f.lookup.InputID = "Hello, World!"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("input_id"), should.Equal, "Hello, World!")
}

func (f *LookupSerializationFixture) TestCountry() {
	f.lookup.Country = "Anvilania"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("country"), should.Equal, "Anvilania")
}

func (f *LookupSerializationFixture) TestGeocodeFalse() {
	f.lookup.Geocode = false
	f.populate()
	f.So(f.query, should.BeEmpty)
}

func (f *LookupSerializationFixture) TestGeocodeTrue() {
	f.lookup.Geocode = true

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("geocode"), should.Equal, "true")
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

func (f *LookupSerializationFixture) TestOrganization() {
	f.lookup.Organization = "SmartyStreets"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("organization"), should.Equal, "SmartyStreets")
}

func (f *LookupSerializationFixture) TestFreeform() {
	f.lookup.Freeform = "free"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("freeform"), should.Equal, "free")
}

func (f *LookupSerializationFixture) TestAddress1() {
	f.lookup.Address1 = "one"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("address1"), should.Equal, "one")
}

func (f *LookupSerializationFixture) TestAddress2() {
	f.lookup.Address2 = "two"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("address2"), should.Equal, "two")
}

func (f *LookupSerializationFixture) TestAddress3() {
	f.lookup.Address3 = "three"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("address3"), should.Equal, "three")
}

func (f *LookupSerializationFixture) TestAddress4() {
	f.lookup.Address4 = "four"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("address4"), should.Equal, "four")
}

func (f *LookupSerializationFixture) TestLocality() {
	f.lookup.Locality = "Provo"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("locality"), should.Equal, "Provo")
}

func (f *LookupSerializationFixture) TestAdministrativeArea() {
	f.lookup.AdministrativeArea = "Admin"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("administrative_area"), should.Equal, "Admin")
}

func (f *LookupSerializationFixture) TestPostalCode() {
	f.lookup.PostalCode = "12345"

	f.populate()

	f.So(f.query, should.HaveLength, 1)
	f.So(f.query.Get("postal_code"), should.Equal, "12345")
}
