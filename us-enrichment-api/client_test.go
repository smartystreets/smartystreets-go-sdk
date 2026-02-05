package us_enrichment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
	"github.com/smartystreets/smartystreets-go-sdk"
)

func TestClientFixture(t *testing.T) {
	gunit.Run(new(ClientFixture), t)
}

type ClientFixture struct {
	*gunit.Fixture

	sender *FakeSender
	client *Client

	input enrichmentLookup
}

func (f *ClientFixture) Setup() {
	f.sender = &FakeSender{}
	f.client = NewClient(f.sender)
	f.input = new(principalLookup)
}

func (f *ClientFixture) TestLookupSerializedAndSentWithContext__ResponseSuggestionsIncorporatedIntoLookup() {
	smartyKey := "123"
	f.sender.response = validPrincipalResponse
	f.input = &principalLookup{Lookup: &Lookup{SmartyKey: smartyKey}}

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.sendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/"+smartyKey+"/"+f.input.getDataSet()+"/"+f.input.getDataSubset())
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	response := f.input.(*principalLookup).Response

	f.So(response, should.Resemble, []*PrincipalResponse{
		{
			SmartyKey:      "123",
			DataSetName:    "property",
			DataSubsetName: "principal",
			Attributes: PrincipalAttributes{
				FirstFloorSqft:        "1st_Floor_Sqft",
				LenderName2:           "Lender_Name_2",
				LenderSellerCarryBack: "Lender_Seller_Carry_Back",
				YearBuilt:             "Year_Built",
				Zoning:                "Zoning",
			},
			Etag: "ABCDEFG",
		},
	})
}

func (f *ClientFixture) TestNilLookupNOP() {
	err := f.client.sendLookup(nil)
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyLookup_NOP() {
	err := f.client.sendLookup(new(principalLookup))
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestSenderErrorPreventsDeserialization() {
	f.sender.err = errors.New("gophers")
	f.sender.response = validPrincipalResponse // would be deserialized if not for the err (above)
	f.input = &principalLookup{Lookup: &Lookup{SmartyKey: "12345"}}

	err := f.client.sendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.(*principalLookup).Response, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input = &principalLookup{Lookup: &Lookup{SmartyKey: "12345"}}

	err := f.client.sendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.(*principalLookup).Response, should.BeEmpty)
}

func (f *ClientFixture) TestUniversalLookupUnmarshallingWithEtag() {
	lookup := universalLookup{
		Response: []byte(validPrincipalResponse),
	}
	httpHeaders := http.Header{"Etag": []string{"ABCDEFG"}}

	_ = lookup.unmarshalResponse(lookup.Response, httpHeaders)

	f.So(lookup.Response, should.Equal, []byte(`[{"eTag": "ABCDEFG","smarty_key":"123","data_set_name":"property","data_subset_name":"principal","attributes":{"1st_floor_sqft":"1st_Floor_Sqft","lender_name_2":"Lender_Name_2","lender_seller_carry_back":"Lender_Seller_Carry_Back","year_built":"Year_Built","zoning":"Zoning"}}]`))
}

func (f *ClientFixture) TestUniversalLookupUnmarshallingWithNoEtag() {
	lookup := universalLookup{
		Response: []byte(validPrincipalResponse),
	}

	httpHeaders := http.Header{"NotAnEtag": []string{"ABC"}}

	_ = lookup.unmarshalResponse(lookup.Response, httpHeaders)

	f.So(lookup.Response, should.Equal, []byte(validPrincipalResponse))
}

func (f *ClientFixture) TestGeoReference() {
	smartyKey := "123"
	f.sender.response = validGeoReferenceResponse
	f.input = &geoReferenceLookup{Lookup: &Lookup{SmartyKey: smartyKey}}

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.sendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/"+smartyKey+"/"+f.input.getDataSet())
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	response := f.input.(*geoReferenceLookup).Response
	var geoRefResponse []*GeoReferenceResponse
	err = json.Unmarshal([]byte(validGeoReferenceResponse), &geoRefResponse)
	geoRefResponse[0].Etag = "ABCDEFG"
	f.So(err, should.BeNil)
	f.So(response, should.Resemble, geoRefResponse)
}

func (f *ClientFixture) TestRiskLookup() {
	smartyKey := "123"
	f.sender.response = validRiskResponse
	f.input = &riskLookup{Lookup: &Lookup{SmartyKey: smartyKey}}

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.sendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/"+smartyKey+"/"+f.input.getDataSet())
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	response := f.input.(*riskLookup).Response
	var riskResponse []*RiskResponse
	err = json.Unmarshal([]byte(validRiskResponse), &riskResponse)
	riskResponse[0].Etag = "ABCDEFG"
	f.So(err, should.BeNil)
	f.So(response, should.Resemble, riskResponse)
}

func (f *ClientFixture) TestSecondaryLookup() {
	smartyKey := "123"
	f.sender.response = validSecondaryResponse
	f.input = &secondaryLookup{Lookup: &Lookup{SmartyKey: smartyKey}}

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.sendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/"+smartyKey+"/"+f.input.getDataSet())
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	response := f.input.(*secondaryLookup).Response
	var secondaryResponse []*SecondaryResponse
	err = json.Unmarshal([]byte(validSecondaryResponse), &secondaryResponse)
	secondaryResponse[0].Etag = "ABCDEFG"
	f.So(err, should.BeNil)
	f.So(response, should.Resemble, secondaryResponse)
}

func (f *ClientFixture) TestSecondaryCount() {
	smartyKey := "123"
	f.sender.response = validSecondaryCountResponse
	f.input = &secondaryCountLookup{Lookup: &Lookup{SmartyKey: smartyKey}}

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.sendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/"+smartyKey+"/"+f.input.getDataSet()+"/"+f.input.getDataSubset())
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	response := f.input.(*secondaryCountLookup).Response
	secondaryCountResponse := []*SecondaryCountResponse{
		{
			SmartyKey: smartyKey,
			Count:     3,
			Etag:      "ABCDEFG",
		},
	}
	f.So(response, should.Resemble, secondaryCountResponse)
}

// Tests for public API methods

func (f *ClientFixture) TestSendPropertyPrincipal() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendPropertyPrincipal(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/property/principal")
	f.So(response, should.NotBeEmpty)
	f.So(response[0].SmartyKey, should.Equal, "123")
}

func (f *ClientFixture) TestSendPropertyPrincipalWithContextAndAuth() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.WithValue(context.Background(), "key", "value")

	err, response := f.client.SendPropertyPrincipalWithContextAndAuth(ctx, lookup, "myAuthID", "myAuthToken")

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Context(), should.Equal, ctx)
	f.So(f.sender.hasBasicAuth, should.BeTrue)
	f.So(f.sender.capturedAuthID, should.Equal, "myAuthID")
	f.So(f.sender.capturedAuthToken, should.Equal, "myAuthToken")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendGeoReferencePublicMethod() {
	f.sender.response = validGeoReferenceResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendGeoReference(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/geo-reference")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendGeoReferenceWithContextAndAuth() {
	f.sender.response = validGeoReferenceResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.WithValue(context.Background(), "key", "value")

	err, response := f.client.SendGeoReferenceWithContextAndAuth(ctx, lookup, "authID", "authToken")

	f.So(err, should.BeNil)
	f.So(f.sender.request.Context(), should.Equal, ctx)
	f.So(f.sender.hasBasicAuth, should.BeTrue)
	f.So(f.sender.capturedAuthID, should.Equal, "authID")
	f.So(f.sender.capturedAuthToken, should.Equal, "authToken")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendGeoReferenceWithVersion() {
	f.sender.response = validGeoReferenceResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendGeoReferenceWithVersion(lookup, "census-2020")

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/geo-reference/census-2020")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendGeoReferenceWithVersionContextAndAuth() {
	f.sender.response = validGeoReferenceResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.WithValue(context.Background(), "key", "value")

	err, response := f.client.SendGeoReferenceWithVersionContextAndAuth(ctx, lookup, "census-2010", "authID", "authToken")

	f.So(err, should.BeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/geo-reference/census-2010")
	f.So(f.sender.request.Context(), should.Equal, ctx)
	f.So(f.sender.hasBasicAuth, should.BeTrue)
	f.So(f.sender.capturedAuthID, should.Equal, "authID")
	f.So(f.sender.capturedAuthToken, should.Equal, "authToken")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendRiskPublicMethod() {
	f.sender.response = validRiskResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendRisk(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/risk")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendRiskWithContextAndAuth() {
	f.sender.response = validRiskResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.WithValue(context.Background(), "key", "value")

	err, response := f.client.SendRiskWithContextAndAuth(ctx, lookup, "authID", "authToken")

	f.So(err, should.BeNil)
	f.So(f.sender.request.Context(), should.Equal, ctx)
	f.So(f.sender.hasBasicAuth, should.BeTrue)
	f.So(f.sender.capturedAuthID, should.Equal, "authID")
	f.So(f.sender.capturedAuthToken, should.Equal, "authToken")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendSecondaryPublicMethod() {
	f.sender.response = validSecondaryResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendSecondary(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/secondary")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendSecondaryWithContextAndAuth() {
	f.sender.response = validSecondaryResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.WithValue(context.Background(), "key", "value")

	err, response := f.client.SendSecondaryWithContextAndAuth(ctx, lookup, "authID", "authToken")

	f.So(err, should.BeNil)
	f.So(f.sender.request.Context(), should.Equal, ctx)
	f.So(f.sender.hasBasicAuth, should.BeTrue)
	f.So(f.sender.capturedAuthID, should.Equal, "authID")
	f.So(f.sender.capturedAuthToken, should.Equal, "authToken")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendSecondaryCountPublicMethod() {
	f.sender.response = validSecondaryCountResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendSecondaryCount(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/secondary/count")
	f.So(response, should.NotBeEmpty)
	f.So(response[0].Count, should.Equal, 3)
}

func (f *ClientFixture) TestSendSecondaryCountWithContextAndAuth() {
	f.sender.response = validSecondaryCountResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.WithValue(context.Background(), "key", "value")

	err, response := f.client.SendSecondaryCountWithContextAndAuth(ctx, lookup, "authID", "authToken")

	f.So(err, should.BeNil)
	f.So(f.sender.request.Context(), should.Equal, ctx)
	f.So(f.sender.hasBasicAuth, should.BeTrue)
	f.So(f.sender.capturedAuthID, should.Equal, "authID")
	f.So(f.sender.capturedAuthToken, should.Equal, "authToken")
	f.So(response, should.NotBeEmpty)
}

// Tests for Universal Lookup methods

func (f *ClientFixture) TestSendUniversalLookup() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendUniversalLookup(lookup, "property", "principal")

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/property/principal")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendUniversalLookupWithContext() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.WithValue(context.Background(), "key", "value")

	err, response := f.client.SendUniversalLookupWithContext(ctx, lookup, "property", "principal")

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/property/principal")
	f.So(f.sender.request.Context(), should.Equal, ctx)
	f.So(f.sender.hasBasicAuth, should.BeFalse)
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendUniversalLookupWithContextAndAuth() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.WithValue(context.Background(), "key", "value")

	err, response := f.client.SendUniversalLookupWithContextAndAuth(ctx, lookup, "property", "principal", "authID", "authToken")

	f.So(err, should.BeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/property/principal")
	f.So(f.sender.request.Context(), should.Equal, ctx)
	f.So(f.sender.hasBasicAuth, should.BeTrue)
	f.So(f.sender.capturedAuthID, should.Equal, "authID")
	f.So(f.sender.capturedAuthToken, should.Equal, "authToken")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestSendUniversalLookupWithoutDataSubset() {
	f.sender.response = validGeoReferenceResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendUniversalLookup(lookup, "geo-reference", "")

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/geo-reference")
	f.So(response, should.NotBeEmpty)
}

// Tests for address search (freeform lookup without SmartyKey)

func (f *ClientFixture) TestAddressSearchWithFreeform() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{Freeform: "123 Main St, Denver CO"}

	err, _ := f.client.SendPropertyPrincipal(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/search/property/principal")
	f.So(f.sender.request.URL.Query().Get("freeform"), should.Equal, "123 Main St, Denver CO")
}

func (f *ClientFixture) TestAddressSearchWithComponents() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{
		Street:  "123 Main St",
		City:    "Denver",
		State:   "CO",
		ZIPCode: "80202",
	}

	err, _ := f.client.SendPropertyPrincipal(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/search/property/principal")
	f.So(f.sender.request.URL.Query().Get("street"), should.Equal, "123 Main St")
	f.So(f.sender.request.URL.Query().Get("city"), should.Equal, "Denver")
	f.So(f.sender.request.URL.Query().Get("state"), should.Equal, "CO")
	f.So(f.sender.request.URL.Query().Get("zipcode"), should.Equal, "80202")
}

// Tests for query parameters (include, exclude, features)

func (f *ClientFixture) TestLookupWithIncludeExcludeFeatures() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{
		SmartyKey: "123",
		Include:   "field1,field2",
		Exclude:   "field3",
		Features:  "feature1",
	}

	err, _ := f.client.SendPropertyPrincipal(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request.URL.Query().Get("include"), should.Equal, "field1,field2")
	f.So(f.sender.request.URL.Query().Get("exclude"), should.Equal, "field3")
	f.So(f.sender.request.URL.Query().Get("features"), should.Equal, "feature1")
}

func (f *ClientFixture) TestLookupWithETag() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{
		SmartyKey: "123",
		ETag:      "my-etag-value",
	}

	err, _ := f.client.SendPropertyPrincipal(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request.Header.Get("Etag"), should.Equal, "my-etag-value")
}

// Tests for IsHTTPErrorCode

func (f *ClientFixture) TestIsHTTPErrorCode_MatchingCode() {
	err := sdk.NewHTTPStatusError(404, []byte("not found"))

	result := f.client.IsHTTPErrorCode(err, 404)

	f.So(result, should.BeTrue)
}

func (f *ClientFixture) TestIsHTTPErrorCode_NonMatchingCode() {
	err := sdk.NewHTTPStatusError(500, []byte("internal error"))

	result := f.client.IsHTTPErrorCode(err, 404)

	f.So(result, should.BeFalse)
}

func (f *ClientFixture) TestIsHTTPErrorCode_NonHTTPError() {
	err := errors.New("some other error")

	result := f.client.IsHTTPErrorCode(err, 404)

	f.So(result, should.BeFalse)
}

func (f *ClientFixture) TestIsHTTPErrorCode_NilError() {
	result := f.client.IsHTTPErrorCode(nil, 404)

	f.So(result, should.BeFalse)
}

// Tests for per-request auth with empty credentials (should not set auth)

func (f *ClientFixture) TestPerRequestAuthEmptyCredentialsNotSet() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.Background()

	f.client.SendPropertyPrincipalWithContextAndAuth(ctx, lookup, "", "")

	f.So(f.sender.hasBasicAuth, should.BeFalse)
}

func (f *ClientFixture) TestPerRequestAuthPartialCredentialsNotSet() {
	f.sender.response = validPrincipalResponse
	lookup := &Lookup{SmartyKey: "123"}
	ctx := context.Background()

	// Only authID provided
	f.client.SendPropertyPrincipalWithContextAndAuth(ctx, lookup, "authID", "")
	f.So(f.sender.hasBasicAuth, should.BeFalse)

	f.sender.Reset()

	// Only authToken provided
	f.client.SendPropertyPrincipalWithContextAndAuth(ctx, lookup, "", "authToken")
	f.So(f.sender.hasBasicAuth, should.BeFalse)
}

// Tests for deprecated methods

func (f *ClientFixture) TestDeprecatedSendPropertyPrincipalLookup() {
	f.sender.response = validPrincipalResponse

	err, response := f.client.SendPropertyPrincipalLookup("123")

	f.So(err, should.BeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/property/principal")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestDeprecatedSendSecondaryLookup() {
	f.sender.response = validSecondaryResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendSecondaryLookup(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/secondary")
	f.So(response, should.NotBeEmpty)
}

func (f *ClientFixture) TestDeprecatedSendSecondaryCountLookup() {
	f.sender.response = validSecondaryCountResponse
	lookup := &Lookup{SmartyKey: "123"}

	err, response := f.client.SendSecondaryCountLookup(lookup)

	f.So(err, should.BeNil)
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/123/secondary/count")
	f.So(response, should.NotBeEmpty)
}

var validPrincipalResponse = `[{"smarty_key":"123","data_set_name":"property","data_subset_name":"principal","attributes":{"1st_floor_sqft":"1st_Floor_Sqft","lender_name_2":"Lender_Name_2","lender_seller_carry_back":"Lender_Seller_Carry_Back","year_built":"Year_Built","zoning":"Zoning"}}]`
var validGeoReferenceResponse = `[{"smarty_key":"123","data_set_name":"geo-reference","data_set_version":"census-2020","attributes":{"census_block":{"accuracy":"block","geoid":"180759630002012"},"census_county_division":{"accuracy":"exact","code":"1807581764","name":"Wayne"},"census_tract":{"code":"9630.00"},"place":{"accuracy":"exact","code":"1861236","name":"Portland","type":"incorporated"}}}]`
var validRiskResponse = `[{"smarty_key":"123","data_set_name":"risk","attributes":{"AGRIVALUE":"data","ALR_NPCTL":"data","ALR_VALA":"data","ALR_VALB":"data","ALR_VALP":"data","ALR_VRA_NPCTL":"data","AREA":"data","AVLN_AFREQ":"data","AVLN_ALRB":"data","AVLN_ALRP":"data","AVLN_ALR_NPCTL":"data","AVLN_EALB":"data","AVLN_EALP":"data","AVLN_EALPE":"data","AVLN_EALR":"data","AVLN_EALS":"data","AVLN_EALT":"data","AVLN_EVNTS":"data","AVLN_EXPB":"data","AVLN_EXPP":"data","AVLN_EXPPE":"data"}}]`
var validSecondaryResponse = `[{"smarty_key":"123","root_address":{"secondary_count":10,"smarty_key":"123","primary_number":"3105","street_name":"National Park Service","street_suffix":"Rd","city_name":"Juneau","state_abbreviation":"AK","zipcode":"99801","plus4_code":"8437"},"aliases":[{"smarty_key":"1882749021","primary_number":"3105","street_name":"National Park","street_suffix":"Rd","city_name":"Juneau","state_abbreviation":"AK","zipcode":"99801","plus4_code":"8437"}],"secondaries":[{"smarty_key":"1785903890","secondary_designator":"Apt","secondary_number":"A5","plus4_code":"8437"},{"smarty_key":"696702050","secondary_designator":"Apt","secondary_number":"B1","plus4_code":"8441"}]}]`
var validSecondaryCountResponse = `[{"smarty_key":"123","count":3}]`

/**************************************************************************/

type FakeSender struct {
	callCount int
	request   *http.Request

	response string
	err      error

	capturedAuthID    string
	capturedAuthToken string
	hasBasicAuth      bool
}

func (f *FakeSender) Send(request *http.Request) ([]byte, error) {
	f.callCount++
	f.request = request
	f.request.Response = &http.Response{Header: http.Header{"Etag": []string{"ABCDEFG"}}}
	f.capturedAuthID, f.capturedAuthToken, f.hasBasicAuth = request.BasicAuth()
	return []byte(f.response), f.err
}

func (f *FakeSender) Reset() {
	f.callCount = 0
	f.request = nil
	f.response = ""
	f.err = nil
	f.capturedAuthID = ""
	f.capturedAuthToken = ""
	f.hasBasicAuth = false
}
