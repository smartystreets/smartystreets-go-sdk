package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/smartystreets/smartystreets-go-sdk/us-street-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"
)

const (
	fieldDelimiter    = '\t'
	reportingInterval = 10000
)

const (
	fieldsPerRecord = 4

	streetField  = 0
	cityField    = 1
	stateField   = 2
	zipCodeField = 3
)

func main() {
	inFile := openFile("input.txt")
	defer inFile.Close()

	bufferedReader := bufio.NewReader(inFile)
	reader := initializeReader(bufferedReader)

	outFile := createFile("output.txt")
	defer outFile.Close()

	bufferedWriter := bufio.NewWriter(outFile)
	defer bufferedWriter.Flush()

	writer := initializeWriter(outFile)
	defer writer.Flush()

	incoming := make(chan *street.Lookup, 1024*48)
	outgoing := make(chan *street.Lookup, 1024*48)
	client := wireup.NewClientBuilder().
		WithSecretKeyCredential(os.Getenv("SMARTY_AUTH_ID"), os.Getenv("SMARTY_AUTH_TOKEN")).
		BuildUSStreetAPIClient()

	// This processes all records provided to the incoming channel
	// in appropriately sized batches and puts the results on the
	// outgoing channel for the writer to stream to disk.
	go client.SendFromChannel(incoming, outgoing)
	go readRecords(reader, incoming)
	writeRecords(writer, outgoing)
}

func openFile(name string) io.ReadCloser {
	file, err := os.Open(name)
	if err != nil {
		log.Panic(err)
	}
	return file
}

func initializeReader(file io.Reader) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = fieldDelimiter
	reader.FieldsPerRecord = fieldsPerRecord
	return reader
}

func createFile(name string) io.WriteCloser {
	file, err := os.Create(name)
	if err != nil {
		log.Panic(err)
	}
	return file
}

func initializeWriter(file io.Writer) *csv.Writer {
	writer := csv.NewWriter(file)
	writer.Comma = fieldDelimiter
	return writer
}

func readRecords(reader *csv.Reader, pipe chan *street.Lookup) {
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}

		// Assuming each line in the input csv file has the following format:
		// Street\tCity\tState\tZIPCode\n
		pipe <- &street.Lookup{
			Street:        record[streetField],
			City:          record[cityField],
			State:         record[stateField],
			ZIPCode:       record[zipCodeField],
			MatchStrategy: street.MatchInvalid,
		}
	}

	close(pipe)
}

func writeRecords(writer *csv.Writer, pipe chan *street.Lookup) {
	sequence := 1

	for lookup := range pipe {
		writer.Write(toCSVRecord(sequence, lookup))
		sequence++

		if sequence%reportingInterval == 0 {
			log.Println("Processed:", sequence)
		}
	}

	log.Println("Finished.")
}

func toCSVRecord(sequence int, lookup *street.Lookup) []string {
	return []string{
		strconv.Itoa(sequence),
		lookup.Street,
		lookup.City,
		lookup.State,
		lookup.ZIPCode,
		"", // This (optional) blank column separates the input from output.
		lookup.Results[0].Addressee,
		lookup.Results[0].DeliveryLine1,
		lookup.Results[0].DeliveryLine2,
		lookup.Results[0].LastLine,
		lookup.Results[0].DeliveryPointBarcode,
		// Include additional fields here if desired.
	}
}
