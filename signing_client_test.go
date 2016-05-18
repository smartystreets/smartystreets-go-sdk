package sdk

import (
	"errors"
	"net/http"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type SigningClientFixture struct {
	*gunit.Fixture

	client     *SigningClient
	inner      *FakeHTTPClient
	credential *FakeCredential
	request    *http.Request
}

func (f *SigningClientFixture) Setup() {
	f.credential = &FakeCredential{}
	f.inner = &FakeHTTPClient{}
	f.client = NewSigningClient(f.inner, f.credential)
	f.request, _ = http.NewRequest("GET", "http://google.com", nil)
}

func (f *SigningClientFixture) TestSuccessfulSigning() {
	expected := &http.Response{}
	f.inner.response = expected
	response, err := f.client.Do(f.request)
	f.So(response, should.Equal, expected)
	f.So(err, should.BeNil)
	f.So(f.request.Header.Get("Auth"), should.Equal, "Success")
}

func (f *SigningClientFixture) TestSigningFailed() {
	f.credential.err = errors.New("GOPHERS!")
	f.inner.response = &http.Response{}
	response, err := f.client.Do(f.request)
	f.So(response, should.BeNil)
	f.So(err, should.NotBeNil)
	f.So(f.request.Header.Get("Auth"), should.BeBlank)
}

/*////////////////////////////////////////////////////////////////////////*/

type FakeCredential struct {
	err error
}

func (f *FakeCredential) Sign(request *http.Request) error {
	if f.err == nil {
		request.Header.Set("Auth", "Success")
	}
	return f.err
}
