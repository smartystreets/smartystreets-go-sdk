package international_autocomplete_api

type Result struct {
	Candidates []*Candidate `json:"candidates"`
}
