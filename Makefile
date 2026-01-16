#!/usr/bin/make -f

VERSION_FILE := version.go
VERSION      := $(shell bumpit -p `git describe`)

test: fmt clean
	go test -short -cover -count=1 ./...

fmt:
	go mod tidy && go fmt ./...

clean:
	@git checkout "$(VERSION_FILE)"

compile: clean
	go build ./...

build: test compile

cover: compile
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out

integrate: compile test
	@go run examples/international-autocomplete-api/main.go > /dev/null
	@go run examples/international-postal-code-api/main.go > /dev/null
	@go run examples/international-street-api/main.go > /dev/null
	@go run examples/us-autocomplete-pro-api/main.go > /dev/null
	@go run examples/us-enrichment-api/address-search/main.go > /dev/null
	@go run examples/us-enrichment-api/geo-reference/main.go > /dev/null
	@go run examples/us-enrichment-api/property-principal/main.go > /dev/null
	@go run examples/us-enrichment-api/risk/main.go > /dev/null
	@go run examples/us-enrichment-api/secondary/main.go > /dev/null
	@go run examples/us-enrichment-api/secondary-count/main.go > /dev/null
	@go run examples/us-enrichment-api/universal/main.go > /dev/null
	@go run examples/us-extract-api/main.go > /dev/null
	@go run examples/us-reverse-geo-api/main.go > /dev/null
	@go run examples/us-street-api/main.go > /dev/null
	@go run examples/us-zipcode-api/main.go > /dev/null

version:
	printf 'package sdk\n\nconst VERSION = "%s"\n' "$(VERSION)" > "$(VERSION_FILE)"

publish: compile test version
	git commit -am "Incremented version."; tagit -p; git push origin master --tags

.PHONY: test fmt clean compile build cover integrate version package publish
