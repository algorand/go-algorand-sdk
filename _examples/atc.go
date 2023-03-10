package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/algorand/go-algorand-sdk/v2/abi"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
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

	appID := deployApp(algodClient, acct1)
	log.Printf("%d", appID)

	// example: ATC_CONTRACT_INIT
	b, err := ioutil.ReadFile("calculator/contract.json")
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

func deployApp(algodClient *algod.Client, creator crypto.Account) uint64 {

	var (
		approvalBinary = make([]byte, 1000)
		clearBinary    = make([]byte, 1000)
	)

	// Compile approval program
	approvalTeal, err := ioutil.ReadFile("calculator/approval.teal")
	if err != nil {
		log.Fatalf("failed to read approval program: %s", err)
	}

	approvalResult, err := algodClient.TealCompile(approvalTeal).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile program: %s", err)
	}

	_, err = base64.StdEncoding.Decode(approvalBinary, []byte(approvalResult.Result))
	if err != nil {
		log.Fatalf("failed to decode compiled program: %s", err)
	}

	// Compile clear program
	clearTeal, err := ioutil.ReadFile("calculator/clear.teal")
	if err != nil {
		log.Fatalf("failed to read clear program: %s", err)
	}

	clearResult, err := algodClient.TealCompile(clearTeal).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile program: %s", err)
	}

	_, err = base64.StdEncoding.Decode(clearBinary, []byte(clearResult.Result))
	if err != nil {
		log.Fatalf("failed to decode compiled program: %s", err)
	}

	// Create application
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	txn, err := transaction.MakeApplicationCreateTx(
		false, approvalBinary, clearBinary,
		types.StateSchema{}, types.StateSchema{},
		nil, nil, nil, nil,
		sp, creator.Address, nil,
		types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	txid, stx, err := crypto.SignTransaction(creator.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	return confirmedTxn.ApplicationIndex
}
