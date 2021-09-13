package templates

import (
	"encoding/base64"
	"fmt"
	"math"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/logic"
	"github.com/algorand/go-algorand-sdk/types"
)

// Split template representation
//
// Deprecated: Use TealCompile source compilation instead.
type Split struct {
	ContractTemplate
	rat1        uint64
	rat2        uint64
	receiverOne types.Address
	receiverTwo types.Address
}

//GetSplitFundsTransaction returns a group transaction array which transfer funds according to the contract's ratio
// the returned byte array is suitable for passing to SendRawTransaction
// contract: the bytecode of the contract to be used
// amount: uint64 total number of algos to be transferred (payment1_amount + payment2_amount)
// params: is typically received from algod, it defines common-to-all-txns arguments like fee and validity period
//
// Deprecated: Use TealCompile source compilation instead.
func GetSplitFundsTransaction(contract []byte, amount uint64, params types.SuggestedParams) ([]byte, error) {
	ints, byteArrays, err := logic.ReadProgram(contract, nil)
	if err != nil {
		return nil, err
	}
	rat1 := ints[6]
	rat2 := ints[5]
	// Convert the byteArrays[0] to receiver
	var receiverOne types.Address //byteArrays[0]
	n := copy(receiverOne[:], byteArrays[1])
	if n != ed25519.PublicKeySize {
		err = fmt.Errorf("address generated from receiver bytes is the wrong size")
		return nil, err
	}
	// Convert the byteArrays[2] to receiverTwo
	var receiverTwo types.Address
	n = copy(receiverTwo[:], byteArrays[2])
	if n != ed25519.PublicKeySize {
		err = fmt.Errorf("address generated from closeRemainderTo bytes is the wrong size")
		return nil, err
	}

	ratio := float64(rat2) / float64(rat1)
	amountForReceiverOneFloat := float64(amount) / (1 + ratio)
	amountForReceiverOne := uint64(math.Round(amountForReceiverOneFloat))
	amountForReceiverTwo := amount - amountForReceiverOne
	if rat2*amountForReceiverOne != rat1*amountForReceiverTwo {
		err = fmt.Errorf("could not split funds in a way that satisfied the contract ratio (%d * %d != %d * %d)", rat2, amountForReceiverOne, rat1, amountForReceiverTwo)
		return nil, err
	}

	from := crypto.AddressFromProgram(contract)
	tx1, err := future.MakePaymentTxn(from.String(), receiverOne.String(), amountForReceiverOne, nil, "", params)
	if err != nil {
		return nil, err
	}
	tx2, err := future.MakePaymentTxn(from.String(), receiverTwo.String(), amountForReceiverTwo, nil, "", params)
	if err != nil {
		return nil, err
	}
	gid, err := crypto.ComputeGroupID([]types.Transaction{tx1, tx2})
	if err != nil {
		return nil, err
	}
	tx1.Group = gid
	tx2.Group = gid

	logicSig, err := crypto.MakeLogicSig(contract, nil, nil, crypto.MultisigAccount{})
	if err != nil {
		return nil, err
	}
	_, stx1, err := crypto.SignLogicsigTransaction(logicSig, tx1)
	if err != nil {
		return nil, err
	}
	_, stx2, err := crypto.SignLogicsigTransaction(logicSig, tx2)
	if err != nil {
		return nil, err
	}

	var signedGroup []byte
	signedGroup = append(signedGroup, stx1...)
	signedGroup = append(signedGroup, stx2...)

	return signedGroup, err
}

// MakeSplit splits money sent to some account to two recipients at some ratio.
// This is a contract account.
//
// This allows either a two-transaction group, for executing a
// split, or single transaction, for closing the account.
//
// Withdrawals from this account are allowed as a group transaction which
// sends receiverOne and receiverTwo amounts with exactly the ratio of
// rat1/rat2.  At least minPay must be sent to receiverOne.
// (CloseRemainderTo must be zero.)
//
// After expiryRound passes, all funds can be refunded to owner.
//
// Split ratio:
// firstRecipient_amount * rat2 == secondRecipient_amount * rat1
// or phrased another way
// firstRecipient_amount == secondRecipient_amount * (rat1/rat2)
//
// Parameters:
//  - owner: the address to refund funds to on timeout
//  - receiverOne: the first recipient in the split account
//  - receiverTwo: the second recipient in the split account
//  - rat1: fraction determines resource split ratio
//  - rat2: fraction determines resource split ratio
//  - expiryRound: the round at which the account expires
//  - minPay: minimum amount to be paid out of the account to receiverOne
//  - maxFee: half of the maximum fee used by each split forwarding group transaction
//
// Deprecated: Use TealCompile source compilation instead.
func MakeSplit(owner, receiverOne, receiverTwo string, rat1, rat2, expiryRound, minPay, maxFee uint64) (Split, error) {
	const referenceProgram = "ASAIAQUCAAYHCAkmAyCztwQn0+DycN+vsk+vJWcsoz/b7NDS6i33HOkvTpf+YiC3qUpIgHGWE8/1LPh9SGCalSN7IaITeeWSXbfsS5wsXyC4kBQ38Z8zcwWVAym4S8vpFB/c0XC6R4mnPi9EBADsPDEQIhIxASMMEDIEJBJAABkxCSgSMQcyAxIQMQglEhAxAiEEDRAiQAAuMwAAMwEAEjEJMgMSEDMABykSEDMBByoSEDMACCEFCzMBCCEGCxIQMwAIIQcPEBA="
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return Split{}, err
	}
	var referenceOffsets = []uint64{ /*fee*/ 4 /*timeout*/, 7 /*rat2*/, 8 /*rat1*/, 9 /*minPay*/, 10 /*owner*/, 14 /*receiver1*/, 47 /*receiver2*/, 80}
	ownerAddr, err := types.DecodeAddress(owner)
	if err != nil {
		return Split{}, err
	}
	receiverOneAddr, err := types.DecodeAddress(receiverOne)
	if err != nil {
		return Split{}, err
	}
	receiverTwoAddr, err := types.DecodeAddress(receiverTwo)
	if err != nil {
		return Split{}, err
	}
	injectionVector := []interface{}{maxFee, expiryRound, rat2, rat1, minPay, ownerAddr, receiverOneAddr, receiverTwoAddr}
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return Split{}, err
	}

	address := crypto.AddressFromProgram(injectedBytes)
	split := Split{
		ContractTemplate: ContractTemplate{
			address: address.String(),
			program: injectedBytes,
		},
		rat1:        rat1,
		rat2:        rat2,
		receiverOne: receiverOneAddr,
		receiverTwo: receiverTwoAddr,
	}
	return split, err
}
