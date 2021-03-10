package sdk

import (
	"errors"
	"net/http"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestCustomHeadersClientFixture(t *testing.T) {
	gunit.Run(new(CustomHeadersClientFixture), t)
}

type CustomHeadersClientFixture struct {
	*gunit.Fixture

	headers  http.Header
	client   *CustomHeadersClient
	inner    *FakeHTTPClient
	request  *http.Request
	response *http.Response
	err      error
}

func (this *CustomHeadersClientFixture) Setup() {
	this.request, _ = http.NewRequest("GET", "/", nil)
	this.err = errors.New("GOPHERS!")
	this.response = &http.Response{StatusCode: http.StatusTeapot}
	this.inner = &FakeHTTPClient{
		err:      this.err,
		response: this.response,
	}
	this.headers = make(http.Header)
	this.client = NewCustomHeadersClient(this.inner, this.headers)
}

func (this *CustomHeadersClientFixture) TestAllCustomHeadersAreAddedToTheRequestBeforeItIsSentToTheInnerHandler() {
	this.headers.Add("A", "1")
	this.headers.Add("A", "2")
	this.headers.Add("B", "1")
	this.headers.Add("Host", "some-domain.com")

	response, err := this.client.Do(this.request)

	this.So(err, should.Equal, this.inner.err)
	this.So(response, should.Equal, this.inner.response)
	this.So(this.request.Header, should.HaveLength, 2)
	this.So(this.request.Header["A"], should.Resemble, []string{"1", "2"})
	this.So(this.request.Header.Get("A"), should.Equal, "1")
	this.So(this.request.Header.Get("B"), should.Equal, "1")
	this.So(this.request.Host, should.Equal, "some-domain.com")
}
