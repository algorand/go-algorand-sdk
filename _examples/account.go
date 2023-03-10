package main

import (
	"log"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
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

}

// Missing ACCOUNT_REKEY in GOSDK examples (in ../docs/get-details/accounts/rekey.md:391)
