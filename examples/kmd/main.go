package main

import (
	"crypto/ed25519"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/examples"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func main() {
	var (
		exampleWalletID          string
		exampleWalletHandleToken string
		genResponse              kmd.GenerateKeyResponse
		initResponse             kmd.InitWalletHandleResponse
	)

	// example: KMD_CREATE_CLIENT
	// Create a new kmd client, configured to connect to out local sandbox
	var kmdAddress = "http://localhost:4002"
	var kmdToken = strings.Repeat("a", 64)
	kmdClient, err := kmd.MakeClient(
		kmdAddress,
		kmdToken,
	)
	// example: KMD_CREATE_CLIENT
	_ = kmdClient
	kmdClient = examples.GetKmdClient()

	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}

	// example: KMD_CREATE_WALLET
	// Create the example wallet, if it doesn't already exist
	createResponse, err := kmdClient.CreateWallet(
		"DemoWallet",
		"password",
		kmd.DefaultWalletDriver,
		types.MasterDerivationKey{},
	)
	if err != nil {
		fmt.Printf("error creating wallet: %s\n", err)
		return
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	exampleWalletID = createResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", createResponse.Wallet.Name, exampleWalletID)
	// example: KMD_CREATE_WALLET

	// example: KMD_CREATE_ACCOUNT
	// Get a wallet handle.
	initResponse, _ = kmdClient.InitWalletHandle(
		exampleWalletID,
		"password",
	)
	exampleWalletHandleToken = initResponse.WalletHandleToken

	// Generate a new address from the wallet handle
	genResponse, err = kmdClient.GenerateKey(exampleWalletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return
	}
	accountAddress := genResponse.Address
	fmt.Printf("New Account: %s\n", accountAddress)
	// example: KMD_CREATE_ACCOUNT

	// example: KMD_EXPORT_ACCOUNT
	// Extract the account sk
	accountKeyResponse, _ := kmdClient.ExportKey(
		exampleWalletHandleToken,
		"password",
		accountAddress,
	)
	accountKey := accountKeyResponse.PrivateKey
	// Convert sk to mnemonic
	mn, err := mnemonic.FromPrivateKey(accountKey)
	if err != nil {
		fmt.Printf("Error getting backup phrase: %s\n", err)
		return
	}
	fmt.Printf("Account Mnemonic: %v ", mn)
	// example: KMD_EXPORT_ACCOUNT

	// example: KMD_IMPORT_ACCOUNT
	account := crypto.GenerateAccount()
	fmt.Println("Account Address: ", account.Address)
	mn, err = mnemonic.FromPrivateKey(account.PrivateKey)
	if err != nil {
		fmt.Printf("Error getting backup phrase: %s\n", err)
		return
	}
	fmt.Printf("Account Mnemonic: %s\n", mn)
	importedAccount, _ := kmdClient.ImportKey(
		exampleWalletHandleToken,
		account.PrivateKey,
	)
	fmt.Println("Account Successfully Imported: ", importedAccount.Address)
	// example: KMD_IMPORT_ACCOUNT

	// Get the MDK for Recovery example
	backupResponse, err := kmdClient.ExportMasterDerivationKey(exampleWalletHandleToken, "password")
	if err != nil {
		fmt.Printf("error exporting mdk: %s\n", err)
		return
	}
	backupPhrase, _ := mnemonic.FromMasterDerivationKey(backupResponse.MasterDerivationKey)
	fmt.Printf("Backup: %s\n", backupPhrase)

	// example: KMD_RECOVER_WALLET
	keyBytes, err := mnemonic.ToKey(backupPhrase)
	if err != nil {
		fmt.Printf("failed to get key: %s\n", err)
		return
	}

	var mdk types.MasterDerivationKey
	copy(mdk[:], keyBytes)
	recoverResponse, err := kmdClient.CreateWallet(
		"RecoveryWallet",
		"password",
		kmd.DefaultWalletDriver,
		mdk,
	)
	if err != nil {
		fmt.Printf("error creating wallet: %s\n", err)
		return
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	exampleWalletID = recoverResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", recoverResponse.Wallet.Name, exampleWalletID)

	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	initResponse, err = kmdClient.InitWalletHandle(
		exampleWalletID,
		"password",
	)
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return
	}

	// Extract the wallet handle
	exampleWalletHandleToken = initResponse.WalletHandleToken
	fmt.Printf("Got wallet handle: '%s'\n", exampleWalletHandleToken)

	// Generate a new address from the wallet handle
	genResponse, err = kmdClient.GenerateKey(exampleWalletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return
	}
	fmt.Printf("Recovered address %s\n", genResponse.Address)
	// example: KMD_RECOVER_WALLET

	// example: ACCOUNT_GENERATE
	newAccount := crypto.GenerateAccount()
	passphrase, err := mnemonic.FromPrivateKey(newAccount.PrivateKey)

	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
	} else {
		fmt.Printf("My address: %s\n", newAccount.Address)
		fmt.Printf("My passphrase: %s\n", passphrase)
	}
	// example: ACCOUNT_GENERATE

	// example: MULTISIG_CREATE
	// Get pre-defined set of keys for example
	_, pks := loadAccounts()
	addr1, _ := types.DecodeAddress(pks[1])
	addr2, _ := types.DecodeAddress(pks[2])
	addr3, _ := types.DecodeAddress(pks[3])

	ma, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		addr1,
		addr2,
		addr3,
	})

	if err != nil {
		panic("invalid multisig parameters")
	}
	fromAddr, _ := ma.Address()
	// Print multisig account
	fmt.Printf("Multisig address : %s \n", fromAddr)
	// example: MULTISIG_CREATE
}

// Accounts to be used through examples
func loadAccounts() (map[int][]byte, map[int]string) {
	// Shown for demonstration purposes. NEVER reveal secret mnemonics in practice.
	// Change these values to use the accounts created previously.
	// Paste in mnemonic phrases for all three accounts
	acc1 := crypto.GenerateAccount()
	mnemonic1, _ := mnemonic.FromPrivateKey(acc1.PrivateKey)
	acc2 := crypto.GenerateAccount()
	mnemonic2, _ := mnemonic.FromPrivateKey(acc2.PrivateKey)
	acc3 := crypto.GenerateAccount()
	mnemonic3, _ := mnemonic.FromPrivateKey(acc3.PrivateKey)

	mnemonics := []string{mnemonic1, mnemonic2, mnemonic3}
	pks := map[int]string{1: "", 2: "", 3: ""}
	var sks = make(map[int][]byte)

	for i, m := range mnemonics {
		var err error
		sk, err := mnemonic.ToPrivateKey(m)
		sks[i+1] = sk
		if err != nil {
			fmt.Printf("Issue with account %d private key conversion.", i+1)
		}
		// derive public address from Secret Key.
		pk := sk.Public()
		var a types.Address
		cpk := pk.(ed25519.PublicKey)
		copy(a[:], cpk[:])
		pks[i+1] = a.String()
		fmt.Printf("Loaded Key %d: %s\n", i+1, pks[i+1])
	}
	return sks, pks
}
