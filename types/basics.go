package types

import (
	"math"

	"golang.org/x/crypto/ed25519"
)

// TxType identifies the type of the transaction
type TxType string

const (
	// PaymentTx is the TxType for payment transactions
	PaymentTx TxType = "pay"
	// KeyRegistrationTx is the TxType for key registration transactions
	KeyRegistrationTx TxType = "keyreg"
	// AssetConfigTx creates, re-configures, or destroys an asset
	AssetConfigTx TxType = "acfg"
	// AssetTransferTx transfers assets between accounts (optionally closing)
	AssetTransferTx TxType = "axfer"
	// AssetFreezeTx changes the freeze status of an asset
	AssetFreezeTx TxType = "afrz"
)

const masterDerivationKeyLenBytes = 32

// MaxTxGroupSize is max number of transactions in a single group
const MaxTxGroupSize = 16

// LogicSigMaxSize is a max TEAL program size (with args)
const LogicSigMaxSize = 1000

// LogicSigMaxCost is a max execution const of a TEAL program
const LogicSigMaxCost = 20000

// MicroAlgos are the base unit of currency in Algorand
type MicroAlgos uint64

// Round represents a round of the Algorand consensus protocol
type Round uint64

// VotePK is the participation public key used in key registration transactions
type VotePK [ed25519.PublicKeySize]byte

// VRFPK is the VRF public key used in key registration transactions
type VRFPK [ed25519.PublicKeySize]byte

// MasterDerivationKey is the secret key used to derive keys in wallets
type MasterDerivationKey [masterDerivationKeyLenBytes]byte

// Digest is a SHA512_256 hash
type Digest [hashLenBytes]byte

const microAlgoConversionFactor = 1e6

// ToAlgos converts amount in microAlgos to Algos
func (microalgos MicroAlgos) ToAlgos() float64 {
	return float64(microalgos) / microAlgoConversionFactor
}

// ToMicroAlgos converts amount in Algos to microAlgos
func ToMicroAlgos(algos float64) MicroAlgos {
	return MicroAlgos(math.Round(algos * microAlgoConversionFactor))
}
