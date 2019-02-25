# go-algorand-sdk

The Algorand golang SDK provides:

- HTTP clients for the algod (agreement) and kmd (key management) APIs
- Standalone functionality for interacting with the Algorand protocol, including transaction signing, message encoding, etc.

# Documentation

Full documentation is available [on godoc](#). You can also self-host the documentation by running `godoc -http=:8099` and visiting `http://localhost:8099/pkg/github.com/algorand/go-algorand-sdk` in your web browser.

# Package overview

In `client/`, the `algod` and `kmd` packages provide HTTP clients for their corresponding APIs. `algod` is the Algorand protocol daemon, responsible for reaching consensus with the network and participating in the Algorand protocol. You can use it to check the status of the blockchain, read a block, look at transactions, look at balances, or submit a signed transaction. `kmd` is the key management daemon. It is responsible for managing spending key material, signing transactions, and managing wallets.

`types` contains the data structures you'll use when interacting with the network, including addresses, transactions, multisig signatures, etc. Some types (like `Transaction`) have their own packages containing constructors (like `MakePaymentTxn`).

`encoding` contains the `json` and `msgpack` packages, which can be used to serialize messages for the algod/kmd APIs and the network.

`mnemonic` contains support for turning 32-byte keys into checksummed, human-readable mnemonics (and going from mnemonics back to keys).

# Quick Start

## algod client

Here is an example that creates an algod client and uses it to fetch node status information, and then a specific block.

```golang
import (
	"encoding/json"
	"fmt"

	"github.com/algorand/go-algorand/sdk/client/algod"
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
import (
	"fmt"

	"github.com/algorand/go-algorand/sdk/client/kmd"
)

// These constants represent the kmdd REST endpoint and the corresponding API
// token. You can retrieve these from the `kmd.net` and `kmd.token` files in
// the kmd data directory.
const kmdAddress = "http://localhost:7833"
const kmdToken = "51ab7f41f250aa6c1b9c873df20f4030cfa8207a93e0ba2b18348ae18c6a2ade"

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

	//Generate a new address from the wallet handle
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
// Create a kmd client
kmdClient, err := kmd.MakeClient(kmdAddress, kmdToken)
if err != nil {
	fmt.Printf("failed to make kmd client: %s\n", err)
	return
}
fmt.Println("Made a kmd client")

//Get the list of wallets
listResponse, err := kmdClient.ListWallets()
if err != nil {
	fmt.Printf("error listing wallets: %s\n", err)
	return
}

//Find our wallet name in the list
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

//Get backup phrase
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
```

To restore a wallet, convert the phrase to a key and pass it to `CreateWallet`:

```golang
backupPhrase := "first umbrella camera pass middle misery copper sniff cargo wealth predict car wise trim guide middle virus bamboo evolve letter prefer panther venue about delay"
keyBytes, err := mnemonic.ToKey(backupPhrase)
if err != nil {
	fmt.Printf("failed to get key: %s\n", err)
	return
}

var mdk types.MasterDerivationKey
copy(mdk[:], keyBytes)
_, err := kmdClient.CreateWallet("testwallet", "testpassword", kmd.DefaultWalletDriver, mdk)
if err != nil {
	fmt.Printf("error creating wallet: %s\n", err)
	return
}
```

## Signing and submitting a transaction

```golang
// Make transaction
tx, err := transaction.MakePaymentTxn(fromAddr, toAddr, 1, 100, 300, 400, nil)
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

// Create an algod client
algodClient, err := algod.MakeClient(algodAddress, algodToken)
if err != nil {
	fmt.Printf("failed to make algod client: %s\n", err)
	return
}

// Broadcast the transaction to the network
sendResponse, err := algodClient.SendRawTransaction(signResponse.SignedTransaction)
if err != nil {
	fmt.Printf("failed to send transaction: %s\n", err)
	return
}

fmt.Printf("Transaction ID: %s\n", sendResponse.TxID)
```
