# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
# Run tests with coverage (also formats code and tidies modules)
make test

# Run a single test
go test -v -run TestName ./path/to/package

# Format code and tidy modules only
make fmt

# Build all packages
make compile

# Full build (test + compile)
make build

# Generate coverage report (opens HTML)
make cover

# Run every example program (requires API credentials)
make examples

# Run one example (target = its path under examples/, slashes replaced with hyphens)
make us-street-api
make us-enrichment-api-address-search

# Compile, test, then run every example program (requires API credentials)
make integrate
```

## Architecture

This is the official Go SDK for SmartyStreets address validation APIs. Compatible with Go 1.27. The SDK uses two key architectural patterns:

### Middleware/Decorator Pattern (internal/sdk/)

HTTP request processing is implemented as a chain of composable clients that wrap each other:

```
custom_query_client → license_client → keep_alive_close_client → base_url_client → custom_headers_client
  → signing_client → retry_client → debug_output_client → tracing_client → http.Client
```

Each middleware client implements the internal `HTTPClient` interface (`Do(*http.Request)`) and wraps another one, adding specific functionality (authentication, retries, base URL injection, custom headers, etc.). The outermost client is wrapped by `HTTPSender`, which implements the root-level `sdk.RequestSender` interface that the API packages consume.

### JSON Encoding (internal/json/)

All library code marshals and unmarshals through `internal/json`, a thin wrapper over `encoding/json/v2` that pins the SDK's option set: case-insensitive field matching, tolerance of invalid UTF-8 and duplicate object names, and deterministic output. Do not import `encoding/json` directly in library packages. Example programs are consumer code and may use either package.

### Builder Pattern (wireup/)

Clients are constructed using `wireup.Build*APIClient()` functions with functional options:

```go
client := wireup.BuildUSStreetAPIClient(
    wireup.SecretKeyCredential(authID, authToken),
    wireup.CustomHeader("X-Custom", "value"),
    wireup.MaxRetry(5),
)
```

### API Client Packages

Ten API packages follow identical structure:
- `us-street-api/`, `us-zipcode-api/`, `us-autocomplete-api`, `us-autocomplete-pro-api/`
- `us-enrichment-api/`, `us-extract-api/`, `us-reverse-geo-api/`
- `international-street-api/`, `international-postal-code-api/`, `international-autocomplete-api/`

Each contains: `client.go` (Client struct with Send methods), `lookup.go` (request struct), and usually `batch.go` for batch operations. Every client accepts a context and per-request credentials through the `*WithContextAndAuth` variants of its Send methods.

### Credential Types

Three authentication strategies defined at root level:
- `SecretKeyCredential` - adds auth-id/auth-token as query parameters
- `BasicAuthCredential` - HTTP Basic Authentication header
- `WebsiteKeyCredential` - for client-side applications with referrer restrictions
- Embedded/website keys are GET-only — not valid for batch (POST) requests or the US Extract API (POST-only): https://www.smarty.com/docs/cloud/authentication
