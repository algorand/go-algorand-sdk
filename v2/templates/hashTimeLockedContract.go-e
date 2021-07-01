package templates

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/logic"
	"github.com/algorand/go-algorand-sdk/types"
	"golang.org/x/crypto/sha3"
)

// HTLC template representation
type HTLC struct {
	ContractTemplate
}

// MakeHTLC allows a user to receive the Algo prior to a deadline (in terms of a round) by proving a knowledge
// of a special value or to forfeit the ability to claim, returning it to the payer.
// This contract is usually used to perform cross-chained atomic swaps
//
// More formally -
// Algos can be transferred under only two circumstances:
// 1. To receiver if hash_function(arg_0) = hash_value
// 2. To owner if txn.FirstValid > expiry_round
// ...
//
// Parameters:
// - owner : string an address that can receive the asset after the expiry round
// - receiver: string address to receive Algos
// - hashFunction : string the hash function to be used (must be either sha256 or keccak256)
// - hashImage : string the hash image in base64
// - expiryRound : uint64 the round on which the assets can be transferred back to owner
// - maxFee : uint64 the maximum fee that can be paid to the network by the account
func MakeHTLC(owner, receiver, hashFunction, hashImage string, expiryRound, maxFee uint64) (HTLC, error) {
	var referenceProgram string
	if hashFunction == "sha256" {
		referenceProgram = "ASAECAEACSYDIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5IH+DsWV/8fxTuS3BgUih1l38LUsfo9Z3KErd0gASbZBpIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITMQEiDjEQIxIQMQcyAxIQMQgkEhAxCSgSLQEpEhAxCSoSMQIlDRAREA=="
	} else if hashFunction == "keccak256" {
		referenceProgram = "ASAECAEACSYDIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5IH+DsWV/8fxTuS3BgUih1l38LUsfo9Z3KErd0gASbZBpIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITMQEiDjEQIxIQMQcyAxIQMQgkEhAxCSgSLQIpEhAxCSoSMQIlDRAREA=="
	} else {
		return HTLC{}, fmt.Errorf("invalid hash function supplied")
	}
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return HTLC{}, err
	}
	ownerAddr, err := types.DecodeAddress(owner)
	if err != nil {
		return HTLC{}, err
	}
	receiverAddr, err := types.DecodeAddress(receiver)
	if err != nil {
		return HTLC{}, err
	}
	//validate hashImage
	_, err = base64.StdEncoding.DecodeString(hashImage)
	if err != nil {
		return HTLC{}, err
	}
	var referenceOffsets = []uint64{ /*fee*/ 3 /*expiryRound*/, 6 /*receiver*/, 10 /*hashImage*/, 42 /*owner*/, 76}
	injectionVector := []interface{}{maxFee, expiryRound, receiverAddr, hashImage, ownerAddr}
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return HTLC{}, err
	}

	address := crypto.AddressFromProgram(injectedBytes)
	htlc := HTLC{
		ContractTemplate: ContractTemplate{
			address: address.String(),
			program: injectedBytes,
		},
	}
	return htlc, err
}

// SignTransactionWithHTLCUnlock accepts a transaction, such as a payment, and builds the HTLC-unlocking signature around that transaction
func SignTransactionWithHTLCUnlock(program []byte, txn types.Transaction, preImageAsBase64 string) (txid string, stx []byte, err error) {
	preImageAsArgument, err := base64.StdEncoding.DecodeString(preImageAsBase64)
	if err != nil {
		return
	}
	hashFunction := program[len(program)-15]
	_, byteArrays, err := logic.ReadProgram(program, nil)
	expectedHashImage := byteArrays[1]
	if err != nil {
		return
	}
	if hashFunction == 1 {
		sha256 := sha256.Sum256(preImageAsArgument)
		if !bytes.Equal(sha256[:], expectedHashImage) {
			err = fmt.Errorf("sha256 hash of preimage failed to match expected hash image")
			return
		}
	} else if hashFunction == 2 {
		keccak256 := sha3.NewLegacyKeccak256()
		keccak256Hash := keccak256.Sum(preImageAsArgument)
		if !bytes.Equal(keccak256Hash[:], expectedHashImage) {
			err = fmt.Errorf("keccak256 hash of preimage failed to match expected hash image")
		}
	} else {
		err = fmt.Errorf("found invalid hash function %d in contract", hashFunction)
		return
	}
	args := make([][]byte, 1)
	args[0] = preImageAsArgument
	var blankMultisig crypto.MultisigAccount
	lsig, err := crypto.MakeLogicSig(program, args, nil, blankMultisig)
	if err != nil {
		return
	}
	txn.Receiver = types.Address{} //txn must have no receiver but MakePayment et al disallow this.
	txid, stx, err = crypto.SignLogicsigTransaction(lsig, txn)
	return
}
