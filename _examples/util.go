package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
)

// add sandbox and other stuff
var (
	ALGOD_ADDRESS = "http://localhost:4001"
	ALGOD_TOKEN   = strings.Repeat("a", 64)

	KMD_ADDRESS         = "http://localhost:4002"
	KMD_TOKEN           = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	KMD_WALLET_NAME     = "unencrypted-default-wallet"
	KMD_WALLET_PASSWORD = ""
)

func getAlgodClient() *algod.Client {
	algodClient, err := algod.MakeClient(
		ALGOD_ADDRESS,
		ALGOD_TOKEN,
	)

	if err != nil {
		log.Fatalf("Failed to create algod client: %s", err)
	}

	return algodClient
}

func getSandboxAccounts() ([]crypto.Account, error) {
	client, err := kmd.MakeClient(KMD_ADDRESS, KMD_TOKEN)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %+v", err)
	}

	resp, err := client.ListWallets()
	if err != nil {
		return nil, fmt.Errorf("Failed to list wallets: %+v", err)
	}

	var walletId string
	for _, wallet := range resp.Wallets {
		if wallet.Name == KMD_WALLET_NAME {
			walletId = wallet.ID
		}
	}

	if walletId == "" {
		return nil, fmt.Errorf("No wallet named %s", KMD_WALLET_NAME)
	}

	whResp, err := client.InitWalletHandle(walletId, KMD_WALLET_PASSWORD)
	if err != nil {
		return nil, fmt.Errorf("Failed to init wallet handle: %+v", err)
	}

	addrResp, err := client.ListKeys(whResp.WalletHandleToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to list keys: %+v", err)
	}

	var accts []crypto.Account
	for _, addr := range addrResp.Addresses {
		expResp, err := client.ExportKey(whResp.WalletHandleToken, KMD_WALLET_PASSWORD, addr)
		if err != nil {
			return nil, fmt.Errorf("Failed to export key: %+v", err)
		}

		acct, err := crypto.AccountFromPrivateKey(expResp.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("Failed to create account from private key: %+v", err)
		}

		accts = append(accts, acct)
	}

	return accts, nil
}
