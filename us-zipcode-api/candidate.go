package zipcode

type (
	// Result contains all output fields defined here:
	// https://smartystreets.com/docs/us-zipcode-api#http-response-output
	Result struct {
		InputID    string `json:"input_id,omitempty"`
		InputIndex int    `json:"input_index"`

		Status     string      `json:"status"`
		Reason     string      `json:"reason"`
		CityStates []CityState `json:"city_states"`
		ZIPCodes   []ZIPCode   `json:"zipcodes"`
	}

	// CityState contains all output fields defined here:
	// https://smartystreets.com/docs/us-zipcode-api#cities
	CityState struct {
		City              string `json:"city"`
		MailableCity      bool   `json:"mailable_city"`
		StateAbbreviation string `json:"state_abbreviation"`
		State             string `json:"state"`
	}

	// ZIPCode contains all output fields defined here:
	// https://smartystreets.com/docs/us-zipcode-api#zipcodes
	ZIPCode struct {
		County
		ZIPCode           string   `json:"zipcode"`
		ZIPCodeType       string   `json:"zipcode_type"`
		DefaultCity       string   `json:"default_city"`
		Latitude          float64  `json:"latitude"`
		Longitude         float64  `json:"longitude"`
		Precision         string   `json:"precision"`
		AlternateCounties []County `json:"alternate_counties"`
	}

	County struct {
		CountyFIPS        string `json:"county_fips"`
		CountyName        string `json:"county_name"`
		StateAbbreviation string `json:"state_abbreviation"`
		State             string `json:"state"`
	}
)
