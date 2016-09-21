package sdk

import "net/http"

type KeepAliveCloseClient struct {
	inner HTTPClient
}

func NewKeepAliveCloseClient(client HTTPClient, close bool) HTTPClient {
	if !close {
		return client
	}
	return &KeepAliveCloseClient{inner: client}
}

func (c *KeepAliveCloseClient) Do(request *http.Request) (*http.Response, error) {
	request.Close = true
	return c.inner.Do(request)
}

