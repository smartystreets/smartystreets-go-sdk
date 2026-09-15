#### SMARTY DISCLAIMER: Subject to the terms of the associated license agreement, this software is freely available for your use. This software is FREE, AS IN PUPPIES, and is a gift. Enjoy your new responsibility. This means that while we may consider enhancement requests, we may or may not choose to entertain requests at our sole and absolute discretion.

[![Test](https://github.com/smartystreets/smartystreets-go-sdk/actions/workflows/test.yml/badge.svg)](https://github.com/smartystreets/smartystreets-go-sdk/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartystreets/smartystreets-go-sdk/v2.svg)](https://pkg.go.dev/github.com/smartystreets/smartystreets-go-sdk/v2)

# smartystreets-go-sdk

The official Go client library for the [Smarty](https://www.smarty.com) address APIs. The SDK covers US and international address validation, autocomplete, address extraction, reverse geocoding, and property enrichment.

- [Requirements](#requirements)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Supported APIs](#supported-apis)
- [Authentication](#authentication)
- [Client options](#client-options)
- [Batch requests](#batch-requests)
- [Errors and retries](#errors-and-retries)
- [Examples](#examples)
- [Development](#development)
- [Changelog](#changelog)
- [License](#license)

## Requirements

- Go 1.27 or later.
- A Smarty account with an active subscription. Create API keys on the [Smarty keys page](https://www.smarty.com/account/keys).

## Installation

```bash
go get github.com/smartystreets/smartystreets-go-sdk/v2
```

Import the `wireup` package and the package for each API that you use:

```go
import (
	"github.com/smartystreets/smartystreets-go-sdk/v2/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/v2/wireup"
)
```

## Quick start

This program validates one US address with the US Street API.

1. Export your secret key pair:

   ```bash
   export SMARTY_AUTH_ID="your-auth-id"
   export SMARTY_AUTH_TOKEN="your-auth-token"
   ```

2. Build a client, send a lookup, and read the candidates:

   ```go
   package main

   import (
   	"context"
   	"fmt"
   	"log"
   	"os"

   	"github.com/smartystreets/smartystreets-go-sdk/v2/us-street-api"
   	"github.com/smartystreets/smartystreets-go-sdk/v2/wireup"
   )

   func main() {
   	client := wireup.BuildUSStreetAPIClient(
   		wireup.BasicAuthCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
   	)

   	lookup := &street.Lookup{
   		Street:        "1600 Amphitheatre Parkway",
   		City:          "Mountain View",
   		State:         "CA",
   		ZIPCode:       "94043",
   		MaxCandidates: 3,
   		MatchStrategy: street.MatchEnhanced,
   	}

   	batch := street.NewBatch()
   	batch.Append(lookup)

   	if err := client.SendBatchWithContext(context.Background(), batch); err != nil {
   		log.Fatal(err)
   	}

   	for _, candidate := range lookup.Results {
   		fmt.Println(candidate.DeliveryLine1)
   		fmt.Println(candidate.LastLine)
   	}
   }
   ```

Each client exposes its send method in up to three forms:

- `SendBatch(batch)` or `SendLookup(lookup)` uses a background context.
- `SendBatchWithContext(ctx, batch)` or `SendLookupWithContext(ctx, lookup)` accepts a context for deadlines and cancellation.
- `SendBatchWithContextAndAuth(ctx, batch, credential)` or `SendLookupWithContextAndAuth(ctx, lookup, credential)` also overrides the client credential for that one request.

The autocomplete clients do not have the `WithContextAndAuth` form.

## Supported APIs

| API | Import path suffix | Package name | Builder | Example |
|---|---|---|---|---|
| US Street | `us-street-api` | `street` | `wireup.BuildUSStreetAPIClient` | [examples/us-street-api](examples/us-street-api) |
| US ZIP Code | `us-zipcode-api` | `zipcode` | `wireup.BuildUSZIPCodeAPIClient` | [examples/us-zipcode-api](examples/us-zipcode-api) |
| US Autocomplete | `us-autocomplete-api` | `autocomplete` | `wireup.BuildUSAutocompleteAPIClient` | [examples/us-autocomplete-api](examples/us-autocomplete-api) |
| US Autocomplete Pro | `us-autocomplete-pro-api` | `autocomplete_pro` | `wireup.BuildUSAutocompleteProAPIClient` | [examples/us-autocomplete-pro-api](examples/us-autocomplete-pro-api) |
| US Extract | `us-extract-api` | `extract` | `wireup.BuildUSExtractAPIClient` | [examples/us-extract-api](examples/us-extract-api) |
| US Reverse Geocoding | `us-reverse-geo-api` | `us_reverse_geo` | `wireup.BuildUSReverseGeocodingAPIClient` | [examples/us-reverse-geo-api](examples/us-reverse-geo-api) |
| US Address Enrichment | `us-enrichment-api` | `us_enrichment` | `wireup.BuildUSEnrichmentAPIClient` | [examples/us-enrichment-api](examples/us-enrichment-api) |
| International Street | `international-street-api` | `street` | `wireup.BuildInternationalStreetAPIClient` | [examples/international-street-api](examples/international-street-api) |
| International Postal Code | `international-postal-code-api` | `international_postal_code` | `wireup.BuildInternationalPostalCodeAPIClient` | [examples/international-postal-code-api](examples/international-postal-code-api) |
| International Autocomplete | `international-autocomplete-api` | `international_autocomplete_api` | `wireup.BuildInternationalAutocompleteAPIClient` | [examples/international-autocomplete-api](examples/international-autocomplete-api) |

The import path for each API is `github.com/smartystreets/smartystreets-go-sdk/v2/` plus the suffix in the table. The US Street and International Street packages share the package name `street`. Use an import alias if you import both.

The [Smarty cloud API documentation](https://www.smarty.com/docs/cloud) describes the input and output fields for each API. Each `Lookup` struct in this SDK maps one-to-one to the documented input fields.

## Authentication

Pass one credential option to the builder.

| Option | How the request is signed | Use it for |
|---|---|---|
| `wireup.SecretKeyCredential(authID, authToken)` | Adds `auth-id` and `auth-token` query parameters | Server-side code |
| `wireup.BasicAuthCredential(authID, authToken)` | Sets the `Authorization` header with HTTP Basic Authentication | Server-side code |
| `wireup.WebsiteKeyCredential(key, hostnameOrIP)` | Adds a `key` query parameter and a `Referer` header | Browser and mobile keys with referrer restrictions |

Embedded (website) keys work only with GET requests. A batch with two or more lookups is sent with POST. The US Extract API is POST-only. Use a secret key for those requests. See the [authentication documentation](https://www.smarty.com/docs/cloud/authentication).

`BasicAuthCredential` panics if the auth ID or the auth token is empty.

To sign a single request with a different credential, call the `WithContextAndAuth` form of the send method. Create the credential with `sdk.NewSecretKeyCredential`, `sdk.NewBasicAuthCredential`, or `sdk.NewWebsiteKeyCredential` from the root package `github.com/smartystreets/smartystreets-go-sdk/v2`.

## Client options

Every `wireup.Build*APIClient` function accepts any number of options. The defaults are 5 retries, a 10-second timeout, and HTTP/2 over TLS.

| Option | Effect |
|---|---|
| `MaxRetry(n)` | Retry a request up to `n` times after a network error or a retryable status code. Set `0` to disable retries. |
| `Timeout(d)` | Set the timeout for each request. |
| `CustomBaseURL(address)` | Send requests to a different scheme, host, and path. Use this for an on-premises deployment. |
| `WithLicenses(licenses...)` | Add the `license` query parameter to select a subscription track. |
| `CustomHeader(key, value)` | Add a header to every request. |
| `AppendedHeader(key, value, separator)` | Append a value to a single-value header such as `User-Agent`. |
| `WithCustomQuery(key, value)` | Add a query parameter to every request. |
| `WithCustomCommaSeparatedQuery(key, value)` | Append a value to a comma-separated query parameter. |
| `WithFeatureComponentAnalysis()` | Turn on the US Street `component-analysis` feature. |
| `WithFeatureIANATimeZone()` | Turn on the US Street `iana-timezone` feature. |
| `ViaProxy(address)` | Send all requests through an HTTP proxy. |
| `DisableKeepAlive()` | Close the connection after each request. Use this if the environment limits open files. |
| `WithMaxIdleConnections(n)` | Set `MaxIdleConnsPerHost` on the transport. |
| `DisableHTTP2()` | Force HTTP/1.1. |
| `EnableCleartextHTTP2()` | Use unencrypted HTTP/2 (h2c). This requires an `http://` base URL. Use it only for local development or a TLS-terminating sidecar. |
| `WithHTTPClient(client)` | Supply your own `*http.Client`. The client then ignores `Timeout`, `ViaProxy`, `DisableHTTP2`, `EnableCleartextHTTP2`, and `WithMaxIdleConnections`. |
| `DebugHTTPOutput()` | Log the full HTTP request and response for each call. |
| `DebugHTTPTracing()` | Log HTTP-level trace events for each call. |

## Batch requests

The US Street and US ZIP Code clients send lookups in batches.

- A `Batch` holds up to 100 lookups (`MaxBatchSize`). `Append` returns `false` when the batch is full.
- A batch with one lookup is sent with GET. A batch with two or more lookups is sent with POST.
- After the send, the `Results` field of each lookup holds its candidates.
- `SendLookups(lookups...)` splits any number of lookups into full batches and sends them one after another.
- `SendFromChannel(input, output)` reads lookups from `input`, sends them in batches, and writes each processed lookup to `output`. Close `input` when you have no more lookups. The SDK closes `output` for you.

The [list example](examples/us-street-api/list) streams a tab-delimited file through `SendFromChannel`.

## Errors and retries

If the API returns a status other than `200 OK`, the send method returns a `*sdk.HTTPStatusError`. The error message contains the status text and the error messages from the API response body.

```go
import (
	"errors"

	sdk "github.com/smartystreets/smartystreets-go-sdk/v2"
)

var statusErr *sdk.HTTPStatusError
if errors.As(err, &statusErr) {
	fmt.Println(statusErr.StatusCode(), statusErr.Content())
}
```

The retry client behaves as follows:

- It retries after a network error and after any status other than `200 OK`, with two exceptions.
- It does not retry a `4xx` response from `400` through `422`, or a `304 Not Modified` response.
- It sleeps a random whole number of seconds between attempts, up to 10 seconds.
- On `429 Too Many Requests`, it waits for the duration in the `Retry-After` header.
- If the context is canceled, it stops and returns the context error.

The US Address Enrichment client supports conditional requests. Set `Lookup.ETag` to the tag from an earlier response. If the record did not change, the API returns `304`, and the send method returns an empty slice with a nil error. After each call, `Lookup.ResponseETag` holds the tag from the response. See the [etag example](examples/us-enrichment-api/etag).

## Examples

Each directory under [examples/](examples) that contains a `main.go` file is a runnable program. The programs read `SMARTY_AUTH_ID` and `SMARTY_AUTH_TOKEN` from the environment. Some programs include a commented `WebsiteKeyCredential` alternative that reads `SMARTY_AUTH_WEB` and `SMARTY_AUTH_REFERER`.

| Example | Shows |
|---|---|
| `us-street-api` | A batch of three US Street lookups |
| `us-street-api/list` | A tab-delimited file streamed through `SendFromChannel` |
| `us-street-api/component-analysis` | The component analysis feature |
| `us-street-api/iana-timezone` | The IANA time zone feature |
| `us-street-api/match-strategy` | The `strict`, `invalid`, and `enhanced` match strategies |
| `us-zipcode-api` | A batch of US ZIP Code lookups |
| `us-autocomplete-api` | US Autocomplete suggestions |
| `us-autocomplete-pro-api` | US Autocomplete Pro suggestions with city and state filters |
| `us-extract-api` | Addresses extracted from free-form text |
| `us-reverse-geo-api` | Addresses near a latitude and longitude |
| `us-enrichment-api/*` | Property principal, geo-reference, secondary, secondary count, business, business name search, address search, ETag, and universal lookups |
| `international-street-api` | An international address lookup |
| `international-postal-code-api` | An international postal code lookup |
| `international-autocomplete-api` | International autocomplete suggestions |

To run one example with `make`, use the path under `examples/` with each slash replaced by a hyphen:

```bash
make us-street-api
make us-street-api-list
make us-enrichment-api-property-principal
```

To run one example without `make`, change to its directory and run it there. Some examples read files from the current directory.

```bash
cd examples/us-street-api && go run .
```

`make examples` runs every example. `make us-enrichment-api` runs every enrichment example.

## Development

```bash
make test      # go mod tidy, go fmt, then go test with coverage
make compile   # go build ./...
make build     # make test, then make compile
make cover     # open an HTML coverage report in the browser
make integrate # compile, test, then run every example (requires credentials)
```

To run one test:

```bash
go test -v -run TestName ./us-street-api
```

CI runs `make test` on Go 1.27 and on the latest stable Go release.

Library packages marshal and unmarshal with `encoding/json/v2`. Do not import `encoding/json` in library packages. Example programs can use either package.

## Changelog

See the [Go SDK changelog](https://github.com/smartystreets/changelog/blob/master/sdk/go.md) in the changelog repository.

## License

[Apache 2.0](LICENSE.md)
