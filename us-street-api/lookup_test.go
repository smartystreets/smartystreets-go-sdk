package street

import (
	"net/url"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
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
	this.So(this.encode(&Lookup{Street: "A"}), should.Resemble, url.Values{"street": {"A"}})
	this.So(this.encode(&Lookup{Street2: "A"}), should.Resemble, url.Values{"street2": {"A"}})
	this.So(this.encode(&Lookup{Secondary: "A"}), should.Resemble, url.Values{"secondary": {"A"}})
	this.So(this.encode(&Lookup{City: "A"}), should.Resemble, url.Values{"city": {"A"}})
	this.So(this.encode(&Lookup{State: "A"}), should.Resemble, url.Values{"state": {"A"}})
	this.So(this.encode(&Lookup{ZIPCode: "A"}), should.Resemble, url.Values{"zipcode": {"A"}})
	this.So(this.encode(&Lookup{LastLine: "A"}), should.Resemble, url.Values{"lastline": {"A"}})
	this.So(this.encode(&Lookup{Addressee: "A"}), should.Resemble, url.Values{"addressee": {"A"}})
	this.So(this.encode(&Lookup{Urbanization: "A"}), should.Resemble, url.Values{"urbanization": {"A"}})
	this.So(this.encode(&Lookup{InputID: "A"}), should.Resemble, url.Values{"input_id": {"A"}})
	this.So(this.encode(&Lookup{MaxCandidates: 2}), should.Resemble, url.Values{"candidates": {"2"}})
	this.So(this.encode(&Lookup{MatchStrategy: MatchInvalid}), should.Resemble, url.Values{"match": {"invalid"}})
}

func (this *LookupFixture) TestQueryStringEncoding_AllFieldsSerialized() {
	this.So(this.encode(&Lookup{
		MatchStrategy: MatchEnhanced,
		MaxCandidates: 0,
		InputID:       "InputID",
		ZIPCode:       "ZIPCode",
		LastLine:      "LastLine",
		Urbanization:  "Urbanization",
		Addressee:     "Addressee",
		Street:        "Street",
		City:          "City",
		Secondary:     "Secondary",
		State:         "State",
		Street2:       "Street2",
	}), should.Resemble, url.Values{
		"match":        {"enhanced"},
		"candidates":   {"5"},
		"input_id":     {"InputID"},
		"zipcode":      {"ZIPCode"},
		"lastline":     {"LastLine"},
		"urbanization": {"Urbanization"},
		"addressee":    {"Addressee"},
		"street":       {"Street"},
		"city":         {"City"},
		"secondary":    {"Secondary"},
		"state":        {"State"},
		"street2":      {"Street2"},
	})
}

func (this *LookupFixture) TestQueryStringEncoding_WithOutputFormat() {
	this.So(this.encode(&Lookup{OutputFormat: FormatDefault}), should.Resemble, url.Values{})
	this.So(this.encode(&Lookup{OutputFormat: FormatProjectUSA}), should.Resemble, url.Values{"format": {"project-usa"}})
}

func (this *LookupFixture) TestQueryStringEncoding_OutputFormatSerialized() {
	this.So(this.encode(&Lookup{
		OutputFormat: FormatProjectUSA,
	}), should.Resemble, url.Values{
		"format": {"project-usa"},
	})
}
