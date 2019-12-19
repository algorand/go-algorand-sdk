package templates

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
	"golang.org/x/crypto/ed25519"
)

// DynamicFee template representation
type DynamicFee struct {
	ContractTemplate
	receiver       types.Address
	closeRemainder types.Address
	firstValid     uint64
	lastValid      uint64
	amount         uint64
	leaseBase64    string
}

// MakeDynamicFee contract allows you to create a transaction without
// specifying the fee. The fee will be determined at the moment of
// transfer.
//
// Parameters:
//  - receiver: address which is authorized to receive withdrawals
//  - amount: the maximum number of funds allowed for a single withdrawal
//  - withdrawWindow: the duration of a withdrawal period
//  - period: the time between a pair of withdrawal periods
//  - expiryRound: the round at which the account expires
//  - maxFee: maximum fee used by the withdrawal transaction
func MakeDynamicFee(receiver, closeRemainder string, amount, firstValid, lastValid uint64) (DynamicFee, error) {
	leaseBytes := make([]byte, 32)
	crypto.RandomBytes(leaseBytes)
	leaseString := base64.StdEncoding.EncodeToString(leaseBytes)
	return MakeDynamicFeeWithLease(receiver, closeRemainder, leaseString, amount, firstValid, lastValid)
}

// MakeDynamicFeeWithLease is as MakeDynamicFee, but the caller can specify the lease (using b64 string)
func MakeDynamicFeeWithLease(receiver, closeRemainder, lease string, amount, firstValid, lastValid uint64) (DynamicFee, error) {
	const referenceProgram = "ASAFAgEHBgUmAyD+vKC7FEpaTqe0OKRoGsgObKEFvLYH/FZTJclWlfaiEyDmmpYeby1feshmB5JlUr6YI17TM2PKiJGLuck4qRW2+SB/g7Flf/H8U7ktwYFIodZd/C1LH6PWdyhK3dIAEm2QaTIEIhIzABAjEhAzAAcxABIQMwAIMQESEDEWIxIQMRAjEhAxBygSEDEJKRIQMQgkEhAxAiUSEDEEIQQSEDEGKhIQ"
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return DynamicFee{}, err
	}
	receiverAddr, err := types.DecodeAddress(receiver)
	if err != nil {
		return DynamicFee{}, err
	}
	closeRemainderAddr, err := types.DecodeAddress(closeRemainder)
	if err != nil {
		return DynamicFee{}, err
	}

	var referenceOffsets = []uint64{ /*amount*/ 5 /*firstValid*/, 6 /*lastValid*/, 7 /*receiver*/, 11 /*closeRemainder*/, 44 /*lease*/, 76}
	injectionVector := []interface{}{amount, firstValid, lastValid, receiverAddr, closeRemainderAddr, lease}
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return DynamicFee{}, err
	}

	address := crypto.AddressFromProgram(injectedBytes)
	dynamicFee := DynamicFee{
		ContractTemplate: ContractTemplate{
			address: address.String(),
			program: injectedBytes,
		},
		receiver:       receiverAddr,
		closeRemainder: closeRemainderAddr,
		firstValid:     firstValid,
		lastValid:      lastValid,
		amount:         amount,
		leaseBase64:    lease,
	}
	return dynamicFee, err
}

// GetDynamicFeeTransactions creates and signs the secondary dynamic fee transaction, updates
// transaction fields, and signs as the fee payer; it returns both
// transactions as bytes suitable for sendRaw.
// Parameters:
// txn - main transaction from payer
// lsig - the signed logic received from the payer
// privateKey - the private key for the account that pays the fee
// fee - fee per byte for both transactions
// firstValid - first protocol round on which both transactions will be valid
// lastValid - last protocol round on which both transactions will be valid
func GetDynamicFeeTransactions(txn types.Transaction, lsig types.LogicSig, privateKey ed25519.PrivateKey, fee, firstValid, lastValid uint64) ([]byte, error) {
	txn.FirstValid = types.Round(firstValid)
	txn.LastValid = types.Round(lastValid)
	txn.Fee = types.MicroAlgos(fee)
	eSize, err := transaction.EstimateSize(txn)
	if err != nil {
		return nil, err
	}
	txn.Fee = types.MicroAlgos(eSize * fee)

	if txn.Fee < transaction.MinTxnFee {
		txn.Fee = transaction.MinTxnFee
	}

	address := types.Address{}
	copy(address[:], privateKey[ed25519.PublicKeySize:])
	genesisHash := make([]byte, 32)
	copy(genesisHash[:], txn.GenesisHash[:])
	feePayTxn, err := transaction.MakePaymentTxn(address.String(), txn.Sender.String(), fee, uint64(txn.Fee), firstValid, lastValid, nil, "", "", genesisHash)
	if err != nil {
		return nil, err
	}
	feePayTxn.AddLease(txn.Lease, fee)

	txnGroup := []types.Transaction{txn, feePayTxn}

	updatedTxns, err := transaction.AssignGroupID(txnGroup, "")

	_, stx1Bytes, err := crypto.SignLogicsigTransaction(lsig, updatedTxns[0])
	if err != nil {
		return nil, err
	}
	_, stx2Bytes, err := crypto.SignTransaction(privateKey, updatedTxns[1])
	if err != nil {
		return nil, err
	}
	return append(stx1Bytes, stx2Bytes...), nil
}

// SignDynamicFee returns the main transaction and signed logic needed to complete the
// transfer. These should be sent to the fee payer, who can use
// GetDynamicFeeTransactions() to update fields and create the auxiliary
// transaction.
func (c DynamicFee) SignDynamicFee(privateKey ed25519.PrivateKey, genesisHash []byte) (txn types.Transaction, lsig types.LogicSig, err error) {
	sender := types.Address{}
	copy(sender[:], privateKey[ed25519.PublicKeySize:])
	fee := uint64(0)
	txn, err = transaction.MakePaymentTxn(sender.String(), c.receiver.String(), fee, c.amount, c.firstValid, c.lastValid, nil, c.closeRemainder.String(), "", genesisHash)
	if err != nil {
		return
	}
	lsig, err = crypto.MakeLogicSig(c.GetProgram(), nil, privateKey, crypto.MultisigAccount{})
	return
}
