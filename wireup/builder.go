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

type clientBuilder struct {
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

func newClientBuilder() *clientBuilder {
	return &clientBuilder{
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

func (b *clientBuilder) withSecretKeyCredential(authID, authToken string) *clientBuilder {
	b.credential = sdk.NewSecretKeyCredential(authID, authToken)
	return b
}

func (b *clientBuilder) withCustomBaseURL(address string) *clientBuilder {
	parsed, err := url.Parse(address)
	if err != nil {
		panic(fmt.Sprint("Could not parse provided address:", err.Error()))
	}
	b.baseURL = parsed
	return b
}

func (b *clientBuilder) withMaxRetry(retries int) *clientBuilder {
	if retries < 0 {
		panic(fmt.Sprintf("Please provide a non-negative number of retry attempts (you supplied %d).", retries))
	}
	b.retries = retries
	return b
}

func (b *clientBuilder) withTimeout(duration time.Duration) *clientBuilder {
	if duration < 0 {
		panic(fmt.Sprintf("Please provide a non-negative duration (you supplied %s).", duration.String()))
	}
	b.timeout = duration
	return b
}

func (b *clientBuilder) withDebugHTTPOutput() *clientBuilder {
	b.debug = true
	return b
}

func (b *clientBuilder) withHTTPRequestTracing() *clientBuilder {
	b.trace = true
	return b
}

func (b *clientBuilder) withCustomHeader(key, value string) *clientBuilder {
	b.headers.Add(key, value)
	return b
}

func (b *clientBuilder) withoutKeepAlive() *clientBuilder {
	b.close = true
	return b
}

func (b *clientBuilder) viaProxy(address string) *clientBuilder {
	proxy, err := url.Parse(address)
	if err != nil {
		panic(fmt.Sprint("Could not parse provided address:", err.Error()))
	}
	b.proxy = proxy
	return b
}

func (b *clientBuilder) buildUSStreetAPIClient() *street.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USStreetAPI)
	return street.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) buildUSZIPCodeAPIClient() *zipcode.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USZIPCodeAPI)
	return zipcode.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) buildUSAutocompleteAPIClient() *autocomplete.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USAutocompleteAPI)
	return autocomplete.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) buildUSExtractAPIClient() *extract.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USExtractAPI)
	return extract.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) buildInternationalStreetAPIClient() *international_street.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_InternationalStreetAPI)
	return international_street.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) ensureBaseURLNotNil(u *url.URL) {
	if b.baseURL == nil {
		b.baseURL = u
	}
}

func (b *clientBuilder) buildHTTPSender() *internal.HTTPSender {
	client := b.buildHTTPClient()
	return internal.NewHTTPSender(client)
}

func (b *clientBuilder) buildHTTPClient() (wrapped internal.HTTPClient) {
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

func (b *clientBuilder) buildTransport() *http.Transport {
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
