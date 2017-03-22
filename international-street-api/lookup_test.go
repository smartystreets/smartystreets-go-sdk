package international

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

func (f *LookupFixture) TestPopulateFromStringFields() {
	f.lookup.InputID = "INPUT_ID"
	f.lookup.Address1 = "1"
	f.lookup.Address2 = "2"
	f.lookup.Address3 = "3"
	f.lookup.Address4 = "4"
	f.lookup.Freeform = "FREEFORM"
	f.lookup.AdministrativeArea = "ADMINISTRATIVE_AREA"
	f.lookup.Country = "COUNTRY"
	f.lookup.Locality = "LOCALITY"
	f.lookup.Organization = "ORGANIZATION"

	f.lookup.populate(f.query)

	f.So(f.query.Get("input_id"), should.Equal, "INPUT_ID")
	f.So(f.query.Get("address1"), should.Equal, "1")
	f.So(f.query.Get("address2"), should.Equal, "2")
	f.So(f.query.Get("address3"), should.Equal, "3")
	f.So(f.query.Get("address4"), should.Equal, "4")
	f.So(f.query.Get("freeform"), should.Equal, "FREEFORM")
	f.So(f.query.Get("administrative_area"), should.Equal, "ADMINISTRATIVE_AREA")
	f.So(f.query.Get("country"), should.Equal, "COUNTRY")
	f.So(f.query.Get("locality"), should.Equal, "LOCALITY")
	f.So(f.query.Get("organization"), should.Equal, "ORGANIZATION")
}

func (f *LookupFixture) TestPopulateGeocodeFalse() {
	f.lookup.Geocode = false
	f.lookup.populate(f.query)
	f.So(f.query.Get("geocode"), should.Equal, "false")
}

func (f *LookupFixture) TestPopulateGeocodeTrue() {
	f.lookup.Geocode = true
	f.lookup.populate(f.query)
	f.So(f.query.Get("geocode"), should.Equal, "true")
}

func (f *LookupFixture) TestPopulateOutputLanguageDefault() {
	f.lookup.OutputLanguage = OutputDefault
	f.lookup.populate(f.query)
	f.So(f.query.Get("language"), should.Equal, "")
}

func (f *LookupFixture) TestPopulateOutputLanguageLatin() {
	f.lookup.OutputLanguage = OutputLatin
	f.lookup.populate(f.query)
	f.So(f.query.Get("language"), should.Equal, "latin")
}

func (f *LookupFixture) TestPopulateOutputLanguageNative() {
	f.lookup.OutputLanguage = OutputNative
	f.lookup.populate(f.query)
	f.So(f.query.Get("language"), should.Equal, "native")
}
