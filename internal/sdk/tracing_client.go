package sdk

import (
	"bytes"
	"fmt"
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
	trace := NewHTTPTraceLogger()
	request = trace.AttachToRequest(request)
	response, err := d.inner.Do(request)
	trace.DumpResponse(response, err)

	if err != nil || response.StatusCode != http.StatusOK {
		d.logger.Println("\n" + trace.String())
	}

	return response, err
}

/**************************************************************************/

type HTTPTraceLogger struct {
	*bytes.Buffer
	trace *httptrace.ClientTrace
}

func NewHTTPTraceLogger() *HTTPTraceLogger {
	trace := &HTTPTraceLogger{Buffer: new(bytes.Buffer)}
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

func (t *HTTPTraceLogger) AttachToRequest(request *http.Request) *http.Request {
	fmt.Fprintln(t, dumpRequest(request))
	return request.WithContext(httptrace.WithClientTrace(request.Context(), t.trace))
}

func (t *HTTPTraceLogger) DumpResponse(response *http.Response, err error) {
	fmt.Fprintln(t, dumpResponse(response, err))
}

func (t *HTTPTraceLogger) GetConn(hostPort string) {
	fmt.Fprintln(t, "GetConn hostPort:", hostPort)
}

func (t *HTTPTraceLogger) GotConn(info httptrace.GotConnInfo) {
	fmt.Fprintf(t, "GotConn info: Conn: %s %#v\n", info.Conn.RemoteAddr(), info)
}

func (t *HTTPTraceLogger) PutIdleConn(err error) {
	fmt.Fprintln(t, "PutIdleConn err:", err)
}

func (t *HTTPTraceLogger) GotFirstResponseByte() {
	fmt.Fprintln(t, "GotFirstResponseByte")
}

func (t *HTTPTraceLogger) Got100Continue() {
	fmt.Fprintln(t, "Got100Continue")
}

func (t *HTTPTraceLogger) DNSStart(info httptrace.DNSStartInfo) {
	fmt.Fprintf(t, "DNSStart info: %#v\n", info)
}

func (t *HTTPTraceLogger) DNSDone(info httptrace.DNSDoneInfo) {
	addresses := new(bytes.Buffer)
	for i, address := range info.Addrs {
		fmt.Fprintf(addresses, "  %d. %s\n", i, address.String())
	}
	fmt.Fprintf(t, "DNSDone info: Coalesced: %t Err: %v Addresses:\n%s",
		info.Coalesced, info.Err, addresses.String())
}

func (t *HTTPTraceLogger) ConnectStart(network, address string) {
	fmt.Fprintln(t, "ConnectStart network:", network, "address:", address)
}

func (t *HTTPTraceLogger) ConnectDone(network, address string, err error) {
	fmt.Fprintln(t, "ConnectDone network:", network, "address:", address, "error:", err)
}

func (t *HTTPTraceLogger) WroteHeaders() {
	fmt.Fprintln(t, "WroteHeaders")
}

func (t *HTTPTraceLogger) Wait100Continue() {
	fmt.Fprintln(t, "Wait100Continue")
}

func (t *HTTPTraceLogger) WroteRequest(info httptrace.WroteRequestInfo) {
	fmt.Fprintf(t, "WroteRequest info: %#v\n", info)
}
