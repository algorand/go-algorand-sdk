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
	cd test && go test --godog.strict=true --godog.format=pretty --godog.tags="@unit.offline,@unit.algod,@unit.indexer,@unit.rekey,@unit.tealsign,@unit.dryrun,@unit.responses,@unit.applications" --test.v .

integration:
	go test $(TEST_SOURCES_NO_CUCUMBER)
	cd test && go test --godog.strict=true --godog.format=pretty --godog.tags="@algod,@assets,@auction,@kmd,@send,@template,@indexer,@rekey,@dryrun,@compile,@applications,@applications.verified" --test.v .

docker-test:
	./test/docker/run_docker.sh
