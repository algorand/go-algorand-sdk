package main

import (
	"context"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
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

	appID := deployApp(algodClient, acct1)

	// example: DEBUG_DRYRUN_DUMP
	var (
		args     [][]byte
		accounts []string
		apps     []uint64
		assets   []uint64
	)

	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get suggested params: %s", err)
	}

	appCallTxn, err := transaction.MakeApplicationNoOpTx(
		appID, args, accounts, apps, assets, sp, acct1.Address,
		nil, types.Digest{}, [32]byte{}, types.Address{},
	)
	if err != nil {
		log.Fatalf("Failed to create app call txn: %+v", err)
	}

	_, stxn, err := crypto.SignTransaction(acct1.PrivateKey, appCallTxn)
	if err != nil {
		log.Fatalf("Failed to sign app txn: %+v", err)
	}

	signedAppCallTxn := types.SignedTxn{}
	msgpack.Decode(stxn, &signedAppCallTxn)

	drr, err := transaction.CreateDryrun(algodClient, []types.SignedTxn{signedAppCallTxn}, nil, context.Background())
	if err != nil {
		log.Fatalf("Failed to create dryrun: %+v", err)
	}

	os.WriteFile("dryrun.msgp", msgpack.Encode(drr), 0666)
	// example: DEBUG_DRYRUN_DUMP

	// example: DEBUG_DRYRUN_SUBMIT
	// Create the dryrun request object
	drReq, err := transaction.CreateDryrun(algodClient, []types.SignedTxn{signedAppCallTxn}, nil, context.Background())
	if err != nil {
		log.Fatalf("Failed to create dryrun: %+v", err)
	}

	// Pass dryrun request to algod server
	dryrunResponse, err := algodClient.TealDryrun(drReq).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to dryrun request: %s", err)
	}

	// Inspect the response to check result
	for _, txn := range dryrunResponse.Txns {
		log.Printf("%+v", txn.AppCallTrace)
	}
	// example: DEBUG_DRYRUN_SUBMIT
	os.Remove("dryrun.msgp")
}
