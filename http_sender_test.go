package sdk

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type HTTPSenderFixture struct {
	*gunit.Fixture

	sender  *HTTPSender
	client  *FakeHTTPClient
	request *http.Request
}

func (this *HTTPSenderFixture) Setup() {
	this.client = &FakeHTTPClient{}
	this.sender = NewHTTPSender(this.client)
	this.request, _ = http.NewRequest("GET", "http://google.com", nil)
}

func (this *HTTPSenderFixture) TestRequestSentToClient_ResponseFromClientReadAndReturned() {
	closer := &Closer{Buffer: bytes.NewBufferString("Hello, World!")}
	this.client.response = &http.Response{
		StatusCode: 200,
		Body:       closer,
	}
	result, err := this.sender.Send(this.request)
	this.So(err, should.BeNil)
	this.So(string(result), should.Equal, "Hello, World!")
	this.So(closer.closed, should.BeTrue)
}

func (this *HTTPSenderFixture) TestErrorWhenClosingResponseBody_ReturnsContentAndError() {
	closer := &Closer{Buffer: bytes.NewBufferString("Hello, World!"), closeError: errors.New("GOPHERS!")}
	this.client.response = &http.Response{
		StatusCode: 200,
		Body:       closer,
	}

	result, err := this.sender.Send(this.request)
	this.So(err, should.NotBeNil)
	this.So(string(result), should.Equal, "Hello, World!")
	this.So(closer.closed, should.BeTrue)
}

func (this *HTTPSenderFixture) TestErrorWhenReadingResponseBody_ReturnsNoContentAndError() {
	body := &Closer{Buffer: bytes.NewBufferString("Hello, World!"), readError: errors.New("GOPHERS!")}
	this.client.response = &http.Response{
		StatusCode: 200,
		Body:       body,
	}

	result, err := this.sender.Send(this.request)
	this.So(err, should.NotBeNil)
	this.So(result, should.BeEmpty)
	this.So(body.closed, should.BeFalse)
}

func (this *HTTPSenderFixture) TestErrorWhenSendingRequest_ReturnsNoContentAndError() {
	this.client.err = errors.New("GOPHERS!")
	result, err := this.sender.Send(this.request)
	this.So(err, should.NotBeNil)
	this.So(result, should.BeEmpty)
}

/*////////////////////////////////////////////////////////////////////////*/

type Closer struct {
	*bytes.Buffer
	closed     bool
	closeError error
	readError  error
}

func (this *Closer) Close() error {
	this.closed = true
	return this.closeError
}

func (this *Closer) Read(p []byte) (int, error) {
	if this.readError != nil {
		return 0, this.readError
	}
	return this.Buffer.Read(p)
}
