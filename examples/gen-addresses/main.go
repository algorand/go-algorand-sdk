package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// CHANGE ME
const (
	exampleWalletName     = "unencrypted-default-wallet"
	exampleWalletPassword = ""
	exampleWalletDriver   = kmd.DefaultWalletDriver
)

var (
	kmdAddress = "http://localhost:4002"
	kmdToken   = strings.Repeat("a", 64)

	algodAddress = "http://localhost:4001"
	algodToken   = strings.Repeat("a", 64)
)

func main() {
	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}
	fmt.Println("Made a kmd client")

	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return
	}

	// Print algod status
	nodeStatus, err := algodClient.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		return
	}
	fmt.Printf("algod last round: %d\n", nodeStatus.LastRound)

	// List existing wallets, and check if our example wallet already exists
	resp0, err := kmdClient.ListWallets()
	if err != nil {
		fmt.Printf("error listing wallets: %s\n", err)
		return
	}
	fmt.Printf("Got %d wallet(s):\n", len(resp0.Wallets))
	var exampleExists bool
	var exampleWalletID string
	for _, wallet := range resp0.Wallets {
		fmt.Printf("ID: %s\tName: %s\n", wallet.ID, wallet.Name)
		if wallet.Name == exampleWalletName {
			exampleWalletID = wallet.ID
			exampleExists = true
		}
	}

	// Create the example wallet, if it doesn't already exist
	if !exampleExists {
		resp1, err := kmdClient.CreateWallet(exampleWalletName, exampleWalletPassword, exampleWalletDriver, types.MasterDerivationKey{})
		if err != nil {
			fmt.Printf("error creating wallet: %s\n", err)
			return
		}
		exampleWalletID = resp1.Wallet.ID
		fmt.Printf("Created wallet '%s' with ID: %s\n", resp1.Wallet.Name, exampleWalletID)
	}

	// Get a wallet handle
	resp2, err := kmdClient.InitWalletHandle(exampleWalletID, exampleWalletPassword)
	if err != nil {
		fmt.Printf("Error initializing wallet: %s\n", err)
		return
	}

	// Extract the wallet handle
	exampleWalletHandleToken := resp2.WalletHandleToken

	// Generate some addresses in the wallet
	fmt.Println("Generating 10 addresses")
	var addresses []string
	for i := 0; i < 10; i++ {
		resp3, err := kmdClient.GenerateKey(exampleWalletHandleToken)
		if err != nil {
			fmt.Printf("Error generating key: %s\n", err)
			return
		}
		fmt.Printf("Generated address %s\n", resp3.Address)
		addresses = append(addresses, resp3.Address)
	}

	// Extract the private key of the first address
	fmt.Printf("Extracting private key for %s\n", addresses[0])
	resp4, err := kmdClient.ExportKey(exampleWalletHandleToken, exampleWalletPassword, addresses[0])
	if err != nil {
		fmt.Printf("Error extracting secret key: %s\n", err)
		return
	}
	privateKey := resp4.PrivateKey

	// Get the suggested transaction parameters
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("error getting suggested tx params: %s\n", err)
		return
	}

	tx, err := transaction.MakePaymentTxn(addresses[0], addresses[1], 100, nil, "", txParams)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}
	fmt.Printf("Made unsigned transaction: %+v\n", tx)
	fmt.Println("Signing transaction with go-algo-sdk library function (not kmd)")

	txid, stx, err := crypto.SignTransaction(privateKey, tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}

	fmt.Printf("Made signed transaction with TxID %s: %x\n", txid, stx)

	// Sign the same transaction with kmd
	fmt.Println("Signing same transaction with kmd")
	resp5, err := kmdClient.SignTransaction(exampleWalletHandleToken, exampleWalletPassword, tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction with kmd: %s\n", err)
		return
	}

	fmt.Printf("kmd made signed transaction with bytes: %x\n", resp5.SignedTransaction)
	if bytes.Equal(resp5.SignedTransaction, stx) {
		fmt.Println("signed transactions match!")
	} else {
		fmt.Println("signed transactions don't match!")
	}
}
