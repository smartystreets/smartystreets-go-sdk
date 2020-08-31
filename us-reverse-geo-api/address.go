package us_reverse_geo

type resultListing struct {
	Listing []Address `json:"results"`
}

// Address fields defined here: https://smartystreets.com/docs/cloud/us-reverse-geo-api#http-response-output
type Address struct {
	Latitude          float64 `json:"latitude,omitempty"`
	Longitude         float64 `json:"longitude,omitempty"`
	License           int     `json:"license,omitempty"`
	Distance          float64 `json:"distance,omitempty"`
	Street            string  `json:"street,omitempty"`
	City              string  `json:"city,omitempty"`
	StateAbbreviation string  `json:"state,omitempty"`
	ZIPCode           string  `json:"zipcode,omitempty"`
}

// License values defined here: https://smartystreets.com/docs/cloud/us-reverse-geo-api#licenses
const (
	LicenseSmartyStreets  = 0
	LicenseGatewaySpatial = 1
)
