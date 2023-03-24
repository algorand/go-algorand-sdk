#!/bin/bash

export ALGOD_PORT="60000"
export INDEXER_PORT="59999"
export KMD_PORT="60001"

go run account/main.go
go run apps/main.go
go run asa/main.go
go run atc/main.go
go run atomic_transactions/main.go
go run codec/main.go
go run debug/main.go
go run indexer/main.go
go run kmd/main.go
go run lsig/main.go
go run participation/main.go
go run overview/main.go