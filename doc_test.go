package sdk

import "net/http"

//go:generate go install github.com/smartystreets/gunit/gunit
//go:generate gunit

type FakeHTTPClient struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (this *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	this.request = request
	return this.response, this.err
}

/*////////////////////////////////////////////////////////////////////////*/
