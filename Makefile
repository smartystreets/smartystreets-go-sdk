#!/usr/bin/make -f

REPO_NAME := smartystreets-go-sdk
REPO_PATH := github.com/smartystreets/$(REPO_NAME)
FULL_PATH := $(GOPATH)/src/$(REPO_PATH)
VERSION_FILE = ./version.go

test:
	go test -short ./...

compile:
	go build ./...

dependencies: gopath
	go get github.com/smartystreets/gunit
	go get github.com/smartystreets/assertions
	go get github.com/smartystreets/clock
	go get github.com/smartystreets/logging

gopath:
	@mkdir -p "$(dir $(FULL_PATH))"
	@test -e "$(FULL_PATH)" || ln -sf "$(PWD)" "$(FULL_PATH)" # gopath compatibility

cover: compile
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

integrate: compile
	@go run examples/international-street-api/main.go > /dev/null
	@go run examples/us-street-api/main.go > /dev/null
	@go run examples/us-autocomplete-api/main.go > /dev/null
	@go run examples/us-extract-api/main.go > /dev/null
	@go run examples/us-zipcode-api/main.go > /dev/null

publish:
	$(eval VERSION := $(shell bumpit -p $(shell git describe) --dirty))
	@echo "package sdk\n\nconst VERSION = \"$(VERSION)\"" > "$(VERSION_FILE)"
	git add "$(VERSION_FILE)"
	git commit -m "Incremented version number to $(VERSION)"
	tagit -p
	git push origin master --tags

#########################################################33

container-test:
	docker-compose run sdk make test
container-compile:
	docker-compose run sdk make compile
container-integrate:
	docker-compose run sdk make integrate
