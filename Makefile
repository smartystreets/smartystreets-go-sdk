#!/usr/bin/make -f

VERSION_FILE := version.go
VERSION      := $(shell tagit -p --dryrun)

clean:
	git checkout "$(VERSION_FILE)"

test: clean
	go test -short ./...

compile:
	go build ./...

cover: compile
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out

integrate: test compile
	find examples/ -type f -name "main.go" -exec go run "{}" > /dev/null \;

package: clean test compile
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
