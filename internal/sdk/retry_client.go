package sdk

import "net/http"

type RetryClient struct {
	inner      httpClient
	maxRetries int
}

func NewRetryClient(client httpClient, maxRetries int) *RetryClient {
	return &RetryClient{
		inner:      client,
		maxRetries: maxRetries,
	}
}

func (this *RetryClient) Do(request *http.Request) (response *http.Response, err error) {
	for attempt := 0; attempt <= this.maxRetries; attempt++ {
		response, err = this.inner.Do(request)
		if err == nil && response.StatusCode == http.StatusOK {
			break
		}
	}
	return response, err
}
