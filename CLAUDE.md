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

# Run integration tests (executes all example programs - requires API credentials)
make integrate
```

## Architecture

This is the official Go SDK for SmartyStreets address validation APIs. The SDK uses two key architectural patterns:

### Middleware/Decorator Pattern (internal/sdk/)

HTTP request processing is implemented as a chain of composable clients that wrap each other:

```
signing_client → retry_client → base_url_client → custom_header_client → http_sender
```

Each middleware client implements the `sdk.RequestSender` interface and wraps another sender, adding specific functionality (authentication, retries, base URL injection, custom headers, etc.).

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

Nine API packages follow identical structure:
- `us-street-api/`, `us-zipcode-api/`, `us-autocomplete-pro-api/`
- `us-enrichment-api/`, `us-extract-api/`, `us-reverse-geo-api/`
- `international-street-api/`, `international-postal-code-api/`, `international-autocomplete-api/`

Each contains: `client.go` (Client struct with Send methods), `lookup.go` (request struct), and usually `batch.go` for batch operations. All clients support context and per-request authentication via `SendBatchWithContextAndAuth()`.

### Credential Types

Three authentication strategies defined at root level:
- `SecretKeyCredential` - adds auth-id/auth-token as query parameters
- `BasicAuthCredential` - HTTP Basic Authentication header
- `WebsiteKeyCredential` - for client-side applications with referrer restrictions
