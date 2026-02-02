package sdk

import (
	"net/http"
	"net/url"
)

type CustomQueryClient struct {
	inner   HTTPClient
	queries url.Values
}

func NewCustomQueryClient(inner HTTPClient, queries url.Values) *CustomQueryClient {
	return &CustomQueryClient{inner: inner, queries: queries}
}

func (this *CustomQueryClient) Do(request *http.Request) (*http.Response, error) {
	if len(this.queries) > 0 {
		values := request.URL.Query()
		for key, v := range this.queries {
			values[key] = append(values[key], v...)
		}
		request.URL.RawQuery = values.Encode()
	}

	return this.inner.Do(request)
}
