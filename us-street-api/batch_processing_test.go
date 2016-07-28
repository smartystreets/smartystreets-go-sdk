package us_street

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BatchProcessingFixture struct {
	*gunit.Fixture
	client *Client
	sender *FakeMultiSender
}

func (f *BatchProcessingFixture) Setup() {
	f.sender = &FakeMultiSender{}
	f.client = NewClient(f.sender)
}

func (f *BatchProcessingFixture) TestManyLookupsSentInBatches() {
	lookups := make([]*Lookup, 250)
	for x := 0; x < len(lookups); x++ {
		lookups[x] = &Lookup{InputID: strconv.Itoa(x)}
	}
	f.client.SendLookups(lookups...)

	if !f.So(f.sender.callCount, should.Equal, 3) {
		return
	}
	f.So(f.sender.requestBodies[0], should.StartWith, `[{"input_id":"0"},`)
	f.So(f.sender.requestBodies[1], should.StartWith, `[{"input_id":"100"},`)
	f.So(f.sender.requestBodies[2], should.StartWith, `[{"input_id":"200"},`)

	f.So(f.sender.requestBodies[0], should.EndWith, `,{"input_id":"99"}]`)
	f.So(f.sender.requestBodies[1], should.EndWith, `,{"input_id":"199"}]`)
	f.So(f.sender.requestBodies[2], should.EndWith, `,{"input_id":"249"}]`)
}

func (f *BatchProcessingFixture) TestErrorPreventsAllLookupsFromBeingBatched() {
	lookups := make([]*Lookup, 250)
	for x := 0; x < len(lookups); x++ {
		lookups[x] = &Lookup{InputID: strconv.Itoa(x)}
	}
	f.sender.err = errors.New("GOPHERS!")
	f.sender.errOnCall = 2
	f.client.SendLookups(lookups...)

	if !f.So(f.sender.callCount, should.Equal, 2) {
		return
	}
	f.So(f.sender.requestBodies[0], should.StartWith, `[{"input_id":"0"},`)
	f.So(f.sender.requestBodies[1], should.StartWith, `[{"input_id":"100"},`)

	f.So(f.sender.requestBodies[0], should.EndWith, `,{"input_id":"99"}]`)
	f.So(f.sender.requestBodies[1], should.EndWith, `,{"input_id":"199"}]`)
}

func (f *BatchProcessingFixture) TestChannelOfLookupsSentInBatches() {
	lookups := make(chan *Lookup, 250)
	for x := 0; x < cap(lookups); x++ {
		lookups <- &Lookup{InputID: strconv.Itoa(x)}
	}
	close(lookups)

	f.client.SendFromChannel(lookups)

	if !f.So(f.sender.callCount, should.Equal, 3) {
		return
	}
	f.So(f.sender.requestBodies[0], should.StartWith, `[{"input_id":"0"},`)
	f.So(f.sender.requestBodies[1], should.StartWith, `[{"input_id":"100"},`)
	f.So(f.sender.requestBodies[2], should.StartWith, `[{"input_id":"200"},`)

	f.So(f.sender.requestBodies[0], should.EndWith, `,{"input_id":"99"}]`)
	f.So(f.sender.requestBodies[1], should.EndWith, `,{"input_id":"199"}]`)
	f.So(f.sender.requestBodies[2], should.EndWith, `,{"input_id":"249"}]`)
}

func (f *BatchProcessingFixture) TestErrorPreventsAllLookupsOnChannelFromBeingBatched() {
	lookups := make([]*Lookup, 250)
	for x := 0; x < len(lookups); x++ {
		lookups[x] = &Lookup{InputID: strconv.Itoa(x)}
	}
	f.sender.err = errors.New("GOPHERS!")
	f.sender.errOnCall = 2
	f.client.SendLookups(lookups...)

	if !f.So(f.sender.callCount, should.Equal, 2) {
		return
	}
	f.So(f.sender.requestBodies[0], should.StartWith, `[{"input_id":"0"},`)
	f.So(f.sender.requestBodies[1], should.StartWith, `[{"input_id":"100"},`)

	f.So(f.sender.requestBodies[0], should.EndWith, `,{"input_id":"99"}]`)
	f.So(f.sender.requestBodies[1], should.EndWith, `,{"input_id":"199"}]`)
}

/*////////////////////////////////////////////////////////////////////////*/

type FakeMultiSender struct {
	callCount     int
	requests      []*http.Request
	requestBodies []string

	err       error
	errOnCall int
}

func (f *FakeMultiSender) Send(request *http.Request) ([]byte, error) {
	f.callCount++
	f.requests = append(f.requests, request)

	if request.Body != nil {
		body, _ := ioutil.ReadAll(request.Body)
		f.requestBodies = append(f.requestBodies, string(body))
	}

	var err error
	if f.errOnCall == f.callCount {
		err = f.err
	}
	return []byte(fmt.Sprintf(multiSenderResponseFormat, f.callCount)), err
}

const multiSenderResponseFormat = `[{"input_index": %d}]`
