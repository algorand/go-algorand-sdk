SRCPATH     := $(shell pwd)
TEST_SOURCES := $(shell cd $(SRCPATH) && go list ./...)
TEST_SOURCES_NO_CUCUMBER := $(shell cd $(SRCPATH) && go list ./... | grep -v test)
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
	cd test && go test -timeout 0s --godog.strict=true --godog.format=pretty --godog.tags="@unit.sourcemap,@unit.offline,@unit.algod,@unit.indexer,@unit.transactions.keyreg,@unit.rekey,@unit.tealsign,@unit.dryrun,@unit.responses,@unit.applications,@unit.transactions,@unit.indexer.rekey,@unit.responses.messagepack,@unit.responses.231,@unit.responses.messagepack.231,@unit.responses.genesis,@unit.feetest,@unit.indexer.logs,@unit.abijson,@unit.abijson.byname,@unit.transactions.payment,@unit.atomic_transaction_composer,@unit.responses.unlimited_assets,@unit.indexer.ledger_refactoring,@unit.algod.ledger_refactoring,@unit.dryrun.trace.application" --test.v .

integration:
	go test $(TEST_SOURCES_NO_CUCUMBER)
	cd test && go test -timeout 0s --godog.strict=true --godog.format=pretty --godog.tags="@algod,@assets,@auction,@kmd,@send,@indexer,@rekey_v1,@send.keyregtxn,@dryrun,@compile,@applications.verified,@indexer.applications,@indexer.231,@abi,@c2c,@compile.sourcemap" --test.v .

harness:
	./test-harness.sh

docker-gosdk-build:
	echo "Building docker image from base $(GO_IMAGE)"
	docker build -t go-sdk-testing --build-arg GO_IMAGE="$(GO_IMAGE)" -f test/docker/Dockerfile $(shell pwd)

docker-gosdk-run:
	docker ps -a
	docker run -it --network host go-sdk-testing:latest

docker-test: harness docker-gosdk-build docker-gosdk-run


.PHONY: test fmt
