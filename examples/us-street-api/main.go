package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.NewClientBuilder().
		WithSecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")).
		//WithDebugHTTPOutput(). // uncomment this line to see detailed HTTP request/response information.
		BuildUSStreetAPIClient()

	lookup1 := &street.Lookup{
		Street:        "1 Rosedale",
		City:          "Baltimore",
		State:         "MD",
		MaxCandidates: 10, // This input produces more than one candidate!
	}
	lookup2 := &street.Lookup{
		Street: "1600 Pennsylvania Avenue",
		City:   "Washington",
		State:  "DC",
	}
	lookup3 := &street.Lookup{
		Street: "1600 Amphitheatre Parkway Mountain View, CA 94043",
	}

	batch := street.NewBatch()
	batch.Append(lookup1)
	batch.Append(lookup2)
	batch.Append(lookup3)

	if err := client.SendBatch(batch); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	for i, input := range batch.Records() {
		fmt.Println("Results for input:", i)
		fmt.Println()
		for j, candidate := range input.Results {
			fmt.Println("  Candidate:", j)
			fmt.Println(" ", candidate.DeliveryLine1)
			fmt.Println(" ", candidate.LastLine)
			fmt.Println()
		}
	}

	log.Println("OK")
}
