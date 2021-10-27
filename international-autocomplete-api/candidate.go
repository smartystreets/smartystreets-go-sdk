package international_autocomplete_api

type Candidate struct {
	Street             string `json:"street"`
	Locality           string `json:"locality"`
	AdministrativeArea string `json:"administrative_area"`
	PostalCode         string `json:"postal_code"`
	CountryIso3        string `json:"country_iso3"`
}
