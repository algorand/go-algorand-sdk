package types

// TxType identifies the type of the transaction
type TxType string

const (
	// PaymentTx is the TxType for payment transactions
	PaymentTx TxType = "pay"
	// KeyRegistrationTx is the TxType for key registration transactions
	KeyRegistrationTx TxType = "keyreg"
)

// KeysLenBytes is the length of the participation public key, VRF public key,
// and master derivation key.
const KeysLenBytes = 32

// MicroAlgos are the base unit of currency in Algorand
type MicroAlgos uint64

// Round represents a round of the Algorand consensus protocol
type Round uint64

// VotePK is the participation public key used in key registration transactions
type VotePK [KeysLenBytes]byte

// VRFPK is the VRF public key used in key registration transactions
type VRFPK [KeysLenBytes]byte

// MasterDerivationKey is the secret key used to derive keys in wallets
type MasterDerivationKey [KeysLenBytes]byte

// Digest is a SHA512_256 hash
type Digest [hashLenBytes]byte

const microAlgoConversionFactor = 1e6

func (microalgos MicroAlgos) ToAlgos() uint64 {
	return uint64(microalgos) / microAlgoConversionFactor
}

func ToMicroAlgos(algos uint64) MicroAlgos {
	return MicroAlgos(algos * microAlgoConversionFactor)
}
