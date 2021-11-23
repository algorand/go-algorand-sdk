SRCPATH     := $(shell pwd)
TEST_SOURCES := $(shell cd $(SRCPATH) && go list ./...)
TEST_SOURCES_NO_CUCUMBER := $(shell cd $(SRCPATH) && go list ./... | grep -v test)

lint:
	golint `go list ./... | grep -v /vendor/`

fmt:
	go fmt ./...

generate:
	cd $(SRCPATH) && go generate ./logic

build: generate
	cd $(SRCPATH) && go test -run xxx_phony_test $(TEST_SOURCES)

test:
	go test $(TEST_SOURCES_NO_CUCUMBER)

unit:
	go test $(TEST_SOURCES_NO_CUCUMBER)
	cd test && go test --godog.strict=true --godog.format=pretty --godog.tags="@unit.transactions.payment" --test.v .

docker-test:
	./test/docker/run_docker.sh

.PHONY: test fmt
