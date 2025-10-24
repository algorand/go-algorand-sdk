package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
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
	appID := appCreate(algodClient, acct1)
	appOptIn(algodClient, appID, acct1)

	// example: APP_READ_STATE
	// grab global state and config of application
	appInfo, err := algodClient.GetApplicationByID(appID).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get app info: %s", err)
	}
	log.Printf("app info: %+v", appInfo)

	// grab local state for an app id for a single account
	acctInfo, err := algodClient.AccountApplicationInformation(
		acct1.Address.String(), appID,
	).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get app info: %s", err)
	}
	log.Printf("app info: %+v", acctInfo)
	// example: APP_READ_STATE

	appNoOp(algodClient, appID, acct1)
	appUpdate(algodClient, appID, acct1)
	appCall(algodClient, appID, acct1)
	appCloseOut(algodClient, appID, acct1)
	appDelete(algodClient, appID, acct1)

}

func appCreate(algodClient *algod.Client, creator crypto.Account) uint64 {
	// example: APP_SCHEMA
	// declare application state storage (immutable)
	var (
		localInts  uint64 = 1
		localBytes uint64 = 1
		globalInts uint64 = 1
		//revive:disable:var-declaration // We want the explicit declaration here for the example
		globalBytes uint64 = 0
	)

	// define schema
	globalSchema := types.StateSchema{NumUint: globalInts, NumByteSlice: globalBytes}
	localSchema := types.StateSchema{NumUint: localInts, NumByteSlice: localBytes}
	// example: APP_SCHEMA

	// example: APP_SOURCE
	approvalTeal, err := os.ReadFile("application/approval.teal")
	if err != nil {
		log.Fatalf("failed to read approval program: %s", err)
	}
	clearTeal, err := os.ReadFile("application/clear.teal")
	if err != nil {
		log.Fatalf("failed to read clear program: %s", err)
	}
	// example: APP_SOURCE

	// example: APP_COMPILE
	approvalResult, err := algodClient.TealCompile(approvalTeal).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile program: %s", err)
	}

	approvalBinary, err := base64.StdEncoding.DecodeString(approvalResult.Result)
	if err != nil {
		log.Fatalf("failed to decode compiled program: %s", err)
	}

	clearResult, err := algodClient.TealCompile(clearTeal).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile program: %s", err)
	}

	clearBinary, err := base64.StdEncoding.DecodeString(clearResult.Result)
	if err != nil {
		log.Fatalf("failed to decode compiled program: %s", err)
	}
	// example: APP_COMPILE

	// example: APP_CREATE
	// Create application
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	txn, err := transaction.MakeApplicationCreateTx(
		false, approvalBinary, clearBinary, globalSchema, localSchema,
		nil, nil, nil, nil, sp, creator.Address, nil,
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
	appID := confirmedTxn.ApplicationIndex
	log.Printf("Created app with id: %d", appID)
	// example: APP_CREATE
	return appID
}

func appOptIn(algodClient *algod.Client, appID uint64, caller crypto.Account) {
	// example: APP_OPTIN
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	// Create a new clawback transaction with the target of the user address and the recipient as the creator
	// address, being sent from the address marked as `clawback` on the asset, in this case the same as creator
	txn, err := transaction.MakeApplicationOptInTx(
		appID, nil, nil, nil, nil, sp,
		caller.Address, nil, types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(caller.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	// Broadcast the transaction to the network
	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	// Wait for confirmation
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	log.Printf("OptIn Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: APP_OPTIN
}

func appNoOp(algodClient *algod.Client, appID uint64, caller crypto.Account) {
	// example: APP_NOOP
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	var (
		appArgs [][]byte
		accts   []string
		apps    []uint64
		assets  []uint64
	)

	// Add an arg to our app call
	appArgs = append(appArgs, []byte("arg0"))

	txn, err := transaction.MakeApplicationNoOpTx(
		appID, appArgs, accts, apps, assets, sp,
		caller.Address, nil, types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	// sign the transaction
	txid, stx, err := crypto.SignTransaction(caller.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	// Broadcast the transaction to the network
	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	// Wait for confirmation
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	log.Printf("NoOp Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: APP_NOOP
}

func appUpdate(algodClient *algod.Client, appID uint64, caller crypto.Account) {
	approvalBinary := examples.CompileTeal(algodClient, "application/approval_refactored.teal")
	clearBinary := examples.CompileTeal(algodClient, "application/clear.teal")

	// example: APP_UPDATE
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	var (
		appArgs [][]byte
		accts   []string
		apps    []uint64
		assets  []uint64
	)

	txn, err := transaction.MakeApplicationUpdateTx(
		appID, appArgs, accts, apps, assets, approvalBinary, clearBinary,
		sp, caller.Address, nil, types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	// sign the transaction
	txid, stx, err := crypto.SignTransaction(caller.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	// Broadcast the transaction to the network
	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	// Wait for confirmation
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	log.Printf("Update Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: APP_UPDATE
}

func appCloseOut(algodClient *algod.Client, appID uint64, caller crypto.Account) {
	// example: APP_CLOSEOUT
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	var (
		appArgs [][]byte
		accts   []string
		apps    []uint64
		assets  []uint64
	)

	txn, err := transaction.MakeApplicationCloseOutTx(
		appID, appArgs, accts, apps, assets, sp,
		caller.Address, nil, types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	// sign the transaction
	txid, stx, err := crypto.SignTransaction(caller.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	// Broadcast the transaction to the network
	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	// Wait for confirmation
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	log.Printf("Closeout Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: APP_CLOSEOUT
}

//nolint:unused // example code for documentation
func appClearState(algodClient *algod.Client, appID uint64, caller crypto.Account) {
	// example: APP_CLEAR
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	var (
		appArgs [][]byte
		accts   []string
		apps    []uint64
		assets  []uint64
	)

	txn, err := transaction.MakeApplicationClearStateTx(
		appID, appArgs, accts, apps, assets, sp,
		caller.Address, nil, types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	// sign the transaction
	txid, stx, err := crypto.SignTransaction(caller.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	// Broadcast the transaction to the network
	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	// Wait for confirmation
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	log.Printf("ClearState Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: APP_CLEAR
}

func appCall(algodClient *algod.Client, appID uint64, caller crypto.Account) {
	// example: APP_CALL
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	var (
		appArgs [][]byte
		accts   []string
		apps    []uint64
		assets  []uint64
	)

	datetime := time.Now().Format("2006-01-02 at 15:04:05")
	appArgs = append(appArgs, []byte(datetime))

	txn, err := transaction.MakeApplicationNoOpTx(
		appID, appArgs, accts, apps, assets, sp,
		caller.Address, nil, types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	// sign the transaction
	txid, stx, err := crypto.SignTransaction(caller.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	// Broadcast the transaction to the network
	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	// Wait for confirmation
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	log.Printf("NoOp Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: APP_CALL
}

func appDelete(algodClient *algod.Client, appID uint64, caller crypto.Account) {
	// example: APP_DELETE
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	var (
		appArgs [][]byte
		accts   []string
		apps    []uint64
		assets  []uint64
	)

	txn, err := transaction.MakeApplicationDeleteTx(
		appID, appArgs, accts, apps, assets, sp,
		caller.Address, nil, types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	// sign the transaction
	txid, stx, err := crypto.SignTransaction(caller.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	// Broadcast the transaction to the network
	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	// Wait for confirmation
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	log.Printf("Delete Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: APP_DELETE
}
