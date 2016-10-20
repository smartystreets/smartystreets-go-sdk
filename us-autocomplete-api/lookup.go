package autocomplete

type (
	// Lookup represents all input fields documented here:
	// https://smartystreets.com/docs/us-autocomplete-api#http-request-input-fields
	Lookup struct {
		Prefix               string   `json:"prefix,omitempty"`
		MaxSuggestions       int      `json:"suggestions,omitempty"`
		CityFilter           []string `json:"city_filter,omitempty"`
		StateFilter          []string `json:"state_filter,omitempty"`
		CityStatePreferences []string `json:"prefer,omitempty"`
		Geolocation          GeolocateMode

		Results []*Suggestion
	}

	GeolocateMode int
)

const (
	GeolocateNone GeolocateMode = 0
	GeolocateState
	GeolocateCity
)
