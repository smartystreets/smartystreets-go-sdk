package sdk

import (
	"net/http"
	"net/url"
	"path"
)

// BaseURLClient amends the http.Request.URL with the configured scheme, host, and path.
// This comes in handy when calling an onsite installation of the us-street-api.
type BaseURLClient struct {
	inner HTTPClient
	base  *url.URL
}

func NewBaseURLClient(inner HTTPClient, baseURL *url.URL) *BaseURLClient {
	return &BaseURLClient{
		inner: inner,
		base:  baseURL,
	}
}

func (c *BaseURLClient) Do(request *http.Request) (*http.Response, error) {
	request.URL.Scheme = c.base.Scheme
	request.URL.Host = c.base.Host
	request.URL.Path = path.Join(c.base.Path, request.URL.Path)
	return c.inner.Do(request)
}
