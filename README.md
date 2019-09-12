# go-algorand-sdk

[![Build Status](https://travis-ci.com/algorand/go-algorand-sdk.svg?branch=master)](https://travis-ci.com/algorand/go-algorand-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/algorand/go-algorand-sdk)](https://goreportcard.com/report/github.com/algorand/go-algorand-sdk)
[![GoDoc](https://godoc.org/github.com/algorand/go-algorand-sdk?status.svg)](https://godoc.org/github.com/algorand/go-algorand-sdk)

The Algorand golang SDK provides:

- HTTP clients for the algod (agreement) and kmd (key management) APIs
- Standalone functionality for interacting with the Algorand protocol, including transaction signing, message encoding, etc.

# Documentation

Full documentation is available [on godoc](https://godoc.org/github.com/algorand/go-algorand-sdk). You can also self-host the documentation by running `godoc -http=:8099` and visiting `http://localhost:8099/pkg/github.com/algorand/go-algorand-sdk` in your web browser.

Additional developer documentation can be found on [developer.algorand.org](https://developer.algorand.org/)

# Package overview

In `client/`, the `algod` and `kmd` packages provide HTTP clients for their corresponding APIs. `algod` is the Algorand protocol daemon, responsible for reaching consensus with the network and participating in the Algorand protocol. You can use it to check the status of the blockchain, read a block, look at transactions, or submit a signed transaction. `kmd` is the key management daemon. It is responsible for managing spending key material, signing transactions, and managing wallets.

`types` contains the data structures you'll use when interacting with the network, including addresses, transactions, multisig signatures, etc. Some types (like `Transaction`) have their own packages containing constructors (like `MakePaymentTxn`).

`encoding` contains the `json` and `msgpack` packages, which can be used to serialize messages for the algod/kmd APIs and the network.

`mnemonic` contains support for turning 32-byte keys into checksummed, human-readable mnemonics (and going from mnemonics back to keys).

# Quick Start
To download the SDK, open a terminal and use the `go get` command.

```command
go get -u github.com/algorand/go-algorand-sdk/...
```

If you are connected to the Algorand network, your algod process should already be running. The kmd process must be started manually, however. Start and stop kmd using `goal kmd start` and `goal kmd stop`:

```command
goal kmd start -d <your-data-directory>
```

Here's a simple example which creates clients for algod and kmd:

```golang
package main

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
)

const algodAddress = "http://localhost:8080"
const kmdAddress = "http://localhost:7833"
const algodToken = "contents-of-algod.token"
const kmdToken = "contents-of-kmd.token"

func main() {
	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return
	}

	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		return
	}

	fmt.Printf("algod: %T, kmd: %T\n", algodClient, kmdClient)
}
```

# Examples

## algod client

Here is an example that creates an algod client and uses it to fetch node status information, and then a specific block.

```golang
package main

import (
	"encoding/json"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/algod"
)

// These constants represent the algod REST endpoint and the corresponding
// API token. You can retrieve these from the `algod.net` and `algod.token`
// files in the algod data directory.
const algodAddress = "http://localhost:8080"
const algodToken = "e48a9bbe064a08f19cde9f0f1b589c1188b24e5059bc661b31bd20b4c8fa4ce7"

func main() {
	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return
	}

	// Print algod status
	nodeStatus, err := algodClient.Status()
	if err != nil {
		fmt.Printf("error getting algod status: %s\n", err)
		return
	}

	fmt.Printf("algod last round: %d\n", nodeStatus.LastRound)
	fmt.Printf("algod time since last round: %d\n", nodeStatus.TimeSinceLastRound)
	fmt.Printf("algod catchup: %d\n", nodeStatus.CatchupTime)
	fmt.Printf("algod latest version: %s\n", nodeStatus.LastVersion)

	// Fetch block information
	lastBlock, err := algodClient.Block(nodeStatus.LastRound)
	if err != nil {
		fmt.Printf("error getting last block: %s\n", err)
		return
	}

	// Print the block information
	fmt.Printf("\n-----------------Block Information-------------------\n")
	blockJSON, err := json.MarshalIndent(lastBlock, "", "\t")
	if err != nil {
		fmt.Printf("Can not marshall block data: %s\n", err)
	}
	fmt.Printf("%s\n", blockJSON)
}
```

## kmd client

The following example creates a wallet, and generates an account within that wallet.

```golang
package main

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/types"
)

// These constants represent the kmdd REST endpoint and the corresponding API
// token. You can retrieve these from the `kmd.net` and `kmd.token` files in
// the kmd data directory.
const kmdAddress = "http://localhost:7833"
const kmdToken = "42b7482737a77d9e5dffb8493ac8899db5f95cbc744d4fcffc0f1c47a6db0c1e"

func main() {
	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}
	fmt.Println("Made a kmd client")

	// Create the example wallet, if it doesn't already exist
	cwResponse, err := kmdClient.CreateWallet("testwallet", "testpassword", kmd.DefaultWalletDriver, types.MasterDerivationKey{})
	if err != nil {
		fmt.Printf("error creating wallet: %s\n", err)
		return
	}

	// We need the wallet ID in order to get a wallet handle, so we can add accounts
	exampleWalletID := cwResponse.Wallet.ID
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, exampleWalletID)

	// Get a wallet handle. The wallet handle is used for things like signing transactions
	// and creating accounts. Wallet handles do expire, but they can be renewed
	initResponse, err := kmdClient.InitWalletHandle(exampleWalletID, "testpassword")
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return
	}

	// Extract the wallet handle
	exampleWalletHandleToken := initResponse.WalletHandleToken

	// Generate a new address from the wallet handle
	genResponse, err := kmdClient.GenerateKey(exampleWalletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return
	}
	fmt.Printf("Generated address %s\n", genResponse.Address)
}
```

This account can now be used to sign transactions, but you will need some funds to get started. If you are on the test network, you can use the [dispenser](https://bank.testnet.algorand.network) to seed your account with some Algos.

## Backing up a Wallet

You can export a master derivation key from the wallet and convert it to a mnemonic phrase in order to back up any generated addresses. This backup phrase will only allow you to recover wallet-generated keys; if you import an external key into a kmd-managed wallet, you'll need to back up that key by itself in order to recover it.

```golang
package main

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

// These constants represent the kmd REST endpoint and the corresponding API
// token. You can retrieve these from the `kmd.net` and `kmd.token` files in
// the kmd data directory.
const kmdAddress = "http://localhost:7833"
const kmdToken = "42b7482737a77d9e5dffb8493ac8899db5f95cbc744d4fcffc0f1c47a6db0c1e"

func main() {
	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}
	fmt.Println("Made a kmd client")

	// Get the list of wallets
	listResponse, err := kmdClient.ListWallets()
	if err != nil {
		fmt.Printf("error listing wallets: %s\n", err)
		return
	}

	// Find our wallet name in the list
	var exampleWalletID string
	fmt.Printf("Got %d wallet(s):\n", len(listResponse.Wallets))
	for _, wallet := range listResponse.Wallets {
		fmt.Printf("ID: %s\tName: %s\n", wallet.ID, wallet.Name)
		if wallet.Name == "testwallet" {
			fmt.Printf("found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			exampleWalletID = wallet.ID
		}
	}

	// Get a wallet handle
	initResponse, err := kmdClient.InitWalletHandle(exampleWalletID, "testpassword")
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return
	}

	// Extract the wallet handle
	exampleWalletHandleToken := initResponse.WalletHandleToken

	// Get the backup phrase
	exportResponse, err := kmdClient.ExportMasterDerivationKey(exampleWalletHandleToken, "testpassword")
	if err != nil {
		fmt.Printf("Error exporting backup phrase: %s\n", err)
		return
	}
	mdk := exportResponse.MasterDerivationKey

	// This string should be kept in a safe place and not shared
	stringToSave, err := mnemonic.FromKey(mdk[:])
	if err != nil {
		fmt.Printf("Error getting backup phrase: %s\n", err)
		return
	}

	fmt.Printf("Backup Phrase: %s\n", stringToSave)
}
```

To restore a wallet, convert the phrase to a key and pass it to `CreateWallet`. This call will fail if the wallet already exists:

```golang
package main

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
)

// These constants represent the kmd REST endpoint and the corresponding API
// token. You can retrieve these from the `kmd.net` and `kmd.token` files in
// the kmd data directory.
const kmdAddress = "http://localhost:7833"
const kmdToken = "42b7482737a77d9e5dffb8493ac8899db5f95cbc744d4fcffc0f1c47a6db0c1e"

func main() {
	// Create a kmd client
	kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
	if err != nil {
		fmt.Printf("failed to make kmd client: %s\n", err)
		return
	}
	backupPhrase := "fire enlist diesel stamp nuclear chunk student stumble call snow flock brush example slab guide choice option recall south kangaroo hundred matrix school above zero"
	keyBytes, err := mnemonic.ToKey(backupPhrase)
	if err != nil {
		fmt.Printf("failed to get key: %s\n", err)
		return
	}

	var mdk types.MasterDerivationKey
	copy(mdk[:], keyBytes)
	cwResponse, err := kmdClient.CreateWallet("testwallet", "testpassword", kmd.DefaultWalletDriver, mdk)
	if err != nil {
		fmt.Printf("error creating wallet: %s\n", err)
		return
	}
	fmt.Printf("Created wallet '%s' with ID: %s\n", cwResponse.Wallet.Name, cwResponse.Wallet.ID)
}
```

## Signing and submitting a transaction

The following example shows how to to use both KMD and Algod when signing and submitting a transaction.  You can also sign a transaction offline, which is shown in the next section of this document.
```golang
package main

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/transaction"
)

// CHANGE ME
const kmdAddress = "http://localhost:7833"
const kmdToken = "42b7482737a77d9e5dffb8493ac8899db5f95cbc744d4fcffc0f1c47a6db0c1e"
const algodAddress = "http://localhost:8080"
const algodToken = "6218386c0d964e371f34bbff4adf543dab14a7d9720c11c6f11970774d4575de"

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
	fmt.Println("Made an algod client")

	// Get the list of wallets
	listResponse, err := kmdClient.ListWallets()
	if err != nil {
		fmt.Printf("error listing wallets: %s\n", err)
		return
	}

	// Find our wallet name in the list
	var exampleWalletID string
	fmt.Printf("Got %d wallet(s):\n", len(listResponse.Wallets))
	for _, wallet := range listResponse.Wallets {
		fmt.Printf("ID: %s\tName: %s\n", wallet.ID, wallet.Name)
		if wallet.Name == "testwallet" {
			fmt.Printf("found wallet '%s' with ID: %s\n", wallet.Name, wallet.ID)
			exampleWalletID = wallet.ID
		}
	}
	// Get a wallet handle
	initResponse, err := kmdClient.InitWalletHandle(exampleWalletID, "testpassword")
	if err != nil {
		fmt.Printf("Error initializing wallet handle: %s\n", err)
		return
	}

	// Extract the wallet handle
	exampleWalletHandleToken := initResponse.WalletHandleToken

	// Generate a new address from the wallet handle
	gen1Response, err := kmdClient.GenerateKey(exampleWalletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return
	}
	fmt.Printf("Generated address 1 %s\n", gen1Response.Address)
	fromAddr := gen1Response.Address

	// Generate a new address from the wallet handle
	gen2Response, err := kmdClient.GenerateKey(exampleWalletHandleToken)
	if err != nil {
		fmt.Printf("Error generating key: %s\n", err)
		return
	}
	fmt.Printf("Generated address 2 %s\n", gen2Response.Address)
	toAddr := gen2Response.Address

	// Get the suggested transaction parameters
	txParams, err := algodClient.SuggestedParams()
        if err != nil {
                fmt.Printf("error getting suggested tx params: %s\n", err)
                return
        }

	// Make transaction
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	tx, err := transaction.MakePaymentTxn(fromAddr, toAddr, 1000, 200000, txParams.LastRound, (txParams.LastRound + 1000), nil, "", genID, genHash)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}

	// Sign the transaction
	signResponse, err := kmdClient.SignTransaction(exampleWalletHandleToken, "testpassword", tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction with kmd: %s\n", err)
		return
	}

	fmt.Printf("kmd made signed transaction with bytes: %x\n", signResponse.SignedTransaction)

	// Broadcast the transaction to the network
	// Note that this transaction will get rejected because the accounts do not have any tokens
	sendResponse, err := algodClient.SendRawTransaction(signResponse.SignedTransaction)
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}

	fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)
}
```
## Sign a transaction offline

The following example shows how to create a transaction and sign it offline. You can also create the transaction online and then sign it offline.
```golang
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
)

func main() {

	account := crypto.GenerateAccount()
	fmt.Printf("account address: %s\n", account.Address)

	m, err := mnemonic.FromPrivateKey(account.PrivateKey)
	fmt.Printf("backup phrase = %s\n", m)

	// Create and sign a sample transaction using this library, *not* kmd
	// This transaction will not be valid as the example parameters will most likely not be valid
	// You can use the algod client to get suggested values for the fee, first and last rounds, and genesisID
	const fee = 1000
	const amount = 20000
	const firstRound = 642715
	const lastRound = firstRound + 1000
	tx, err := transaction.MakePaymentTxn(
		account.Address.String(), "4MYUHDWHWXAKA5KA7U5PEN646VYUANBFXVJNONBK3TIMHEMWMD4UBOJBI4",
		fee, amount, firstRound, lastRound, nil, "", "", []byte("JgsgCaCTqIaLeVhyL6XlRu3n7Rfk2FxMeK+wRSaQ7dI="),
	)
	if err != nil {
		fmt.Printf("Error creating transaction: %s\n", err)
		return
	}
	fmt.Printf("Made unsigned transaction: %+v\n", tx)
	fmt.Println("Signing transaction with go-algo-sdk library function (not kmd)")

	// Sign the Transaction
	txid, bytes, err := crypto.SignTransaction(account.PrivateKey, tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}

	// Save the signed object to disk
	fmt.Printf("Made signed transaction with TxID %s\n", txid)
	filename := "./signed.tx"
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		fmt.Printf("Failed in saving transaction to file %s, error %s\n", filename, err)
		return
	}
	fmt.Printf("Saved signed transaction to file: %s\n", filename)
}
```
## Submit the transaction from a file

This example takes the output from the previous example (file containing signed transaction) and submits it to Algod process of a node.
```golang
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/algorand/go-algorand-sdk/client/algod"
)

// CHANGE ME
const algodAddress = "http://localhost:8080"
const algodToken = "f1dee49e36a82face92fdb21cd3d340a1b369925cd12f3ee7371378f1665b9b1"

func main() {

	rawTx, err := ioutil.ReadFile("./signed.tx")
	if err != nil {
		fmt.Printf("failed to open signed transaction: %s\n", err)
		return
	}

	// Create an algod client
	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("failed to make algod client: %s\n", err)
		return
	}

	// Broadcast the transaction to the network
	sendResponse, err := algodClient.SendRawTransaction(rawTx)
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}

	fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)
}
```

## Manipulating multisig transactions

Here, we first create a simple multisig payment transaction,
with three public identities and a threshold of 2:

```golang
addr1, _ := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
addr2, _ := types.DecodeAddress("BFRTECKTOOE7A5LHCF3TTEOH2A7BW46IYT2SX5VP6ANKEXHZYJY77SJTVM")
addr3, _ := types.DecodeAddress("47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU")
ma, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
	addr1,
	addr2,
	addr3,
})
if err != nil {
	panic("invalid multisig parameters")
}
fromAddr, _ := ma.Address()
txn, err := transaction.MakePaymentTxn(
	fromAddr.String(),
	"INSERTTOADDRESHERE",
	10,     // fee per byte
	10000,  // amount
	100000, // first valid round
	101000, // last valid round
	nil,    // note
	"",     // closeRemainderTo
	"",     // genesisID
	[]byte, // genesisHash (Cannot be empty in practice)
)
txid, txBytes, err := crypto.SignMultisigTransaction(secretKey, ma, txn)
if err != nil {
	panic("could not sign multisig transaction")
}
fmt.Printf("Made partially-signed multisig transaction with TxID %s: %x\n", txid, txBytes)

```

Now, we can write the returned bytes to disk:
```golang
_ := ioutil.WriteFile("./arbitrary_file.tx", txBytes, 0644)
```

And read them back in:
```golang
readTxBytes, _ := ioutil.ReadFile("./arbitrary_file.tx")
```

Now, we can append another signature, to hit the threshold. Note that
this SDK forces new signers to know the parameters of the multisig -
after all, we don't want to sign things without knowing the identity
of the multi-signature.
```golang
// as before
addr1, _ := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
addr2, _ := types.DecodeAddress("BFRTECKTOOE7A5LHCF3TTEOH2A7BW46IYT2SX5VP6ANKEXHZYJY77SJTVM")
addr3, _ := types.DecodeAddress("47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU")
ma, _ := crypto.MultisigAccountWithParams(1, 2, []types.Address{
	addr1,
	addr2,
	addr3,
})
// append our signature to readTxBytes
txid, twoOfThreeTxBytes, err := crypto.AppendMultisigTransaction(secretKey, ma, readTxBytes)
if err != nil {
	panic("could not append signature to multisig transaction")
}
fmt.Printf("Made 2-out-of-3 multisig transaction with TxID %s: %x\n", txid, twoOfThreeTxBytes)

```

We can also merge raw, partially-signed multisig transactions:
```golang
otherTxBytes := ... // generate another raw multisig transaction somehow
txid, mergedTxBytes, err := crypto.MergeMultisigTransactions(twoOfThreeTxBytes, otherTxBytes)
```
