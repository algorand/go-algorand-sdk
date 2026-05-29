package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
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
	alice := accts[0]
	bob := accts[1]

	// Check versions
	versions, err := algodClient.Versions().Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get versions: %s", err)
	}
	// Verify we have at least version 4.7
	if versions.Build.Major < 4 {
		log.Fatalf("Need algod version >= 4.7")
	}
	if versions.Build.Major == 4 && versions.Build.Minor < 7 {
		log.Fatalf("Need algod version >= 4.7")
	}

	// Create app with box storage
	appID := createAppWithBoxes(algodClient, alice)
	fmt.Printf("Created app %d\n", appID)

	// Fund the app with some algos for box storage
	appAddr := crypto.GetApplicationAddress(appID)
	fundApp(algodClient, alice, appAddr)

	// Set boxes for Alice and Bob
	setBoxes(algodClient, appID, alice, "Alice", "red")
	setBoxes(algodClient, appID, bob, "Bob", "blue")

	// Get box keys without values
	boxKeys, err := algodClient.GetApplicationBoxes(appID).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get box keys: %s", err)
	}
	fmt.Printf("Box Keys count: %d\n", len(boxKeys.Boxes))
	if len(boxKeys.Boxes) == 0 {
		log.Fatal("Expected at least one box")
	}
	if len(boxKeys.Boxes[0].Value) != 0 {
		log.Fatalf("Expected box value to be empty (nil), got: %v", boxKeys.Boxes[0].Value)
	}

	// Get boxes with values
	boxValues, err := algodClient.GetApplicationBoxes(appID).Include([]string{"values"}).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get box values: %s", err)
	}
	fmt.Printf("Box Values count: %d\n", len(boxValues.Boxes))
	if len(boxValues.Boxes) != 4 {
		log.Fatalf("Expected 4 boxes, got %d", len(boxValues.Boxes))
	}
	if len(boxValues.Boxes[0].Value) == 0 {
		log.Fatalf("Expected box value to be populated, got empty: %v", boxValues.Boxes[0].Value)
	}

	// Get Alice's boxes by prefix
	alicePrefix := fmt.Sprintf("b64:%s", base64.StdEncoding.EncodeToString(alice.Address[:]))
	aliceValues, err := algodClient.GetApplicationBoxes(appID).
		Include([]string{"values"}).
		Prefix(alicePrefix).
		Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get alice's boxes: %s", err)
	}
	fmt.Printf("Alice's values count: %d\n", len(aliceValues.Boxes))
	if len(aliceValues.Boxes) != 2 {
		log.Fatalf("Expected 2 boxes for Alice, got %d", len(aliceValues.Boxes))
	}

	// Get Bob's boxes by prefix
	bobPrefix := fmt.Sprintf("b64:%s", base64.StdEncoding.EncodeToString(bob.Address[:]))
	bobValues, err := algodClient.GetApplicationBoxes(appID).
		Include([]string{"values"}).
		Prefix(bobPrefix).
		Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get bob's boxes: %s", err)
	}
	fmt.Printf("Bob's values count: %d\n", len(bobValues.Boxes))
	if len(bobValues.Boxes) != 2 {
		log.Fatalf("Expected 2 boxes for Bob, got %d", len(bobValues.Boxes))
	}

	// Paginate through boxes
	var pages []models.BoxesResponse
	var nextToken string
	for {
		page, err := algodClient.GetApplicationBoxes(appID).
			Include([]string{"values"}).
			Limit(1).
			Next(nextToken).
			Do(context.Background())
		if err != nil {
			log.Fatalf("failed to get paginated boxes: %s", err)
		}
		pages = append(pages, page)
		nextToken = page.NextToken
		if nextToken == "" {
			break
		}
	}
	fmt.Printf("Pages count: %d\n", len(pages))
	if len(pages) != 4 {
		log.Fatalf("Expected 4 pages, got %d", len(pages))
	}

	fmt.Println("All assertions passed!")
}

func createAppWithBoxes(algodClient *algod.Client, creator crypto.Account) uint64 {
	// Approval program that stores box data
	approvalProgram := `#pragma version 12

txn ApplicationID
bz return

txn Sender
byte "name"
concat
txna ApplicationArgs 0
box_put

txn Sender
byte "favoriteColor"
concat
txna ApplicationArgs 1
box_put

return:
  int 1
  return
`
	approvalResult, err := algodClient.TealCompile([]byte(approvalProgram)).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile approval program: %s", err)
	}

	approvalBinary, err := base64.StdEncoding.DecodeString(approvalResult.Result)
	if err != nil {
		log.Fatalf("failed to decode compiled approval program: %s", err)
	}

	// Clear program - always return 1
	clearProgram := `#pragma version 12
int 1
return
`
	clearResult, err := algodClient.TealCompile([]byte(clearProgram)).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile clear program: %s", err)
	}

	clearBinary, err := base64.StdEncoding.DecodeString(clearResult.Result)
	if err != nil {
		log.Fatalf("failed to decode compiled clear program: %s", err)
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
		log.Fatalf("error waiting for confirmation: %s", err)
	}

	return confirmedTxn.ApplicationIndex
}

func fundApp(algodClient *algod.Client, sender crypto.Account, appAddr types.Address) {
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	txn, err := transaction.MakePaymentTxn(
		sender.Address.String(), appAddr.String(),
		1000000, // 1 algo
		nil, "", sp,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	txid, stx, err := crypto.SignTransaction(sender.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	_, err = transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation: %s", err)
	}
}

func setBoxes(algodClient *algod.Client, appID uint64, sender crypto.Account, name string, favoriteColor string) {
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("error getting suggested tx params: %s", err)
	}

	// Create box names by concatenating address with field name
	nameBox := append(sender.Address[:], []byte("name")...)
	colorBox := append(sender.Address[:], []byte("favoriteColor")...)

	appArgs := [][]byte{
		[]byte(name),
		[]byte(favoriteColor),
	}

	boxRefs := []types.AppBoxReference{
		{AppID: appID, Name: nameBox},
		{AppID: appID, Name: colorBox},
	}

	txn, err := transaction.MakeApplicationNoOpTxWithBoxes(
		appID, appArgs, nil, nil, nil, boxRefs,
		sp, sender.Address, nil,
		types.Digest{}, [32]byte{}, types.ZeroAddress,
	)
	if err != nil {
		log.Fatalf("failed to make txn: %s", err)
	}

	txid, stx, err := crypto.SignTransaction(sender.PrivateKey, txn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	_, err = algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	_, err = transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("error waiting for confirmation: %s", err)
	}

	fmt.Printf("Set %s to %s\n", sender.Address.String(), map[string]string{"name": name, "favoriteColor": favoriteColor})
}
