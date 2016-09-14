package wireup

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	sdk "github.com/smartystreets/smartystreets-go-sdk"
	internal "github.com/smartystreets/smartystreets-go-sdk/internal/sdk"
	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
)

// ClientBuilder is responsible for accepting credentials and other configuration options to combine
// all components necessary to assemble a fully functional Client for use in an application.
type ClientBuilder struct {
	credential sdk.Credential
	baseURL    *url.URL
	retries    int
	timeout    time.Duration
	debug      bool
	headers    http.Header
}

// NewClientBuilder creates a new client builder, ready to receive calls to its chain-able methods.
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		credential: &internal.NopCredential{},
		retries:    4,
		timeout:    time.Second * 10,
		headers:    initializeHeadersWithUserAgent(),
	}
}

func initializeHeadersWithUserAgent() http.Header {
	headers := make(http.Header)
	headers.Add("User-Agent", fmt.Sprintf("smartystreets (sdk:go@%s)", sdk.VERSION))
	return headers
}

// WithSecretKeyCredential allows the caller to set the authID and authToken for use with the client.
// In all but very few cases calling this method with a valid authID and authToken is required.
func (b *ClientBuilder) WithSecretKeyCredential(authID, authToken string) *ClientBuilder {
	b.credential = &sdk.SecretKeyCredential{AuthID: authID, AuthToken: authToken}
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

// WithCustomHeader ensures the provided header is added to every API request made with the resulting client.
func (b *ClientBuilder) WithCustomHeader(key, value string) *ClientBuilder {
	b.headers.Add(key, value)
	return b
}

// BuildUSStreetAPIClient builds the us-street-api client using the provided
// configuration details provided by other methods on the ClientBuilder.
func (b *ClientBuilder) BuildUSStreetAPIClient() *street.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USStreetAPI)
	return street.NewClient(b.buildHTTPSender())
}

// BuildUSZIPCodeAPIClient builds the us-zipcode-api client using the provided
// configuration details provided by other methods on the ClientBuilder.
func (b *ClientBuilder) BuildUSZIPCodeAPIClient() *zipcode.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USZIPCodeAPI)
	return zipcode.NewClient(b.buildHTTPSender())
}

func (b *ClientBuilder) ensureBaseURLNotNil(u *url.URL) {
	if b.baseURL == nil {
		b.baseURL = u
	}
}

func (b *ClientBuilder) buildHTTPSender() *internal.HTTPSender {
	client := b.buildHTTPClient()
	return internal.NewHTTPSender(client)
}

func (b *ClientBuilder) buildHTTPClient() (wrapped internal.HTTPClient) {
	wrapped = &http.Client{Timeout: b.timeout}
	wrapped = internal.NewDebugOutputClient(wrapped, b.debug)
	wrapped = internal.NewRetryClient(wrapped, b.retries)
	wrapped = internal.NewSigningClient(wrapped, b.credential)
	wrapped = internal.NewBaseURLClient(wrapped, b.baseURL)
	wrapped = internal.NewCustomHeadersClient(wrapped, b.headers)
	return wrapped
}

var (
	defaultBaseURL_USStreetAPI, _  = url.Parse("https://us-street.api.smartystreets.com")
	defaultBaseURL_USZIPCodeAPI, _ = url.Parse("https://us-zipcode.api.smartystreets.com")
)
