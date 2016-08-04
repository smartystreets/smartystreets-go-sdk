package street

// SendLookups is a high-level convenience function that leverages an internal reusable Batch
// to send all lookups provided in serial, blocking fashion.
func (c *Client) SendLookups(lookups ...*Lookup) error {
	stream := make(chan *Lookup)
	go load(stream, lookups)
	return c.SendFromChannel(stream)
}

func load(stream chan *Lookup, lookups []*Lookup) {
	for _, lookup := range lookups {
		stream <- lookup
	}
	close(stream)
}

// SendFromChannel is a high-level convenience function that leverages an internal reusable Batch
// to send everything received from the provided lookups channel in serial, blocking fashion.
func (c *Client) SendFromChannel(lookups <-chan *Lookup) error {
	processor := newBatchProcessor(lookups, c)
	processor.ProcessAll()
	return processor.Error()
}

type batchProcessor struct {
	lookups <-chan *Lookup
	client  *Client
	batch   *Batch
	err     error
}

func newBatchProcessor(lookups <-chan *Lookup, client *Client) *batchProcessor {
	return &batchProcessor{
		lookups: lookups,
		client:  client,
		batch:   NewBatch(),
	}
}

func (c *batchProcessor) ProcessAll() {
	c.readLookups()
	c.sendLastBatch()
}

func (c *batchProcessor) readLookups() {
	for lookup := range c.lookups {
		c.processLookup(lookup)
		if c.errorOccurred() {
			break
		}
	}
}

func (c *batchProcessor) processLookup(lookup *Lookup) {
	c.batch.Append(lookup)
	if c.batch.IsFull() {
		c.sendBatch()
	}
}

func (c *batchProcessor) sendBatch() {
	c.err = c.client.SendBatch(c.batch)
	c.batch.Clear()
}

func (c *batchProcessor) sendLastBatch() {
	if c.errorOccurred() || c.batch.isEmpty() {
		return
	}
	c.sendBatch()
}

func (c *batchProcessor) errorOccurred() bool {
	return c.err != nil
}

func (c *batchProcessor) Error() error {
	return c.err
}
