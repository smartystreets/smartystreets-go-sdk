package extract

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestLookupFixture(t *testing.T) {
	gunit.Run(new(LookupFixture), t)
}

type LookupFixture struct {
	*gunit.Fixture

	lookup  *Lookup
	request *http.Request
}

func (f *LookupFixture) Setup() {
	f.lookup = new(Lookup)
	f.request, _ = http.NewRequest("POST", "/?hello=world", nil)
}

func (f *LookupFixture) query() url.Values {
	return f.request.URL.Query()
}

func readBody(request *http.Request) string {
	bytes, _ := ioutil.ReadAll(request.Body)
	return string(bytes)
}

func (f *LookupFixture) TestPopulate_TextCopiedToBody() {
	body := "Hello, World!"
	f.lookup.Text = body

	f.lookup.populate(f.request)

	f.So(readBody(f.request), should.Equal, body)
	f.So(f.request.ContentLength, should.Equal, len(body))
	f.So(f.request.Header.Get("Content-Type"), should.Equal, "text/plain")
}

func (f *LookupFixture) TestPopulate_NothingSet_NothingAddedToQueryStringOrBody() {
	f.lookup.populate(f.request)
	f.So(f.query().Encode(), should.Equal, "hello=world")
	f.So(f.request.Body, should.BeNil)
	f.So(f.request.ContentLength, should.Equal, 0)
}

func (f *LookupFixture) TestPopulate_HTMLYes() {
	f.lookup.HTML = HTMLYes
	f.lookup.populate(f.request)
	f.So(f.query().Get("html"), should.Equal, "true")
}

func (f *LookupFixture) TestPopulate_HTMLNo() {
	f.lookup.HTML = HTMLNo
	f.lookup.populate(f.request)
	f.So(f.query().Get("html"), should.Equal, "false")
}

func (f *LookupFixture) TestPopulate_SimpleParameters_Set() {
	f.lookup.AddressesPerLine = 42
	f.lookup.AddressesWithLineBreaks = true
	f.lookup.Aggressive = true

	f.lookup.populate(f.request)

	f.So(f.query().Get("aggressive"), should.Equal, "true")
	f.So(f.query().Get("addr_line_breaks"), should.Equal, "true")
	f.So(f.query().Get("addr_per_line"), should.Equal, "42")
}
