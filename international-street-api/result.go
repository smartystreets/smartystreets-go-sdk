package international

type Result struct {
	InputID      string `json:"input_id"`
	Organization string `json:"organization"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	Address3     string `json:"address3"`
	Address4     string `json:"address4"`
	Address5     string `json:"address5"`
	Address6     string `json:"address6"`
	Address7     string `json:"address7"`
	Address8     string `json:"address8"`
	Address9     string `json:"address9"`
	Address10    string `json:"address10"`
	Address11    string `json:"address11"`
	Address12    string `json:"address12"`

	Components Components `json:"components"`
	Metadata   Metadata   `json:"metadata"`
	Analysis   Analysis   `json:"analysis"`
}

type Components struct {
	CountryISO3                        string `json:"country_iso_3"`
	SuperAdministrativeArea            string `json:"super_administrative_area"`
	AdministrativeArea                 string `json:"administrative_area"`
	SubAdministrativeArea              string `json:"sub_administrative_area"`
	DependentLocality                  string `json:"dependent_locality"`
	DependentLocalityName              string `json:"dependent_locality_name"`
	DoubleDependentLocality            string `json:"double_dependent_locality"`
	Locality                           string `json:"locality"`
	PostalCode                         string `json:"postal_code"`
	PostalCodeShort                    string `json:"postal_code_short"`
	PostalCodeExtra                    string `json:"postal_code_extra"`
	Premise                            string `json:"premise"`
	PremiseExtra                       string `json:"premise_extra"`
	PremiseNumber                      string `json:"premise_number"`
	PremiseType                        string `json:"premise_type"`
	Thoroughfare                       string `json:"thoroughfare"`
	ThoroughfarePredirection           string `json:"thoroughfare_predirection"`
	ThoroughfarePostdirection          string `json:"thoroughfare_postdirection"`
	ThoroughfareName                   string `json:"thoroughfare_name"`
	ThoroughfareTrailingType           string `json:"thoroughfare_trailing_type"`
	ThoroughfareType                   string `json:"thoroughfare_type"`
	DependentThoroughfare              string `json:"dependent_thoroughfare"`
	DependentThoroughfarePredirection  string `json:"dependent_thoroughfare_predirection"`
	DependentThoroughfarePostdirection string `json:"dependent_thoroughfare_postdirection"`
	DependentThoroughfareName          string `json:"dependent_thoroughfare_name"`
	DependentThoroughfareTrailingType  string `json:"dependent_thoroughfare_trailing_type"`
	DependentThoroughfareType          string `json:"dependent_thoroughfare_type"`
	Building                           string `json:"building"`
	BuildingLeadingType                string `json:"building_leading_type"`
	BuildingName                       string `json:"building_name"`
	BuildingTrailingType               string `json:"building_trailing_type"`
	SubBuildingType                    string `json:"sub_building_type"`
	SubBuildingNumber                  string `json:"sub_building_number"`
	SubBuildingName                    string `json:"sub_building_name"`
	SubBuilding                        string `json:"sub_building"`
	PostBox                            string `json:"post_box"`
	PostBoxType                        string `json:"post_box_type"`
	PostBoxNumber                      string `json:"post_box_number"`
}

type Metadata struct {
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	GeocodePrecision    string  `json:"geocode_precision"`
	MaxGeocodePrecision string  `json:"max_geocode_precision"`
}

type Analysis struct {
	VerificationStatus  string `json:"verification_status"`
	AddressPrecision    string `json:"address_precision"`
	MaxAddressPrecision string `json:"max_address_precision"`
}
