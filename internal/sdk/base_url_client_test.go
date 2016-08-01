package sdk

import (
	"github.com/smartystreets/gunit"
	"net/http"
	"net/url"
	"github.com/smartystreets/assertions/should"
)

type BaseURLClientFixture struct {
	*gunit.Fixture
}

func (this *BaseURLClientFixture) TestProvidedURLOverridesRequestURL() {
	inner := &FakeHTTPClient{}
	inner.response = &http.Response{StatusCode: 123}
	original, _ := http.NewRequest("GET", "http://www.google.com/the/path/stays", nil)
	override, _ := url.Parse("https://smartystreets.com/the/path/is/ignored")
	client := NewBaseURLClient(inner, override)

	response, err := client.Do(original)

	this.So(err, should.BeNil)
	this.So(response, should.Equal, inner.response)
	this.So(original.URL.String(), should.Equal, override.Scheme + "://" + override.Host + original.URL.Path)
	this.So(original.URL.Scheme, should.Equal, override.Scheme)
	this.So(original.URL.Host, should.Equal, override.Host)
}
