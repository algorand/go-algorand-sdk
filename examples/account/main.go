package main

import (
	"context"
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/examples"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

func main() {
	// example: ACCOUNT_GENERATE
	account := crypto.GenerateAccount()
	mn, err := mnemonic.FromPrivateKey(account.PrivateKey)

	if err != nil {
		log.Fatalf("failed to generate account: %s", err)
	}

	log.Printf("Address: %s\n", account.Address)
	log.Printf("Mnemonic: %s\n", mn)
	// example: ACCOUNT_GENERATE

	// example: ACCOUNT_RECOVER_MNEMONIC
	k, err := mnemonic.ToPrivateKey(mn)
	if err != nil {
		log.Fatalf("failed to parse mnemonic: %s", err)
	}

	recovered, err := crypto.AccountFromPrivateKey(k)
	if err != nil {
		log.Fatalf("failed to recover account from key: %s", err)
	}

	log.Printf("%+v", recovered)
	// example: ACCOUNT_RECOVER_MNEMONIC
	accts, err := examples.GetSandboxAccounts()
	if err != nil {
		log.Fatalf("failed to get sandbox accounts: %s", err)
	}
	rekeyAccount(examples.GetAlgodClient(), accts[0], accts[1])
}

func rekeyAccount(algodClient *algod.Client, acct crypto.Account, rekeyTarget crypto.Account) {
	// example: ACCOUNT_REKEY
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get suggested params: %s", err)
	}

	addr := acct.Address.String()
	// here we create a payment transaction but rekey is valid
	// on any transaction type
	rktxn, err := transaction.MakePaymentTxn(addr, addr, 0, nil, "", sp)
	if err != nil {
		log.Fatalf("failed to creating transaction: %s\n", err)
	}
	// Set the rekey parameter
	rktxn.RekeyTo = rekeyTarget.Address

	_, stxn, err := crypto.SignTransaction(acct.PrivateKey, rktxn)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
	}

	txID, err := algodClient.SendRawTransaction(stxn).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}

	result, err := transaction.WaitForConfirmation(algodClient, txID, 4, context.Background())
	if err != nil {
		fmt.Printf("Error waiting for confirmation on txID: %s\n", txID)
		return
	}

	fmt.Printf("Confirmed Transaction: %s in Round %d\n", txID, result.ConfirmedRound)
	// example: ACCOUNT_REKEY

	// rekey back
	rktxn, _ = transaction.MakePaymentTxn(addr, addr, 0, nil, "", sp)
	rktxn.RekeyTo = acct.Address
	_, stxn, _ = crypto.SignTransaction(rekeyTarget.PrivateKey, rktxn)
	txID, _ = algodClient.SendRawTransaction(stxn).Do(context.Background())
	result, _ = transaction.WaitForConfirmation(algodClient, txID, 4, context.Background())
}
