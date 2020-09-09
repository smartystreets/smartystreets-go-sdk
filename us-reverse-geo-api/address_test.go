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
	this.So(CoordinateLicenseSmartyStreets.String(), should.Equal, "SmartyStreets")
	this.So(CoordinateLicenseGatewaySpatial.String(), should.Equal, "Gateway Spatial, LLC")
	this.So(CoordinateLicense(42).String(), should.Equal, "SmartyStreets")
}
