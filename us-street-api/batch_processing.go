package street

// SendLookups is a high-level convenience function that leverages an internal reusable Batch
// to send all lookups provided in serial, blocking fashion.
func (c *Client) SendLookups(lookups ...*Lookup) error {
	stream := make(chan *Lookup)
	go loadAll(stream, lookups)
	return c.SendFromChannel(stream, nil)
}

func loadAll(stream chan *Lookup, lookups []*Lookup) {
	load(stream, lookups)
	close(stream)
}

func load(stream chan *Lookup, lookups []*Lookup) {
	for _, lookup := range lookups {
		stream <- lookup
	}
}

// SendFromChannel is a high-level convenience function that leverages an
// internal reusable Batch to send everything received from the provided
// input channel in serial, blocking fashion. The same lookups are sent
// on the output channel as they are processed. Whereas it is the caller's
// job to close the input channel when all inputs have been sent, the output
// channel will be closed for the caller.
func (c *Client) SendFromChannel(input, output chan *Lookup) error {
	processor := newBatchProcessor(input, output, c)
	processor.ProcessAll()
	return processor.Error()
}

type batchProcessor struct {
	input  chan *Lookup
	output chan *Lookup
	client *Client
	batch  *Batch
	err    error
}

func newBatchProcessor(input, output chan *Lookup, client *Client) *batchProcessor {
	return &batchProcessor{
		input:  input,
		output: output,
		client: client,
		batch:  NewBatch(),
	}
}

func (c *batchProcessor) ProcessAll() {
	c.readLookups()
	c.sendLastBatch()
	c.cleanup()
}

func (c *batchProcessor) readLookups() {
	for lookup := range c.input {
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
	c.routeOutput()
	c.batch.Clear()
}

func (c *batchProcessor) routeOutput() {
	if c.output != nil {
		load(c.output, c.batch.Records())
	}
}

func (c *batchProcessor) sendLastBatch() {
	if !c.errorOccurred() && !c.batch.isEmpty() {
		c.sendBatch()
	}
}

func (c *batchProcessor) errorOccurred() bool {
	return c.err != nil
}

func (c *batchProcessor) cleanup() {
	if c.output != nil {
		close(c.output)
	}
}
func (c *batchProcessor) Error() error {
	return c.err
}
