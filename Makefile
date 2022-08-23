SRCPATH     := $(shell pwd)
TEST_SOURCES := $(shell cd $(SRCPATH) && go list ./...)
TEST_SOURCES_NO_CUCUMBER := $(shell cd $(SRCPATH) && go list ./... | grep -v test)
UNIT_TAGS :=  "$(shell awk '{print $2}' test/unit.tags | paste -s -d, -)"
INTEGRATIONS_TAGS := "$(shell awk '{print $2}' test/integration.tags | paste -s -d, -)"
GO_IMAGE := golang:$(subst go,,$(shell go version | cut -d' ' -f 3 | cut -d'.' -f 1,2))-stretch

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
	cd test && go test -timeout 0s --godog.strict=true --godog.format=pretty --godog.tags=$(UNIT_TAGS) --test.v .

integration:
	go test $(TEST_SOURCES_NO_CUCUMBER)
	cd test && go test -timeout 0s --godog.strict=true --godog.format=pretty --godog.tags="@algod,@assets,@auction,@kmd,@send,@indexer,@rekey_v1,@send.keyregtxn,@dryrun,@compile,@applications.verified,@indexer.applications,@indexer.231,@abi,@c2c,@compile.sourcemap" --test.v .


.PHONY: test fmt
