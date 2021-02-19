package street

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	sdk "github.com/smartystreets/smartystreets-go-sdk"
)

func TestClientFixture(t *testing.T) {
	gunit.Run(new(ClientFixture), t)
}

type ClientFixture struct {
	*gunit.Fixture

	sender *FakeSender
	client *Client
	batch  *Batch
}

func (f *ClientFixture) Setup() {
	f.sender = &FakeSender{}
	f.client = NewClient(f.sender)
	f.batch = NewBatch()
}

func (f *ClientFixture) TestSingleAddressBatchWithContext_SentInQueryStringAsGET() {
	f.sender.response = `[{"input_index": 0, "input_id": "42"}]`
	input := &Lookup{InputID: "42"}
	f.batch.Append(input)

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendBatchWithContext(ctx, f.batch)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/street-address")
	f.So(f.sender.requestBody, should.BeNil)
	f.So(f.sender.request.ContentLength, should.Equal, 0)
	f.So(f.sender.request.URL.String(), should.StartWith, verifyURL)
	f.So(f.sender.request.URL.Query(), should.Resemble, url.Values{"input_id": {"42"}})
	f.So(f.sender.request.Context(), should.Resemble, ctx)
}

func (f *ClientFixture) TestAddressBatchSerializedAndSent__ResponseCandidatesIncorporatedIntoBatch() {
	f.sender.response = `[
		{"input_index": 0, "input_id": "42"},
		{"input_index": 2, "input_id": "44"},
		{"input_index": 2, "input_id": "44", "candidate_index": 1}
	]`
	input0 := &Lookup{InputID: "42"}
	input1 := &Lookup{InputID: "43"}
	input2 := &Lookup{InputID: "44"}
	f.batch.Append(input0)
	f.batch.Append(input1)
	f.batch.Append(input2)

	err := f.client.SendBatch(f.batch)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "POST")
	f.So(f.sender.request.URL.Path, should.Equal, "/street-address")
	f.So(f.sender.request.ContentLength, should.Equal, len(f.sender.requestBody))
	f.So(string(f.sender.requestBody), should.Equal, `[{"input_id":"42"},{"input_id":"43"},{"input_id":"44"}]`)
	f.So(f.sender.request.URL.String(), should.Equal, verifyURL)

	f.So(input0.Results, should.Resemble, []*Candidate{{InputID: "42"}})
	f.So(input1.Results, should.BeEmpty)
	f.So(input2.Results, should.Resemble, []*Candidate{{InputID: "44", InputIndex: 2}, {InputID: "44", InputIndex: 2, CandidateIndex: 1}})
}

func (f *ClientFixture) TestNilBatchNOP() {
	err := f.client.SendBatch(nil)
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyBatch_NOP() {
	err := f.client.SendBatch(new(Batch))
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestSenderErrorPreventsDeserialization() {
	f.sender.err = errors.New("GOPHERS!")
	f.sender.response = `[
		{"input_index": 0, "input_id": "42"},
		{"input_index": 2, "input_id": "44"},
		{"input_index": 2, "input_id": "44", "candidate_index": 1}
	]` // would be deserialized if not for the err (above)

	input := new(Lookup)
	f.batch.Append(input)

	err := f.client.SendBatch(f.batch)

	f.So(err, should.NotBeNil)
	f.So(input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	input := new(Lookup)
	f.batch.Append(input)

	err := f.client.SendBatch(f.batch)

	f.So(err, should.NotBeNil)
	f.So(input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestNullCandidatesWithinResponseArrayAreIgnoredAfterDeserialization() {
	f.sender.response = `[null]`
	lookup := new(Lookup)
	f.batch.Append(lookup)
	f.So(func() { f.client.SendBatch(f.batch) }, should.NotPanic)
	f.So(lookup.Results, should.BeEmpty)
}

func (f *ClientFixture) TestOutOfRangeCandidatesWithinResponseArrayAreIgnoredAfterDeserialization() {
	f.sender.response = `[{"input_index": 9999999}]`
	lookup := new(Lookup)
	f.batch.Append(lookup)
	f.So(func() { f.client.SendBatch(f.batch) }, should.NotPanic)
	f.So(lookup.Results, should.BeEmpty)
}

func (f *ClientFixture) TestFullJSONResponseDeserialization() {
	f.sender.response = `[
  {
	"input_id": "blah",
    "input_index": 0,
    "candidate_index": 4242,
	"addressee": "John Smith",
    "delivery_line_1": "3214 N University Ave # 409",
    "delivery_line_2": "blah blah",
    "last_line": "Provo UT 84604-4405",
    "delivery_point_barcode": "846044405140",
    "components": {
      "primary_number": "3214",
      "street_predirection": "N",
      "street_postdirection": "Q",
      "street_name": "University",
      "street_suffix": "Ave",
      "secondary_number": "409",
      "secondary_designator": "#",
      "extra_secondary_number": "410",
      "extra_secondary_designator": "Apt",
      "pmb_number": "411",
      "pmb_designator": "Box",
      "city_name": "Provo",
      "default_city_name": "Provo",
      "state_abbreviation": "UT",
      "zipcode": "84604",
      "plus4_code": "4405",
      "delivery_point": "14",
      "delivery_point_check_digit": "0",
      "urbanization": "urbanization"
    },
    "metadata": {
      "record_type": "S",
      "zip_type": "Standard",
      "county_fips": "49049",
      "county_name": "Utah",
      "carrier_route": "C016",
      "congressional_district": "03",
	  "building_default_indicator": "hi",
      "rdi": "Commercial",
      "elot_sequence": "0016",
      "elot_sort": "A",
      "latitude": 40.27658,
      "longitude": -111.65759,
      "coordinate_license": 1,
      "precision": "Rooftop",
      "time_zone": "Mountain",
      "utc_offset": -7,
      "dst": true,
	  "ews_match": true
    },
    "analysis": {
      "dpv_match_code": "S",
      "dpv_footnotes": "AACCRR",
      "dpv_cmra": "Y",
      "dpv_vacant": "N",
      "dpv_no_stat": "N",
      "active": "Y",
      "footnotes": "footnotes",
      "lacslink_code": "lacslink_code",
      "lacslink_indicator": "lacslink_indicator",
      "suitelink_match": true
    }
  }
]`
	lookup := new(Lookup)
	f.batch.Append(lookup)
	err := f.client.SendBatch(f.batch)
	f.So(err, should.BeNil)
	f.So(lookup.Results, should.Resemble, []*Candidate{
		{
			InputID:              "blah",
			InputIndex:           0,
			CandidateIndex:       4242,
			Addressee:            "John Smith",
			DeliveryLine1:        "3214 N University Ave # 409",
			DeliveryLine2:        "blah blah",
			LastLine:             "Provo UT 84604-4405",
			DeliveryPointBarcode: "846044405140",
			Components: Components{
				PrimaryNumber:            "3214",
				StreetPredirection:       "N",
				StreetName:               "University",
				StreetPostdirection:      "Q",
				StreetSuffix:             "Ave",
				SecondaryNumber:          "409",
				SecondaryDesignator:      "#",
				ExtraSecondaryNumber:     "410",
				ExtraSecondaryDesignator: "Apt",
				PMBNumber:                "411",
				PMBDesignator:            "Box",
				CityName:                 "Provo",
				DefaultCityName:          "Provo",
				StateAbbreviation:        "UT",
				ZIPCode:                  "84604",
				Plus4Code:                "4405",
				DeliveryPoint:            "14",
				DeliveryPointCheckDigit:  "0",
				Urbanization:             "urbanization",
			},
			Metadata: Metadata{
				RecordType:               "S",
				ZIPType:                  "Standard",
				CountyFIPS:               "49049",
				CountyName:               "Utah",
				CarrierRoute:             "C016",
				CongressionalDistrict:    "03",
				BuildingDefaultIndicator: "hi",
				RDI:                      "Commercial",
				ELOTSequence:             "0016",
				ELOTSort:                 "A",
				Latitude:                 40.27658,
				Longitude:                -111.65759,
				CoordinateLicense:        sdk.CoordinateLicenseSmartyStreetsProprietary,
				Precision:                "Rooftop",
				TimeZone:                 "Mountain",
				UTCOffset:                -7,
				DST:                      true,
				EWSMatch:                 true,
			},
			Analysis: Analysis{
				DPVMatchCode:      "S",
				DPVFootnotes:      "AACCRR",
				DPVCMRACode:       "Y",
				DPVVacantCode:     "N",
				DPVNoStat:         "N",
				Active:            "Y",
				Footnotes:         "footnotes",
				LACSLinkCode:      "lacslink_code",
				LACSLinkIndicator: "lacslink_indicator",
				SuiteLinkMatch:    true,
				EWSMatch:          false,
			},
		},
	})
}

/*////////////////////////////////////////////////////////////////////////*/

type FakeSender struct {
	callCount int

	request     *http.Request
	requestBody []byte

	response string
	err      error
}

func (f *FakeSender) Send(request *http.Request) ([]byte, error) {
	f.callCount++
	f.request = request
	if request != nil && request.Body != nil {
		f.requestBody, _ = ioutil.ReadAll(request.Body)
	}
	return []byte(f.response), f.err
}
