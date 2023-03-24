#!/bin/bash

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

#go run account/main.go
#go run apps/main.go
#go run asa/main.go
#go run atc/main.go
#go run atomic_transactions/main.go
#go run codec/main.go
#go run debug/main.go
#go run indexer/main.go
#go run kmd/main.go
#go run lsig/main.go
#go run participation/main.go
#go run overview/main.go