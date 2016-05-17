package us_street

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type Batch struct {
	Error   error
	records []*Input
}

func NewBatch() *Batch {
	return &Batch{}
}

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

func (b *Batch) Length() int { return len(b.records) }

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

var emptyBatchError = errors.New("The batch was nil or had no records.")
