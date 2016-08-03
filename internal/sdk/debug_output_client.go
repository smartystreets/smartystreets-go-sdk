package sdk

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/smartystreets/logging"
)

// DebugOutputClient makes use of functions from net/http/httputil to log http request/response details.
type DebugOutputClient struct {
	inner  HTTPClient
	logger *logging.Logger
}

func NewDebugOutputClient(inner HTTPClient, dump bool) HTTPClient {
	if !dump {
		return inner
	}
	return &DebugOutputClient{inner: inner}
}

func (d *DebugOutputClient) Do(request *http.Request) (*http.Response, error) {
	d.logger.Println(dumpRequest(request))
	response, err := d.inner.Do(request)
	d.logger.Println(dumpResponse(response, err))
	return response, err
}

func dumpRequest(request *http.Request) string {
	dump, err := httputil.DumpRequestOut(request, true)
	return composeDump("request", string(dump), err)
}

func dumpResponse(response *http.Response, err error) string {
	if err != nil {
		return composeDump("err", err.Error(), nil)
	}
	dump, err := httputil.DumpResponse(response, true)
	return composeDump("response", string(dump), err)
}

func composeDump(title string, dump string, err error) string {
	if err != nil {
		return fmt.Sprintf("Could not dump HTTP %s: %s\n", title, err.Error())
	} else {
		return fmt.Sprintf("HTTP %s:\n%s\n", strings.Title(title), dump)
	}
}
