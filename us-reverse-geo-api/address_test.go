package us_reverse_geo

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestAddressFixture(t *testing.T) {
    gunit.Run(new(AddressFixture), t)
}

type AddressFixture struct {
    *gunit.Fixture
}

func (this *AddressFixture) TestLicenseString() {
	this.So(LicenseSmartyStreets.String(), should.Equal, "SmartyStreets")
	this.So(LicenseGatewaySpatial.String(), should.Equal, "Gateway Spatial, LLC")
	this.So(License(42).String(), should.Equal, "SmartyStreets")
}
