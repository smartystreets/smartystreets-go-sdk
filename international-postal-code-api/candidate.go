package international_postal_code

type (
	// Candidate contains all output fields defined here:
	// https://smartystreets.com/docs/international-postal-code-api
	Candidate struct {
		InputID                 string `json:"input_id,omitempty"`
		AdministrativeArea      string `json:"administrative_area,omitempty"`
		SubAdministrativeArea   string `json:"sub_administrative_area,omitempty"`
		SuperAdministrativeArea string `json:"super_administrative_area,omitempty"`
		CountryIso3             string `json:"country_iso_3,omitempty"`
		Locality                string `json:"locality,omitempty"`
		DependentLocality       string `json:"dependent_locality,omitempty"`
		DependentLocalityName   string `json:"dependent_locality_name,omitempty"`
		DoubleDependentLocality string `json:"double_dependent_locality,omitempty"`
		//PostalCode              string `json:"postal_code,omitempty"`
		PostalCodeShort string `json:"postal_code,omitempty"`
		PostalCodeExtra string `json:"postal_code_extra,omitempty"`
	}
)
