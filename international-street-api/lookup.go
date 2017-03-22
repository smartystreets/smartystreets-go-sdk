package international

import "net/url"

type Lookup struct {
	InputID        string
	Geocode        bool
	OutputLanguage OutputLanguage

	Freeform           string
	Address1           string
	Address2           string
	Address3           string
	Address4           string
	Organization       string
	Locality           string
	AdministrativeArea string
	PostalCode         string
	Country            string

	Results []*Result
}

func (l *Lookup) populate(query url.Values) {
	query.Set("input_id", l.InputID)
	query.Set("freeform", l.Freeform)
	query.Set("address1", l.Address1)
	query.Set("address2", l.Address2)
	query.Set("address3", l.Address3)
	query.Set("address4", l.Address4)
	query.Set("organization", l.Organization)
	query.Set("locality", l.Locality)
	query.Set("administrative_area", l.AdministrativeArea)
	query.Set("postal_code", l.PostalCode)
	query.Set("country", l.Country)
	query.Set("geocode", boolString(l.Geocode))
	query.Set("language", string(l.OutputLanguage))
}
func boolString(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

type OutputLanguage string

const (
	OutputDefault OutputLanguage = ""
	OutputNative  OutputLanguage = "native"
	OutputLatin   OutputLanguage = "latin"
)
