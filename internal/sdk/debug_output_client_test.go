package sdk

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"github.com/smartystreets/logging"
)

type DebugOutputClientFixture struct {
	*gunit.Fixture
	inner   *FakeHTTPClient
	client  *DebugOutputClient
	request *http.Request
}

func (f *DebugOutputClientFixture) Setup() {
	f.inner = &FakeHTTPClient{}
	f.inner.response = &http.Response{
		ProtoMajor: 1, ProtoMinor: 1,
		StatusCode: http.StatusTeapot,
		Body:       ioutil.NopCloser(strings.NewReader("Goodbye, World!")),
	}
	f.client = NewDebugOutputClient(f.inner, true).(*DebugOutputClient)
	f.client.logger = logging.Capture()
	f.request, _ = http.NewRequest("PUT", "https://www.google.com/search", strings.NewReader("Hello, World!"))
}

func (f *DebugOutputClientFixture) TestDumpFalse_ReturnInnerClientInstead() {
	f.So(NewDebugOutputClient(f.inner, false), should.Equal, f.inner)
}

func (f *DebugOutputClientFixture) TestRequestDumped() {
	f.client.Do(f.request)

	log := f.client.logger.Log.String()
	f.assertRequestDumped(log)
	f.assertResponseDumped(log)
}
func (f *DebugOutputClientFixture) assertRequestDumped(log string) {
	f.So(log, should.ContainSubstring, "HTTP Request:\n")
	f.So(log, should.ContainSubstring, "PUT /search HTTP/1.1")
	f.So(log, should.ContainSubstring, "Host: www.google.com")
	f.So(log, should.ContainSubstring, "User-Agent: Go-http-client/1.1")
	f.So(log, should.ContainSubstring, "Content-Length: 13")
	f.So(log, should.ContainSubstring, "Accept-Encoding: gzip")
	f.So(log, should.ContainSubstring, "\n\n")
	f.So(log, should.ContainSubstring, "Hello, World!")
}
func (f *DebugOutputClientFixture) assertResponseDumped(log string) {
	f.So(log, should.ContainSubstring, "HTTP Response:")
	f.So(log, should.ContainSubstring, "HTTP/1.1 418 I'm a teapot")
	f.So(log, should.ContainSubstring, "Connection: close")
	f.So(log, should.ContainSubstring, "\n\n")
	f.So(log, should.ContainSubstring, "Goodbye, World!")
}

func (f *DebugOutputClientFixture) TestErrorDumpedIfResponseNil() {
	f.inner.err = errors.New("GOPHERS!")

	f.client.Do(f.request)

	log := f.client.logger.Log.String()
	f.assertResponseNOTDumped(log)
	f.assertErrorDumped(log)
}
func (f *DebugOutputClientFixture) assertResponseNOTDumped(log string) {
	f.So(log, should.NotContainSubstring, "HTTP Response:")
	f.So(log, should.NotContainSubstring, "HTTP/1.1 418 I'm a teapot")
	f.So(log, should.NotContainSubstring, "Connection: close")
	f.So(log, should.NotContainSubstring, "Goodbye, World!")
}
func (f *DebugOutputClientFixture) assertErrorDumped(log string) {
	f.So(log, should.ContainSubstring, "HTTP Err:")
	f.So(log, should.ContainSubstring, "GOPHERS!")
}

func (f *DebugOutputClientFixture) TestComposeDumpInErrorCase() {
	dump := composeDump("stuff", "I won't be include", errors.New("I will be included."))
	f.So(dump, should.Equal, "Could not dump HTTP stuff: I will be included.\n")
}
