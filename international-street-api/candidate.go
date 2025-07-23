package street

type (
	// Candidate contains all output fields defined here:
	// https://smartystreets.com/docs/international-street-api#http-response-output
	Candidate struct {
		InputID string `json:"input_id,omitempty"`
		RootLevel
		Components Components `json:"components,omitempty"`
		Metadata   Metadata   `json:"metadata,omitempty"`
		Analysis   Analysis   `json:"analysis,omitempty"`
	}

	RootLevel struct {
		Organization string `json:"organization,omitempty"`
		Address1     string `json:"address1,omitempty"`
		Address2     string `json:"address2,omitempty"`
		Address3     string `json:"address3,omitempty"`
		Address4     string `json:"address4,omitempty"`
		Address5     string `json:"address5,omitempty"`
		Address6     string `json:"address6,omitempty"`
		Address7     string `json:"address7,omitempty"`
		Address8     string `json:"address8,omitempty"`
		Address9     string `json:"address9,omitempty"`
		Address10    string `json:"address10,omitempty"`
		Address11    string `json:"address11,omitempty"`
		Address12    string `json:"address12,omitempty"`
	}

	// Components contains all output fields defined here:
	// https://smartystreets.com/docs/international-street-api#components
	Components struct {
		SuperAdministrativeArea            string `json:"super_administrative_area,omitempty"`
		AdministrativeArea                 string `json:"administrative_area,omitempty"`
		AdministrativeAreaISO2             string `json:"administrative_area_iso2,omitempty"`
		AdministrativeAreaShort            string `json:"administrative_area_short,omitempty"`
		AdministrativeAreaLong             string `json:"administrative_area_long,omitempty"`
		SubAdministrativeArea              string `json:"sub_administrative_area,omitempty"`
		Building                           string `json:"building,omitempty"`
		DependentLocality                  string `json:"dependent_locality,omitempty"`
		DependentLocalityName              string `json:"dependent_locality_name,omitempty"`
		DoubleDependentLocality            string `json:"double_dependent_locality,omitempty"`
		CountryISO3                        string `json:"country_iso_3,omitempty"`
		Locality                           string `json:"locality,omitempty"`
		PostalCode                         string `json:"postal_code,omitempty"`
		PostalCodeShort                    string `json:"postal_code_short,omitempty"`
		PostalCodeExtra                    string `json:"postal_code_extra,omitempty"`
		Premise                            string `json:"premise,omitempty"`
		PremiseExtra                       string `json:"premise_extra,omitempty"`
		PremiseNumber                      string `json:"premise_number,omitempty"`
		PremiseType                        string `json:"premise_type,omitempty"`
		PremisePrefixNumber                string `json:"premise_prefix_number,omitempty"`
		Thoroughfare                       string `json:"thoroughfare,omitempty"`
		ThoroughfarePredirection           string `json:"thoroughfare_predirection,omitempty"`
		ThoroughfarePostdirection          string `json:"thoroughfare_postdirection,omitempty"`
		ThoroughfareName                   string `json:"thoroughfare_name,omitempty"`
		ThoroughfareTrailingType           string `json:"thoroughfare_trailing_type,omitempty"`
		ThoroughfareType                   string `json:"thoroughfare_type,omitempty"`
		DependentThoroughfare              string `json:"dependent_thoroughfare,omitempty"`
		DependentThoroughfarePredirection  string `json:"dependent_thoroughfare_predirection,omitempty"`
		DependentThoroughfarePostdirection string `json:"dependent_thoroughfare_postdirection,omitempty"`
		DependentThoroughfareName          string `json:"dependent_thoroughfare_name,omitempty"`
		DependentThoroughfareTrailingType  string `json:"dependent_thoroughfare_trailing_type,omitempty"`
		DependentThoroughfareType          string `json:"dependent_thoroughfare_type,omitempty"`
		BuildingLeadingType                string `json:"building_leading_type,omitempty"`
		BuildingName                       string `json:"building_name,omitempty"`
		BuildingTrailingType               string `json:"building_trailing_type,omitempty"`
		SubBuildingType                    string `json:"sub_building_type,omitempty"`
		SubBuildingNumber                  string `json:"sub_building_number,omitempty"`
		SubBuildingName                    string `json:"sub_building_name,omitempty"`
		SubBuilding                        string `json:"sub_building,omitempty"`
		LevelType                          string `json:"level_type,omitempty"`
		LevelNumber                        string `json:"level_number,omitempty"`
		PostBox                            string `json:"post_box,omitempty"`
		PostBoxType                        string `json:"post_box_type,omitempty"`
		PostBoxNumber                      string `json:"post_box_number,omitempty"`
		AdditionalContent                  string `json:"additional_content,omitempty"`
		DeliveryInstallation               string `json:"delivery_installation,omitempty"`
		DeliveryInstallationType           string `json:"delivery_installation_type,omitempty"`
		DeliveryInstallationQualifierName  string `json:"delivery_installation_qualifier_name,omitempty"`
		Route                              string `json:"route,omitempty"`
		RouteNumber                        string `json:"route_number,omitempty"`
		RouteType                          string `json:"route_type,omitempty"`
	}

	// Metadata contains all output fields defined here:
	// https://smartystreets.com/docs/international-street-api#metadata
	Metadata struct {
		Latitude            float64 `json:"latitude,omitempty"`
		Longitude           float64 `json:"longitude,omitempty"`
		GeocodePrecision    string  `json:"geocode_precision,omitempty"`
		MaxGeocodePrecision string  `json:"max_geocode_precision,omitempty"`
		AddressFormat       string  `json:"address_format,omitempty"`
		OccupantUse         string  `json:"occupant_use,omitempty"`
	}

	// Analysis contains all output fields defined here:
	// https://smartystreets.com/docs/international-street-api#analysis
	Analysis struct {
		VerificationStatus  string  `json:"verification_status,omitempty"`
		AddressPrecision    string  `json:"address_precision,omitempty"`
		MaxAddressPrecision string  `json:"max_address_precision,omitempty"`
		Changes             Changes `json:"changes,omitempty"`
	}

	// Changes contains all output fields defined here:
	// https://smartystreets.com/docs/international-street-api#changes
	Changes struct {
		RootLevel
		Components Components `json:"components,omitempty"`
	}
)
