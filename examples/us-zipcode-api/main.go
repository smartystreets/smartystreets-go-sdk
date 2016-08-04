package main

import (
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime)

	client := wireup.NewClientBuilder().
		WithSecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")).
		BuildUSZIPCodeAPIClient()

	if err := client.Ping(); err != nil {
		fmt.Println("Ping failed:", err)
		os.Exit(1)
	} else {
		fmt.Println("Ping successful; service is reachable and responding.")
	}

	batch := us_zipcode.NewBatch()
	for batch.Append(&us_zipcode.Lookup{City: "PROVO", State: "UT", ZIPCode: "84604"}) {
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
}
