package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"
	"log"

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
	seedAcct := accts[0]

	// example: LSIG_COMPILE
	teal, err := ioutil.ReadFile("lsig/simple.teal")
	if err != nil {
		log.Fatalf("failed to read approval program: %s", err)
	}

	result, err := algodClient.TealCompile(teal).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to compile program: %s", err)
	}

	lsigBinary, err := base64.StdEncoding.DecodeString(result.Result)
	if err != nil {
		log.Fatalf("failed to decode compiled program: %s", err)
	}
	// example: LSIG_COMPILE

	// example: LSIG_INIT
	lsig := crypto.LogicSigAccount{
		Lsig: types.LogicSig{Logic: lsigBinary, Args: nil},
	}
	// example: LSIG_INIT
	_ = lsig

	// example: LSIG_PASS_ARGS
	encodedArg := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedArg, 123)

	lsigWithArgs := crypto.LogicSigAccount{
		Lsig: types.LogicSig{Logic: lsigBinary, Args: [][]byte{encodedArg}},
	}
	// example: LSIG_PASS_ARGS
	_ = lsigWithArgs

	// seed lsig so the pay from the lsig works
	lsa, err := lsig.Address()
	if err != nil {
		log.Fatalf("failed to get lsig address: %s", err)
	}
	seedAddr := seedAcct.Address.String()
	sp, _ := algodClient.SuggestedParams().Do(context.Background())
	txn, _ := transaction.MakePaymentTxn(seedAddr, lsa.String(), 1000000, nil, "", sp)
	txid, stx, _ := crypto.SignTransaction(seedAcct.PrivateKey, txn)
	algodClient.SendRawTransaction(stx).Do(context.Background())
	transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())

	// example: LSIG_SIGN_FULL
	sp, err = algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		log.Fatalf("failed to get suggested params: %s", err)
	}

	lsigAddr, err := lsig.Address()
	if err != nil {
		log.Fatalf("failed to get lsig address: %s", err)
	}
	ptxn, err := transaction.MakePaymentTxn(lsigAddr.String(), seedAddr, 10000, nil, "", sp)
	if err != nil {
		log.Fatalf("failed to make transaction: %s", err)
	}
	txid, stxn, err := crypto.SignLogicSigAccountTransaction(lsig, ptxn)
	if err != nil {
		log.Fatalf("failed to sign transaction with lsig: %s", err)
	}
	_, err = algodClient.SendRawTransaction(stxn).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	payResult, err := transaction.WaitForConfirmation(algodClient, txid, 4, context.Background())
	if err != nil {
		log.Fatalf("failed while waiting for transaction: %s", err)
	}
	log.Printf("Lsig pay confirmed in round: %d", payResult.ConfirmedRound)
	// example: LSIG_SIGN_FULL

	// example: LSIG_DELEGATE_FULL
	// account signs the logic, and now the logic may be passed instead
	// of a signature for a transaction
	var args [][]byte
	delSig, err := crypto.MakeLogicSigAccountDelegated(lsigBinary, args, seedAcct.PrivateKey)
	if err != nil {
		log.Fatalf("failed to make delegate lsig: %s", err)
	}

	delSigPay, err := transaction.MakePaymentTxn(seedAddr, lsigAddr.String(), 10000, nil, "", sp)
	if err != nil {
		log.Fatalf("failed to make transaction: %s", err)
	}

	delTxId, delStxn, err := crypto.SignLogicSigAccountTransaction(delSig, delSigPay)
	if err != nil {
		log.Fatalf("failed to sign with delegate sig: %s", err)
	}

	_, err = algodClient.SendRawTransaction(delStxn).Do(context.Background())
	if err != nil {
		log.Fatalf("failed to send transaction: %s", err)
	}

	delPayResult, err := transaction.WaitForConfirmation(algodClient, delTxId, 4, context.Background())
	if err != nil {
		log.Fatalf("failed while waiting for transaction: %s", err)
	}

	log.Printf("Delegated Lsig pay confirmed in round: %d", delPayResult.ConfirmedRound)
	// example: LSIG_DELEGATE_FULL
}
