package sdk

type CoordinateLicense uint16

// CoordinateLicense values and associated details defined here: https://smartystreets.com/docs/cloud/us-reverse-geo-api#licenses
const (
	CoordinateLicenseSmartyStreets            CoordinateLicense = 0
	CoordinateLicenseSmartyStreetsProprietary CoordinateLicense = 1
)

func (this CoordinateLicense) String() string {
	switch this {
	case CoordinateLicenseSmartyStreetsProprietary:
		return "SmartyStreets Proprietary"
	default:
		return "SmartyStreets"
	}
}
