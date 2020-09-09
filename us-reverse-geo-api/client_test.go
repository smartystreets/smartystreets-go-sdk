package us_reverse_geo

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestClientFixture(t *testing.T) {
	gunit.Run(new(ClientFixture), t)
}

type ClientFixture struct {
	*gunit.Fixture

	sender *FakeSender
	client *Client

	input *Lookup
}

func (f *ClientFixture) Setup() {
	f.sender = &FakeSender{}
	f.client = NewClient(f.sender)
	f.input = new(Lookup)
}

func (f *ClientFixture) TestAddressLookupSerializedAndSentWithContext__ResponseSuggestionsIncorporatedIntoLookup() {
	f.sender.response = validResponseJSON
	f.input.Latitude = 40
	f.input.Longitude = -111

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, lookupURL)
	f.So(f.sender.request.URL.Query().Get("latitude"), should.Equal, "40.00000000")
	f.So(f.sender.request.URL.Query().Get("longitude"), should.Equal, "-111.00000000")
	f.So(f.sender.request.URL.String(), should.Equal, lookupURL+"?latitude=40.00000000&longitude=-111.00000000")
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	f.So(f.input.Results, should.Resemble, []Address{
		{
			Latitude:          40.209549,
			Longitude:         -110.840134,
			CoordinateLicense: CoordinateLicenseGatewaySpatial,
			Distance:          26977.779297,
			Street:            "6186 S 45000 W",
			City:              "Fruitland",
			StateAbbreviation: "UT",
			ZIPCode:           "84027",
		},
		{
			Latitude:          39.123411,
			Longitude:         -110.872123,
			CoordinateLicense: CoordinateLicenseSmartyStreets,
			Distance:          34721.824219,
			Street:            "340 Hardscrabble Rd",
			City:              "Helper",
			StateAbbreviation: "UT",
			ZIPCode:           "84526",
		},
	})
}

func (f *ClientFixture) TestNilLookupNOP() {
	err := f.client.SendLookup(nil)
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyLookup_NOP() {
	err := f.client.SendLookup(new(Lookup))
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestSenderErrorPreventsDeserialization() {
	f.sender.err = errors.New("GOPHERS!")
	f.sender.response = validResponseJSON // would be deserialized if not for the err (above)
	f.input.Latitude = 40
	f.input.Longitude = -111

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.Latitude = 40
	f.input.Longitude = -111

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeEmpty)
}

var validResponseJSON = `{
  "results": [
    {
      "latitude": 40.209549,
      "longitude": -110.840134,
      "coordinate_license": 1,
      "distance": 26977.779297,
      "street": "6186 S 45000 W",
      "city": "Fruitland",
      "state_abbreviation": "UT",
      "zipcode": "84027"
    },
    {
      "latitude": 39.123411,
      "longitude": -110.872123,
      "coordinate_license": 0,
      "distance": 34721.824219,
      "street": "340 Hardscrabble Rd",
      "city": "Helper",
      "state_abbreviation": "UT",
      "zipcode": "84526"
    }
  ]
}`

/**************************************************************************/

type FakeSender struct {
	callCount int
	request   *http.Request

	response string
	err      error
}

func (f *FakeSender) Send(request *http.Request) ([]byte, error) {
	f.callCount++
	f.request = request
	return []byte(f.response), f.err
}
