package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

type addressCase struct {
	label, street, city, state, zip string
}

type lookupCase struct {
	label, address string
	strategy       street.MatchStrategy
}

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	// You don't have to store your keys in environment variables, but we recommend it.
	client := wireup.BuildUSStreetAPIClient(
		wireup.BasicAuthCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
	)

	// Each address is run through all three match strategies so you can compare how
	// "strict", "enhanced", and "invalid" each handle a valid, an invalid, and an
	// ambiguous address.
	//   - strict:   only returns candidates that are valid, mailable addresses.
	//   - enhanced: returns a more comprehensive dataset (requires a US Core or Rooftop license).
	//   - invalid:  most permissive; always returns at least one candidate (a best-guess standardization).
	// Documentation for input fields: https://smartystreets.com/docs/us-street-api#input-fields
	addresses := []addressCase{
		{"valid (real, deliverable)", "1600 Amphitheatre Pkwy", "Mountain View", "CA", "94043"},
		{"invalid (no such address)", "9999 W 1150 S", "Provo", "UT", "84601"},
		{"ambiguous (missing ZIP/unit)", "1 Rosedale St", "Baltimore", "MD", ""},
	}
	strategies := []street.MatchStrategy{street.MatchStrict, street.MatchEnhanced, street.MatchInvalid}

	batch := street.NewBatch()
	var cases []lookupCase // parallel metadata for each lookup, in the order they are appended

	for _, a := range addresses {
		for _, strategy := range strategies {
			batch.Append(&street.Lookup{
				Street:        a.street,
				City:          a.city,
				State:         a.state,
				ZIPCode:       a.zip,
				MatchStrategy: strategy,
				MaxCandidates: 10, // allow ambiguous addresses to return more than one match
			})
			cases = append(cases, lookupCase{a.label, fmt.Sprintf("%s, %s, %s", a.street, a.city, a.state), strategy})
		}
	}

	if err := client.SendBatchWithContext(context.Background(), batch); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	separator := strings.Repeat("=", 70)
	lastAddress := ""

	for i, input := range batch.Records() {
		c := cases[i]

		if c.address != lastAddress {
			fmt.Printf("\n%s\n Address: %s  [%s]\n%s\n", separator, c.address, c.label, separator)
			lastAddress = c.address
		}

		fmt.Printf("\n--- '%s' strategy ---\n", c.strategy)

		if len(input.Results) == 0 {
			fmt.Println("  0 candidates - no match returned under this strategy.")
			continue
		}

		fmt.Printf("  %d candidate(s):\n", len(input.Results))
		for _, candidate := range input.Results {
			fmt.Printf("    [%d] %s  %s\n", candidate.CandidateIndex, candidate.DeliveryLine1, candidate.LastLine)
		}
	}
}
