package extract

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Lookup represents all input fields documented here:
// https://smartystreets.com/docs/cloud/us-extract-api#http-request-input-fields
type Lookup struct {
	Text string

	HTML                    HTMLPayload
	Aggressive              bool
	AddressesWithLineBreaks bool
	AddressesPerLine        int

	Result *Result
}

type HTMLPayload string

const (
	HTMLUnspecified HTMLPayload = ""      // Indicates that the server may decide if Lookup.Text is HTML or not.
	HTMLYes         HTMLPayload = "true"  // Indicates that the Lookup.Text is known to be HTML.
	HTMLNo          HTMLPayload = "false" // Indicates that the Lookup.Text is known to NOT be HTML.
)

func (l *Lookup) populate(request *http.Request) {
	l.setQuery(request)
	l.setBody(request)
	l.setHeaders(request)
}
func (l *Lookup) setQuery(request *http.Request) {
	query := request.URL.Query()

	if l.HTML != HTMLUnspecified {
		query.Set("html", string(l.HTML))
	}

	if l.Aggressive {
		query.Set("aggressive", "true")
	}

	if l.AddressesWithLineBreaks {
		query.Set("addr_line_breaks", "true")
	}

	if l.AddressesPerLine > 0 {
		query.Set("addr_per_line", strconv.Itoa(l.AddressesPerLine))
	}

	request.URL.RawQuery = query.Encode()
}
func (l *Lookup) setBody(request *http.Request) {
	if len(l.Text) == 0 {
		return
	}

	body := strings.NewReader(l.Text)
	request.Body = ioutil.NopCloser(body)
	request.ContentLength = int64(body.Len())
}
func (l *Lookup) setHeaders(request *http.Request) {
	request.Header.Set("Content-Type", "text/plain")
}
