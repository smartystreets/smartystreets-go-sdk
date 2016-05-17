package sdk

import (
	"github.com/smartystreets/gunit"
	"net/http"
	"github.com/smartystreets/assertions/should"
	"errors"
)


type SigningClientFixture struct {
	*gunit.Fixture

	client *SigningClient
	inner *FakeHTTPClient
	credential *FakeCredential
	request *http.Request
}

func (this *SigningClientFixture) Setup() {
	this.credential = &FakeCredential{}
	this.inner = &FakeHTTPClient{}
	this.client = NewSigningClient(this.inner, this.credential)
	this.request, _ = http.NewRequest("GET", "http://google.com", nil)
}

func (this *SigningClientFixture) TestSuccessfulSigning() {
	expected := &http.Response{}
	this.inner.response = expected
	response, err := this.client.Do(this.request)
	this.So(response, should.Equal, expected)
	this.So(err, should.BeNil)
	this.So(this.request.Header.Get("Auth"), should.Equal, "Success")
}

func (this *SigningClientFixture) TestSigningFailed() {
	this.credential.err = errors.New("GOPHERS!")
	this.inner.response = &http.Response{}
	response, err := this.client.Do(this.request)
	this.So(response, should.BeNil)
	this.So(err, should.NotBeNil)
	this.So(this.request.Header.Get("Auth"), should.BeBlank)
}

/*////////////////////////////////////////////////////////////////////////*/

type FakeCredential struct {
	err error
}

func (this *FakeCredential) Sign(request *http.Request) error {
	if this.err == nil {
		request.Header.Set("Auth", "Success")
	}
	return this.err
}