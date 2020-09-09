package us_reverse_geo

type resultListing struct {
	Listing []Address `json:"results"`
}

// Address fields defined here: https://smartystreets.com/docs/cloud/us-reverse-geo-api#http-response-output
type Address struct {
	Latitude          float64           `json:"latitude,omitempty"`
	Longitude         float64           `json:"longitude,omitempty"`
	CoordinateLicense CoordinateLicense `json:"coordinate_license,omitempty"`
	Distance          float64           `json:"distance,omitempty"`
	Street            string            `json:"street,omitempty"`
	City              string            `json:"city,omitempty"`
	StateAbbreviation string            `json:"state_abbreviation,omitempty"`
	ZIPCode           string            `json:"zipcode,omitempty"`
}

type CoordinateLicense uint16

// CoordinateLicense values and associated details defined here: https://smartystreets.com/docs/cloud/us-reverse-geo-api#licenses
const (
	CoordinateLicenseSmartyStreets  CoordinateLicense = 0
	CoordinateLicenseGatewaySpatial CoordinateLicense = 1
)

func (this CoordinateLicense) String() string {
	switch this {
	case CoordinateLicenseGatewaySpatial:
		return "Gateway Spatial, LLC"
	default:
		return "SmartyStreets"
	}
}
