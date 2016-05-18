package us_street

import (
	"net/http"
	"net/url"

	"bitbucket.org/smartystreets/smartystreets-go-sdk"
)

// ClientBuilder is responsible for accepting credentials and other configuration options to combine
// all components necessary to assemble a fully-functional Client for use in an application.
type ClientBuilder struct {
	credential *sdk.SecretKeyCredential
	baseURL    string
	err        error
}

// NewClientBuilder creates a new client builder, ready to receive calls to its chain-able methods.
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{}
}

// WithSecretKeyCredential allows the caller to set the authID and authToken for use with the client.
// In all but very few cases calling this method with a valid authID and authToken is required.
func (b *ClientBuilder) WithSecretKeyCredential(authID, authToken string) *ClientBuilder {
	b.credential = &sdk.SecretKeyCredential{AuthID: authID, AuthToken: authToken}
	return b
}

// WithSecretKeyCredential allows the caller to specify the url that the client will use.
// In all but very few use cases the default value is sufficient and this method should not be called.
func (b *ClientBuilder) WithCustomBaseURL(uri string) *ClientBuilder {
	_, b.err = url.Parse(uri)
	b.baseURL = uri
	return b
}

// Builds the client using the provided configuration details provided by other methods on the ClientBuilder.
func (b *ClientBuilder) Build() (*Client, error) {
	if b.err != nil {
		return nil, b.err
	}

	var sender requestSender

	if b.credential != nil {
		signingClient := sdk.NewSigningClient(http.DefaultClient, b.credential)
		sender = sdk.NewHTTPSender(signingClient)
	} else {
		sender = sdk.NewHTTPSender(http.DefaultClient)
	}
	return NewClient(sender), nil
}
