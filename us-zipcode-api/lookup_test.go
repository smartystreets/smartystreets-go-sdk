package zipcode

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
}

func (this *LookupFixture) encode(lookup *Lookup) url.Values {
	var values = make(url.Values)
	lookup.encodeQueryString(values)
	return values
}

func (this *LookupFixture) TestQueryStringEncoding_OnlySerializeNonDefaultFields() {
	this.So(this.encode(&Lookup{City: "A"}), should.Resemble, url.Values{"city": {"A"}})
	this.So(this.encode(&Lookup{State: "A"}), should.Resemble, url.Values{"state": {"A"}})
	this.So(this.encode(&Lookup{ZIPCode: "A"}), should.Resemble, url.Values{"zipcode": {"A"}})
	this.So(this.encode(&Lookup{InputID: "A"}), should.Resemble, url.Values{"input_id": {"A"}})
}

func (this *LookupFixture) TestQueryStringEncoding_AllFieldsSerialized() {
	this.So(this.encode(&Lookup{
		InputID: "InputID",
		ZIPCode: "ZIPCode",
		City:    "City",
		State:   "State",
	}), should.Resemble, url.Values{
		"input_id": {"InputID"},
		"zipcode":  {"ZIPCode"},
		"city":     {"City"},
		"state":    {"State"},
	})
}
