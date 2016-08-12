package sdk

import "net/http"

type CustomHeadersClient struct {
	inner   HTTPClient
	headers http.Header
}

func NewCustomHeadersClient(inner HTTPClient, headers http.Header) *CustomHeadersClient {
	return &CustomHeadersClient{
		inner:   inner,
		headers: headers,
	}
}

func (this *CustomHeadersClient) Do(request *http.Request) (*http.Response, error) {
	this.addHeaders(request.Header)
	return this.inner.Do(request)
}

func (this *CustomHeadersClient) addHeaders(headers http.Header) {
	for key, values := range this.headers {
		for _, value := range values {
			headers.Add(key, value)
		}
	}
}
