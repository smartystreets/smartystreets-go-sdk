package sdk

import (
	"net/http"
	"strings"
)

type CustomHeadersClient struct {
	inner         HTTPClient
	headers       http.Header
	appendHeaders map[string]string
}

func NewCustomHeadersClient(inner HTTPClient, headers http.Header, appendHeaders map[string]string) *CustomHeadersClient {
	return &CustomHeadersClient{
		inner:         inner,
		headers:       headers,
		appendHeaders: appendHeaders,
	}
}

func (this *CustomHeadersClient) Do(request *http.Request) (*http.Response, error) {
	this.addHeaders(request)
	return this.inner.Do(request)
}

func (this *CustomHeadersClient) addHeaders(request *http.Request) {
	headers := request.Header

	for key, values := range this.headers {
		if separator, ok := this.appendHeaders[key]; ok {
			headers.Set(key, strings.Join(values, separator))
		} else {
			for _, value := range values {
				if key == "Host" {
					request.Host = value
				} else {
					headers.Add(key, value)
				}
			}
		}
	}
}
