package main

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/smartystreets/smartystreets-go-sdk/us-street-api"
)

func main() {
	log.SetFlags(log.Ltime)

	client := us_street.NewClientBuilder().
		WithSecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")).
		Build()

	batch := us_street.NewBatch()
	for batch.Append(&us_street.Lookup{Street: "3214 N University ave", LastLine: "Provo UT 84604"}) {
		fmt.Print(".")
	}
	fmt.Println("\nBatch full, preparing to send inputs:", batch.Length())

	if err := client.SendBatch(batch); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	for i, input := range batch.Records() {
		fmt.Println("Input:", i)
		for j, candidate := range input.Results {
			fmt.Println("Candidate:", j)
			fmt.Println(candidate.DeliveryLine1)
			fmt.Println(candidate.LastLine)
			fmt.Println()
		}
	}
}
