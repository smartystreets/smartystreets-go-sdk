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

package: compile test
	echo "package sdk\n\nconst VERSION = \"$(VERSION)\"" > "$(VERSION_FILE)"

##########################################################

workspace:
	docker-compose run sdk /bin/sh

release:
	docker-compose run sdk make package \
		&& git commit -am "Incremented version to $(VERSION)" \
		&& tagit -p \
		&& git push origin master --tags

.PHONY: clean test compile cover integrate package workspace release
