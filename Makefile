SRCPATH := $(GOPATH)/src/github.com/algorand/go-algorand-sdk
TEST_SOURCES := $(shell cd $(SRCPATH) && go list ./...)

lint:
	golint `go list ./... | grep -v /vendor/`

build:
	cd $(SRCPATH) && go test -run xxx_phony_test $(TEST_SOURCES)

android:
	gomobile bind -target=android github.com/algorand/go-algorand-sdk/crypto
	gomobile bind -target=android github.com/algorand/go-algorand-sdk/auction
	gomobile bind -target=android github.com/algorand/go-algorand-sdk/transaction
	gomobile bind -target=android github.com/algorand/go-algorand-sdk/mnemonic
	gomobile bind -target=android github.com/algorand/go-algorand-sdk/encoding/msgpack

ios:
	gomobile bind -target=ios github.com/algorand/go-algorand-sdk/crypto
	gomobile bind -target=ios github.com/algorand/go-algorand-sdk/auction
	gomobile bind -target=ios github.com/algorand/go-algorand-sdk/transaction
	gomobile bind -target=ios github.com/algorand/go-algorand-sdk/mnemonic
	gomobile bind -target=ios github.com/algorand/go-algorand-sdk/encoding/msgpack
