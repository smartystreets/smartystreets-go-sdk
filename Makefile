#!/usr/bin/make -f

VERSION_FILE := version.go
VERSION      := $(shell tagit -p --dryrun)

clean:
	@git checkout "$(VERSION_FILE)"

test: clean
	go test -short ./...

compile: clean
	go build ./...

cover: compile
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out

integrate: compile test
	@go run examples/international-street-api/main.go > /dev/null
	@go run examples/us-street-api/main.go > /dev/null
	@go run examples/us-autocomplete-api/main.go > /dev/null
	@go run examples/us-extract-api/main.go > /dev/null
	@go run examples/us-zipcode-api/main.go > /dev/null

version:
	printf 'package sdk\n\nconst VERSION = "%s"\n' "$(VERSION)" > "$(VERSION_FILE)"

publish: compile test version
		&& git commit -am "Incremented version to $(VERSION)" \
		&& tagit -p \
		&& git push origin master --tags

.PHONY: clean test compile cover integrate version package publish
