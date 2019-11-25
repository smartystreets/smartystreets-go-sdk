package street

import (
	"errors"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
	"net/http"
	"testing"
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

func (f *ClientFixture) TestAddressLookupSerializedAndSent__ResponseSuggestionsIncorporatedIntoLookup() {
	f.sender.response = `[
		{"address1": "1"},
		{"address1": "2"},
		{"address1": "3"}
	]`
	f.input.Freeform = "42"

	err := f.client.SendLookup(f.input)

	f.So(err, should.BeNil)
	f.So(f.sender.request, should.NotBeNil)
	f.So(f.sender.request.Method, should.Equal, "GET")
	f.So(f.sender.request.URL.Path, should.Equal, verifyURL)
	f.So(string(f.sender.request.URL.Query().Get("freeform")), should.Equal, "42")
	f.So(f.sender.request.URL.String(), should.Equal, verifyURL+"?freeform=42")
	f.So(f.input.Results, should.Resemble, []*Candidate{
		{RootLevel: RootLevel{Address1: "1"}},
		{RootLevel: RootLevel{Address1: "2"}},
		{RootLevel: RootLevel{Address1: "3"}}})
	f.So(f.input.Results, should.Resemble, []*Candidate{
		{RootLevel: RootLevel{Address1: "1"}},
		{RootLevel: RootLevel{Address1: "2"}},
		{RootLevel: RootLevel{Address1: "3"}},
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
	f.sender.err = errors.New("GOPHERS!")
	f.sender.response = `{"suggestions":[
		{"text": "1"},
		{"text": "2"},
		{"text": "3"}
	]}` // would be deserialized if not for the err (above)
	f.input.Freeform = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeEmpty)
}

func (f *ClientFixture) TestDeserializationErrorPreventsDeserialization() {
	f.sender.response = `I can't haz JSON`
	f.input.Freeform = "HI"

	err := f.client.SendLookup(f.input)

	f.So(err, should.NotBeNil)
	f.So(f.input.Results, should.BeEmpty)
}

func (f *ClientFixture) FocusTestFullJSONResponseDeserialization() {
	f.sender.response = `[
{
	"input_id": "blah",
	"organization": "blah",
	"address1": "Rua Antônio Loes Marin, 121",
	"address2": "Casa Verde",
	"address3": "São Paulo - SP",
	"address4": "02516-050",
	"address5": "blank",
	"address6": "also empty",
	"address7": "there",
	"address8": "is",
	"address9": "nothing",
	"address10": "to",
	"address11": "show",
	"address12": "here",
	"components": {
	  "super_administrative_area": "super_blah",
	  "administrative_area": "SP",
	  "sub_administrative_area": "sub_blah",
	  "building": "building",
	  "dependent_locality": "Casa Verde",
	  "dependent_locality_name": "dependent_locality_name",
	  "double_dependent_locality": "x2",
	  "country_iso_3": "BRA",
	  "locality": "São Paulo",
	  "postal_code": "02516-050",
	  "postal_code_short": "02516-050",
	  "postal_code_extra": "ditto++",
	  "premise": "121",
	  "premise_extra": "premise_extra",
	  "premise_number": "121",
	  "premise_type": "premise_type",
	  "premise_prefix_number": "prefix",
	  "thoroughfare": "Rua Antônio Lopes Marin",
	  "thoroughfare_predirection": "Q",
	  "thoroughfare_postdirection": "K",
	  "thoroughfare_name": "Rua Antônio Lopes Marin",
	  "thoroughfare_trailing_type": "Rua",
	  "thoroughfare_type": "empty",
	  "dependent_thoroughfare": "dependent",
	  "dependent_thoroughfare_predirection": "before_dependent",
	  "dependent_thoroughfare_postdirection": "after_dependent",
	  "dependent_thoroughfare_name": "dependent_name",
	  "dependent_thoroughfare_trailing_type": "dependent_trail",
	  "dependent_thoroughfare_type": "dependent_type",
	  "building_leading_type": "leading_type",
	  "building_name": "building_name",
	  "building_trailing_type": "building_trail",
	  "sub_building_type": "almost_bldg_type",
	  "sub_building_number": "almost_bldg_number",
	  "sub_building_name": "almost_bldg_name",
	  "sub_building": "almost_bldg",
	  "post_box": "box",
	  "post_box_type": "cube",
	  "post_box_number": "blank"
	},
	"metadata": {
	  "latitude": -23.509659,
	  "longitude": -46.659711,
	  "geocode_precision": "Premise",
	  "max_geocode_precision": "DeliveryPoint",
	  "address_format": "thoroughfare, premise|dependent_locality|locality - administrative_area|postal_code"
	},
    "analysis": {
	  "verification_status": "Verified",
	  "address_precision": "Premise",
	  "max_address_precision": "DeliveryPoint",
	  "changes": {
		"organization": "blank",
		"address1": "Verified-AliasChange",
		"address2": "Verified-AliasChange",
		"address3": "Verified-AliasChange",
		"address4": "Verified-AliasChange",
		"address5": "5",
		"address6": "6",
		"address7": "7",
		"address8": "8",
		"address9": "9",
		"address10": "10",
		"address11": "11",
		"address12": "12",
		"components": {
	 	  "super_administrative_area": "blank",
		  "administrative_area": "Verified-NoChange",
		  "sub_administrative_area": "blank",
		  "building": "blank",
		  "dependent_locality": "Added",
		  "dependent_locality_name": "blank",
		  "double_dependent_locality": "blank",
		  "country_iso_3": "Added",
		  "locality": "Verified-AliasChange",
		  "postal_code": "Verified-SmallChange",
		  "postal_code_short": "Verified-SmallChange",
		  "postal_code_extra": "blank",
		  "premise": "Verified-NoChange",
		  "premise_extra": "blank",
		  "premise_number": "Verified-NoChange",
		  "premise_type": "blank",
		  "premise_prefix_number": "blank",
		  "thoroughfare": "Verified-SmallChange",
		  "thoroughfare_predirection": "blank",
		  "thoroughfare_postdirection": "blank",
		  "thoroughfare_name": "Identified-ContextChange",
		  "thoroughfare_trailing_type": "blank",
		  "thoroughfare_type": "Identified-AliasChange",
		  "dependent_thoroughfare": "blank",
		  "dependent_thoroughfare_predirection": "blank",
		  "dependent_thoroughfare_postdirection": "blank",
		  "dependent_thoroughfare_name": "blank",
		  "dependent_thoroughfare_trailing_type": "blank",
		  "dependent_thoroughfare_type": "blank",
		  "building_leading_type": "blank",
		  "building_name": "blank",
		  "building_trailing_type": "blank",
		  "sub_building_type": "blank",
		  "sub_building_number": "blank",
		  "sub_building_name": "blank",
		  "sub_building": "blank",
		  "post_box": "blank",
		  "post_box_type": "blank",
		  "post_box_number": "blank"
		}
	  }
	}
}
]`
	lookup := new(Lookup)
	response := []byte(f.sender.response)
	err := deserializeResponse(response, lookup)
	candidate := lookup.Results[0]
	component := candidate.Components
	metadata := candidate.Metadata
	analysis := candidate.Analysis
	changes := analysis.Changes
	ccomponents := changes.Components
	f.So(err, should.BeNil)
	f.So(candidate.InputID, should.Equal, "blah")
	f.So(candidate.Organization, should.Equal, "blah")
	f.So(candidate.Address1, should.Equal, "Rua Antônio Loes Marin, 121")
	f.So(candidate.Address2, should.Equal, "Casa Verde")
	f.So(candidate.Address3, should.Equal, "São Paulo - SP")
	f.So(candidate.Address4, should.Equal, "02516-050")
	f.So(candidate.Address5, should.Equal, "blank")
	f.So(candidate.Address6, should.Equal, "also empty")
	f.So(candidate.Address7, should.Equal, "there")
	f.So(candidate.Address8, should.Equal, "is")
	f.So(candidate.Address9, should.Equal, "nothing")
	f.So(candidate.Address10, should.Equal, "to")
	f.So(candidate.Address11, should.Equal, "show")
	f.So(candidate.Address12, should.Equal, "here")
	f.So(component.SuperAdministrativeArea, should.Equal, "super_blah")
	f.So(component.AdministrativeArea, should.Equal, "SP")
	f.So(component.SubAdministrativeArea, should.Equal, "sub_blah")
	f.So(component.Building, should.Equal, "building")
	f.So(component.DependentLocality, should.Equal, "Casa Verde")
	f.So(component.DependentLocalityName, should.Equal, "dependent_locality_name")
	f.So(component.DoubleDependentLocality, should.Equal, "x2")
	f.So(component.CountryISO3, should.Equal, "BRA")
	f.So(component.Locality, should.Equal, "São Paulo")
	f.So(component.PostalCode, should.Equal, "02516-050")
	f.So(component.PostalCodeShort, should.Equal, "02516-050")
	f.So(component.PostalCodeExtra, should.Equal, "ditto++")
	f.So(component.Premise, should.Equal, "121")
	f.So(component.PremiseExtra, should.Equal, "premise_extra")
	f.So(component.PremiseNumber, should.Equal, "121")
	f.So(component.PremiseType, should.Equal, "premise_type")
	f.So(component.PremisePrefixNumber, should.Equal, "prefix")
	f.So(component.Thoroughfare, should.Equal, "Rua Antônio Lopes Marin")
	f.So(component.ThoroughfarePredirection, should.Equal, "Q")
	f.So(component.ThoroughfarePostdirection, should.Equal, "K")
	f.So(component.ThoroughfareName, should.Equal, "Rua Antônio Lopes Marin")
	f.So(component.ThoroughfareTrailingType, should.Equal, "Rua")
	f.So(component.ThoroughfareType, should.Equal, "empty")
	f.So(component.DependentThoroughfare, should.Equal, "dependent")
	f.So(component.DependentThoroughfarePredirection, should.Equal, "before_dependent")
	f.So(component.DependentThoroughfarePostdirection, should.Equal, "after_dependent")
	f.So(component.DependentThoroughfareName, should.Equal, "dependent_name")
	f.So(component.DependentThoroughfareTrailingType, should.Equal, "dependent_trail")
	f.So(component.DependentThoroughfareType, should.Equal, "dependent_type")
	f.So(component.BuildingLeadingType, should.Equal, "leading_type")
	f.So(component.BuildingName, should.Equal, "building_name")
	f.So(component.BuildingTrailingType, should.Equal, "building_trail")
	f.So(component.SubBuildingType, should.Equal, "almost_bldg_type")
	f.So(component.SubBuildingNumber, should.Equal, "almost_bldg_number")
	f.So(component.SubBuildingName, should.Equal, "almost_bldg_name")
	f.So(component.SubBuilding, should.Equal, "almost_bldg")
	f.So(component.PostBox, should.Equal, "box")
	f.So(component.PostBoxType, should.Equal, "cube")
	f.So(component.PostBoxNumber, should.Equal, "blank")
	f.So(metadata.Latitude, should.Equal, -23.509659)
	f.So(metadata.Longitude, should.Equal, -46.659711)
	f.So(metadata.GeocodePrecision, should.Equal, "Premise")
	f.So(metadata.MaxGeocodePrecision, should.Equal, "DeliveryPoint")
	f.So(metadata.AddressFormat, should.Equal, "thoroughfare, premise|dependent_locality|locality - administrative_area|postal_code")
	f.So(analysis.VerificationStatus, should.Equal, "Verified")
	f.So(analysis.AddressPrecision, should.Equal, "Premise")
	f.So(analysis.MaxAddressPrecision, should.Equal, "DeliveryPoint")
	f.So(changes.Organization, should.Equal, "blank")
	f.So(changes.Address1, should.Equal, "Verified-AliasChange")
	f.So(changes.Address2, should.Equal, "Verified-AliasChange")
	f.So(changes.Address3, should.Equal, "Verified-AliasChange")
	f.So(changes.Address4, should.Equal, "Verified-AliasChange")
	f.So(changes.Address5, should.Equal, "5")
	f.So(changes.Address6, should.Equal, "6")
	f.So(changes.Address7, should.Equal, "7")
	f.So(changes.Address8, should.Equal, "8")
	f.So(changes.Address9, should.Equal, "9")
	f.So(changes.Address10, should.Equal, "10")
	f.So(changes.Address11, should.Equal, "11")
	f.So(changes.Address12, should.Equal, "12")
	f.So(ccomponents.SuperAdministrativeArea, should.Equal, "blank")
	f.So(ccomponents.AdministrativeArea, should.Equal, "Verified-NoChange")
	f.So(ccomponents.SubAdministrativeArea, should.Equal, "blank")
	f.So(ccomponents.Building, should.Equal, "blank")
	f.So(ccomponents.DependentLocality, should.Equal, "Added")
	f.So(ccomponents.DependentLocalityName, should.Equal, "blank")
	f.So(ccomponents.DoubleDependentLocality, should.Equal, "blank")
	f.So(ccomponents.CountryISO3, should.Equal, "Added")
	f.So(ccomponents.Locality, should.Equal, "Verified-AliasChange")
	f.So(ccomponents.PostalCode, should.Equal, "Verified-SmallChange")
	f.So(ccomponents.PostalCodeShort, should.Equal, "Verified-SmallChange")
	f.So(ccomponents.PostalCodeExtra, should.Equal, "blank")
	f.So(ccomponents.Premise, should.Equal, "Verified-NoChange")
	f.So(ccomponents.PremiseExtra, should.Equal, "blank")
	f.So(ccomponents.PremiseNumber, should.Equal, "Verified-NoChange")
	f.So(ccomponents.PremiseType, should.Equal, "blank")
	f.So(ccomponents.PremisePrefixNumber, should.Equal, "blank")
	f.So(ccomponents.Thoroughfare, should.Equal, "Verified-SmallChange")
	f.So(ccomponents.ThoroughfarePredirection, should.Equal, "blank")
	f.So(ccomponents.ThoroughfarePostdirection, should.Equal, "blank")
	f.So(ccomponents.ThoroughfareName, should.Equal, "Identified-ContextChange")
	f.So(ccomponents.ThoroughfareTrailingType, should.Equal, "blank")
	f.So(ccomponents.ThoroughfareType, should.Equal, "Identified-AliasChange")
	f.So(ccomponents.DependentThoroughfare, should.Equal, "blank")
	f.So(ccomponents.DependentThoroughfarePredirection, should.Equal, "blank")
	f.So(ccomponents.DependentThoroughfarePostdirection, should.Equal, "blank")
	f.So(ccomponents.DependentThoroughfareName, should.Equal, "blank")
	f.So(ccomponents.DependentThoroughfareTrailingType, should.Equal, "blank")
	f.So(ccomponents.DependentThoroughfareType, should.Equal, "blank")
	f.So(ccomponents.BuildingLeadingType, should.Equal, "blank")
	f.So(ccomponents.BuildingName, should.Equal, "blank")
	f.So(ccomponents.BuildingTrailingType, should.Equal, "blank")
	f.So(ccomponents.SubBuildingType, should.Equal, "blank")
	f.So(ccomponents.SubBuildingNumber, should.Equal, "blank")
	f.So(ccomponents.SubBuildingName, should.Equal, "blank")
	f.So(ccomponents.SubBuilding, should.Equal, "blank")
	f.So(ccomponents.PostBox, should.Equal, "blank")
	f.So(ccomponents.PostBoxType, should.Equal, "blank")
	f.So(ccomponents.PostBoxNumber, should.Equal, "blank")
}

/*////////////////////////////////////////////////////////////////////////*/

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
