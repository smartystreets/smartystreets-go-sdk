package international_autocomplete_api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
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
				"administrative_area_short": "3.0",
				"administrative_area_long": "3.1",
				"postal_code": "4",
				"country_iso3": "5"
			},
			{
				"street": "6",
				"locality": "7",
				"administrative_area": "8",
				"administrative_area_short": "8.0",
				"administrative_area_long": "8.1",
				"postal_code": "9",
				"country_iso3": "10"
			}
		]
	}`
	f.input.Search = "42"
	f.input.Country = "FRA"

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, suggestURL)
	f.So(f.sender.request.URL.Query().Get("search"), should.Equal, "42")
	f.So(f.sender.request.URL.String(), should.Equal, suggestURL+"?country=FRA&max_group_results=100&max_results=5&search=42")
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	f.So(f.input.Result, should.Resemble, &Result{Candidates: []*Candidate{
		{
			Street:                  "1",
			Locality:                "2",
			AdministrativeArea:      "3",
			AdministrativeAreaShort: "3.0",
			AdministrativeAreaLong:  "3.1",
			PostalCode:              "4",
			CountryIso3:             "5",
		},
		{
			Street:                  "6",
			Locality:                "7",
			AdministrativeArea:      "8",
			AdministrativeAreaShort: "8.0",
			AdministrativeAreaLong:  "8.1",
			PostalCode:              "9",
			CountryIso3:             "10",
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
	f.input.Country = "FRA"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Result, should.BeNil)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.Search = "HI"
	f.input.Country = "FRA"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Result, should.BeNil)
}

func (f *ClientFixture) TestAddressIDAppendsToURL() {
	f.input.Country = "FRA"
	f.input.AddressID = "thisisid"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)

	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, suggestURL+"/thisisid")
	f.So(f.sender.request.URL.String(), should.Equal, suggestURL+"/thisisid?country=FRA&max_group_results=100&max_results=5")
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
