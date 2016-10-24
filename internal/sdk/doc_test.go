package sdk

import (
	"io/ioutil"
	"net/http"
)

//go:generate go install github.com/smartystreets/gunit/gunit
//go:generate gunit

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
	requests  []*http.Request
	bodies    []string
	responses []*http.Response
	errors    []error
	call      int
}

func (f *FakeMultiHTTPClient) Do(request *http.Request) (*http.Response, error) {
	defer f.increment()
	f.simulateServerReadingRequestBody(request)
	f.requests = append(f.requests, request)
	return f.responses[f.call], f.errors[f.call]
}

func (f *FakeMultiHTTPClient) simulateServerReadingRequestBody(request *http.Request) {
	if request.Body != nil {
		body, _ := ioutil.ReadAll(request.Body)
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
