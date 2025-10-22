package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/v2/abi"
	"github.com/algorand/go-algorand-sdk/v2/examples"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func main() {
	algodClient := examples.GetAlgodClient()
	accts, err := examples.GetSandboxAccounts()
	if err != nil {
		log.Fatalf("failed to get sandbox accounts: %s", err)
	}

	acct1 := accts[0]

	appID := examples.DeployApp(algodClient, acct1)
	log.Printf("%d", appID)

	// example: ATC_CONTRACT_INIT
	b, err := os.ReadFile("calculator/contract.json")
	if err != nil {
		log.Fatalf("failed to read contract file: %s", err)
	}

	contract := &abi.Contract{}
	if err := json.Unmarshal(b, contract); err != nil {
		log.Fatalf("failed to unmarshal contract: %s", err)
	}
	// example: ATC_CONTRACT_INIT

	// example: ATC_CREATE
	// Create the atc we'll use to compose our transaction group
	var atc = transaction.AtomicTransactionComposer{}
	// example: ATC_CREATE

	// example: ATC_ADD_TRANSACTION
	// Get suggested params and make a transaction as usual
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	txn, err := transaction.MakePaymentTxn(acct1.Address.String(), acct1.Address.String(), 10000, nil, "", sp)
	if err != nil {
		log.Fatalf("failed to make transaction: %s", err)
	}

	// Construct a TransactionWithSigner and pass it to the atc
	signer := transaction.BasicAccountTransactionSigner{Account: acct1}
	atc.AddTransaction(transaction.TransactionWithSigner{Txn: txn, Signer: signer})
	// example: ATC_ADD_TRANSACTION

	// example: ATC_ADD_METHOD_CALL
	// Grab the method from out contract object
	addMethod, err := contract.GetMethodByName("add")
	if err != nil {
		log.Fatalf("failed to get add method: %s", err)
	}

	// Set up method call params
	mcp := transaction.AddMethodCallParams{
		AppID:           appID,
		Sender:          acct1.Address,
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
		Method:          addMethod,
		MethodArgs:      []interface{}{1, 1},
	}
	if err := atc.AddMethodCall(mcp); err != nil {
		log.Fatalf("failed to add method call: %s", err)
	}
	// example: ATC_ADD_METHOD_CALL

	// example: ATC_RESULTS
	result, err := atc.Execute(algodClient, context.Background(), 4)
	if err != nil {
		log.Fatalf("failed to get add method: %s", err)
	}

	for _, r := range result.MethodResults {
		log.Printf("%s => %v", r.Method.Name, r.ReturnValue)
	}
	// example: ATC_RESULTS

	// example: ATC_BOX_REF
	boxName := "coolBoxName"
	//nolint:ineffassign,staticcheck // example code demonstrating BoxReferences structure
	mcp = transaction.AddMethodCallParams{
		AppID:           appID,
		Sender:          acct1.Address,
		SuggestedParams: sp,
		OnComplete:      types.NoOpOC,
		Signer:          signer,
		Method:          addMethod,
		MethodArgs:      []interface{}{1, 1},
		// Here we're passing a box reference so our app
		// can reference it during evaluation
		BoxReferences: []types.AppBoxReference{
			{AppID: appID, Name: []byte(boxName)},
		},
	}
	// ...
	// example: ATC_BOX_REF

}
