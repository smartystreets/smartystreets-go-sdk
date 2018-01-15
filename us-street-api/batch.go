package street

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Batch stores input records and settings related to a group of addresses to be verified in a batch.
type Batch struct {
	lookups []*Lookup
}

// NewBatch creates a new, empty batch.
func NewBatch() *Batch {
	return &Batch{}
}

// Append includes the record in the collection to be sent if there is still room (max: 100).
func (b *Batch) Append(record *Lookup) bool {
	if b.IsFull() {
		return false
	}
	b.lookups = append(b.lookups, record)
	return true
}

func (b *Batch) attach(candidates []*Candidate) {
	for _, candidate := range candidates {
		if candidate != nil {
			if i := candidate.InputIndex; i < len(b.lookups) {
				lookup := b.lookups[i]
				lookup.Results = append(lookup.Results, candidate)
			}
		}
	}
}

// IsFull returns true when the batch has 100 lookups, false in every other case.
func (b *Batch) IsFull() bool {
	return b.Length() == MaxBatchSize
}

func (b *Batch) isEmpty() bool {
	return b.Length() == 0
}

// Length returns
func (b *Batch) Length() int {
	return len(b.lookups)
}

// Records returns the internal records collection.
func (b *Batch) Records() []*Lookup {
	return b.lookups
}

// Clear clears the internal collection.
func (b *Batch) Clear() {
	b.lookups = nil
}

func (b *Batch) buildRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, verifyURL, nil) // We control the method and the URL. This is safe.
	if b.Length() == 1 {
		b.serializeGET(request)
	} else {
		b.serializePOST(request)
	}
	return request
}
func (b *Batch) serializeGET(request *http.Request) {
	request.Method = http.MethodGet
	query := request.URL.Query()
	b.lookups[0].encodeQueryString(query)
	request.URL.RawQuery = query.Encode()
}
func (b *Batch) serializePOST(request *http.Request) {
	request.Method = http.MethodPost
	payload, _ := json.Marshal(b.lookups) // We control the types being serialized. This is safe.
	request.Body = ioutil.NopCloser(bytes.NewReader(payload))
	request.ContentLength = int64(len(payload))
	request.Header.Set("Content-Type", "application/json")
}

const verifyURL = "/street-address" // Remaining parts will be completed later by the sdk.BaseURLClient.
const MaxBatchSize = 100
