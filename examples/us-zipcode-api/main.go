package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSZIPCodeAPIClient(
		wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// For complete list of input fields, refer to:
	// https://smartystreets.com/docs/us-street-api#input-fields

	lookup1 := &zipcode.Lookup{
		InputID: "dfc33cb6-829e-4fea-aa1b-b6d6580f0817", // Optional ID from your system
		City:    "PROVO",
		State:   "UT",
		ZIPCode: "84604",
	}

	lookup2 := &zipcode.Lookup{
		InputID: "01189998819991197253",
		ZIPCode: "90210",
	}

	batch := zipcode.NewBatch()
	for batch.Append(lookup1) {
		fmt.Print(".")
	}
	for batch.Append(lookup2) {
		fmt.Print(".")
	}
	fmt.Println("\nBatch full, preparing to send inputs:", batch.Length())

	if err := client.SendBatch(batch); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	for i, input := range batch.Records() {
		fmt.Println("Input:", i, input.City, input.State, input.ZIPCode)
		fmt.Printf("%#v\n", input.Result)
		fmt.Println()
	}

	log.Println("OK")
}
