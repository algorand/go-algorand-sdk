package main

import (
  "fmt"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

func main() {
  // example: GOSDK_ACCOUNT_GENERATE
  account := crypto.GenerateAccount()
  passphrase, err := mnemonic.FromPrivateKey(account.PrivateKey)

  if err != nil {
  	fmt.Printf("Error with private key: %s\n", err)
  } else {
  	fmt.Printf("My address: %s\n", account.Address)
  	fmt.Printf("My passphrase: %s\n", passphrase)
  }
  // example: GOSDK_ACCOUNT_GENERATE
}
