package us_street

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
	c.readLookupsInBatches()
	c.sendLastBatch()
}

func (c *batchProcessor) readLookupsInBatches() {
	for lookup := range c.lookups {
		c.processLookup(lookup)
		if c.haltedPrematurely() {
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
	if c.haltedPrematurely() || c.batch.isEmpty() {
		return
	}
	c.err = c.client.SendBatch(c.batch)
}

func (c *batchProcessor) haltedPrematurely() bool {
	return c.err != nil
}

func (c *batchProcessor) Error() error {
	return c.err
}

func load(stream chan *Lookup, lookups []*Lookup) {
	for _, lookup := range lookups {
		stream <- lookup
	}
	close(stream)
}
