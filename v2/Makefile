SRCPATH     := $(shell pwd)
TEST_SOURCES := $(shell cd $(SRCPATH) && go list ./...)
TEST_SOURCES_NO_CUCUMBER := $(shell cd $(SRCPATH) && go list ./... | grep -v test)

lint:
	golint `go list ./... | grep -v /vendor/`

generate:
	cd $(SRCPATH) && go generate ./logic

build: generate
	cd $(SRCPATH) && go test -run xxx_phony_test $(TEST_SOURCES)

unit:
	go test $(TEST_SOURCES_NO_CUCUMBER)
	cd test && go test --godog.strict=true --godog.format=pretty --godog.tags="@unit.offline,@unit.algod,@unit.indexer,@unit.rekey,@unit.tealsign,@unit.dryrun,@unit.responses,@unit.applications,@unit.transactions,@unit.indexer.rekey,@unit.responses.messagepack,@unit.responses.231,@unit.responses.messagepack.231,@unit.responses.genesis,@unit.feetest" --test.v .

integration:
	go test $(TEST_SOURCES_NO_CUCUMBER)
	cd test && go test --godog.strict=true --godog.format=pretty --godog.tags="@algod,@assets,@auction,@kmd,@send,@template,@indexer,@rekey,@dryrun,@compile,@applications.verified,@indexer.applications,@indexer.231" --test.v .

docker-test:
	./test/docker/run_docker.s
