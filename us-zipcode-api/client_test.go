package zipcode

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
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

	batch *Batch
}

func (f *ClientFixture) Setup() {
	f.sender = &FakeSender{}
	f.client = NewClient(f.sender)
	f.batch = NewBatch()
}

func (f *ClientFixture) TestSingleLookupSerializedInQueryStringGET() {
	f.sender.response = `[{"input_index": 0, "input_id": "42"}]`
	input := &Lookup{InputID: "42", ZIPCode: "10001", City: "NYC", State: "NY"}
	f.batch.Append(input)

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendBatchWithContext(ctx, f.batch)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup")
	f.So(f.sender.requestBody, should.BeNil)
	f.So(f.sender.request.URL.String(), should.StartWith, placeholderURL)
	f.So(f.sender.request.URL.Query(), should.Resemble, url.Values{"input_id": {"42"}, "zipcode": {"10001"}, "city": {"NYC"}, "state": {"NY"}})
	f.So(f.sender.request.Context(), should.Resemble, ctx)
}

func (f *ClientFixture) TestLookupBatchSerializedAndSent__ResultsIncorporatedBackIntoBatch() {
	f.sender.response = `[
		{"input_index": 0, "input_id": "42"},
		{"input_index": 1, "input_id": "43"},
		{"input_index": 2, "input_id": "44"}
	]`
	input0 := &Lookup{InputID: "42", ZIPCode: "10001"}
	input1 := &Lookup{InputID: "43", State: "NY", City: "NYC"}
	input2 := &Lookup{InputID: "44"}
	f.batch.Append(input0)
	f.batch.Append(input1)
	f.batch.Append(input2)

	err := f.client.SendBatch(f.batch)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "POST")
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup")
	f.So(string(f.sender.requestBody), should.Equal, `[{"zipcode":"10001","input_id":"42"},{"city":"NYC","state":"NY","input_id":"43"},{"input_id":"44"}]`)
	f.So(f.sender.request.URL.String(), should.Equal, placeholderURL)
	f.So(f.sender.request.Header.Get("Content-Type"), should.Equal, "application/json")

	f.So(input0.Result, should.Resemble, &Result{InputID: "42"})
	f.So(input1.Result, should.Resemble, &Result{InputID: "43", InputIndex: 1})
	f.So(input2.Result, should.Resemble, &Result{InputID: "44", InputIndex: 2})
}

func (f *ClientFixture) TestNilBatchNOP() {
	err := f.client.SendBatch(nil)
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyBatchCausesSerializationError__PreventsBatchBeingSent() {
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
	f.So(input.Result, should.BeNil)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	input := new(Lookup)
	f.batch.Append(input)

	err := f.client.SendBatch(f.batch)

	f.So(err, should.NotBeNil)
	f.So(input.Result, should.BeNil)
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
	if request.Body != nil {
		f.requestBody, _ = io.ReadAll(request.Body)
	}
	return []byte(f.response), f.err
}
