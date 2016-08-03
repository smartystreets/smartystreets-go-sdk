package sdk

import (
	"net/http"
	"bitbucket.org/smartystreets/smartystreets-go-sdk"
)

// Signing client signs each request it receives using the credential
// provided and forwards the request to the inner client.
type SigningClient struct {
	inner      HTTPClient
	credential smarty_sdk.Credential
}

func NewSigningClient(inner HTTPClient, credential smarty_sdk.Credential) *SigningClient {
	return &SigningClient{
		inner:      inner,
		credential: credential,
	}
}

func (c *SigningClient) Do(request *http.Request) (*http.Response, error) {
	err := c.credential.Sign(request)
	if err != nil {
		return nil, err
	}
	return c.inner.Do(request)
}
