package sdk

import (
	"errors"
	"net/http"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestLicenseClientFixture(t *testing.T) {
	gunit.Run(new(LicenseClientFixture), t)
}

type LicenseClientFixture struct {
	*gunit.Fixture

	client   *LicenseClient
	inner    *FakeHTTPClient
	request  *http.Request
	response *http.Response
	err      error
}

func (this *LicenseClientFixture) Setup() {
	this.request, _ = http.NewRequest("GET", "/", nil)
	this.err = errors.New("GOPHERS!")
	this.response = &http.Response{StatusCode: http.StatusTeapot}
	this.inner = &FakeHTTPClient{
		err:      this.err,
		response: this.response,
	}
}

func (this *LicenseClientFixture) TestLicenseIsAddedToTheRequestBeforeItIsSentToTheInnerHandler() {
	this.client = NewLicenseClient(this.inner, "0", "1", "", "2", "3")

	response, err := this.client.Do(this.request)

	this.So(err, should.Equal, this.inner.err)
	this.So(response, should.Equal, this.inner.response)
	this.So(this.request.URL.Query().Get("license"), should.Equal, "0,1,2,3")
}
func (this *LicenseClientFixture) TestLicenseNotAddedToQueryIfNotSpecified() {
	this.client = NewLicenseClient(this.inner, "")

	response, err := this.client.Do(this.request)

	this.So(err, should.Equal, this.inner.err)
	this.So(response, should.Equal, this.inner.response)
	this.So(this.request.URL.Query(), should.HaveLength, 0)
}
