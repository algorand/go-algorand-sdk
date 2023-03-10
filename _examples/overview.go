package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
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
}

// Missing MULTISIG_SIGN in GOSDK examples (in ../docs/get-details/transactions/signatures.md:317)
// Missing TRANSACTION_FEE_OVERRIDE in GOSDK examples (in ../docs/get-details/transactions/index.md:786)
// Missing TRANSACTION_KEYREG_OFFLINE_CREATE in GOSDK examples (in ../docs/run-a-node/participate/offline.md:38)
// Missing CONST_MIN_FEE in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/guidelines.md:71)
// Missing SP_MIN_FEE in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/guidelines.md:106)
