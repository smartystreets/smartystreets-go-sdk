package us_enrichment

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/smarty/assertions/should"
	"github.com/smarty/gunit"
)

func TestClientFixture(t *testing.T) {
	gunit.Run(new(ClientFixture), t)
}

type ClientFixture struct {
	*gunit.Fixture

	sender *FakeSender
	client *Client

	input *Lookup
}

func (f *ClientFixture) Setup() {
	f.sender = &FakeSender{}
	f.client = NewClient(f.sender)
	f.input = new(Lookup)
}

func (f *ClientFixture) TestLookupSerializedAndSentWithContext__ResponseSuggestionsIncorporatedIntoLookup() {
	f.sender.response = validFinancialResponse
	f.input.SmartyKey = "12345"
	f.input.DataSet = "property"
	f.input.DataSubSet = "financial"

	ctx := context.WithValue(context.Background(), "key", "value")
	err := f.client.SendLookupWithContext(ctx, f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, "/lookup/"+f.input.SmartyKey+"/"+f.input.DataSet+"/"+f.input.DataSubSet)
	f.So(f.sender.request.Context(), should.Resemble, ctx)

	f.So(f.input.FinancialResponse, should.Resemble, []FinancialResponse{
		{
			SmartyKey:      "7",
			DataSetName:    "property",
			DataSubsetName: "financial",
			Attributes: FinancialAttributes{
				AssessedImprovementPercent: "Assessed_Improvement_Percent",
				VeteranTaxExemption:        "Veteran_Tax_Exemption",
				WidowTaxExemption:          "Widow_Tax_Exemption",
			},
		},
	})
}

func (f *ClientFixture) TestNilLookupNOP() {
	err := f.client.SendLookup(nil)
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestEmptyLookup_NOP() {
	err := f.client.SendLookup(new(Lookup))
	f.So(err, should.BeNil)
	f.So(f.sender.request, should.BeNil)
}

func (f *ClientFixture) TestSenderErrorPreventsDeserialization() {
	f.sender.err = errors.New("gophers")
	f.sender.response = validPrincipalResponse // would be deserialized if not for the err (above)
	f.input.SmartyKey = "12345"
	f.input.DataSet = "property"
	f.input.DataSubSet = "principal"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.PrincipalResponse, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.SmartyKey = "12345"
	f.input.DataSet = "property"
	f.input.DataSubSet = "principal"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.PrincipalResponse, should.BeEmpty)
}

var validFinancialResponse = `[{"smarty_key":"7","data_set_name":"property","data_subset_name":"financial","attributes":{"assessed_improvement_percent":"Assessed_Improvement_Percent","veteran_tax_exemption":"Veteran_Tax_Exemption","widow_tax_exemption":"Widow_Tax_Exemption"}}]`
var validPrincipalResponse = `[{"smarty_key":"7","data_set_name":"property","data_subset_name":"principal","attributes":{"1st_floor_sqft":"1st_Floor_Sqft",lender_name_2":"Lender_Name_2","lender_seller_carry_back":"Lender_Seller_Carry_Back","year_built":"Year_Built","zoning":"Zoning"}}]`

/**************************************************************************/

type FakeSender struct {
	callCount int
	request   *http.Request

	response string
	err      error
}

func (f *FakeSender) Send(request *http.Request) ([]byte, error) {
	f.callCount++
	f.request = request
	return []byte(f.response), f.err
}
