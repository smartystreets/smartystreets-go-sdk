#!/usr/bin/make -f

test: build
	go test -v -short ./...

build:
	go build ./...

cover:
	go build
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
