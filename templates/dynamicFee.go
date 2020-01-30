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

// DynamicFee template representation
type DynamicFee struct {
	ContractTemplate
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
	return makeDynamicFeeWithLease(receiver, closeRemainder, leaseString, amount, firstValid, lastValid)
}

// makeDynamicFeeWithLease is as MakeDynamicFee, but the caller can specify the lease (using b64 string)
func makeDynamicFeeWithLease(receiver, closeRemainder, lease string, amount, firstValid, lastValid uint64) (DynamicFee, error) {
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
func GetDynamicFeeTransactions(txn types.Transaction, lsig types.LogicSig, privateKey ed25519.PrivateKey, fee uint64) ([]byte, error) {
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
	feePayTxn, err := transaction.MakePaymentTxn(address.String(), txn.Sender.String(), fee, uint64(txn.Fee), uint64(txn.FirstValid), uint64(txn.LastValid), nil, "", "", genesisHash)
	if err != nil {
		return nil, err
	}
	feePayTxn.AddLease(txn.Lease, fee)

	txnGroup := []types.Transaction{feePayTxn, txn}

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

// SignDynamicFee takes in the contract bytes and returns the main transaction and signed logic needed to complete the
// transfer. These should be sent to the fee payer, who can use
// GetDynamicFeeTransactions() to update fields and create the auxiliary
// transaction.
// Parameters:
// contract - the bytearray representing the contract in question
// privateKey - the privateKey that will sign the delegated LogicSig
// genesisHash - the bytearray representing the network for the txns
func SignDynamicFee(contract []byte, privateKey ed25519.PrivateKey, genesisHash []byte) (txn types.Transaction, lsig types.LogicSig, err error) {
	ints, byteArrays, err := logic.ReadProgram(contract, nil)
	if err != nil {
		return
	}

	// Convert the byteArrays[0] to receiver
	var receiver types.Address //byteArrays[0]
	n := copy(receiver[:], byteArrays[0])
	if n != ed25519.PublicKeySize {
		err = fmt.Errorf("address generated from receiver bytes is the wrong size")
		return
	}
	// Convert the byteArrays[1] to closeRemainderTo
	var closeRemainderTo types.Address
	n = copy(closeRemainderTo[:], byteArrays[1])
	if n != ed25519.PublicKeySize {
		err = fmt.Errorf("address generated from closeRemainderTo bytes is the wrong size")
		return
	}
	contractLease := byteArrays[2]
	amount, firstValid, lastValid := ints[2], ints[3], ints[4]

	sender := types.Address{}
	copy(sender[:], privateKey[ed25519.PublicKeySize:])
	fee := uint64(0)
	txn, err = transaction.MakePaymentTxn(sender.String(), receiver.String(), fee, amount, firstValid, lastValid, nil, closeRemainderTo.String(), "", genesisHash)
	if err != nil {
		return
	}
	lease := [32]byte{}
	copy(lease[:], contractLease) // convert from []byte to [32]byte
	txn.AddLease(lease, fee)
	lsig, err = crypto.MakeLogicSig(contract, nil, privateKey, crypto.MultisigAccount{})
	return
}
