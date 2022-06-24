#!/usr/bin/env bash

set -e

# reset test harness
rm -rf test-harness
rm -rf test/features
git clone --single-branch --branch sourcemap https://github.com/algorand/algorand-sdk-testing.git test-harness
#copy feature files into project
mv test-harness/features test/features

exit
GO_VERSION=$(go version | cut -d' ' -f 3 | cut -d'.' -f 1,2)
GO_IMAGE=golang:${GO_VERSION:2}-stretch

echo "Building docker image from base \"$GO_IMAGE\""

#build test environment
docker build -t go-sdk-testing --build-arg GO_IMAGE="$GO_IMAGE" -f test/docker/Dockerfile "$(pwd)"

# Start test harness environment
./test-harness/scripts/up.sh -p

docker run -it \
     --network host \
     go-sdk-testing:latest
