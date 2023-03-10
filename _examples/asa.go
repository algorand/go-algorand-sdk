package main

import (
	"context"
	"log"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

func main() {

	algodClient := getAlgodClient()
	accts, err := getSandboxAccounts()
	if err != nil {
		log.Fatalf("failed to get sandbox accounts: %s", err)
	}

	creator := accts[0]
	user := accts[1]

	assetID := createAsset(algodClient, creator)
	// example: ASSET_INFO
	info, err := algodClient.GetAssetByID(assetID).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get asset info: %s", err)
	}
	log.Printf("Asset info for %d: %+v", assetID, info)
	// example: ASSET_INFO

	configureAsset(algodClient, assetID, creator)
	optInAsset(algodClient, assetID, user)
	xferAsset(algodClient, assetID, creator, user)
	freezeAsset(algodClient, assetID, creator, user)
	clawbackAsset(algodClient, assetID, creator, user)
	deleteAsset(algodClient, assetID, creator)
}

func createAsset(algodClient *algod.Client, creator crypto.Account) uint64 {
	// example: ASSET_CREATE
	// Configure parameters for asset creation
	var (
		creatorAddr       = creator.Address.String()
		assetName         = "Really Useful Gift"
		unitName          = "rug"
		assetURL          = "https://path/to/my/asset/details"
		assetMetadataHash = "thisIsSomeLength32HashCommitment"
		defaultFrozen     = false
		decimals          = uint32(0)
		totalIssuance     = uint64(1000)

		manager  = creatorAddr
		reserve  = creatorAddr
		freeze   = creatorAddr
		clawback = creatorAddr

		note []byte
	)

	// Get network-related transaction parameters and assign
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	// Construct the transaction
	txn, err := transaction.MakeAssetCreateTxn(
		creatorAddr, note, txParams, totalIssuance, decimals,
		defaultFrozen, manager, reserve, freeze, clawback,
		unitName, assetName, assetURL, assetMetadataHash,
	)

	if err != nil {
		log.Fatalf("failed to make transaction: %s", err)
	}

	// sign the transaction
	txid, stx, err := crypto.SignTransaction(creator.PrivateKey, txn)
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

	log.Printf("Create Transaction: %s confirmed in Round %d with new asset id: %d\n",
		txid, confirmedTxn.ConfirmedRound, confirmedTxn.AssetIndex)
	// example: ASSET_CREATE
	return confirmedTxn.AssetIndex
}

func configureAsset(algodClient *algod.Client, assetID uint64, creator crypto.Account) {
	// example: ASSET_CONFIG
	creatorAddr := creator.Address.String()
	var (
		newManager  = creatorAddr
		newFreeze   = creatorAddr
		newClawback = creatorAddr
		newReserve  = ""

		strictAddrCheck = false
		note            []byte
	)

	// Get network-related transaction parameters and assign
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	txn, err := transaction.MakeAssetConfigTxn(creatorAddr, note, sp, assetID, newManager, newReserve, newFreeze, newClawback, strictAddrCheck)
	if err != nil {
		log.Fatalf("failed to make  txn: %s", err)
	}
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(creator.PrivateKey, txn)
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

	log.Printf("Asset Config Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: ASSET_CONFIG
}

func optInAsset(algodClient *algod.Client, assetID uint64, user crypto.Account) {
	// example: ASSET_OPTIN
	userAddr := user.Address.String()

	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	txn, err := transaction.MakeAssetAcceptanceTxn(userAddr, nil, sp, assetID)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(user.PrivateKey, txn)
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
	// example: ASSET_OPTIN
}

func xferAsset(algodClient *algod.Client, assetID uint64, creator crypto.Account, user crypto.Account) {
	// example: ASSET_XFER
	var (
		creatorAddr = creator.Address.String()
		userAddr    = user.Address.String()
	)

	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	txn, err := transaction.MakeAssetTransferTxn(creatorAddr, userAddr, 1, nil, sp, "", assetID)
	if err != nil {
		log.Fatalf("failed to make asset txn: %s", err)
	}
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(creator.PrivateKey, txn)
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

	log.Printf("Asset Transfer Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: ASSET_XFER
}

func freezeAsset(algodClient *algod.Client, assetID uint64, creator crypto.Account, user crypto.Account) {
	// example: ASSET_FREEZE
	var (
		creatorAddr = creator.Address.String()
		userAddr    = user.Address.String()
	)

	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	// Create a freeze asset transaction with the target of the user address
	// and the new freeze setting of `true`
	txn, err := transaction.MakeAssetFreezeTxn(creatorAddr, nil, sp, assetID, userAddr, true)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(creator.PrivateKey, txn)
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

	log.Printf("Freeze Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: ASSET_FREEZE
}

func clawbackAsset(algodClient *algod.Client, assetID uint64, creator crypto.Account, user crypto.Account) {
	// example: ASSET_CLAWBACK
	var (
		creatorAddr = creator.Address.String()
		userAddr    = user.Address.String()
	)

	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	// Create a new clawback transaction with the target of the user address and the recipient as the creator
	// address, being sent from the address marked as `clawback` on the asset, in this case the same as creator
	txn, err := transaction.MakeAssetRevocationTxn(creatorAddr, userAddr, 1, creatorAddr, nil, sp, assetID)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(creator.PrivateKey, txn)
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

	log.Printf("Clawback Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: ASSET_CLAWBACK
}

func deleteAsset(algodClient *algod.Client, assetID uint64, creator crypto.Account) {
	// example: ASSET_DELETE
	var (
		creatorAddr = creator.Address.String()
	)

	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	// Create a new clawback transaction with the target of the user address and the recipient as the creator
	// address, being sent from the address marked as `clawback` on the asset, in this case the same as creator
	txn, err := transaction.MakeAssetDestroyTxn(creatorAddr, nil, sp, assetID)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}
	// sign the transaction
	txid, stx, err := crypto.SignTransaction(creator.PrivateKey, txn)
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

	log.Printf("Destroy Transaction: %s confirmed in Round %d\n", txid, confirmedTxn.ConfirmedRound)
	// example: ASSET_DELETE
}
