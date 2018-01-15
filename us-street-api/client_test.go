package street

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
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
	batch  *Batch
}

func (f *ClientFixture) Setup() {
	f.sender = &FakeSender{}
	f.client = NewClient(f.sender)
	f.batch = NewBatch()
}

func (f *ClientFixture) TestSingleAddressBatch_SentInQueryStringAsGET() {
	f.sender.response = `[{"input_index": 0, "input_id": "42"}]`
	input := &Lookup{InputID: "42"}
	f.batch.Append(input)

	err := f.client.SendBatch(f.batch)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/street-address")
	f.So(f.sender.requestBody, should.BeNil)
	f.So(f.sender.request.ContentLength, should.Equal, 0)
	f.So(f.sender.request.URL.String(), should.StartWith, verifyURL)
	f.So(f.sender.request.URL.Query(), should.Resemble, url.Values{"input_id": {"42"}})
}

func (f *ClientFixture) TestAddressBatchSerializedAndSent__ResponseCandidatesIncorporatedIntoBatch() {
	f.sender.response = `[
		{"input_index": 0, "input_id": "42"},
		{"input_index": 2, "input_id": "44"},
		{"input_index": 2, "input_id": "44", "candidate_index": 1}
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
	f.So(f.sender.request.URL.Path, should.Equal, "/street-address")
	f.So(f.sender.request.ContentLength, should.Equal, len(f.sender.requestBody))
	f.So(string(f.sender.requestBody), should.Equal, `[{"input_id":"42"},{"input_id":"43"},{"input_id":"44"}]`)
	f.So(f.sender.request.URL.String(), should.Equal, verifyURL)

	f.So(input0.Results, should.Resemble, []*Candidate{{InputID: "42"}})
	f.So(input1.Results, should.BeEmpty)
	f.So(input2.Results, should.Resemble, []*Candidate{{InputID: "44", InputIndex: 2}, {InputID: "44", InputIndex: 2, CandidateIndex: 1}})
}

func (f *ClientFixture) TestNilBatchNOP() {
	err := f.client.SendBatch(nil)
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyBatch_NOP() {
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
	f.So(input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	input := new(Lookup)
	f.batch.Append(input)

	err := f.client.SendBatch(f.batch)

	f.So(err, should.NotBeNil)
	f.So(input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestNullCandidatesWithinResponseArrayAreIgnoredAfterDeserialization() {
	f.sender.response = `[null]`
	lookup := new(Lookup)
	f.batch.Append(lookup)
	f.So(func() { f.client.SendBatch(f.batch) }, should.NotPanic)
	f.So(lookup.Results, should.BeEmpty)
}

func (f *ClientFixture) TestOutOfRangeCandidatesWithinResponseArrayAreIgnoredAfterDeserialization() {
	f.sender.response = `[{"input_index": 9999999}]`
	lookup := new(Lookup)
	f.batch.Append(lookup)
	f.So(func() { f.client.SendBatch(f.batch) }, should.NotPanic)
	f.So(lookup.Results, should.BeEmpty)
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
	if request != nil && request.Body != nil {
		f.requestBody, _ = ioutil.ReadAll(request.Body)
	}
	return []byte(f.response), f.err
}
