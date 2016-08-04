package us_zipcode

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

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

func (f *ClientFixture) TestLookupBatchSerializedAndSent__ResultsIncorporatedBackIntoBatch() {
	f.sender.response = `[
		{"input_index": 0, "input_id": "42"},
		{"input_index": 1, "input_id": "43"},
		{"input_index": 2, "input_id": "44"}
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
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup")
	f.So(string(f.sender.requestBody), should.Equal, `[{"input_id":"42"},{"input_id":"43"},{"input_id":"44"}]`)
	f.So(f.sender.request.URL.String(), should.Equal, placeholderURL)

	f.So(input0.Result, should.Resemble, &Result{InputID: "42"})
	f.So(input1.Result, should.Resemble, &Result{InputID: "43", InputIndex: 1})
	f.So(input2.Result, should.Resemble, &Result{InputID: "44", InputIndex: 2})
}

func (f *ClientFixture) TestNilBatchCausesSerializationError__PreventsBatchBeingSent() {
	err := f.client.SendBatch(nil)
	f.So(err, should.NotBeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyBatchCausesSerializationError__PreventsBatchBeingSent() {
	err := f.client.SendBatch(new(Batch))
	f.So(err, should.NotBeNil)
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

func (f *ClientFixture) TestPingReturnsTrueWhenServiceIsUp() {
	f.sender.response = "OK"
	up, err := f.client.Ping()
	f.So(up, should.BeTrue)
	f.So(err, should.BeNil)
}

func (f *ClientFixture) TestPingReturnsFalseWhenServiceIsDown() {
	f.sender.response = "Not OK"
	up, err := f.client.Ping()
	f.So(up, should.BeFalse)
	f.So(err, should.BeNil)
}

func (f *ClientFixture) TestPingReturnsErrorAndFalseWhenServiceIsUnreachable() {
	f.sender.response = "OK" // This wouldn't happen, but having an error next should trump even an OK in the response.
	f.sender.err = errors.New("HOT POCKETS!")
	up, err := f.client.Ping()
	f.So(up, should.BeFalse)
	f.So(err, should.Equal, f.sender.err)
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
		f.requestBody, _ = ioutil.ReadAll(request.Body)
	}
	return []byte(f.response), f.err
}
