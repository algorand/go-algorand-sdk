SRCPATH     := $(shell pwd)
TEST_SOURCES := $(shell cd $(SRCPATH) && go list ./...)

lint:
	golint `go list ./... | grep -v /vendor/`

generate:
	cd $(SRCPATH) && go generate ./logic

build: generate
	cd $(SRCPATH) && go test -run xxx_phony_test $(TEST_SOURCES)
