package wireup

import (
	"time"

	international_street "github.com/smartystreets/smartystreets-go-sdk/international-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-autocomplete-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
)

// BuildUSStreetAPIClient builds a client for the US Street API using the provided options.
func BuildUSStreetAPIClient(options ...Option) *street.Client {
	return configure(options...).BuildUSStreetAPIClient()
}

// BuildUSZIPCodeAPIClient builds a client for the US ZIP Code API using the provided options.
func BuildUSZIPCodeAPIClient(options ...Option) *zipcode.Client {
	return configure(options...).BuildUSZIPCodeAPIClient()
}

// BuildUSAutocompleteAPIClient builds a client for the US Autocomplete API using the provided options.
func BuildUSAutocompleteAPIClient(options ...Option) *autocomplete.Client {
	return configure(options...).BuildUSAutocompleteAPIClient()
}

// BuildUSExtractAPIClient builds a client for the US Extract API using the provided options.
func BuildUSExtractAPIClient(options ...Option) *extract.Client {
	return configure(options...).BuildUSExtractAPIClient()
}

// BuildInternationalStreetAPIClient builds a client for the International Street API using the provided options.
func BuildInternationalStreetAPIClient(options ...Option) *international_street.Client {
	return configure(options...).BuildInternationalStreetAPIClient()
}
func configure(options ...Option) *ClientBuilder {
	builder := NewClientBuilder()
	for _, option := range options {
		option(builder)
	}
	return builder
}

type Option func(builder *ClientBuilder)

// SecretKeyCredential sets the authID and authToken for use with the client.
// In all but very few cases calling this method with a valid authID and authToken is required.
func SecretKeyCredential(authID, authToken string) Option {
	return func(builder *ClientBuilder) {
		builder.WithSecretKeyCredential(authID, authToken)
	}
}

// CustomBaseURL specifies the url that the client will use.
// In all but very few use cases the default value is sufficient and this method should not be called.
// The address provided should be a url that consists of only the scheme and host. Any other elements
// (such as a path, query string, or fragment) will be ignored.
func CustomBaseURL(address string) Option {
	return func(builder *ClientBuilder) {
		builder.WithCustomBaseURL(address)
	}
}

// MaxRetry specifies the number of times an API request will be resent in the
// case of network errors or unexpected results.
func MaxRetry(retries int) Option {
	return func(builder *ClientBuilder) {
		builder.WithMaxRetry(retries)
	}
}

// Timeout specifies the timeout for all API requests.
func Timeout(duration time.Duration) Option {
	return func(builder *ClientBuilder) {
		builder.WithTimeout(duration)
	}
}

// DebugHTTPOutput engages detailed HTTP request/response logging using functions from net/http/httputil.
func DebugHTTPOutput() Option {
	return func(builder *ClientBuilder) {
		builder.WithDebugHTTPOutput()
	}
}

// DebugHTTPTracing engages additional HTTP-level tracing for each API request.
func DebugHTTPTracing() Option {
	return func(builder *ClientBuilder) {
		builder.WithHTTPRequestTracing()
	}
}

// CustomHeader ensures the provided header is added to every API request made with the resulting client.
func CustomHeader(key, value string) Option {
	return func(builder *ClientBuilder) {
		builder.WithCustomHeader(key, value)
	}
}

// DisableKeepAlive disables keep-alive for API requests.
// This is helpful if your environment limits the number of open files.
func DisableKeepAlive() Option {
	return func(builder *ClientBuilder) {
		builder.WithoutKeepAlive()
	}
}

// ViaProxy saves the address of your proxy server through which to send all requests.
func ViaProxy(address string) Option {
	return func(builder *ClientBuilder) {
		builder.ViaProxy(address)
	}
}
