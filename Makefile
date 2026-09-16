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

# Every directory under examples/ that holds a main.go is an example program. Its target name is the
# path under examples/ with slashes replaced by hyphens (examples/us-street-api/list -> us-street-api-list).
# Each example runs from its own directory so that relative paths such as input.txt resolve.
EXAMPLE_DIRS    := $(patsubst %/main.go,%,$(wildcard examples/*/main.go examples/*/*/main.go))
example-target   = $(subst /,-,$(patsubst examples/%,%,$(1)))
EXAMPLE_TARGETS := $(foreach dir,$(EXAMPLE_DIRS),$(call example-target,$(dir)))

define EXAMPLE_RULE
$(call example-target,$(1)):
	cd $(1) && go run .
endef
$(foreach dir,$(EXAMPLE_DIRS),$(eval $(call EXAMPLE_RULE,$(dir))))

# The enrichment API has no top-level example, only sub-examples; this aggregate runs them all.
us-enrichment-api: $(filter us-enrichment-api-%,$(EXAMPLE_TARGETS))

examples: $(EXAMPLE_TARGETS)

integrate: compile test examples

version:
	printf 'package sdk\n\nconst VERSION = "%s"\n' "$(VERSION)" > "$(VERSION_FILE)"

publish: compile test version
	git commit -am "Incremented version."; tagit -p; git push origin master --tags

.PHONY: test fmt clean compile build cover integrate version package publish examples us-enrichment-api $(EXAMPLE_TARGETS)
