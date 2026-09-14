package international_postal_code

import (
	"context"
	"errors"
	"net/http"

	"github.com/smartystreets/smartystreets-go-sdk"
	"github.com/smartystreets/smartystreets-go-sdk/internal/json"
)

// Client is responsible for sending batches of addresses to the international-street-api.
type Client struct {
	sender sdk.RequestSender
}

// NewClient creates a client with the provided sender.
func NewClient(sender sdk.RequestSender) *Client {
	return &Client{sender: sender}
}

// SendLookup sends the lookup, populating the output if the request was successful.
func (c *Client) SendLookup(lookup *Lookup) error {
	return c.SendLookupWithContext(context.Background(), lookup)
}

func (c *Client) SendLookupWithContext(ctx context.Context, lookup *Lookup) error {
	if lookup == nil {
		return errors.New("lookup cannot be nil")
	}

	request, err := buildRequest(ctx, lookup)
	if err != nil {
		return err
	}
	response, err := c.sender.Send(request)
	if err != nil {
		return err
	}
	return deserializeResponse(response, lookup)
}

// SendLookupWithContextAndAuth sends a lookup with the provided context and per-request credentials.
// If credential is non-nil, it will be used to sign this request instead of the client-level credentials.
// This is useful for multi-tenant scenarios where different requests require different credentials.
func (c *Client) SendLookupWithContextAndAuth(ctx context.Context, lookup *Lookup, credential sdk.Credential) error {
	if lookup == nil {
		return errors.New("lookup cannot be nil")
	}

	request, err := buildRequest(ctx, lookup)
	if err != nil {
		return err
	}
	if credential != nil {
		if err := credential.Sign(request); err != nil {
			return err
		}
	}
	response, err := c.sender.Send(request)
	if err != nil {
		return err
	}
	return deserializeResponse(response, lookup)
}

func deserializeResponse(response []byte, lookup *Lookup) error {
	var candidates []*Candidate
	err := json.Unmarshal(response, &candidates)
	if err != nil {
		return err
	}
	lookup.Results = candidates
	return nil
}

func buildRequest(ctx context.Context, lookup *Lookup) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", lookupUrl, nil)
	if err != nil {
		return nil, err
	}
	query := request.URL.Query()
	lookup.populate(query)
	request.URL.RawQuery = query.Encode()
	return request, nil
}

const lookupUrl = "/lookup" // Remaining parts will be completed later by the sdk.BaseURLClient.
