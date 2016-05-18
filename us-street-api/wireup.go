package us_street

import (
	"net/http"
	"net/url"

	"bitbucket.org/smartystreets/smartystreets-go-sdk"
)

type ClientBuilder struct {
	credential *sdk.SecretKeyCredential
	baseURL    string
	err        error
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

func (b *ClientBuilder) WithSecretKeyCredential(authID, authToken string) *ClientBuilder {
	b.credential = &sdk.SecretKeyCredential{AuthID: authID, AuthToken: authToken}
	return b
}

func (b *ClientBuilder) WithCustomBaseURL(uri string) *ClientBuilder {
	_, b.err = url.Parse(uri)
	b.baseURL = uri
	return b
}

func (b *ClientBuilder) Build() (*Client, error) {
	if b.err != nil {
		return nil, b.err
	}
	client := http.DefaultClient
	signingClient := sdk.NewSigningClient(client, b.credential)
	sender := sdk.NewHTTPSender(signingClient)
	return NewClient(sender), nil
}
