#!/bin/env/bin bash

export ALGOD_PORT="60000"
export INDEXER_PORT="59999"
export KMD_PORT="60001"

# Loop through each directory in the current working directory
for dir in */; do
  # Check if main.go exists in the directory
  if [ -f "${dir}main.go" ]; then
    # Run the "go run" command with the relative path
    go run "${dir}main.go"
  fi
done