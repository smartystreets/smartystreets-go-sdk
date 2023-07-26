package sdk

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
)

func TestSecretKeySigning(t *testing.T) {
	t.Parallel()

	credential := NewSecretKeyCredential("my id", "my token")
	request, _ := http.NewRequest("GET", "http://google.com", nil)
	err := credential.Sign(request)
	assertions.New(t).So(err, should.BeNil)
	assertions.New(t).So(request.URL.String(), should.Equal, "http://google.com?auth-id=my+id&auth-token=my+token")
}

func TestURLEncodedSecretKeySigning_WeShouldNOTDoubleURLEncodeTheValue(t *testing.T) {
	t.Parallel()

	var (
		AuthToken                 = "Hello, World!"
		Base64EncodedAuthToken    = base64.StdEncoding.EncodeToString([]byte(AuthToken)) // SGVsbG8sIFdvcmxkIQ==
		URLEncodedBase64AuthToken = url.QueryEscape(Base64EncodedAuthToken)              // SGVsbG8sIFdvcmxkIQ%3D%3D
	)

	credential := NewSecretKeyCredential("auth-id", URLEncodedBase64AuthToken)
	request, _ := http.NewRequest("GET", "http://google.com", nil)
	err := credential.Sign(request)
	assertions.New(t).So(err, should.BeNil)
	assertions.New(t).So(request.URL.String(), should.Equal, "http://google.com?auth-id=auth-id&auth-token="+URLEncodedBase64AuthToken)
}
