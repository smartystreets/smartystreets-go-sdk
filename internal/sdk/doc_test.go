package sdk

import (
	"io"
	"net/http"
	"strconv"
)

type FakeHTTPClient struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (f *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	f.request = request
	return f.response, f.err
}

/*////////////////////////////////////////////////////////////////////////*/

type FakeMultiHTTPClient struct {
	requests      []*http.Request
	headers       []*http.Header
	bodies        []string
	responses     []*http.Response
	errors        []error
	call          int
	headerKey     string
	rateLimitTime int
}

func (f *FakeMultiHTTPClient) Do(request *http.Request) (*http.Response, error) {
	defer f.increment()
	f.simulateServerReadingRequestBody(request)
	f.requests = append(f.requests, request)
	response := f.responses[f.call]
	if response.StatusCode == 429 {
		response.Header = http.Header{}
		response.Header.Set(f.headerKey, strconv.Itoa(f.rateLimitTime))
	}
	return f.responses[f.call], f.errors[f.call]
}

func (f *FakeMultiHTTPClient) simulateServerReadingRequestBody(request *http.Request) {
	if request.Body != nil {
		body, _ := io.ReadAll(request.Body)
		f.bodies = append(f.bodies, string(body))
	} else {
		f.bodies = append(f.bodies, request.URL.Query().Get("body"))
	}
}

func (f *FakeMultiHTTPClient) increment() {
	f.call++
}

func NewFailingHTTPClient(statusCodes ...int) *FakeMultiHTTPClient {
	client := &FakeMultiHTTPClient{}
	for _, statusCode := range statusCodes {
		client.responses = append(client.responses, &http.Response{StatusCode: statusCode})
		client.errors = append(client.errors, nil)
	}
	return client
}

func NewErringHTTPClient(errors ...error) *FakeMultiHTTPClient {
	client := &FakeMultiHTTPClient{}
	for _, err := range errors {
		client.responses = append(client.responses, &http.Response{StatusCode: http.StatusOK})
		client.errors = append(client.errors, err)
	}
	return client
}
