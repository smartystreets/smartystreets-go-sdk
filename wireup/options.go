package wireup

import (
	"net/http"
	"time"

	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
)

// BuildUSStreetAPIClient builds a client for the US Street API using the provided options.
func BuildUSStreetAPIClient(options ...Option) *street.Client {
	return configure(options...).buildUSStreetAPIClient()
}

// BuildUSZIPCodeAPIClient builds a client for the US ZIP Code API using the provided options.
func BuildUSZIPCodeAPIClient(options ...Option) *zipcode.Client {
	return configure(options...).buildUSZIPCodeAPIClient()
}

// BuildUSAutocompleteAPIClient builds a client for the US Autocomplete API using the provided options.
func BuildUSAutocompleteAPIClient(options ...Option) *autocomplete.Client {
	return configure(options...).buildUSAutocompleteAPIClient()
}

// BuildUSExtractAPIClient builds a client for the US Extract API using the provided options.
func BuildUSExtractAPIClient(options ...Option) *extract.Client {
	return configure(options...).buildUSExtractAPIClient()
}

// BuildInternationalStreetAPIClient builds a client for the International Street API using the provided options.
func BuildInternationalStreetAPIClient(options ...Option) *international_street.Client {
	return configure(options...).buildInternationalStreetAPIClient()
}

func configure(options ...Option) *clientBuilder {
	builder := newClientBuilder()
	for _, option := range options {
		if option != nil {
			option(builder)
		}
	}
	return builder
}

type Option func(builder *clientBuilder)

// SecretKeyCredential sets the authID and authToken for use with the client.
// In all but very few cases calling this method with a valid authID and authToken is required.
func SecretKeyCredential(authID, authToken string) Option {
	return func(builder *clientBuilder) {
		builder.withSecretKeyCredential(authID, authToken)
	}
}

// WebsiteKeyCredential sets the key and hostnameOrIP for use with the client.
// This kind of authentication is generally only used for client-side applications but it
// included here for completeness.
func WebsiteKeyCredential(key, hostnameOrIP string) Option {
	return func(builder *clientBuilder) {
		builder.withWebsiteKeyCredential(key, hostnameOrIP)
	}
}

// CustomBaseURL specifies the url that the client will use.
// In all but very few use cases the default value is sufficient and this method should not be called.
// The address provided should be a url that consists of only the scheme and host. Any other elements
// (such as a path, query string, or fragment) will be ignored.
func CustomBaseURL(address string) Option {
	return func(builder *clientBuilder) {
		builder.withCustomBaseURL(address)
	}
}

// MaxRetry specifies the number of times an API request will be resent in the
// case of network errors or unexpected results.
func MaxRetry(retries int) Option {
	return func(builder *clientBuilder) {
		builder.withMaxRetry(retries)
	}
}

// Timeout specifies the timeout for all API requests.
func Timeout(duration time.Duration) Option {
	return func(builder *clientBuilder) {
		builder.withTimeout(duration)
	}
}

// DebugHTTPOutput engages detailed HTTP request/response logging using functions from net/http/httputil.
func DebugHTTPOutput() Option {
	return func(builder *clientBuilder) {
		builder.withDebugHTTPOutput()
	}
}

// DebugHTTPTracing engages additional HTTP-level tracing for each API request.
func DebugHTTPTracing() Option {
	return func(builder *clientBuilder) {
		builder.withHTTPRequestTracing()
	}
}

// CustomHeader ensures the provided header is added to every API request made with the resulting client.
func CustomHeader(key, value string) Option {
	return func(builder *clientBuilder) {
		builder.withCustomHeader(key, value)
	}
}

// DisableKeepAlive disables keep-alive for API requests.
// This is helpful if your environment limits the number of open files.
func DisableKeepAlive() Option {
	return func(builder *clientBuilder) {
		builder.withoutKeepAlive()
	}
}

// ViaProxy saves the address of your proxy server through which to send all requests.
func ViaProxy(address string) Option {
	return func(builder *clientBuilder) {
		builder.viaProxy(address)
	}
}

// WithMaxIdleConnections sets MaxIdleConnsPerHost on the http.Transport used to send requests.
// Docs for http.Transport.MaxIdleConnsPerHost: https://golang.org/pkg/net/http/#Transport
// Also see: https://stackoverflow.com/questions/22881090/golang-about-maxidleconnsperhost-in-the-http-clients-transport
func WithMaxIdleConnections(max int) Option {
	return func(builder *clientBuilder) {
		builder.withMaxIdleConnections(max)
	}
}

// DisableHTTP2 prevents clients from making use of the http2 protocol. This is achieved by following the instructions
// from the http package documentation (see: https://golang.org/pkg/net/http):
// > "Programs that must disable HTTP/2 can do so by setting Transport.TLSNextProto to a non-nil, empty map."
func DisableHTTP2() Option {
	return func(builder *clientBuilder) {
		builder.disableHTTP2()
	}
}

// WithHTTPClient allows the caller to supply their own *http.Client. This is useful if you want full
// control over the http client and its properties, but keep in mind that it reduces the following
// options to no-ops (you would need to specify any of those details on the *http.Client you provide):
//
// - DisableHTTP2
// - WithMaxIdleConnections
// - ViaProxy
// - Timeout
//
func WithHTTPClient(client *http.Client) Option {
	return func(builder *clientBuilder) {
		builder.client = client
	}
}