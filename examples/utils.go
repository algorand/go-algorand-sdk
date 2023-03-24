package examples

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

var (
	ALGOD_ADDRESS = "http://localhost"
	ALGOD_PORT    = "4001"
	ALGOD_URL     = ""
	ALGOD_TOKEN   = strings.Repeat("a", 64)

	INDEXER_ADDRESS = "http://localhost"
	INDEXER_PORT    = "8980"
	INDEXER_TOKEN   = strings.Repeat("a", 64)
	INDEXER_URL     = ""

	KMD_ADDRESS = "http://localhost"
	KMD_PORT    = "4002"
	KMD_TOKEN   = strings.Repeat("a", 64)
	KMD_URL     = ""

	KMD_WALLET_NAME     = "unencrypted-default-wallet"
	KMD_WALLET_PASSWORD = ""
)

func init() {
	if aport, ok := os.LookupEnv("ALGOD_PORT"); ok {
		ALGOD_PORT = aport
	}
	ALGOD_URL = fmt.Sprintf("%s:%s", ALGOD_ADDRESS, ALGOD_PORT)

	if iport, ok := os.LookupEnv("INDEXER_PORT"); ok {
		INDEXER_PORT = iport
	}
	INDEXER_URL = fmt.Sprintf("%s:%s", INDEXER_ADDRESS, INDEXER_PORT)

	if kport, ok := os.LookupEnv("KMD_PORT"); ok {
		KMD_PORT = kport
	}
	KMD_URL = fmt.Sprintf("%s:%s", KMD_ADDRESS, KMD_PORT)
}

func GetAlgodClient() *algod.Client {
	algodClient, err := algod.MakeClient(
		ALGOD_URL,
		ALGOD_TOKEN,
	)

	if err != nil {
		log.Fatalf("Failed to create algod client: %s", err)
	}

	return algodClient
}

func GetKmdClient() kmd.Client {
	kmdClient, err := kmd.MakeClient(
		KMD_URL,
		KMD_TOKEN,
	)

	if err != nil {
		log.Fatalf("Failed to create kmd client: %s", err)
	}

	return kmdClient
}

func GetIndexerClient() *indexer.Client {
	indexerClient, err := indexer.MakeClient(
		INDEXER_URL,
		INDEXER_TOKEN,
	)

	if err != nil {
		log.Fatalf("Failed to create indexer client: %s", err)
	}

	return indexerClient
}

func GetSandboxAccounts() ([]crypto.Account, error) {
	client := GetKmdClient()

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

func CompileTeal(algodClient *algod.Client, path string) []byte {
	teal, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read approval program: %s", err)
	}

	result, err := algodClient.TealCompile(teal).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile program: %s", err)
	}

	bin, err := base64.StdEncoding.DecodeString(result.Result)
	if err != nil {
		log.Fatalf("failed to decode compiled program: %s", err)
	}
	return bin
}

func DeployApp(algodClient *algod.Client, creator crypto.Account) uint64 {

	var (
		approvalBinary = make([]byte, 1000)
		clearBinary    = make([]byte, 1000)
	)

	// Compile approval program
	approvalTeal, err := ioutil.ReadFile("calculator/approval.teal")
	if err != nil {
		log.Fatalf("failed to read approval program: %s", err)
	}

	approvalResult, err := algodClient.TealCompile(approvalTeal).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile program: %s", err)
	}

	_, err = base64.StdEncoding.Decode(approvalBinary, []byte(approvalResult.Result))
	if err != nil {
		log.Fatalf("failed to decode compiled program: %s", err)
	}

	// Compile clear program
	clearTeal, err := ioutil.ReadFile("calculator/clear.teal")
	if err != nil {
		log.Fatalf("failed to read clear program: %s", err)
	}

	clearResult, err := algodClient.TealCompile(clearTeal).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile program: %s", err)
	}

	_, err = base64.StdEncoding.Decode(clearBinary, []byte(clearResult.Result))
	if err != nil {
		log.Fatalf("failed to decode compiled program: %s", err)
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
		log.Fatalf("error waiting for confirmation:  %s", err)
	}

	return confirmedTxn.ApplicationIndex
}
