package sdk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"

	"github.com/smartystreets/logging"
)

// TODO: tests? (The tracing constructs and request context weren't designed with testing in mind...)

type TracingClient struct {
	inner  HTTPClient
	logger *logging.Logger
}

func NewTracingClient(inner HTTPClient, trace bool) HTTPClient {
	if !trace {
		return inner
	}
	return &TracingClient{inner: inner}
}

func (d *TracingClient) Do(request *http.Request) (*http.Response, error) {
	output := new(bytes.Buffer)
	trace := NewHTTPTrace()
	trace.SetLogOutput(output)
	request = trace.AttachToRequest(request)
	response, err := d.inner.Do(request)
	trace.DumpResponse(response, err)

	if err != nil || response.StatusCode != http.StatusOK {
		d.logger.Println("\n" + output.String())
	}

	return response, err
}

/**************************************************************************/

type HTTPTrace struct {
	trace  *httptrace.ClientTrace
	logger *logging.Logger
}

func NewHTTPTrace() *HTTPTrace {
	trace := new(HTTPTrace)
	trace.trace = &httptrace.ClientTrace{
		ConnectDone:          trace.ConnectDone,
		ConnectStart:         trace.ConnectStart,
		DNSDone:              trace.DNSDone,
		DNSStart:             trace.DNSStart,
		GetConn:              trace.GetConn,
		Got100Continue:       trace.Got100Continue,
		GotConn:              trace.GotConn,
		GotFirstResponseByte: trace.GotFirstResponseByte,
		PutIdleConn:          trace.PutIdleConn,
		Wait100Continue:      trace.Wait100Continue,
		WroteHeaders:         trace.WroteHeaders,
		WroteRequest:         trace.WroteRequest,
	}
	return trace
}

func (t *HTTPTrace) SetLogOutput(writer io.Writer) {
	t.logger.SetOutput(writer)
}

func (t *HTTPTrace) AttachToRequest(request *http.Request) *http.Request {
	t.logger.Println(dumpRequest(request))
	return request.WithContext(httptrace.WithClientTrace(request.Context(), t.trace))
}

func (t *HTTPTrace) DumpResponse(response *http.Response, err error) {
	t.logger.Println(dumpResponse(response, err))
}

func (t *HTTPTrace) GetConn(hostPort string) {
	t.logger.Println("GetConn hostPort:", hostPort)
}

func (t *HTTPTrace) GotConn(info httptrace.GotConnInfo) {
	t.logger.Printf("GotConn info: Conn: %s %#v\n", info.Conn.RemoteAddr(), info)
}

func (t *HTTPTrace) PutIdleConn(err error) {
	t.logger.Println("PutIdleConn err:", err)
}

func (t *HTTPTrace) GotFirstResponseByte() {
	t.logger.Println("GotFirstResponseByte")
}

func (t *HTTPTrace) Got100Continue() {
	t.logger.Println("Got100Continue")
}

func (t *HTTPTrace) DNSStart(info httptrace.DNSStartInfo) {
	t.logger.Printf("DNSStart info: %#v\n", info)
}

func (t *HTTPTrace) DNSDone(info httptrace.DNSDoneInfo) {
	addresses := new(bytes.Buffer)
	for i, address := range info.Addrs {
		fmt.Fprintf(addresses, "  %d. %s\n", i, address.String())
	}
	t.logger.Printf("DNSDone info: Coalesced: %t Err: %v Addresses:\n%s",
		info.Coalesced, info.Err, addresses.String())
}

func (t *HTTPTrace) ConnectStart(network, address string) {
	t.logger.Println("ConnectStart network:", network, "address:", address)
}

func (t *HTTPTrace) ConnectDone(network, address string, err error) {
	t.logger.Println("ConnectDone network:", network, "address:", address, "error:", err)
}

func (t *HTTPTrace) WroteHeaders() {
	t.logger.Println("WroteHeaders")
}

func (t *HTTPTrace) Wait100Continue() {
	t.logger.Println("Wait100Continue")
}

func (t *HTTPTrace) WroteRequest(info httptrace.WroteRequestInfo) {
	t.logger.Printf("WroteRequest info: %#v\n", info)
}
