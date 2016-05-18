package sdk

import "net/http"

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
