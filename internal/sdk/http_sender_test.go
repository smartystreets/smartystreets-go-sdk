package sdk

import (
	"bytes"
	"errors"
	"net/http"

	"bitbucket.org/smartystreets/smartystreets-go-sdk"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type HTTPSenderFixture struct {
	*gunit.Fixture

	sender  *HTTPSender
	client  *FakeHTTPClient
	request *http.Request
}

func (f *HTTPSenderFixture) Setup() {
	f.client = &FakeHTTPClient{}
	f.sender = NewHTTPSender(f.client)
	f.request, _ = http.NewRequest("GET", "http://google.com", nil)
}

func (f *HTTPSenderFixture) TestRequestSentToClient_ResponseFromClientReadAndReturned() {
	closer := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!")}
	f.client.response = &http.Response{
		StatusCode: 200,
		Body:       closer,
	}
	result, err := f.sender.Send(f.request)
	f.So(err, should.BeNil)
	f.So(string(result), should.Equal, "Hello, World!")
	f.So(closer.closed, should.BeTrue)
}

func (f *HTTPSenderFixture) TestErrorWhenClosingResponseBody_ReturnsContentAndError() {
	closer := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!"), closeError: errors.New("GOPHERS!")}
	f.client.response = &http.Response{
		StatusCode: 200,
		Body:       closer,
	}

	result, err := f.sender.Send(f.request)
	f.So(err, should.NotBeNil)
	f.So(string(result), should.Equal, "Hello, World!")
	f.So(closer.closed, should.BeTrue)
}

func (f *HTTPSenderFixture) TestErrorWhenReadingResponseBody_ReturnsNoContentAndError() {
	body := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!"), readError: errors.New("GOPHERS!")}
	f.client.response = &http.Response{
		StatusCode: 200,
		Body:       body,
	}

	result, err := f.sender.Send(f.request)
	f.So(err, should.NotBeNil)
	f.So(result, should.BeEmpty)
	f.So(body.closed, should.BeTrue)
}

func (f *HTTPSenderFixture) TestHTTP400() {
	body := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!")}
	f.client.response = &http.Response{StatusCode: 400, Body: body}
	result, err := f.sender.Send(f.request)
	f.So(result, should.BeNil)
	f.So(err, should.Equal, smarty_sdk.StatusBadRequest)
}
func (f *HTTPSenderFixture) TestHTTP401() {
	body := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!")}
	f.client.response = &http.Response{StatusCode: 401, Body: body}
	result, err := f.sender.Send(f.request)
	f.So(result, should.BeNil)
	f.So(err, should.Equal, smarty_sdk.StatusUnauthorized)
}
func (f *HTTPSenderFixture) TestHTTP402() {
	body := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!")}
	f.client.response = &http.Response{StatusCode: 402, Body: body}
	result, err := f.sender.Send(f.request)
	f.So(result, should.BeNil)
	f.So(err, should.Equal, smarty_sdk.StatusPaymentRequired)
}
func (f *HTTPSenderFixture) TestHTTP413() {
	body := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!")}
	f.client.response = &http.Response{StatusCode: 413, Body: body}
	result, err := f.sender.Send(f.request)
	f.So(result, should.BeNil)
	f.So(err, should.Equal, smarty_sdk.StatusTooLarge)
}
func (f *HTTPSenderFixture) TestHTTP429() {
	body := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!")}
	f.client.response = &http.Response{StatusCode: 429, Body: body}
	result, err := f.sender.Send(f.request)
	f.So(result, should.BeNil)
	f.So(err, should.Equal, smarty_sdk.StatusTooManyRequests)
}

func (f *HTTPSenderFixture) TestNon200StatusCode_ReturnsNoContentAndCustomError() {
	body := &ErrorProneReadCloser{Buffer: bytes.NewBufferString("Hello, World!")}
	f.client.response = &http.Response{
		StatusCode: 500,
		Body:       body,
	}

	result, err := f.sender.Send(f.request)
	f.So(err, should.NotBeNil)
	f.So(result, should.BeEmpty)
	f.So(body.closed, should.BeTrue)
}

func (f *HTTPSenderFixture) TestErrorWhenSendingRequest_ReturnsNoContentAndError() {
	f.client.err = errors.New("GOPHERS!")
	result, err := f.sender.Send(f.request)
	f.So(err, should.NotBeNil)
	f.So(result, should.BeEmpty)
}

/*////////////////////////////////////////////////////////////////////////*/

type ErrorProneReadCloser struct {
	*bytes.Buffer
	closed     bool
	closeError error
	readError  error
}

func (e *ErrorProneReadCloser) Close() error {
	e.closed = true
	return e.closeError
}

func (e *ErrorProneReadCloser) Read(p []byte) (int, error) {
	if e.readError != nil {
		return 0, e.readError
	}
	return e.Buffer.Read(p)
}
