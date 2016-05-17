package us_street

type Batch struct {
	records []*Input
}

func NewBatch() *Batch {
	return &Batch{}
}

func (this *Batch) Append(record *Input) bool {
	hasSpace := len(this.records) < 100
	if hasSpace {
		this.records = append(this.records, record)
	}
	return hasSpace
}

func (this *Batch) Length() int {
	return len(this.records)
}
