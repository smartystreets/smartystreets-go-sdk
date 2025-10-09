package wireup

import (
	"net/http"
	"time"

	international_autocomplete "github.com/smartystreets/smartystreets-go-sdk/international-autocomplete-api"
	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	autocomplete_pro "github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-pro-api"
	us_enrichment "github.com/smartystreets/smartystreets-go-sdk/us-enrichment-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	us_reverse_geo "github.com/smartystreets/smartystreets-go-sdk/us-reverse-geo-api"
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

// BuildUSAutocompleteProAPIClient builds a client for the US Autocomplete API using the provided options.
func BuildUSAutocompleteProAPIClient(options ...Option) *autocomplete_pro.Client {
	return configure(options...).buildUSAutocompleteProAPIClient()
}

// BuildUSEnrichmentAPIClient builds a client for the US Enrichment API using the provided options.
func BuildUSEnrichmentAPIClient(options ...Option) *us_enrichment.Client {
	return configure(options...).buildUSEnrichmentAPIClient()
}

// BuildUSExtractAPIClient builds a client for the US Extract API using the provided options.
func BuildUSExtractAPIClient(options ...Option) *extract.Client {
	return configure(options...).buildUSExtractAPIClient()
}

// BuildInternationalStreetAPIClient builds a client for the International Street API using the provided options.
func BuildInternationalStreetAPIClient(options ...Option) *international_street.Client {
	return configure(options...).buildInternationalStreetAPIClient()
}

// BuildInternationalAutocompleteAPIClient builds a client for the International Autocomplete API using the provided options.
func BuildInternationalAutocompleteAPIClient(options ...Option) *international_autocomplete.Client {
	return configure(options...).buildInternationalAutocompleteAPIClient()
}

// BuildUSReverseGeocodingAPIClient builds a client for the US Reverse Geocoding API using the provided options.
func BuildUSReverseGeocodingAPIClient(options ...Option) *us_reverse_geo.Client {
	return configure(options...).buildUSReverseGeocodingAPIClient()
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
// The address provided will be consulted for scheme, host, and path values. Any other URL components
// such as the query string or fragment will be ignored.
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
func WithHTTPClient(client *http.Client) Option {
	return func(builder *clientBuilder) {
		builder.client = client
	}
}

// WithLicenses allows the caller to specify the subscription license (aka "track") they wish to use.
func WithLicenses(licenses ...string) Option {
	return func(builder *clientBuilder) {
		builder.licenses = append(builder.licenses, licenses...)
	}
}

// WithCustomQuery allows the caller to specify key and value pair that is added to the request query.
func WithCustomQuery(key, value string) Option {
	return func(builder *clientBuilder) {
		builder.customQueries.Set(key, value)
	}
}

// WithCustomCommaSeparatedQuery allows the caller to specify key and value pair and appends the value to the current
// value associated with the key, separated by a comma.
func WithCustomCommaSeparatedQuery(key, value string) Option {
	return func(builder *clientBuilder) {
		v := builder.customQueries.Get(key)
		if v == "" {
			v = value
		} else {
			v += "," + value
		}
		builder.customQueries.Set(key, value)
	}
}

// WithFeatureComponentAnalysis adds to the request query to use the component analysis feature.
func WithFeatureComponentAnalysis() Option {
	return WithCustomCommaSeparatedQuery("features", "component-analysis")
}
