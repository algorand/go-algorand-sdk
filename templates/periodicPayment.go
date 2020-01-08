package templates

import (
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/logic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	"golang.org/x/crypto/ed25519"
)

// PeriodicPayment template representation
type PeriodicPayment struct {
	ContractTemplate
}

// GetPeriodicPaymentWithdrawalTransaction returns a signed transaction extracting funds from the contract
// contract: the bytearray defining the contract, received from the payer
// firstValid: the first round on which the txn will be valid
// genesisHash: the hash representing the network for the txn
func GetPeriodicPaymentWithdrawalTransaction(contract []byte, firstValid uint64, genesisHash []byte) ([]byte, error) {
	address := crypto.AddressFromProgram(contract)
	ints, byteArrays, err := logic.ReadProgram(contract, nil)
	if err != nil {
		return nil, err
	}
	// Convert the byteArrays[0] to receiver
	var receiver types.Address
	n := copy(receiver[:], byteArrays[0])
	if n != ed25519.PublicKeySize {
		return nil, fmt.Errorf("address generated from receiver bytes is the wrong size")
	}
	contractLease := byteArrays[1]
	fee, period, withdrawWindow, amount := ints[1], ints[2], ints[4], ints[5]
	if firstValid%period != 0 {
		return nil, fmt.Errorf("firstValid round %d was not a multiple of the contract period %d", firstValid, period)
	}
	lastValid := firstValid + withdrawWindow

	txn, err := transaction.MakePaymentTxn(address.String(), receiver.String(), fee, amount, firstValid, lastValid, nil, "", "", genesisHash)
	if err != nil {
		return nil, err
	}
	lease := [32]byte{}
	copy(lease[:], contractLease) // convert from []byte to [32]byte
	txn.AddLease(lease, fee)

	logicSig, err := crypto.MakeLogicSig(contract, nil, nil, crypto.MultisigAccount{})
	if err != nil {
		return nil, err
	}
	_, signedTxn, err := crypto.SignLogicsigTransaction(logicSig, txn)
	return signedTxn, err
}

// MakePeriodicPayment allows some account to execute periodic withdrawal of funds.
// This is a contract account.
//
// This allows receiver to withdraw amount every
// period rounds for withdrawWindow after every multiple
// of period.
//
// After expiryRound, all remaining funds in the escrow
// are available to receiver.
//
// Parameters:
//  - receiver: address which is authorized to receive withdrawals
//  - amount: the maximum number of funds allowed for a single withdrawal
//  - withdrawWindow: the duration of a withdrawal period
//  - period: the time between a pair of withdrawal periods
//  - expiryRound: the round at which the account expires
//  - maxFee: maximum fee used by the withdrawal transaction
func MakePeriodicPayment(receiver string, amount, withdrawWindow, period, expiryRound, maxFee uint64) (PeriodicPayment, error) {
	leaseBytes := make([]byte, 32)
	crypto.RandomBytes(leaseBytes)
	leaseString := base64.StdEncoding.EncodeToString(leaseBytes)
	return makePeriodicPaymentWithLease(receiver, leaseString, amount, withdrawWindow, period, expiryRound, maxFee)
}

// makePeriodicPaymentWithLease is as MakePeriodicPayment, but the caller can specify the lease (using b64 string)
func makePeriodicPaymentWithLease(receiver, lease string, amount, withdrawWindow, period, expiryRound, maxFee uint64) (PeriodicPayment, error) {
	const referenceProgram = "ASAHAQYFAAQDByYCIAECAwQFBgcIAQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIIJKvkYTkEzwJf2arzJOxERsSogG9nQzKPkpIoc4TzPTFMRAiEjEBIw4QMQIkGCUSEDEEIQQxAggSEDEGKBIQMQkyAxIxBykSEDEIIQUSEDEJKRIxBzIDEhAxAiEGDRAxCCUSEBEQ"
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return PeriodicPayment{}, err
	}
	receiverAddr, err := types.DecodeAddress(receiver)
	if err != nil {
		return PeriodicPayment{}, err
	}

	var referenceOffsets = []uint64{ /*fee*/ 4 /*period*/, 5 /*withdrawWindow*/, 7 /*amount*/, 8 /*expiryRound*/, 9 /*lease*/, 12 /*receiver*/, 46}
	injectionVector := []interface{}{maxFee, period, withdrawWindow, amount, expiryRound, lease, receiverAddr}
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return PeriodicPayment{}, err
	}

	address := crypto.AddressFromProgram(injectedBytes)
	periodicPayment := PeriodicPayment{
		ContractTemplate: ContractTemplate{
			address: address.String(),
			program: injectedBytes,
		},
	}
	return periodicPayment, err
}
