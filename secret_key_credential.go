package sdk

import (
	"net/http"
	"net/url"
	"strings"
)

type secretKeyCredential struct {
	authID    string
	authToken string
}

func NewSecretKeyCredential(authID, authToken string) *secretKeyCredential {
	if oldStyleBase64AuthTokenIsAlreadyURLEncoded(authToken) {
		authToken, _ = url.QueryUnescape(authToken)
	}
	return &secretKeyCredential{
		authID:    authID,
		authToken: authToken,
	}
}

func oldStyleBase64AuthTokenIsAlreadyURLEncoded(authToken string) bool {
	return strings.HasSuffix(authToken, "%3D")
}

func (c secretKeyCredential) Sign(request *http.Request) error {
	query := request.URL.Query()
	query.Set("auth-id", c.authID)
	query.Set("auth-token", c.authToken)
	request.URL.RawQuery = query.Encode()
	return nil
}

// FUTURE: implement HTML key credential
