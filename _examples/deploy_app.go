package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/types"
)

func main() {
	// example: ALGOD_CREATE_CLIENT
	// Create a new algod client, configured to connect to out local sandbox
	var algodAddress = "http://localhost:4001"
	var algodToken = strings.Repeat("a", 64)
	algodClient, err := algod.MakeClient(
		algodAddress,
		algodToken,
	)

	// Or, if necessary, pass alternate headers

	var algodHeader common.Header
	algodHeader.Key = "X-API-Key"
	algodHeader.Value = algodToken
	algodClientWithHeaders, err := algod.MakeClientWithHeaders(
		algodAddress,
		algodToken,
		[]*common.Header{&algodHeader},
	)
	// example: ALGOD_CREATE_CLIENT

	// Suppress `algodClientWithHeaders declared but not used`
	_ = algodClientWithHeaders

	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return
	}

	nodeStatus, err := algodClient.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("Failed to get status: %s\n", err)
		return
	}

	fmt.Printf("Last Round: %d\n", nodeStatus.LastRound)

	// example: GO_STATE_SCHEMA
	// declare application state storage (immutable)
	const localInts = 1
	const localBytes = 1
	const globalInts = 1
	const globalBytes = 0

	// define schema
	globalSchema := types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}
	localSchema := types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)}
	// example: GO_STATE_SCHEMA
}
