package autocomplete

import (
	"errors"
	"net/http"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

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

func (f *ClientFixture) TestAddressLookupSerializedAndSent__ResponseSuggestionsIncorporatedIntoLookup() {
	f.sender.response = `{"suggestions":[
		{"text": "1"},
		{"text": "2"},
		{"text": "3"}
	]}`
	f.input.Prefix = "42"

	err := f.client.SendLookup(f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, suggestURL)
	f.So(f.sender.request.Header.Get("Content-Type"), should.Equal, "text/plain")
	f.So(string(f.sender.request.URL.Query().Get("prefix")), should.Equal, "42")
	f.So(f.sender.request.URL.String(), should.Equal, suggestURL+"?prefix=42")

	f.So(f.input.Results, should.Resemble, []*Suggestion{{Text: "1"}, {Text: "2"}, {Text: "3"}})
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
	f.sender.response = `{"suggestions":[
		{"text": "1"},
		{"text": "2"},
		{"text": "3"}
	]}` // would be deserialized if not for the err (above)
	f.input.Prefix = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.Prefix = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeEmpty)
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
