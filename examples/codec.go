package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"log"
	"os"

	"github.com/algorand/go-algorand-sdk/v2/abi"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func main() {
	algodClient := getAlgodClient()
	accts, err := getSandboxAccounts()
	if err != nil {
		log.Fatalf("failed to get sandbox accounts: %s", err)
	}

	acct1 := accts[0]

	// example: CODEC_TRANSACTION_UNSIGNED
	// Error handling omitted for brevity
	sp, _ := algodClient.SuggestedParams().Do(context.Background())
	ptxn, _ := transaction.MakePaymentTxn(
		acct1.Address.String(), acct1.Address.String(), 10000, nil, "", sp,
	)

	// Encode the txn as bytes,
	// if sending over the wire (like to a frontend) it should also be b64 encoded
	encodedTxn := msgpack.Encode(ptxn)
	os.WriteFile("pay.txn", encodedTxn, 0655)

	var recoveredPayTxn = types.Transaction{}

	msgpack.Decode(encodedTxn, &recoveredPayTxn)
	log.Printf("%+v", recoveredPayTxn)
	// example: CODEC_TRANSACTION_UNSIGNED

	// example: CODEC_TRANSACTION_SIGNED
	// Assuming we already have a pay transaction `ptxn`

	// Sign the transaction
	_, signedTxn, err := crypto.SignTransaction(acct1.PrivateKey, ptxn)
	if err != nil {
		log.Fatalf("failed to sign transaction: %s", err)
	}

	// Save the signed transaction to file
	os.WriteFile("pay.stxn", signedTxn, 0644)

	signedPayTxn := types.SignedTxn{}
	err = msgpack.Decode(signedTxn, &signedPayTxn)
	if err != nil {
		log.Fatalf("failed to decode signed transaction: %s", err)
	}
	// example: CODEC_TRANSACTION_SIGNED

	// cleanup
	os.Remove("pay.stxn")
	os.Remove("pay.txn")

	// example: CODEC_ADDRESS
	address := "4H5UNRBJ2Q6JENAXQ6HNTGKLKINP4J4VTQBEPK5F3I6RDICMZBPGNH6KD4"
	pk, _ := types.DecodeAddress(address)
	addr := pk.String()
	// example: CODEC_ADDRESS
	_ = addr

	// example: CODEC_BASE64
	encoded := "SGksIEknbSBkZWNvZGVkIGZyb20gYmFzZTY0"
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	reencoded := base64.StdEncoding.EncodeToString(decoded)
	// example: CODEC_BASE64
	_ = reencoded

	// example: CODEC_UINT64
	val := 1337
	encodedInt := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedInt, uint64(val))

	decodedInt := binary.BigEndian.Uint64(encodedInt)
	// decodedInt == val
	// example: CODEC_UINT64
	_ = decodedInt

	// example: CODEC_ABI
	tupleCodec, _ := abi.TypeOf("(string,string)")

	tupleVal := []string{"hello", "world"}
	encodedTuple, _ := tupleCodec.Encode(tupleVal)
	log.Printf("%x", encodedTuple)

	decodedTuple, _ := tupleCodec.Decode(encodedTuple)
	log.Printf("%v", decodedTuple) // [hello world]

	arrCodec, _ := abi.TypeOf("uint64[]")
	arrVal := []uint64{1, 2, 3, 4, 5}
	encodedArr, _ := arrCodec.Encode(arrVal)
	log.Printf("%x", encodedArr)

	decodedArr, _ := arrCodec.Decode(encodedArr)
	log.Printf("%v", decodedArr) // [1 2 3 4 5]
	// example: CODEC_ABI
}
