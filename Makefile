#!/usr/bin/make -f

VERSION_FILE := version.go
VERSION      := $(shell tagit -p --dryrun)

test: clean
	go test -short ./...

fmt:
	go fmt ./...

clean:
	@git checkout "$(VERSION_FILE)"

compile: clean
	go build ./...

build: test compile

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
	git commit -am "Incremented version."; tagit -p; git push origin master --tags

.PHONY: test fmt clean compile build cover integrate version package publish
