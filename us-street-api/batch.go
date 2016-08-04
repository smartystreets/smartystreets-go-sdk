package street

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
		i := candidate.InputIndex
		b.lookups[i].Results = append(b.lookups[i].Results, candidate)
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

const MaxBatchSize = 100
