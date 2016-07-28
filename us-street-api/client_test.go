package us_street

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

func (f *ClientFixture) TestSingleAddressBatchSerializedAndSent__ResponseCandidatesIncorporatedIntoBatch() {
	f.sender.response = `[{"input_index": 0, "input_id": "42"}]`
	input := &Lookup{
		InputID:       "42",
		Addressee:     "addressee",
		Street:        "street",
		Street2:       "street2",
		Secondary:     "secondary",
		LastLine:      "lastline",
		Urbanization:  "urbanization",
		ZIPCode:       "zipcode",
		MaxCandidates: 7,
	}
	f.batch.Append(input)

	err := f.client.SendBatch(f.batch)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET") // single address input
	f.So(f.sender.requestBody, should.BeNil)
	f.So(f.sender.request.URL.String(), should.Equal,
		defaultAPIURL+
			"?addressee=addressee"+
			"&candidates=7"+
			"&input_id=42"+
			"&lastline=lastline"+
			"&secondary=secondary"+
			"&street=street"+
			"&street2=street2"+
			"&urbanization=urbanization"+
			"&zipcode=zipcode",
	)
	f.So(input.Results, should.Resemble, []*Candidate{{InputID: "42"}})
}

func (f *ClientFixture) TestMultipleAddressBatchSerializedAndSent__ResponseCandidatesIncorporatedIntoBatch() {
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
	f.So(string(f.sender.requestBody), should.Equal, `[{"input_id":"42"},{"input_id":"43"},{"input_id":"44"}]`)
	f.So(f.sender.request.URL.String(), should.Equal, defaultAPIURL)

	f.So(input0.Results, should.Resemble, []*Candidate{{InputID: "42"}})
	f.So(input1.Results, should.BeEmpty)
	f.So(input2.Results, should.Resemble, []*Candidate{{InputID: "44", InputIndex: 2}, {InputID: "44", InputIndex: 2, CandidateIndex: 1}})
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

func (f *ClientFixture) TestXStandardizeOnlyHeaderAddedWhenSpecified() {
	f.batch.StandardizeOnly(true)
	input := new(Lookup)
	f.batch.Append(input)

	f.client.SendBatch(f.batch)

	f.So(f.sender.request.Header.Get(xStandardizeOnlyHeader), should.Equal, "true")
}

func (f *ClientFixture) TestXStandardizeOnlyHeaderNOTAddedWhenNOTSpecified() {
	f.batch.StandardizeOnly(false)
	input := new(Lookup)
	f.batch.Append(input)

	f.client.SendBatch(f.batch)

	f.So(f.sender.request.Header.Get(xStandardizeOnlyHeader), should.BeBlank)
}

func (f *ClientFixture) TestXIncludeInvalidHeaderAddedWhenSpecified() {
	f.batch.IncludeInvalid(true)
	input := new(Lookup)
	f.batch.Append(input)

	f.client.SendBatch(f.batch)

	f.So(f.sender.request.Header.Get(xIncludeInvalidHeader), should.Equal, "true")
}

func (f *ClientFixture) TestXIncludeInvalidHeaderNOTAddedWhenNOTSpecified() {
	f.batch.IncludeInvalid(false)
	input := new(Lookup)
	f.batch.Append(input)

	f.client.SendBatch(f.batch)

	f.So(f.sender.request.Header.Get(xIncludeInvalidHeader), should.BeBlank)
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
