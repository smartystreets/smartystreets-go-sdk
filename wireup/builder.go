package wireup

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/smartystreets/smartystreets-go-sdk"
	internal "github.com/smartystreets/smartystreets-go-sdk/internal/sdk"
	international_autocomplete "github.com/smartystreets/smartystreets-go-sdk/international-autocomplete-api"
	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	autocomplete_pro "github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-pro-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	us_reverse_geo "github.com/smartystreets/smartystreets-go-sdk/us-reverse-geo-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
)

// clientBuilder is responsible for accepting credentials and other configuration options to combine
// all components necessary to assemble a fully functional Client for use in an application.
type clientBuilder struct {
	credential    sdk.Credential
	baseURL       *url.URL
	proxy         *url.URL
	retries       int
	timeout       time.Duration
	debug         bool
	close         bool
	trace         bool
	headers       http.Header
	idleConns     int
	http2Disabled bool
	client        *http.Client
	licenses      []string
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

func (b *clientBuilder) withWebsiteKeyCredential(key, hostnameOrIP string) *clientBuilder {
	b.credential = sdk.NewWebsiteKeyCredential(key, hostnameOrIP)
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

func (b *clientBuilder) disableHTTP2() *clientBuilder {
	b.http2Disabled = true
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

func (b *clientBuilder) withMaxIdleConnections(max int) *clientBuilder {
	b.idleConns = max
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

func (b *clientBuilder) buildUSAutocompleteProAPIClient() *autocomplete_pro.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USAutocompleteProAPI)
	return autocomplete_pro.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) buildUSExtractAPIClient() *extract.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USExtractAPI)
	return extract.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) buildInternationalStreetAPIClient() *international_street.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_InternationalStreetAPI)
	return international_street.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) buildInternationalAutocompleteAPIClient() *international_autocomplete.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_InternationalAutocompleteAPI)
	return international_autocomplete.NewClient(b.buildHTTPSender())
}

func (b *clientBuilder) buildUSReverseGeocodingAPIClient() *us_reverse_geo.Client {
	b.ensureBaseURLNotNil(defaultBaseURL_USReverseGeocodingAPI)
	return us_reverse_geo.NewClient(b.buildHTTPSender())
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
	wrapped = b.buildClient()
	wrapped = internal.NewTracingClient(wrapped, b.trace)
	wrapped = internal.NewDebugOutputClient(wrapped, b.debug)
	wrapped = internal.NewRetryClient(wrapped, b.retries, time.Sleep)
	wrapped = internal.NewSigningClient(wrapped, b.credential)
	wrapped = internal.NewCustomHeadersClient(wrapped, b.headers)
	wrapped = internal.NewBaseURLClient(wrapped, b.baseURL)
	wrapped = internal.NewKeepAliveCloseClient(wrapped, b.close)
	wrapped = internal.NewLicenseClient(wrapped, b.licenses...)
	// outer-most
	return wrapped
}

func (b *clientBuilder) buildClient() *http.Client {
	if b.client != nil {
		return b.client
	}
	return &http.Client{Timeout: b.timeout, Transport: b.buildTransport()}
}

func (b *clientBuilder) buildTransport() *http.Transport {
	transport := &http.Transport{}
	if b.proxy != nil {
		transport.Proxy = http.ProxyURL(b.proxy)
	}
	if b.idleConns > 0 {
		transport.MaxIdleConnsPerHost = b.idleConns
	}
	if b.http2Disabled { // https://golang.org/pkg/net/http/ ("Programs that must disable HTTP/2 can do so by setting Transport.TLSNextProto to a non-nil, empty map.")
		transport.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper, 0)
	}

	transport.DialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext

	return transport
}

var (
	defaultBaseURL_InternationalStreetAPI       = &url.URL{Scheme: "https", Host: "international-street.api.smarty.com"}
	defaultBaseURL_InternationalAutocompleteAPI = &url.URL{Scheme: "https", Host: "international-autocomplete.api.smarty.com"}
	defaultBaseURL_USStreetAPI                  = &url.URL{Scheme: "https", Host: "us-street.api.smarty.com"}
	defaultBaseURL_USZIPCodeAPI                 = &url.URL{Scheme: "https", Host: "us-zipcode.api.smarty.com"}
	defaultBaseURL_USAutocompleteAPI            = &url.URL{Scheme: "https", Host: "us-autocomplete.api.smarty.com"}
	defaultBaseURL_USExtractAPI                 = &url.URL{Scheme: "https", Host: "us-extract.api.smarty.com"}
	defaultBaseURL_USReverseGeocodingAPI        = &url.URL{Scheme: "https", Host: "us-reverse-geo.api.smarty.com"}
	defaultBaseURL_USAutocompleteProAPI         = &url.URL{Scheme: "https", Host: "us-autocomplete-pro.api.smarty.com"}
)
