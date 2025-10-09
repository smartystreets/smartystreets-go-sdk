package sdk

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestCustomQueryClientFixture(t *testing.T) {
	gunit.Run(new(CustomQueryClientFixture), t)
}

type CustomQueryClientFixture struct {
	*gunit.Fixture

	inner    *FakeHTTPClient
	request  *http.Request
	response *http.Response
	err      error
}

func (this *CustomQueryClientFixture) Setup() {
	this.request, _ = http.NewRequest("GET", "/", nil)
	this.err = errors.New("GOPHERS!")
	this.response = &http.Response{StatusCode: http.StatusTeapot}
	this.inner = &FakeHTTPClient{
		err:      this.err,
		response: this.response,
	}
}

func (this *CustomQueryClientFixture) TestCustomQueryIsAddedToTheRequestBeforeItIsSentToTheInnerHandler() {
	query := url.Values{}
	query.Add("test", "test")
	client := NewCustomQueryClient(this.inner, query)
	response, err := client.Do(this.request)

	this.So(err, should.Equal, this.inner.err)
	this.So(response, should.Equal, this.inner.response)
	this.So(this.request.URL.Query(), should.HaveLength, 1)
	this.So(this.request.URL.Query().Get("test"), should.Equal, "test")
}

func (this *CustomQueryClientFixture) TestAllCustomQueriesAreAddedToTheRequestBeforeItIsSentToTheInnerHandler() {
	query := url.Values{}
	query.Add("test-key1", "test-value1")
	query.Add("test-key1", "test-value1.2")
	query.Add("test-key2", "test-value2")
	query.Add("test-key3", "test-value3")
	client := NewCustomQueryClient(this.inner, query)
	response, err := client.Do(this.request)

	this.So(err, should.Equal, this.inner.err)
	this.So(response, should.Equal, this.inner.response)
	this.So(this.request.URL.Query(), should.HaveLength, 3)
	this.So(this.request.URL.Query().Get("test-key1"), should.Equal, "test-value1")
}

func (this *CustomQueryClientFixture) TestEmptyCustomQueriesNothingAddedToTheRequestBeforeItIsSentToTheInnerHandler() {
	query := url.Values{}
	client := NewCustomQueryClient(this.inner, query)
	response, err := client.Do(this.request)

	this.So(err, should.Equal, this.inner.err)
	this.So(response, should.Equal, this.inner.response)
	this.So(this.request.URL.Query(), should.HaveLength, 0)
	this.So(this.request.URL.Query().Get("test"), should.Equal, "")
}
