package main

import (
	"fmt"
	"log"
	"os"

	us_enrichment "github.com/smartystreets/smartystreets-go-sdk/us-enrichment-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	client := wireup.BuildUSEnrichmentAPIClient(
		wireup.BasicAuthCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
	)

	smartyKey := "1962995076"

	first := &us_enrichment.Lookup{SmartyKey: smartyKey}
	err, results := client.SendBusinessSummary(first)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("First call: %d result(s), Etag=%s\n", len(results), first.ResponseETag)

	second := &us_enrichment.Lookup{SmartyKey: smartyKey, ETag: first.ResponseETag}
	err, results = client.SendBusinessSummary(second)
	if err != nil {
		log.Fatal(err)
	}
	if len(results) == 0 {
		fmt.Printf("Second call: not modified, Etag=%s\n", second.ResponseETag)
	} else {
		fmt.Printf("Second call: modified, %d result(s), Etag=%s\n", len(results), second.ResponseETag)
	}
}
