package us_street

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// Batch stores input records and settings related to a group of addresses to be verified in a batch.
type Batch struct {
	records []*Input

	standardizeOnly bool
	includeInvalid  bool
}

// NewBatch creates a new, empty batch.
func NewBatch() *Batch {
	return &Batch{}
}

// StandardizeOnly sets the X-Standardize-Only header value.
func (b *Batch) StandardizeOnly(on bool) {
	b.standardizeOnly = on
}

// IncludeInvalid sets the X-Include-Invalid header value.
func (b *Batch) IncludeInvalid(on bool) {
	b.includeInvalid = on
}

// Append includes the record in the collection to be sent if there is still room (max: 100).
func (b *Batch) Append(record *Input) bool {
	hasSpace := len(b.records) < 100
	if hasSpace {
		b.records = append(b.records, record)
	}
	return hasSpace
}

func (b *Batch) attach(candidate Candidate) {
	i := candidate.InputIndex
	b.records[i].Results = append(b.records[i].Results, candidate)
}

// Length returns
func (b *Batch) Length() int {
	return len(b.records)
}

// Records returns the internal records collection.
func (b *Batch) Records() []*Input {
	return b.records
}

func (b *Batch) marshalJSON() ([]byte, error) {
	return json.Marshal(b.records)
}

func (b *Batch) marshalQueryString() string {
	record := b.records[0]
	query := make(url.Values)
	query.Set("input_id", record.InputID)
	query.Set("addressee", record.Addressee)
	query.Set("street", record.Street)
	query.Set("street2", record.Street2)
	query.Set("secondary", record.Secondary)
	query.Set("lastline", record.LastLine)
	query.Set("urbanization", record.Urbanization)
	query.Set("zipcode", record.ZIPCode)
	query.Set("candidates", strconv.Itoa(record.MaxCandidates))
	return query.Encode()
}

// Clear clears the internal collection.
func (b *Batch) Clear() {
	b.records = nil
}

// Reset clears the internal collection and resets settings.
func (b *Batch) Reset() {
	b.Clear()
	b.standardizeOnly = false
	b.includeInvalid = false
}