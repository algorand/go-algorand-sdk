package main

import (
	"context"
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func main() {
	algodClient := getAlgodClient()
	accts, err := getSandboxAccounts()
	if err != nil {
		log.Fatalf("failed to get sandbox accounts: %s", err)
	}

	acct1 := accts[0]
	acct2 := accts[1]

	// example: ATOMIC_CREATE_TXNS
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get suggested params: %s", err)
	}

	tx1, err := transaction.MakePaymentTxn(acct1.Address.String(), acct2.Address.String(), 100000, nil, "", sp)
	if err != nil {
		log.Fatalf("failed creating transaction: %s", err)
	}

	// from account 2 to account 1
	tx2, err := transaction.MakePaymentTxn(acct2.Address.String(), acct1.Address.String(), 100000, nil, "", sp)
	if err != nil {
		log.Fatalf("failed creating transaction: %s", err)
	}
	// example: ATOMIC_CREATE_TXNS

	// example: ATOMIC_GROUP_TXNS
	// compute group id and put it into each transaction
	gid, err := crypto.ComputeGroupID([]types.Transaction{tx1, tx2})
	tx1.Group = gid
	tx2.Group = gid
	// example: ATOMIC_GROUP_TXNS

	// example: ATOMIC_GROUP_SIGN
	_, stx1, err := crypto.SignTransaction(acct1.PrivateKey, tx1)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}
	_, stx2, err := crypto.SignTransaction(acct2.PrivateKey, tx2)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
	}
	// example: ATOMIC_GROUP_SIGN

	// example: ATOMIC_GROUP_ASSEMBLE
	var signedGroup []byte
	signedGroup = append(signedGroup, stx1...)
	signedGroup = append(signedGroup, stx2...)

	// example: ATOMIC_GROUP_ASSEMBLE

	// example: ATOMIC_GROUP_SEND
	pendingTxID, err := algodClient.SendRawTransaction(signedGroup).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, pendingTxID, 4, context.Background())
	if err != nil {
		fmt.Printf("Error waiting for confirmation on txID: %s\n", pendingTxID)
		return
	}
	fmt.Printf("Confirmed Transaction: %s in Round %d\n", pendingTxID, confirmedTxn.ConfirmedRound)
	// example: ATOMIC_GROUP_SEND

}
