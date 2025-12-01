package international_postal_code

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

func (f *ClientFixture) TestLookupSerializedAndSent__ResponseSuggestionsIncorporatedIntoLookup() {
	f.sender.response = `[
		{"input_id": "1"},
		{"administrative_area": "2"},
		{"locality": "3"},
		{"postal_code": "4"}
	]`
	f.input.AdministrativeArea = "42"

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, lookupUrl)
	f.So(f.sender.request.Context(), should.Resemble, ctx)
	f.So(f.sender.request.URL.Query().Get("administrative_area"), should.Equal, "42")
	f.So(f.sender.request.URL.String(), should.Equal, lookupUrl+"?administrative_area=42")
	f.So(f.input.Results, should.Resemble, []*Candidate{
		{InputID: "1"},
		{AdministrativeArea: "2"},
		{Locality: "3"},
		{PostalCode: "4"},
	})
}

func (f *ClientFixture) TestNilLookupNOP() {
	err := f.client.SendLookup(nil)
	f.So(err.Error(), should.Equal, "lookup cannot be nil")
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyLookup_NOP() {
	err := f.client.SendLookup(new(Lookup))
	f.So(err.Error(), should.Equal, "unexpected end of JSON input")
}

func (f *ClientFixture) TestSenderErrorPreventsDeserialization() {
	f.sender.err = errors.New("GOPHERS!")
	f.sender.response = `[
		{"text": "1"},
		{"text": "2"},
		{"text": "3"}
	]` // would be deserialized if not for the err (above)
	f.input.Locality = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I have no JSON`
	f.input.Locality = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestFullJSONResponseDeserialization() {
	f.sender.response = `[
{
	"input_id": "1",
  "country_iso_3": "2",
	"locality": "3",
	"administrative_area": "4",
	"sub_administrative_area": "5",
	"super_administrative_area": "6",
	"postal_code": "7"
}
]`
	lookup := new(Lookup)
	response := []byte(f.sender.response)
	err := deserializeResponse(response, lookup)
	f.So(err, should.BeNil)
	candidate := lookup.Results[0]
	f.So(candidate.InputID, should.Equal, "1")
	f.So(candidate.CountryIso3, should.Equal, "2")
	f.So(candidate.Locality, should.Equal, "3")
	f.So(candidate.AdministrativeArea, should.Equal, "4")
	f.So(candidate.SubAdministrativeArea, should.Equal, "5")
	f.So(candidate.SuperAdministrativeArea, should.Equal, "6")
	f.So(candidate.PostalCode, should.Equal, "7")
}

/*////////////////////////////////////////////////////////////////////////*/

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
