package templates

import (
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// Split template representation
type Split struct {
	ContractTemplate
	ratn        uint64
	ratd        uint64
	receiverOne types.Address
	receiverTwo types.Address
}

//GetSendFundsTransaction returns a group transaction array which transfer funds according to the contract's ratio
// the returned byte array is suitable for passing to SendRawTransaction
// amount: uint64 number of assets to be transferred total
func (contract Split) GetSendFundsTransaction(amount uint64, firstRound, lastRound, fee uint64, genesisHash []byte) ([]byte, error) {
	ratio := contract.ratn / contract.ratd
	amountForReceiverOne := amount * ratio
	amountForReceiverTwo := amount * (1 - ratio)
	remainder := amount - amountForReceiverOne - amountForReceiverTwo
	if remainder != 0 {
		return nil, fmt.Errorf("could not precisely divide funds between the two accounts")
	}

	from := contract.address
	tx1, err := transaction.MakePaymentTxn(from, contract.receiverOne.String(), fee, amountForReceiverOne, firstRound, lastRound, nil, "", "", genesisHash)
	if err != nil {
		return nil, err
	}
	tx2, err := transaction.MakePaymentTxn(from, contract.receiverTwo.String(), fee, amountForReceiverTwo, firstRound, lastRound, nil, "", "", genesisHash)
	if err != nil {
		return nil, err
	}
	gid, err := crypto.ComputeGroupID([]types.Transaction{tx1, tx2})
	if err != nil {
		return nil, err
	}
	tx1.Group = gid
	tx2.Group = gid

	logicSig, err := crypto.MakeLogicSig(contract.program, nil, nil, crypto.MultisigAccount{})
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
// ratn/ratd.  At least minPay must be sent to receiverOne.
// (CloseRemainderTo must be zero.)
//
// After expiryRound passes, all funds can be refunded to owner.
//
// Parameters:
//  - owner: the address to refund funds to on timeout
//  - receiverOne: the first recipient in the split account
//  - receiverTwo: the second recipient in the split account
//  - ratn: fraction of money to be paid to the first recipient (numerator)
//  - ratd: fraction of money to be paid to the first recipient (denominator)
//  - expiryRound: the round at which the account expires
//  - minPay: minimum amount to be paid out of the account
//  - maxFee: half of the maximum fee used by each split forwarding group transaction
func MakeSplit(owner, receiverOne, receiverTwo string, ratn, ratd, expiryRound, minPay, maxFee uint64) (Split, error) {
	const referenceProgram = "ASAIAQUCAAYHCAkmAyCztwQn0+DycN+vsk+vJWcsoz/b7NDS6i33HOkvTpf+YiC3qUpIgHGWE8/1LPh9SGCalSN7IaITeeWSXbfsS5wsXyC4kBQ38Z8zcwWVAym4S8vpFB/c0XC6R4mnPi9EBADsPDEQIhIxASMMEDIEJBJAABkxCSgSMQcyAxIQMQglEhAxAiEEDRAiQAAuMwAAMwEAEjEJMgMSEDMABykSEDMBByoSEDMACCEFCzMBCCEGCxIQMwAIIQcPEBA="
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return Split{}, err
	}

	var referenceOffsets = []uint64{ /*fee*/ 4 /*timeout*/, 7 /*ratn*/, 8 /*ratd*/, 9 /*minPay*/, 10 /*owner*/, 14 /*receiver1*/, 47 /*receiver2*/, 80}
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
	injectionVector := []interface{}{maxFee, expiryRound, ratn, ratd, minPay, ownerAddr, receiverOneAddr, receiverTwoAddr}
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
		ratn:        ratn,
		ratd:        ratd,
		receiverOne: receiverOneAddr,
		receiverTwo: receiverTwoAddr,
	}
	return split, err
}
