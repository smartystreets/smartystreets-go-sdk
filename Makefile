#!/usr/bin/make -f

default: test

test: build
	go generate ./...
	go test -v -short ./...

build:
	go build ./...

cover:
	go build
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
