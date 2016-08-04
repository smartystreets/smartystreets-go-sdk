package main

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/smartystreets/smartystreets-go-sdk/us-zipcode-api"
	"bitbucket.org/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime)

	client := wireup.NewClientBuilder().
		WithSecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")).
		BuildUSZIPCodeAPIClient()

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
