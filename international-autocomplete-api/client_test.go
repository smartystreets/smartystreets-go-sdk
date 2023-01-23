package international_autocomplete_api

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
	f.sender.response = `{
		"candidates": [
			{
				"street": "1",
				"locality": "2",
				"administrative_area": "3",
                "super_administrative_area": "4",
                "sub_administrative_area": "5",
				"postal_code": "6",
				"country_iso3": "7"
			},
			{
				"street": "8",
				"locality": "9",
				"administrative_area": "10",
                "super_administrative_area": "11",
                "sub_administrative_area": "12",
				"postal_code": "13",
				"country_iso3": "14"
			}
		]
	}`
	f.input.Search = "42"

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, suggestURL)
	f.So(f.sender.request.URL.Query().Get("search"), should.Equal, "42")
	f.So(f.sender.request.URL.String(), should.Equal, suggestURL+"?distance=5&max_results=5&search=42")
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	f.So(f.input.Result, should.Resemble, &Result{Candidates: []*Candidate{
		{
			Street:                  "1",
			Locality:                "2",
			AdministrativeArea:      "3",
			SuperAdministrativeArea: "4",
			SubAdministrativeArea:   "5",
			PostalCode:              "6",
			CountryIso3:             "7",
		},
		{
			Street:                  "8",
			Locality:                "9",
			AdministrativeArea:      "10",
			SuperAdministrativeArea: "11",
			SubAdministrativeArea:   "12",
			PostalCode:              "13",
			CountryIso3:             "14",
		},
	}})
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
	f.sender.response = `{"candidates":[
		{"text": "1"},
		{"text": "2"},
		{"text": "3"}
	]}` // would be deserialized if not for the err (above)
	f.input.Search = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Result, should.BeNil)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.Search = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Result, should.BeNil)
}

//////////////////////////////////////////////////////////////////

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
