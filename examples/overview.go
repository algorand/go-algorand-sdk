package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

func main() {

	algodClient := getAlgodClient()

	nodeStatus, err := algodClient.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("Failed to get status: %s\n", err)
		return
	}

	fmt.Printf("Last Round: %d\n", nodeStatus.LastRound)

	// example: SP_MIN_FEE
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Printf("failed to %s", err)
	}
	// example: SP_MIN_FEE

	// example: CONST_MIN_FEE
	log.Printf("Min fee const: %d", transaction.MinTxnFee)
	// example: CONST_MIN_FEE

	// example: TRANSACTION_FEE_OVERRIDE
	// by using fee pooling and setting our fee to 2x min tx fee
	// we can cover the fee for another transaction in the group
	sp.Fee = 2 * transaction.MinTxnFee
	sp.FlatFee = true
	// ...
	// example: TRANSACTION_FEE_OVERRIDE

}
func exampleAlgod() {
	// example: ALGOD_CREATE_CLIENT
	// Create a new algod client, configured to connect to out local sandbox
	var algodAddress = "http://localhost:4001"
	var algodToken = strings.Repeat("a", 64)
	algodClient, _ := algod.MakeClient(
		algodAddress,
		algodToken,
	)

	// Or, if necessary, pass alternate headers

	var algodHeader common.Header
	algodHeader.Key = "X-API-Key"
	algodHeader.Value = algodToken
	algodClientWithHeaders, _ := algod.MakeClientWithHeaders(
		algodAddress,
		algodToken,
		[]*common.Header{&algodHeader},
	)
	// example: ALGOD_CREATE_CLIENT

	_ = algodClientWithHeaders
	_ = algodClient
}
