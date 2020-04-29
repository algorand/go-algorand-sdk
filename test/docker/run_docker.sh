#!/usr/bin/env bash

set -e

# reset test harness
rm -rf test-harness
rm -rf test/features
git clone --single-branch --branch develop https://github.com/algorand/algorand-sdk-testing.git test-harness
#copy feature files into project
mv test-harness/features test/features

#build test environment
docker build -t go-sdk-testing -f test/docker/Dockerfile "$(pwd)"

# Start test harness environment
./test-harness/scripts/up.sh

docker run -it \
     --network host \
     go-sdk-testing:latest
