package main

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/examples"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

func main() {
	markOnline()
}

func markOnline() {
	// setup connection
	algodClient := examples.GetAlgodClient()

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
