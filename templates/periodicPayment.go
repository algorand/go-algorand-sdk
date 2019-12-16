package templates

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

// PeriodicPayment template representation
type PeriodicPayment struct {
	ContractTemplate
	amount      uint64
	receiver    types.Address
	leaseBase64 string
}

// GetWithdrawalTransaction returns a signed transaction extracting funds from the contract
// fee: the fee to pay in microAlgos
// firstValid: the first round on which the txn will be valid
// lastValid: the final round on which the txn will be valid
// genesisHash: the hash representing the network for the txn
func (c PeriodicPayment) GetWithdrawalTransaction(fee, firstValid, lastValid uint64, genesisHash []byte) ([]byte, error) {
	txn, err := transaction.MakePaymentTxn(c.GetAddress(), c.receiver.String(), fee, c.amount, firstValid, lastValid, nil, "", "", genesisHash)
	if err != nil {
		return nil, err
	}

	leaseBytes, err := base64.StdEncoding.DecodeString(c.leaseBase64)
	if err != nil {
		return nil, err
	}
	lease := [32]byte{}
	copy(lease[:], leaseBytes)
	txn.AddLease(lease, fee)

	logicSig, err := crypto.MakeLogicSig(c.program, nil, nil, crypto.MultisigAccount{})
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
	return MakePeriodicPaymentWithLease(receiver, leaseString, amount, withdrawWindow, period, expiryRound, maxFee)
}

// MakePeriodicPaymentWithLease is as MakePeriodicPayment, but the caller can specify the lease (using b64 string)
func MakePeriodicPaymentWithLease(receiver, lease string, amount, withdrawWindow, period, expiryRound, maxFee uint64) (PeriodicPayment, error) {
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
		amount:      amount,
		leaseBase64: lease,
		receiver:    receiverAddr,
	}
	return periodicPayment, err
}
