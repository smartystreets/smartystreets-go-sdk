package sdk

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestBaseURLClientFixture(t *testing.T) {
	gunit.Run(new(BaseURLClientFixture), t)
}

type BaseURLClientFixture struct {
	*gunit.Fixture

	inner    *FakeHTTPClient
	original *http.Request
	override *url.URL
}

func (f *BaseURLClientFixture) Setup() {
	f.inner = &FakeHTTPClient{}
	f.inner.response = &http.Response{StatusCode: 123}
}
func (f *BaseURLClientFixture) do() {
	client := NewBaseURLClient(f.inner, f.override)
	response, err := client.Do(f.original)
	f.So(err, should.BeNil)
	f.So(response, should.Equal, f.inner.response)
}
func (f *BaseURLClientFixture) assertFinalURL(expected string) {
	f.So(f.inner.request.URL.String(), should.Equal, expected)
}

func (f *BaseURLClientFixture) TestProvidedURLOverridesRequestURL() {
	f.original = httptest.NewRequest("GET", "http://original.com/original", nil)
	f.override, _ = url.Parse( /*********/ "https://override.com/override")

	f.do()

	f.assertFinalURL("https://override.com/override/original")
}
func (f *BaseURLClientFixture) TestHostWithPortInOverridingAddress() {
	f.original = httptest.NewRequest("GET", "https://us-street.api.smartystreets.com/street-address", nil)
	f.override, _ = url.Parse( /**********/ "https://1.2.3.4:5/override")

	f.do()

	f.assertFinalURL( /*******************/ "https://1.2.3.4:5/override/street-address")
}
func (f *BaseURLClientFixture) TestNoPathSpecifiedInOverridingAddress() {
	f.original = httptest.NewRequest("GET", "https://us-street.api.smartystreets.com/street-address", nil)
	f.override, _ = url.Parse( /**********/ "https://1.2.3.4:5/")

	f.do()

	f.assertFinalURL( /*******************/ "https://1.2.3.4:5/street-address")
}
