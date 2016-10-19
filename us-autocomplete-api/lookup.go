package autocomplete

type (
	// Lookup represents all input fields documented here:
	// https://smartystreets.com/docs/us-autocomplete-api#http-request-input-fields
	Lookup struct {
		Prefix               string
		MaxSuggestions       int
		CityFilter           []string
		StateFilter          []string
		CityStatePreferences []string
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
