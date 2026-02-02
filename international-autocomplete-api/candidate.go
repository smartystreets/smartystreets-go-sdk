package international_autocomplete_api

type Candidate struct {
	Street                  string `json:"street"`
	Locality                string `json:"locality"`
	AdministrativeArea      string `json:"administrative_area"`
	AdministrativeAreaShort string `json:"administrative_area_short"`
	AdministrativeAreaLong  string `json:"administrative_area_long"`
	PostalCode              string `json:"postal_code"`
	CountryIso3             string `json:"country_iso3"`

	Entries     int    `json:"entries"`
	AddressText string `json:"address_text"`
	AddressID   string `json:"address_id"`
}
