package sdk

import (
	"net/http"
	"net/url"
)

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
	return c.inner.Do(request)
}
