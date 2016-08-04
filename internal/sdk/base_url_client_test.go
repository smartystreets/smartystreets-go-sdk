package sdk

import (
	"net/http"
	"net/url"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BaseURLClientFixture struct {
	*gunit.Fixture
}

func (f *BaseURLClientFixture) TestProvidedURLOverridesRequestURL() {
	inner := &FakeHTTPClient{}
	inner.response = &http.Response{StatusCode: 123}
	original, _ := http.NewRequest("GET", "http://www.google.com/the/path/stays", nil)
	override, _ := url.Parse("https://smartystreets.com/the/path/is/ignored")
	client := NewBaseURLClient(inner, override)

	response, err := client.Do(original)

	f.So(err, should.BeNil)
	f.So(response, should.Equal, inner.response)
	f.So(original.URL.String(), should.Equal, override.Scheme+"://"+override.Host+original.URL.Path)
	f.So(original.URL.Scheme, should.Equal, override.Scheme)
	f.So(original.URL.Host, should.Equal, override.Host)
}
