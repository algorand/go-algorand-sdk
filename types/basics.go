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
)

const masterDerivationKeyLenBytes = 32

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

func (microalgos MicroAlgos) ToAlgos() float64 {
	return float64(microalgos) / microAlgoConversionFactor
}

func ToMicroAlgos(algos float64) MicroAlgos {
	return MicroAlgos(math.Round(algos * microAlgoConversionFactor))
}
