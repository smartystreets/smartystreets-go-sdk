package us_enrichment

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

func (f *ClientFixture) TestLookupSerializedAndSentWithContext__ResponseSuggestionsIncorporatedIntoLookup() {
	f.sender.response = validResponseJSON
	f.input.SmartyKey = "12345"
	f.input.DataSet = "property"

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/"+f.input.SmartyKey+"/"+f.input.DataSet)
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	f.So(f.input.Response, should.Resemble, []*Response{
		{
			SmartyKey:   "12345",
			DataSetName: "property",
			Attributes: []*Attribute{
				{Key: "PA1", Value: "67890"},
				{Key: "PA2", Value: "unknown"},
			}},
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
	f.input.SmartyKey = "12345"
	f.input.DataSet = "property"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Response, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.SmartyKey = "12345"
	f.input.DataSet = "property"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Response, should.BeEmpty)
}

var validResponseJSON = `[
{
	"smarty-key": "12345",
	"data-set-name": "property",
	"attributes": [{
			"key": "PA1",
			"value": "67890"
		},
		{
			"key": "PA2",
			"value": "unknown"
		}
	]
}]`

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
