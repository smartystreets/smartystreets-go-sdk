package sdk

import "net/http"

type Credential interface {
	Sign(*http.Request) error
}

type SigningClient struct {
	inner      httpClient
	credential Credential
}

func NewSigningClient(inner httpClient, credential Credential) *SigningClient {
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
