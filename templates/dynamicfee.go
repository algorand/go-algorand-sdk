package templates

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
)

type DynamicFee struct {
	ContractTemplate
	leaseAsBase64 string
}

// MakeDynamicFee implements a payment transaction with an undetermined fee.
// This is delegate logic.
//
// This must be present on the first of two transactions.
//
// The first transaction should send money to this account.
// It must send an amount equal to txn.Fee.
//
// The second transaction should be from this account.
// The lease is mandatory!
//
// Parameters:
//  - receiver: the payment receiver
//  - closeTo: the account to close the payment to
//  - amount: the amount of the payment
//  - firstValid: the first valid round of the transaction
//  - lastValid: the last valid round of the transaction
//  - TMPL_LEASE: string to use for the transaction lease // TODO evan: generate and store this and provide getter if needed
func MakeDynamicFee(receiver, closeTo string, amount, firstValid, lastValid uint64) DynamicFee {
	lease := make([]byte, 32)
	crypto.RandomBytes(lease)
	leaseAsBase64 := base64.StdEncoding.EncodeToString(lease)

	return DynamicFee{leaseAsBase64: leaseAsBase64}
}

// returns a TransactionGroup containing
// a payment from
func GetTransactions(fee, firstValid, lastValid uint64, genesisHash []byte) {

}
