SRCPATH     := $(shell pwd)
TEST_SOURCES := $(shell cd $(SRCPATH) && go list ./...)
TEST_SOURCES_NO_CUCUMBER := $(shell cd $(SRCPATH) && go list ./... | grep -v test)
UNIT_TAGS :=  "$(shell awk '{print $2}' test/unit.tags | paste -s -d, -)"
INTEGRATIONS_TAGS := "$(shell awk '{print $2}' test/integration.tags | paste -s -d, -)"
GO_IMAGE := golang:$(subst go,,$(shell go version | cut -d' ' -f 3 | cut -d'.' -f 1,2))-stretch

lint:
	golangci-lint run -c .golangci.yml
	go vet ./...

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
	cd test && go test -timeout 0s --godog.strict=true --godog.format=pretty --godog.tags=$(INTEGRATIONS_TAGS) --test.v .

display-all-go-steps:
	find test -name "*.go" | xargs grep "github.com/cucumber/godog" 2>/dev/null | cut -d: -f1 | sort | uniq | xargs grep -Eo "Step[(].[^\`]+" | awk '{sub(/:Step\(./,":")} 1' | sed -E 's/", [a-zA-Z0-9]+\)//g'

harness:
	./test-harness.sh up

harness-down:
	./test-harness.sh down

docker-gosdk-build:
	echo "Building docker image from base $(GO_IMAGE)"
	docker build -t go-sdk-testing --build-arg GO_IMAGE="$(GO_IMAGE)" -f test/docker/Dockerfile $(shell pwd)

docker-gosdk-run:
	docker ps -a
	docker run -it --network host go-sdk-testing:latest

smoke-test-examples:
	cd "$(SRCPATH)/examples" && bash smoke_test.sh && cd -

docker-test: harness docker-gosdk-build docker-gosdk-run


.PHONY: test fmt
