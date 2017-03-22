package international

import (
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

func (f *ClientFixture) TestLookupSerializedAndSent__ResponseSuggestionsIncorporatedIntoLookup() {
	f.sender.response = `[{"address1":"42"}]`
	f.input.Freeform = "42"
	f.input.Country = "43"

	err := f.client.SendLookup(f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, verifyURL)
	f.So(f.sender.request.URL.Query().Get("freeform"), should.Equal, "42")
	f.So(f.sender.request.URL.Query().Get("country"), should.Equal, "43")

	f.So(f.input.Results, should.Resemble, []*Result{{Address1: "42"}})
}

func (f *ClientFixture) TestNilLookupNOP() {
	err := f.client.SendLookup(nil)
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyLookup_NOP() {
	err := f.client.SendLookup(nil)
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestSenderErrorPreventsDeserialization() {
	f.sender.err = errors.New("GOPHERS!")
	f.sender.response = `[{"address1":"42"}]` // would be deserialized if not for the err (above)
	f.input.Freeform = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeNil)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.Freeform = "42"
	f.input.Country = "43"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeNil)
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
