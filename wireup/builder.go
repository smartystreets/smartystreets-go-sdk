package wireup

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"bitbucket.org/smartystreets/smartystreets-go-sdk"
	"bitbucket.org/smartystreets/smartystreets-go-sdk/internal/sdk"
	"bitbucket.org/smartystreets/smartystreets-go-sdk/us-street-api"
)

// ClientBuilder is responsible for accepting credentials and other configuration options to combine
// all components necessary to assemble a fully functional Client for use in an application.
type ClientBuilder struct {
	credential smarty_sdk.Credential
	baseURL    *url.URL
	retries    int
	timeout    time.Duration
	debug      bool
}

// NewClientBuilder creates a new client builder, ready to receive calls to its chain-able methods.
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		credential: &sdk.NopCredential{},
		retries:    4,
		timeout:    time.Second * 10,
	}
}

// WithSecretKeyCredential allows the caller to set the authID and authToken for use with the client.
// In all but very few cases calling this method with a valid authID and authToken is required.
func (b *ClientBuilder) WithSecretKeyCredential(authID, authToken string) *ClientBuilder {
	b.credential = &smarty_sdk.SecretKeyCredential{AuthID: authID, AuthToken: authToken}
	return b
}

// WithSecretKeyCredential allows the caller to specify the url that the client will use.
// In all but very few use cases the default value is sufficient and this method should not be called.
// The address provided should be a url that consists of only the scheme and host. Any other elements
// (such as a path, query string, or fragment) will be ignored.
func (b *ClientBuilder) WithCustomBaseURL(address string) *ClientBuilder {
	parsed, err := url.Parse(address)
	if err != nil {
		panic(fmt.Sprint("Could not parse provided address:", err.Error()))
	}
	b.baseURL = parsed
	return b
}

// WithMaxRetry allows the caller to specify the number of times an API request will be resent in the
// case of network errors or unexpected results.
func (b *ClientBuilder) WithMaxRetry(retries int) *ClientBuilder {
	if retries < 0 {
		panic(fmt.Sprintf("Please provide a non-negative number of retry attempts (you supplied %d).", retries))
	}
	b.retries = retries
	return b
}

// WithTimeout allows the caller to specify the timeout for all API requests.
func (b *ClientBuilder) WithTimeout(duration time.Duration) *ClientBuilder {
	if duration < 0 {
		panic(fmt.Sprintf("Please provide a non-negative duration (you supplied %s).", duration.String()))
	}
	b.timeout = duration
	return b
}

// WithDebugHTTPOutput enables detailed HTTP request/response logging using functions from net/http/httputil.
func (b *ClientBuilder) WithDebugHTTPOutput() *ClientBuilder {
	b.debug = true
	return b
}

// BuildUSStreetAPIClient builds the client using the provided configuration details provided by other methods on the ClientBuilder.
func (b *ClientBuilder) BuildUSStreetAPIClient() *us_street.Client {
	if b.baseURL == nil {
		b.baseURL = defaultBaseURL_USStreetAPI
	}
	return us_street.NewClient(b.buildHTTPSender())
}

func (b *ClientBuilder) buildHTTPSender() *sdk.HTTPSender {
	return sdk.NewHTTPSender(b.buildHTTPClient())
}

func (b *ClientBuilder) buildHTTPClient() (wrapped sdk.HTTPClient) {
	wrapped = &http.Client{Timeout: b.timeout}
	wrapped = sdk.NewDebugOutputClient(wrapped, b.debug)
	wrapped = sdk.NewRetryClient(wrapped, b.retries)
	wrapped = sdk.NewSigningClient(wrapped, b.credential)
	wrapped = sdk.NewBaseURLClient(wrapped, b.baseURL)
	return wrapped
}

var defaultBaseURL_USStreetAPI, _ = url.Parse("https://api.smartystreets.com")
