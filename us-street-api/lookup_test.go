package street

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"slices"
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
	this.So(this.encode(&Lookup{Street: "A"}), should.Resemble, url.Values{"street": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{Street2: "A"}), should.Resemble, url.Values{"street2": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{Secondary: "A"}), should.Resemble, url.Values{"secondary": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{City: "A"}), should.Resemble, url.Values{"city": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{State: "A"}), should.Resemble, url.Values{"state": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{ZIPCode: "A"}), should.Resemble, url.Values{"zipcode": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{LastLine: "A"}), should.Resemble, url.Values{"lastline": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{Addressee: "A"}), should.Resemble, url.Values{"addressee": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{Urbanization: "A"}), should.Resemble, url.Values{"urbanization": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{InputID: "A"}), should.Resemble, url.Values{"input_id": {"A"}, "match": {"enhanced"}, "candidates": {"5"}})
	this.So(this.encode(&Lookup{MaxCandidates: 2}), should.Resemble, url.Values{"candidates": {"2"}, "match": {"enhanced"}})
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
	this.So(this.encode(&Lookup{OutputFormat: FormatDefault}), should.Resemble, url.Values{"match": {"enhanced"}, "candidates": {"5"}, "format": {"default"}})
	this.So(this.encode(&Lookup{OutputFormat: FormatProjectUSA}), should.Resemble, url.Values{"format": {"project-usa"}, "match": {"enhanced"}, "candidates": {"5"}})
}

func (this *LookupFixture) TestQueryStringEncoding_OutputFormatSerialized() {
	this.So(this.encode(&Lookup{
		OutputFormat: FormatProjectUSA,
	}), should.Resemble, url.Values{
		"format":     {"project-usa"},
		"match":      {"enhanced"},
		"candidates": {"5"},
	})
}

func (this *LookupFixture) TestQueryStringEncoding_CountySourceSerialized() {
	this.So(this.encode(&Lookup{
		CountySource: GeographicCounty,
	}), should.Resemble, url.Values{
		"county_source": {"geographic"},
		"match":         {"enhanced"},
		"candidates":    {"5"},
	})
}

func (this *LookupFixture) TestQueryStringEncoding_CustomParameters() {
	lookup := &Lookup{}
	lookup.AddCustomParameter("test_parameter", "hello")
	this.So(this.encode(lookup), should.Resemble, url.Values{
		"test_parameter": {"hello"},
		"match":          {"enhanced"},
		"candidates":     {"5"},
	})
}

func (this *LookupFixture) TestQueryStringEncoding_MultipleCustomParameters() {
	lookup := &Lookup{}
	lookup.AddCustomParameter("test_parameter_1", "hello_1")
	lookup.AddCustomParameter("test_parameter_2", "hello_2")
	this.So(this.encode(lookup), should.Resemble, url.Values{
		"test_parameter_1": {"hello_1"},
		"test_parameter_2": {"hello_2"},
		"match":            {"enhanced"},
		"candidates":       {"5"},
	})
}

func (this *LookupFixture) TestQueryStringEncoding_DefaultMatchStrategyIsEnhanced() {
	this.So(this.encode(&Lookup{}), should.Resemble, url.Values{
		"match":      {"enhanced"},
		"candidates": {"5"},
	})
}

func (this *LookupFixture) TestQueryStringEncoding_ExplicitMatchStrict() {
	this.So(this.encode(&Lookup{MatchStrategy: MatchStrict}), should.Resemble, url.Values{"match": {"strict"}})
}

func (this *LookupFixture) TestQueryStringEncoding_ExplicitMatchStrictWithCandidates() {
	this.So(this.encode(&Lookup{MatchStrategy: MatchStrict, MaxCandidates: 3}), should.Resemble, url.Values{
		"candidates": {"3"}, "match": {"strict"},
	})
}

func (this *LookupFixture) marshalJSON(lookup *Lookup) map[string]any {
	raw, err := json.Marshal(lookup)
	if err != nil {
		this.Errorf("failed to marshal lookup: %s", err)
	}
	var result map[string]any
	err = json.Unmarshal(raw, &result)
	if err != nil {
		this.Errorf("failed to unmarshal lookup: %s", err)
	}
	return result
}

func (this *LookupFixture) TestJSONEncoding_DefaultMatchStrategyIsEnhancedWithCandidates() {
	result := this.marshalJSON(&Lookup{})
	this.So(result["match"], should.Equal, "enhanced")
	this.So(result["candidates"], should.Equal, float64(5))
}

func (this *LookupFixture) TestJSONEncoding_ExplicitMatchEnhancedWithCandidates() {
	result := this.marshalJSON(&Lookup{MatchStrategy: MatchEnhanced, MaxCandidates: 1})
	this.So(result["match"], should.Equal, "enhanced")
	this.So(result["candidates"], should.Equal, float64(1))
}

func (this *LookupFixture) TestJSONEncoding_ExplicitMatchStrict() {
	result := this.marshalJSON(&Lookup{MatchStrategy: MatchStrict})
	this.So(result["match"], should.Equal, "strict")
	this.So(result, should.NotContainKey, "candidates")
}

func (this *LookupFixture) TestJSONEncoding_ExplicitMatchStrictWithCandidates() {
	result := this.marshalJSON(&Lookup{MatchStrategy: MatchStrict, MaxCandidates: 3})
	this.So(result["match"], should.Equal, "strict")
	this.So(result["candidates"], should.Equal, float64(3))
}

func (this *LookupFixture) TestJSONEncoding_MatchInvalid() {
	result := this.marshalJSON(&Lookup{MatchStrategy: MatchInvalid})
	this.So(result["match"], should.Equal, "invalid")
	this.So(result, should.NotContainKey, "candidates")
}

func (this *LookupFixture) TestJSONEncoding_ExplicitMatchInvalidWithCandidates() {
	result := this.marshalJSON(&Lookup{MatchStrategy: MatchInvalid, MaxCandidates: 3})
	this.So(result["match"], should.Equal, "invalid")
	this.So(result["candidates"], should.Equal, float64(3))
}

func (this *LookupFixture) TestJSONEncoding_CustomParameters() {
	lookup := &Lookup{}
	lookup.AddCustomParameter("test_parameter", "hello")
	result := this.marshalJSON(lookup)
	this.So(result["test_parameter"], should.Equal, "hello")
	this.So(result, should.NotContainKey, "custom_parameters")
	this.So(result["match"], should.Equal, "enhanced")
	this.So(result["candidates"], should.Equal, float64(5))
}

func (this *LookupFixture) TestJSONEncoding_MultipleCustomParameters() {
	lookup := &Lookup{}
	lookup.AddCustomParameter("test_parameter_1", "hello_1")
	lookup.AddCustomParameter("test_parameter_2", "hello_2")

	result := this.marshalJSON(lookup)
	this.So(result["test_parameter_1"], should.Equal, "hello_1")
	this.So(result["test_parameter_2"], should.Equal, "hello_2")
	this.So(result, should.NotContainKey, "custom_parameters")
	this.So(result["match"], should.Equal, "enhanced")
	this.So(result["candidates"], should.Equal, float64(5))
}

func (this *LookupFixture) TestQueryStringEncodingToMatchJSONEncoding_CustomParameters() {
	lookup := &Lookup{}
	lookup.AddCustomParameter("test_parameter", "hello")
	jsonResult := this.marshalJSON(lookup)
	queryResult := this.encode(lookup)
	this.compareJSONandQuery(jsonResult, queryResult)
}

func (this *LookupFixture) TestQueryStringEncodingToMatchJSONEncoding_MultipleCustomParameters() {
	lookup := &Lookup{}
	lookup.AddCustomParameter("test_parameter_1", "hello_1")
	lookup.AddCustomParameter("test_parameter_2", "hello_2")
	jsonResult := this.marshalJSON(lookup)
	queryResult := this.encode(lookup)
	this.compareJSONandQuery(jsonResult, queryResult)
}

func (this *LookupFixture) TestJSONFieldNamesAndValuesMatchQueryStringKeyNamesAndValues() {
	lookupFieldsToSkip := []string{
		"CustomParameters", // tested in TestJSONEncoding_CustomParameters
		"Results",          // not in the request
	}
	l := &Lookup{}
	lv := reflect.ValueOf(l).Elem()
	lt := lv.Type()
	for i := range lt.NumField() {
		field := lt.Field(i)
		if slices.Contains(lookupFieldsToSkip, field.Name) {
			continue
		}
		fv := lv.Field(i)
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(fmt.Sprintf("test_value_%d", i))
		case reflect.Int:
			fv.SetInt(int64(i))
		default:
			this.Errorf("unknown field type: %s", field.Name)
		}
	}

	jsonResult := this.marshalJSON(l)
	query := make(url.Values)
	l.encodeQueryString(query)
	this.compareJSONandQuery(jsonResult, query)
}

func (this *LookupFixture) compareJSONandQuery(jsonResult map[string]any, queryResult url.Values) {
	for key, jsonValue := range jsonResult {
		this.So(queryResult, should.ContainKey, key)
		this.So(queryResult.Get(key), should.Equal, fmt.Sprintf("%v", jsonValue))
	}
	for key := range queryResult {
		this.So(jsonResult, should.ContainKey, key)
		this.So(fmt.Sprintf("%v", jsonResult[key]), should.Equal, queryResult.Get(key))
	}
}
