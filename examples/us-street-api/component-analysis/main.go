package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSStreetAPIClient(
		wireup.HeaderCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		wireup.WithFeatureComponentAnalysis(), // To add component analysis feature you need to specify when you
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
		MatchStrategy: street.MatchEnhanced, // Enhanced matching is required to return component analysis results.
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

			componentAnalysis := &candidate.Analysis.Components // Component Analysis results are found in the Analysis
			// object.
			jsonData, err := json.MarshalIndent(componentAnalysis, "\t", "\t")
			if err != nil {
				log.Fatal("Error marshaling candidate:", err)
			}
			fmt.Println("\tComponent Analysis Results:")
			fmt.Println("\t", string(jsonData))
			fmt.Println()
		}
	}
	log.Println("OK")
}
