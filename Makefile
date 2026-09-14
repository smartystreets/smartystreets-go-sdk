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

international-autocomplete-api:
	go run ./examples/international-autocomplete-api

international-postal-code-api:
	go run ./examples/international-postal-code-api

international-street-api:
	go run ./examples/international-street-api

us-autocomplete-api:
	go run ./examples/us-autocomplete-api

us-autocomplete-pro-api:
	go run ./examples/us-autocomplete-pro-api

us-enrichment-api: us-enrichment-api-address-search us-enrichment-api-business us-enrichment-api-business-name-search us-enrichment-api-etag us-enrichment-api-geo-reference us-enrichment-api-property-principal us-enrichment-api-secondary us-enrichment-api-secondary-count us-enrichment-api-universal

us-enrichment-api-address-search:
	go run ./examples/us-enrichment-api/address-search

us-enrichment-api-business:
	go run ./examples/us-enrichment-api/business

us-enrichment-api-business-name-search:
	go run ./examples/us-enrichment-api/business-name-search

us-enrichment-api-etag:
	go run ./examples/us-enrichment-api/etag

us-enrichment-api-geo-reference:
	go run ./examples/us-enrichment-api/geo-reference

us-enrichment-api-property-principal:
	go run ./examples/us-enrichment-api/property-principal

us-enrichment-api-secondary:
	go run ./examples/us-enrichment-api/secondary

us-enrichment-api-secondary-count:
	go run ./examples/us-enrichment-api/secondary-count

us-enrichment-api-universal:
	go run ./examples/us-enrichment-api/universal

us-extract-api:
	go run ./examples/us-extract-api

us-reverse-geo-api:
	go run ./examples/us-reverse-geo-api

us-street-api:
	go run ./examples/us-street-api

us-street-api-component-analysis:
	go run ./examples/us-street-api/component-analysis

us-street-api-iana-timezone:
	go run ./examples/us-street-api/iana-timezone

us-street-api-list:
	go run ./examples/us-street-api/list

us-street-api-match-strategy:
	go run ./examples/us-street-api/match-strategy

us-zipcode-api:
	go run ./examples/us-zipcode-api

examples: international-autocomplete-api international-postal-code-api international-street-api us-autocomplete-api us-autocomplete-pro-api us-enrichment-api us-extract-api us-reverse-geo-api us-street-api us-street-api-component-analysis us-street-api-iana-timezone us-street-api-list us-street-api-match-strategy us-zipcode-api

integrate: compile test examples

version:
	printf 'package sdk\n\nconst VERSION = "%s"\n' "$(VERSION)" > "$(VERSION_FILE)"

publish: compile test version
	git commit -am "Incremented version."; tagit -p; git push origin master --tags

.PHONY: test fmt clean compile build cover integrate version package publish examples \
	international-autocomplete-api international-postal-code-api international-street-api \
	us-autocomplete-api us-autocomplete-pro-api us-extract-api us-reverse-geo-api us-zipcode-api \
	us-enrichment-api us-enrichment-api-address-search us-enrichment-api-business \
	us-enrichment-api-business-name-search us-enrichment-api-etag us-enrichment-api-geo-reference \
	us-enrichment-api-property-principal us-enrichment-api-secondary us-enrichment-api-secondary-count \
	us-enrichment-api-universal us-street-api us-street-api-component-analysis \
	us-street-api-iana-timezone us-street-api-list us-street-api-match-strategy
