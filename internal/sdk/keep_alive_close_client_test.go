package sdk

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestKeepAliveCloseClientFixture(t *testing.T) {
	gunit.Run(new(KeepAliveCloseClientFixture), t)
}

type KeepAliveCloseClientFixture struct {
	*gunit.Fixture
	inner   *FakeHTTPClient
	request *http.Request
}

func (f *KeepAliveCloseClientFixture) Setup() {
	f.inner = &FakeHTTPClient{}
	f.inner.response = &http.Response{
		ProtoMajor: 1, ProtoMinor: 1,
		StatusCode: http.StatusTeapot,
		Body:       ioutil.NopCloser(strings.NewReader("Goodbye, World!")),
	}
	f.request = httptest.NewRequest("GET", "/", nil)
}

func (f *KeepAliveCloseClientFixture) TestCloseFalse_ReturnInnerClientInstead() {
	f.So(NewKeepAliveCloseClient(f.inner, false), should.Equal, f.inner)
}

func (f *KeepAliveCloseClientFixture) TestRequestClosesConnection() {
	client := NewKeepAliveCloseClient(f.inner, true)
	response, err := client.Do(f.request)
	f.So(f.inner.request.Close, should.BeTrue)
	f.So(response, should.Equal, f.inner.response)
	f.So(err, should.BeNil)
}
