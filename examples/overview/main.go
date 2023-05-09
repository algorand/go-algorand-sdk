package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/examples"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

func main() {

	// example: ALGOD_CREATE_CLIENT
	// Create a new algod client, configured to connect to out local sandbox
	var algodAddress = "http://localhost:4001"
	var algodToken = strings.Repeat("a", 64)
	algodClient, _ := algod.MakeClient(
		algodAddress,
		algodToken,
	)

	// Or, if necessary, pass alternate headers

	var algodHeader common.Header
	algodHeader.Key = "X-API-Key"
	algodHeader.Value = algodToken
	algodClientWithHeaders, _ := algod.MakeClientWithHeaders(
		algodAddress,
		algodToken,
		[]*common.Header{&algodHeader},
	)

	// Or, for better performance, pass a custom Transport, in this case allowing
	// up to 100 simultaneous connections to the same host (ie: an algod node)
	// Clone Go's default transport settings but increase connection values
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.MaxIdleConns = 100
	customTransport.MaxConnsPerHost = 100
	customTransport.MaxIdleConnsPerHost = 100

	algodClientWithCustomTransport, _ := algod.MakeClientWithTransport(
		algodAddress,
		algodToken,
		nil, // accepts additional headers like MakeClientWithHeaders
		customTransport,
	)
	// example: ALGOD_CREATE_CLIENT

	_ = algodClientWithCustomTransport
	_ = algodClientWithHeaders
	_ = algodClient

	// Override with the version that has correct port
	algodClient = examples.GetAlgodClient()

	accts, err := examples.GetSandboxAccounts()
	if err != nil {
		log.Fatalf("failed to get accounts: %s", err)
	}
	acct := accts[0]

	// example: ALGOD_FETCH_ACCOUNT_INFO
	acctInfo, err := algodClient.AccountInformation(acct.Address.String()).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch account info: %s", err)
	}
	log.Printf("Account balance: %d microAlgos", acctInfo.Amount)
	// example: ALGOD_FETCH_ACCOUNT_INFO

	// example: TRANSACTION_PAYMENT_CREATE
	sp, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get suggested params: %s", err)
	}
	// payment from account to itself
	ptxn, err := transaction.MakePaymentTxn(acct.Address.String(), acct.Address.String(), 100000, nil, "", sp)
	if err != nil {
		log.Fatalf("failed creating transaction: %s", err)
	}
	// example: TRANSACTION_PAYMENT_CREATE

	// example: TRANSACTION_PAYMENT_SIGN
	_, sptxn, err := crypto.SignTransaction(acct.PrivateKey, ptxn)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}
	// example: TRANSACTION_PAYMENT_SIGN

	// example: TRANSACTION_PAYMENT_SUBMIT
	pendingTxID, err := algodClient.SendRawTransaction(sptxn).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}
	confirmedTxn, err := transaction.WaitForConfirmation(algodClient, pendingTxID, 4, context.Background())
	if err != nil {
		fmt.Printf("Error waiting for confirmation on txID: %s\n", pendingTxID)
		return
	}
	fmt.Printf("Confirmed Transaction: %s in Round %d\n", pendingTxID, confirmedTxn.ConfirmedRound)
	// example: TRANSACTION_PAYMENT_SUBMIT

	// example: SP_MIN_FEE
	suggestedParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("failed to %s", err)
	}
	log.Printf("Min fee from suggested params: %d", suggestedParams.MinFee)
	// example: SP_MIN_FEE

	// example: CONST_MIN_FEE
	log.Printf("Min fee const: %d", transaction.MinTxnFee)
	// example: CONST_MIN_FEE

	// example: TRANSACTION_FEE_OVERRIDE
	// by using fee pooling and setting our fee to 2x min tx fee
	// we can cover the fee for another transaction in the group
	sp.Fee = 2 * transaction.MinTxnFee
	sp.FlatFee = true
	// ...
	// example: TRANSACTION_FEE_OVERRIDE

}
func exampleAlgod() {
}
