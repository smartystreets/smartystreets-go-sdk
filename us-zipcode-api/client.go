package zipcode

import (
	"context"
	"encoding/json"

	"github.com/smartystreets/smartystreets-go-sdk"
)

// Client is responsible for sending batches of addresses to the us-zipcode-api.
type Client struct {
	sender sdk.RequestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

// SendBatch sends the batch of inputs, populating the output for each input if the batch was successful.
func (c *Client) SendBatch(batch *Batch) error {
	return c.SendBatchWithContextAndAuth(context.Background(), batch, "", "")
}

func (c *Client) SendBatchWithContext(ctx context.Context, batch *Batch) error {
	return c.SendBatchWithContextAndAuth(ctx, batch, "", "")
}

func (c *Client) SendBatchWithContextAndAuth(ctx context.Context, batch *Batch, authID, authToken string) error {
	if batch == nil || batch.Length() == 0 {
		return nil
	}
	request := batch.buildRequest()
	request = request.WithContext(ctx)
	if len(authID) > 0 && len(authToken) > 0 {
		request.SetBasicAuth(authID, authToken)
	}

	response, err := c.sender.Send(request)
	if err != nil {
		return err
	}
	return deserializeResponse(response, batch)
}

func deserializeResponse(response []byte, batch *Batch) error {
	var results []*Result
	err := json.Unmarshal(response, &results)
	if err == nil {
		batch.attach(results)
	}
	return err
}
