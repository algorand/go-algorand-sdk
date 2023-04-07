package main

import (
	"bytes"
	"crypto/ed25519"
	"log"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

var txidPrefix = []byte("TX")

func main() {
	account := crypto.GenerateAccount()

	sp := types.SuggestedParams{
		Fee:             0,
		GenesisID:       "blah",
		GenesisHash:     []byte("blah"),
		FirstRoundValid: 0,
		LastRoundValid:  1,
		MinFee:          1000,
	}

	tx1, _ := transaction.MakePaymentTxn(account.Address.String(), account.Address.String(), 100000, nil, "", sp)
	_, stxn, _ := crypto.SignTransaction(account.PrivateKey, tx1)

	// example: OFFLINE_VERIFY_SIG
	signedTxn := types.SignedTxn{}
	msgpack.Decode(stxn, &signedTxn)

	from := signedTxn.Txn.Sender[:]

	encodedTx := msgpack.Encode(signedTxn.Txn)

	msgParts := [][]byte{txidPrefix, encodedTx}
	msg := bytes.Join(msgParts, nil)

	valid := ed25519.Verify(from, msg, signedTxn.Sig[:])

	log.Printf("Valid? %t", valid)
	// example: OFFLINE_VERIFY_SIG
}
