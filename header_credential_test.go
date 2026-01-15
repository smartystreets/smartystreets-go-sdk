package sdk

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestHeaderCredentialFixture(t *testing.T) {
	gunit.Run(new(HeaderCredentialFixture), t)
}

type HeaderCredentialFixture struct {
	*gunit.Fixture
}

func (this *HeaderCredentialFixture) TestSignWithValidCredentials() {
	cred := NewHeaderCredential("myID", "myToken")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)
	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte("myID:myToken")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func (this *HeaderCredentialFixture) TestSignWithEmptyCredentials() {
	cred := NewHeaderCredential("", "")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)
	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(":")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func (this *HeaderCredentialFixture) TestSignWithCredentialsContainingColon() {
	cred := NewHeaderCredential("id:with:colons", "token:also")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)
	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte("id:with:colons:token:also")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func (this *HeaderCredentialFixture) TestSignWithSpecialCharacters() {
	cred := NewHeaderCredential("user@domain.com", "p@ssw0rd!")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)
	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte("user@domain.com:p@ssw0rd!")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func TestSignRequestFixture(t *testing.T) {
	gunit.Run(new(SignRequestFixture), t)
}

type SignRequestFixture struct {
	*gunit.Fixture
}

func (this *SignRequestFixture) TestSignRequestWithBasicCredentials() {
	req, _ := http.NewRequest("POST", "http://api.example.com/endpoint", nil)

	SignRequest(req, "testUser", "testPass")

	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte("testUser:testPass")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func (this *SignRequestFixture) TestSignRequestWithEmptyAuthID() {
	req, _ := http.NewRequest("POST", "http://api.example.com/endpoint", nil)

	SignRequest(req, "", "password")

	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(":password")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func (this *SignRequestFixture) TestSignRequestWithEmptyAuthToken() {
	req, _ := http.NewRequest("POST", "http://api.example.com/endpoint", nil)

	SignRequest(req, "username", "")

	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte("username:")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func (this *SignRequestFixture) TestSignRequestWithUnicodeCharacters() {
	req, _ := http.NewRequest("POST", "http://api.example.com/endpoint", nil)

	SignRequest(req, "用户", "密码")

	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte("用户:密码")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func (this *SignRequestFixture) TestSignRequestWithLongCredentials() {
	req, _ := http.NewRequest("POST", "http://api.example.com/endpoint", nil)

	authID := "verylongusernamethatexceedsnormallength"
	authToken := "verylongpasswordthatexceedsnormallength"
	SignRequest(req, authID, authToken)

	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte(authID+":"+authToken)))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)
}

func (this *SignRequestFixture) TestSignRequestOverwritesExistingHeader() {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer oldtoken")

	SignRequest(req, "newID", "newToken")

	expectedAuth := fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte("newID:newToken")))
	this.So(req.Header.Get("Authorization"), should.Equal, expectedAuth)

	// Ensure old header was replaced, not appended
	authHeaders := req.Header.Values("Authorization")
	this.So(len(authHeaders), should.Equal, 1)
}

func (this *SignRequestFixture) TestSignRequestPreservesOtherHeaders() {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-Header", "custom-value")

	SignRequest(req, "id", "token")

	this.So(req.Header.Get("Content-Type"), should.Equal, "application/json")
	this.So(req.Header.Get("X-Custom-Header"), should.Equal, "custom-value")
}

func (this *SignRequestFixture) TestSignRequestWithNilRequestPanics() {
	defer func() {
		r := recover()
		this.So(r, should.NotBeNil)
	}()

	SignRequest(nil, "id", "token")
}
