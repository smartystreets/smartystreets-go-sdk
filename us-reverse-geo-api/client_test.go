package us_reverse_geo

import (
	"context"
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

	input *Lookup
}

func (f *ClientFixture) Setup() {
	f.sender = &FakeSender{}
	f.client = NewClient(f.sender)
	f.input = new(Lookup)
}

func (f *ClientFixture) TestAddressLookupSerializedAndSentWithContext__ResponseSuggestionsIncorporatedIntoLookup() {
	f.sender.response = validResponseJSON
	f.input.Latitude = 40.123456789
	f.input.Longitude = -111

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, lookupURL)
	f.So(f.sender.request.URL.Query().Get("latitude"), should.Equal, "40.12345679")
	f.So(f.sender.request.URL.Query().Get("longitude"), should.Equal, "-111.00000000")
	f.So(f.sender.request.URL.String(), should.Equal, lookupURL+"?latitude=40.12345679&longitude=-111.00000000&source=")
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	f.So(f.input.Response, should.Resemble, Response{
		Results: []Result{
			{
				Coordinate: Coordinate{
					Latitude:  40.209549,
					Longitude: -110.840134,
					Accuracy:  "Rooftop",
					License:   sdk.CoordinateLicenseSmartyStreetsProprietary,
				},
				Distance: 26977.779297,
				Address: Address{
					Street:            "6186 S 45000 W",
					City:              "Fruitland",
					StateAbbreviation: "UT",
					ZIPCode:           "84027",
				},
			},
			{
				Coordinate: Coordinate{
					Latitude:  39.123411,
					Longitude: -110.872123,
					Accuracy:  "Zip9",
					License:   sdk.CoordinateLicenseSmartyStreets,
				},
				Distance: 34721.824219,
				Address: Address{
					Street:            "340 Hardscrabble Rd",
					City:              "Helper",
					StateAbbreviation: "UT",
					ZIPCode:           "84526",
				},
			},
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
	f.sender.err = errors.New("gophers")
	f.sender.response = validResponseJSON // would be deserialized if not for the err (above)
	f.input.Latitude = 40
	f.input.Longitude = -111

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Response.Results, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.Latitude = 40
	f.input.Longitude = -111

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Response.Results, should.BeEmpty)
}

var validResponseJSON = `{
  "results": [
    {
      "coordinate": {
        "latitude": 40.209549,
        "longitude": -110.840134,
        "accuracy": "Rooftop",
        "license": 1
      },
      "distance": 26977.779297,
      "address": {
        "street": "6186 S 45000 W",
        "city": "Fruitland",
        "state_abbreviation": "UT",
        "zipcode": "84027"
      }
    },
    {
      "coordinate": {
        "latitude": 39.123411,
        "longitude": -110.872123,
        "accuracy": "Zip9",
        "license": 0
      },
      "distance": 34721.824219,
      "address": {
        "street": "340 Hardscrabble Rd",
        "city": "Helper",
        "state_abbreviation": "UT",
        "zipcode": "84526"
      }
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

func (f *FakeSender) SendAndReturnHeaders(request *http.Request) ([]byte, http.Header, error) {
	return []byte{}, nil, f.err
}
