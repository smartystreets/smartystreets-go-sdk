package sdk

import (
	"net/http/httptest"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestWebsiteKeyCredentialFixture(t *testing.T) {
	gunit.Run(new(WebsiteKeyCredentialFixture), t)
}

type WebsiteKeyCredentialFixture struct {
	*gunit.Fixture
}

func (this *WebsiteKeyCredentialFixture) TestWebsiteKeySigning() {
	credential := NewWebsiteKeyCredential("12345", "abc.com")
	request := httptest.NewRequest("GET", "/", nil)

	credential.Sign(request)

	this.So(request.URL.Query().Get("key"), should.Equal, "12345")
	this.So(request.Header.Get("Referer"), should.Equal, "http://abc.com")
}

func (this *WebsiteKeyCredentialFixture) TestHostAlreadyHasHTTPScheme() {
	credential := NewWebsiteKeyCredential("12345", "http://abc.com")
	request := httptest.NewRequest("GET", "/", nil)

	credential.Sign(request)

	this.So(request.URL.Query().Get("key"), should.Equal, "12345")
	this.So(request.Header.Get("Referer"), should.Equal, "http://abc.com")
}

func (this *WebsiteKeyCredentialFixture) TestHostAlreadyHasHTTPSScheme() {
	credential := NewWebsiteKeyCredential("12345", "https://abc.com")
	request := httptest.NewRequest("GET", "/", nil)

	credential.Sign(request)

	this.So(request.URL.Query().Get("key"), should.Equal, "12345")
	this.So(request.Header.Get("Referer"), should.Equal, "https://abc.com")
}
