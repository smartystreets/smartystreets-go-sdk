package wireup

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/smartystreets/smartystreets-go-sdk"
	internal "github.com/smartystreets/smartystreets-go-sdk/internal/sdk"
	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
)

// ClientBuilder is responsible for accepting credentials and other configuration options to combine
// all components necessary to assemble a fully functional Client for use in an application.
//
// Deprecated: This type (and all associated functions and methods) will be unexported in the future.
//
// Instead of this kind of wireup:
//
// 	client := NewClientBuilder().
// 		WithSecretKeyCredential("auth-id", "auth-token").
// 		WithTimeout(time.Second*20).
// 		BuildUSStreetAPIClient()
//
// Please migrate to this approach instead:
//
// 	client := BuildUSStreetAPIClient(
// 		SecretKeyCredential("auth-id", "auth-token"),
// 		Timeout(time.Second*20),
// 	)
type ClientBuilder struct {
	credential sdk.Credential
	baseURL    *url.URL
	proxy      *url.URL
	retries    int
	timeout    time.Duration
	debug      bool
	close      bool
	trace      bool
	headers    http.Header
}

// Deprecated: (see ClientBuilder godoc for details)
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

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) WithSecretKeyCredential(authID, authToken string) *ClientBuilder {
	b.credential = sdk.NewSecretKeyCredential(authID, authToken)
	return b
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) WithCustomBaseURL(address string) *ClientBuilder {
	parsed, err := url.Parse(address)
	if err != nil {
		panic(fmt.Sprint("Could not parse provided address:", err.Error()))
	}
	b.baseURL = parsed
	return b
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) WithMaxRetry(retries int) *ClientBuilder {
	if retries < 0 {
		panic(fmt.Sprintf("Please provide a non-negative number of retry attempts (you supplied %d).", retries))
	}
	b.retries = retries
	return b
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) WithTimeout(duration time.Duration) *ClientBuilder {
	if duration < 0 {
		panic(fmt.Sprintf("Please provide a non-negative duration (you supplied %s).", duration.String()))
	}
	b.timeout = duration
	return b
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) WithDebugHTTPOutput() *ClientBuilder {
	b.debug = true
	return b
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) WithHTTPRequestTracing() *ClientBuilder {
	b.trace = true
	return b
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) WithCustomHeader(key, value string) *ClientBuilder {
	b.headers.Add(key, value)
	return b
}

func (b *ClientBuilder) WithoutKeepAlive() *ClientBuilder {
	b.close = true
	return b
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) ViaProxy(address string) *ClientBuilder {
	proxy, err := url.Parse(address)
	if err != nil {
		panic(fmt.Sprint("Could not parse provided address:", err.Error()))
	}
	b.proxy = proxy
	return b
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) BuildUSStreetAPIClient() *street.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USStreetAPI)
	return street.NewClient(b.buildHTTPSender())
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) BuildUSZIPCodeAPIClient() *zipcode.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USZIPCodeAPI)
	return zipcode.NewClient(b.buildHTTPSender())
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) BuildUSAutocompleteAPIClient() *autocomplete.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USAutocompleteAPI)
	return autocomplete.NewClient(b.buildHTTPSender())
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) BuildUSExtractAPIClient() *extract.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USExtractAPI)
	return extract.NewClient(b.buildHTTPSender())
}

// Deprecated: (see ClientBuilder godoc for details)
func (b *ClientBuilder) BuildInternationalStreetAPIClient() *international_street.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_InternationalStreetAPI)
	return international_street.NewClient(b.buildHTTPSender())
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
	// inner-most
	wrapped = &http.Client{Timeout: b.timeout, Transport: b.buildTransport()}
	wrapped = internal.NewTracingClient(wrapped, b.trace)
	wrapped = internal.NewDebugOutputClient(wrapped, b.debug)
	wrapped = internal.NewRetryClient(wrapped, b.retries)
	wrapped = internal.NewSigningClient(wrapped, b.credential)
	wrapped = internal.NewBaseURLClient(wrapped, b.baseURL)
	wrapped = internal.NewCustomHeadersClient(wrapped, b.headers)
	wrapped = internal.NewKeepAliveCloseClient(wrapped, b.close)
	// outer-most
	return wrapped
}

func (b *ClientBuilder) buildTransport() *http.Transport {
	transport := &http.Transport{}
	if b.proxy != nil {
		transport.Proxy = http.ProxyURL(b.proxy)
	}
	return transport
}

var (
	defaultBaseURL_InternationalStreetAPI, _ = url.Parse("https://international-street.api.smartystreets.com")
	defaultBaseURL_USStreetAPI, _            = url.Parse("https://us-street.api.smartystreets.com")
	defaultBaseURL_USZIPCodeAPI, _           = url.Parse("https://us-zipcode.api.smartystreets.com")
	defaultBaseURL_USAutocompleteAPI, _      = url.Parse("https://us-autocomplete.api.smartystreets.com")
	defaultBaseURL_USExtractAPI, _           = url.Parse("https://us-extract.api.smartystreets.com")
)
