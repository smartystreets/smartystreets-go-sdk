package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSExtractAPIClient(
		wireup.WebsiteKeyCredential(os.Getenv("SMARTY_AUTH_WEB"), os.Getenv("SMARTY_AUTH_REFERER")),
		//wireup.SecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/cloud/us-extract-api#http-request-input-fields

	lookup := &extract.Lookup{
		Text: "Meet me at 3214 N University Ave Provo UT 84604 just after 3pm. " +
			"Also, here's a non-postal that will show up with enhanced match! 808 County Road 408 Brady, Tx. " +
			"is a beautiful place!",
		Aggressive:              true,
		AddressesWithLineBreaks: false,
		AddressesPerLine:        2,
		MatchStrategy:           street.MatchEnhanced,
	}

	if err := client.SendLookupWithContext(context.Background(), lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	fmt.Println(DumpJSON(lookup))

	log.Println("OK")
}

func DumpJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}

	var indent bytes.Buffer
	err = json.Indent(&indent, b, "", "  ")
	if err != nil {
		return err.Error()
	}
	return indent.String()
}
