package runners

import (
	"cmd/main/contracts"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	street "github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

type USStreetRunner struct {
	client *street.Client
}

func (this *USStreetRunner) Help(writer io.Writer) {

}

func (this *USStreetRunner) Print(writer io.Writer, requests []*street.Lookup) {
	for i, request := range requests {
		fmt.Fprintln(writer, "Results for input:", i)
		fmt.Fprintln(writer)
		for j, candidate := range request.Results {
			fmt.Fprintln(writer, "  Candidate:", j)
			fmt.Fprintln(writer, " Input ID: ", candidate.InputID)
			fmt.Fprintln(writer, " ", candidate.DeliveryLine1)
			fmt.Fprintln(writer, " ", candidate.LastLine)
			fmt.Fprintln(writer)
		}
	}
}

func (this *USStreetRunner) Run(requests []*street.Lookup) (err error) {
	batch := street.NewBatch()
	for _, request := range requests {
		batch.Append(request)
	}

	if err := this.client.SendBatchWithContext(context.Background(), batch); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	return nil
}

func (this *USStreetRunner) Setup(settings contracts.Settings) (err error) {
	var secretStrategy wireup.Option
	if settings.AuthIDVariable != "" && settings.AuthTokenVariable != "" {
		secretStrategy = wireup.SecretKeyCredential(
			os.Getenv(settings.AuthIDVariable),
			os.Getenv(settings.AuthTokenVariable),
		)
	} else if settings.AuthWebVariable != "" && settings.AuthRefererVariable != "" {
		secretStrategy = wireup.WebsiteKeyCredential(
			os.Getenv(settings.AuthWebVariable),
			os.Getenv(settings.AuthRefererVariable),
		)
	} else {
		return contracts.ErrNoAuthorization
	}
	
	this.client = wireup.BuildUSStreetAPIClient(
		secretStrategy,
		wireup.CustomBaseURL(settings.CustomBaseURL),
		wireup.DebugHTTPOutput(),
	)

	return nil
}
