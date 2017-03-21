package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

func main() {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.NewClientBuilder().
		WithSecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")).
		//WithDebugHTTPOutput(). // uncomment this line to see detailed HTTP request/response information.
		BuildUSExtractAPIClient()

	lookup := &extract.Lookup{Text: "Meet me at 3214 N University Ave Provo UT 84604 just after 3pm."}

	if err := client.SendLookup(lookup); err != nil {
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
