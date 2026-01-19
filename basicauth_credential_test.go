package sdk

import (
	"net/http"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestBasicAuthCredentialFixture(t *testing.T) {
	gunit.Run(new(BasicAuthCredentialFixture), t)
}

type BasicAuthCredentialFixture struct {
	*gunit.Fixture
}

func (this *BasicAuthCredentialFixture) TestNewBasicAuthCredentialWithValidCredentials() {
	cred := NewBasicAuthCredential("testID", "testToken")

	this.So(cred, should.NotBeNil)
	this.So(cred.authID, should.Equal, "testID")
	this.So(cred.authToken, should.Equal, "testToken")
}

func (this *BasicAuthCredentialFixture) TestNewBasicAuthCredentialWithEmptyAuthID() {
	cred := NewBasicAuthCredential("", "testToken")

	this.So(cred, should.NotBeNil)
	this.So(cred.authID, should.Equal, "")
	this.So(cred.authToken, should.Equal, "testToken")
}

func (this *BasicAuthCredentialFixture) TestNewBasicAuthCredentialWithEmptyAuthToken() {
	cred := NewBasicAuthCredential("testID", "")

	this.So(cred, should.NotBeNil)
	this.So(cred.authID, should.Equal, "testID")
	this.So(cred.authToken, should.Equal, "")
}

func (this *BasicAuthCredentialFixture) TestNewBasicAuthCredentialWithBothEmpty() {
	cred := NewBasicAuthCredential("", "")

	this.So(cred, should.NotBeNil)
	this.So(cred.authID, should.Equal, "")
	this.So(cred.authToken, should.Equal, "")
}

func (this *BasicAuthCredentialFixture) TestNewBasicAuthCredentialWithSpecialCharacters() {
	cred := NewBasicAuthCredential("test@id#123", "token!@#$%^&*()")

	this.So(cred, should.NotBeNil)
	this.So(cred.authID, should.Equal, "test@id#123")
	this.So(cred.authToken, should.Equal, "token!@#$%^&*()")
}

func TestSignMethodFixture(t *testing.T) {
	gunit.Run(new(SignMethodFixture), t)
}

type SignMethodFixture struct {
	*gunit.Fixture
}

func (this *SignMethodFixture) TestSignWithValidCredentials() {
	cred := NewBasicAuthCredential("myID", "myToken")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)

	// Verify the Authorization header is set
	username, password, ok := req.BasicAuth()
	this.So(ok, should.BeTrue)
	this.So(username, should.Equal, "myID")
	this.So(password, should.Equal, "myToken")
}

func (this *SignMethodFixture) TestSignWithEmptyCredentials() {
	cred := NewBasicAuthCredential("", "")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)

	username, password, ok := req.BasicAuth()
	this.So(ok, should.BeTrue)
	this.So(username, should.Equal, "")
	this.So(password, should.Equal, "")
}

func (this *SignMethodFixture) TestSignWithPasswordContainingColon() {
	// Note: Per RFC 2617, userid must NOT contain colons, but password can
	cred := NewBasicAuthCredential("validUserID", "password:with:colons")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)

	username, password, ok := req.BasicAuth()
	this.So(ok, should.BeTrue)
	this.So(username, should.Equal, "validUserID")
	this.So(password, should.Equal, "password:with:colons")
}

func (this *SignMethodFixture) TestSignWithSpecialCharacters() {
	cred := NewBasicAuthCredential("user@domain.com", "p@ssw0rd!")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)

	username, password, ok := req.BasicAuth()
	this.So(ok, should.BeTrue)
	this.So(username, should.Equal, "user@domain.com")
	this.So(password, should.Equal, "p@ssw0rd!")
}

func (this *SignMethodFixture) TestSignWithUnicodeCharacters() {
	cred := NewBasicAuthCredential("用户", "密码")
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	err := cred.Sign(req)

	this.So(err, should.BeNil)

	username, password, ok := req.BasicAuth()
	this.So(ok, should.BeTrue)
	this.So(username, should.Equal, "用户")
	this.So(password, should.Equal, "密码")
}

func (this *SignMethodFixture) TestSignOverwritesExistingHeader() {
	cred := NewBasicAuthCredential("newID", "newToken")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer oldtoken")

	err := cred.Sign(req)

	this.So(err, should.BeNil)

	username, password, ok := req.BasicAuth()
	this.So(ok, should.BeTrue)
	this.So(username, should.Equal, "newID")
	this.So(password, should.Equal, "newToken")

	// Ensure old header was replaced, not appended
	authHeaders := req.Header.Values("Authorization")
	this.So(len(authHeaders), should.Equal, 1)
}
