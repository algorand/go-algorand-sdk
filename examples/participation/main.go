package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

func main() {
	markOnline()
}

func setupConnection() (c *algod.Client, err error) {
	algod_token := strings.Repeat("a", 64)
	algod_server := "http://127.0.0.1:4001"
	algod_client, err := algod.MakeClient(algod_server, algod_token)
	c = (*algod.Client)(algod_client)
	return
}

func markOnline() {
	// setup connection
	algodClient, err := setupConnection()
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
		return
	}

	// get network suggested parameters
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
		return
	}

	// Mark Account as "Online" (participating)
	// example: TRANSACTION_KEYREG_ONLINE_CREATE
	fromAddr := "MWAPNXBDFFD2V5KWXAHWKBO7FO4JN36VR4CIBDKDDE7WAUAGZIXM3QPJW4"
	voteKey := "87iBW46PP4BpTDz6+IEGvxY6JqEaOtV0g+VWcJqoqtc="
	selKey := "1V2BE2lbFvS937H7pJebN0zxkqe1Nrv+aVHDTPbYRlw="
	sProofKey := "f0CYOA4yXovNBFMFX+1I/tYVBaAl7VN6e0Ki5yZA3H6jGqsU/LYHNaBkMQ/rN4M4F3UmNcpaTmbVbq+GgDsrhQ=="
	voteFirst := uint64(16532750)
	voteLast := uint64(19532750)
	keyDilution := uint64(1732)
	nonpart := false
	tx, err := transaction.MakeKeyRegTxnWithStateProofKey(
		fromAddr,
		[]byte{},
		sp,
		voteKey,
		selKey,
		sProofKey,
		voteFirst,
		voteLast,
		keyDilution,
		nonpart,
	)
	// example: TRANSACTION_KEYREG_ONLINE_CREATE
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}
	_ = tx

	// disabled example: TRANSACTION_KEYREG_OFFLINE_CREATE
	// disabled example: TRANSACTION_KEYREG_OFFLINE_CREATE
}
