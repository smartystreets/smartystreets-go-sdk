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

func (this *ClientBuilder) WithSecretKeyCredential(authID, authToken string) *ClientBuilder {
	this.credential = &sdk.SecretKeyCredential{AuthID: authID, AuthToken: authToken}
	return this
}

func (this *ClientBuilder) WithCustomBaseURL(uri string) *ClientBuilder {
	_, this.err = url.Parse(uri)
	this.baseURL = uri
	return this
}

func (this *ClientBuilder) Build() (*Client, error) {
	if this.err != nil {
		return nil, this.err
	}
	client := http.DefaultClient
	signingClient := sdk.NewSigningClient(client, this.credential)
	sender := sdk.NewHTTPSender(signingClient)
	return NewClient(sender), nil
}
