#!/usr/bin/env bash

set -e

rm -rf temp
rm -rf test/features
git clone --single-branch --branch michelle/test https://github.com/algorand/algorand-sdk-testing.git temp

cp test/docker/sdk.py temp/docker
mv temp/features test/features

docker build -t sdk-testing -f test/docker/Dockerfile "$(pwd)"

docker run -it \
     -v "$(pwd)":/opt/go/src/github.com/algorand/go-algorand-sdk \
     sdk-testing:latest 
