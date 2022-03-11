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
	#cd test && go test -timeout 0s --godog.strict=true --godog.format=pretty --godog.tags="@unit.offline,@unit.algod,@unit.indexer,@unit.transactions.keyreg,@unit.rekey,@unit.tealsign,@unit.dryrun,@unit.responses,@unit.applications,@unit.transactions,@unit.indexer.rekey,@unit.responses.messagepack,@unit.responses.231,@unit.responses.messagepack.231,@unit.responses.genesis,@unit.feetest,@unit.indexer.logs,@unit.abijson,@unit.transactions.payment,@unit.atomic_transaction_composer,@unit.dryrun.trace.application" --test.v .
	cd test && go test -timeout 0s --godog.strict=true --godog.format=pretty --godog.tags="@unit.dryrun.trace.application" --test.v .

integration:
	go test $(TEST_SOURCES_NO_CUCUMBER)
	cd test && go test -timeout 0s --godog.strict=true --godog.format=pretty --godog.tags="@algod,@assets,@auction,@kmd,@send,@indexer,@rekey,@send.keyregtxn,@dryrun,@compile,@applications.verified,@indexer.applications,@indexer.231,@abi,@c2c" --test.v .

docker-test:
	./test/docker/run_docker.sh

.PHONY: test fmt