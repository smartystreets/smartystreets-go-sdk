package main

import (
	"context"
	"fmt"
	"log"
	"os"

	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSStreetAPIClient(
		wireup.BasicAuthCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		wireup.WithFeatureIANATimeZone(), // To add IANA timezone feature you need to specify when you
		// create the client.
	)

	lookup := &street.Lookup{
		InputID:       "24601",
		Addressee:     "John Doe",
		Street:        "1 Rosedale",
		Street2:       "closet under the stairs",
		Secondary:     "APT 2",
		City:          "Baltimore",
		State:         "MD",
		ZIPCode:       "21229",
		MatchStrategy: street.MatchEnhanced, // Enhanced matching is required to return IANA timezone results.
	}

	batch := street.NewBatch()
	batch.Append(lookup)

	if err := client.SendBatchWithContext(context.Background(), batch); err != nil {
		log.Fatal("Error sending batch:", err)
	}
	for i, input := range batch.Records() {
		fmt.Println("Results for input:", i)
		fmt.Println()
		for j, candidate := range input.Results {
			fmt.Println("\tCandidate:", j)
			fmt.Printf("\tInput ID: %s\n", candidate.InputID)
			fmt.Printf("\tDelivery Line 1:\n")
			fmt.Printf("\t\t%s\n", candidate.DeliveryLine1)
			fmt.Printf("\tLast Line:\n")
			fmt.Printf("\t\t%s\n", candidate.LastLine)
			fmt.Printf("\tTimezone: %s\n", candidate.Metadata.TimeZone)
			fmt.Printf("\tUTC Offset: %0.1f\n", candidate.Metadata.UTCOffset)
			fmt.Printf("\tDST: %t\n", candidate.Metadata.DST)
			fmt.Printf("\tIANA Timezone: %s\n", candidate.Metadata.IANATimeZone)
			fmt.Printf("\tIANA UTC Offset: %0.1f\n", candidate.Metadata.IANAUTCOffset)
			fmt.Printf("\tIANA DST: %t\n", candidate.Metadata.IANADST)
			fmt.Println()
		}
	}
	log.Println("OK")
}
