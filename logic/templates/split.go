package templates

import (
	"crypto/sha512"
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/types"
)

type Split struct {
	address string
	program string
}

const referenceProgram = "ASAIAQUCAAYHCAkmAyCztwQn0+DycN+vsk+vJWcsoz/b7NDS6i33HOkvTpf+YiC3qUpIgHGWE8/1LPh9SGCalSN7IaITeeWSXbfsS5wsXyC4kBQ38Z8zcwWVAym4S8vpFB/c0XC6R4mnPi9EBADsPDEQIhIxASMMEDIEJBJAABkxCSgSMQcyAxIQMQglEhAxAiEEDRAiQAAuMwAAMwEAEjEJMgMSEDMABykSEDMBByoSEDMACCEFCzMBCCEGCxIQMwAIIQcPEBA="

var referenceOffsets = []uint64{ /*fee*/ 4 /*timeout*/, 7 /*ratn*/, 8 /*ratd*/, 9 /*minPay*/, 10 /*owner*/, 14 /*receiver1*/, 15 /*receiver2*/, 80}

// GetAddress returns the contract address
func (contract Split) GetAddress() string {
	return contract.address
}

// GetProgram returns b64-encoded version of the program
func (contract Split) GetProgram() string {
	return contract.program
}

//GetSendFundsTransaction returns a group transactions array which transfer funds according to the contract's ratio
// amount: uint64 number of assets to be transferred
// precise: handles rounding error. When False, the amount will be divided as closely as possible but one account will get
// 			slightly more. When true, returns an error.
func (contract Split) GetSendFundsTransaction(amount uint64, precise bool) ([]types.Transaction, error) {
	return nil, nil
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
//  - receiverOne: the first recipient in the split account
//  - receiverTwo: the second recipient in the split account
//  - ratn: fraction of money to be paid to the first recipient (numerator)
//  - ratd: fraction of money to be paid to the first recipient (denominator)
//  - minPay: minimum amount to be paid out of the account
//  - expiryRound: the round at which the account expires
//  - owner: the address to refund funds to on timeout
//  - maxFee: half of the maximum fee used by each split forwarding group transaction
func MakeSplit(owner, receiverOne, receiverTwo string, ratn, ratd, expiryRound, minPay, maxFee uint64) (Split, error) {
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return Split{}, err
	}
	injectionVector := []interface{}{maxFee, expiryRound, ratn, ratd, minPay, owner, receiverOne, receiverTwo} // TODO ordering
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return Split{}, err
	}
	injectedProgram := base64.StdEncoding.EncodeToString(injectedBytes)
	addressBytes := sha512.Sum512_256(injectedBytes)
	address := types.Address(addressBytes)
	return Split{address: address.String(), program: injectedProgram}, err
}
