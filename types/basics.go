package types

import (
	"encoding/base32"
	"encoding/binary"

	"golang.org/x/crypto/ed25519"
)

// TxType identifies the type of the transaction
type TxType string

const (
	// PaymentTx is the TxType for payment transactions
	PaymentTx TxType = "pay"
	// KeyRegistrationTx is the TxType for key registration transactions
	KeyRegistrationTx TxType = "keyreg"
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

// String returns the digest in a human-readable Base32 string
func (d Digest) String() string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(d[:])
}

// TrimUint64 returns the top 64 bits of the digest and converts to uint64
func (d Digest) TrimUint64() uint64 {
	return binary.LittleEndian.Uint64(d[:8])
}

// IsZero return true if the digest contains only zeros, false otherwise
func (d Digest) IsZero() bool {
	return d == Digest{}
}

const microAlgoConversionFactor = 1e6

// ToAlgos converts micro algos into algos
func (microalgos MicroAlgos) ToAlgos() uint64 {
	return uint64(microalgos) / microAlgoConversionFactor
}

// ToMicroAlgos converts algos into MicroAlgos
func ToMicroAlgos(algos uint64) MicroAlgos {
	return MicroAlgos(algos * microAlgoConversionFactor)
}
