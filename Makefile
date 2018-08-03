#!/usr/bin/make -f

REPO_NAME := smartystreets-go-sdk
REPO_PATH := github.com/smartystreets/$(REPO_NAME)
FULL_PATH := $(GOPATH)/src/$(REPO_PATH)

test: compile
	go test -short ./...

compile:
	go build ./...

dependencies: gopath
	go get github.com/smartystreets/gunit
	go get github.com/smartystreets/assertions
	go get github.com/smartystreets/clock
	go get github.com/smartystreets/logging

gopath:
	@mkdir -p $(dir $(FULL_PATH))
	@test -e $(FULL_PATH) || ln -sf $(PWD) $(FULL_PATH) # gopath compatibility

cover: compile
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

integrate: compile
	@go run examples/international-street-api/main.go > /dev/null
	@go run examples/us-street-api/main.go > /dev/null
	@go run examples/us-autocomplete-api/main.go > /dev/null
	@go run examples/us-extract-api/main.go > /dev/null
	@go run examples/us-zipcode-api/main.go > /dev/null

#########################################################33

container-test:
	docker-compose run sdk make test
container-compile:
	docker-compose run sdk make compile
container-integrate:
	docker-compose run sdk make integrate
