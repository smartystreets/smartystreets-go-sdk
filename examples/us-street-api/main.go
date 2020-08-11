package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	// You don't have to store your keys in environment variables, but we recommend it.
	client := wireup.BuildUSStreetAPIClient(
		wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// wireup.WithLicense("us-core-cloud"),
		// wireup.ViaProxy("https://my-proxy.my-company.com"), // uncomment this line to point to the specified proxy.
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
		// ...or maybe you want to supply your own http client:
		// wireup.WithHTTPClient(&http.Client{Timeout: time.Second * 30})
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/us-street-api#input-fields

	lookup1 := &street.Lookup{
		InputID:       "24601", // Optional ID from your system
		Addressee:     "John Doe",
		Street:        "1 Rosedale",
		Street2:       "closet under the stairs",
		Secondary:     "APT 2",
		Urbanization:  "", // Only applies to Puerto Rico addresses
		City:          "Baltimore",
		State:         "MD",
		ZIPCode:       "21229",
		MaxCandidates: 3,
		MatchStrategy: street.MatchInvalid, // "invalid" is the most permissive match,
		// this will always return at least one result even if the address is invalid.
		// Refer to the documentation for additional MatchStrategy options.
	}
	lookup2 := &street.Lookup{
		Street:        "1600 Pennsylvania Avenue",
		LastLine:      "Washington, DC",
		MaxCandidates: 5,
	}
	lookup3 := &street.Lookup{
		InputID:       "8675309",
		Street:        "1600 Amphitheatre Parkway Mountain View, CA 94043",
		MaxCandidates: 1,
	}

	batch := street.NewBatch()
	batch.Append(lookup1)
	batch.Append(lookup2)
	batch.Append(lookup3)

	if err := client.SendBatchWithContext(context.Background(), batch); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	for i, input := range batch.Records() {
		fmt.Println("Results for input:", i)
		fmt.Println()
		for j, candidate := range input.Results {
			fmt.Println("  Candidate:", j)
			fmt.Println(" Input ID: ", candidate.InputID)
			fmt.Println(" ", candidate.DeliveryLine1)
			fmt.Println(" ", candidate.LastLine)
			fmt.Println()
		}
	}

	log.Println("OK")
}
