package types

import (
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"math"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
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
	// ApplicationCallTx allows creating, deleting, and interacting with an application
	ApplicationCallTx TxType = "appl"
	// StateProofTx records a state proof
	StateProofTx TxType = "stpf"
)

const masterDerivationKeyLenBytes = 32

// MaxTxGroupSize is max number of transactions in a single group
const MaxTxGroupSize = 16

// LogicSigMaxSize is a max TEAL program size (with args)
const LogicSigMaxSize = 1000

// LogicSigMaxCost is a max execution const of a TEAL program
const LogicSigMaxCost = 20000

// KeyStoreRootSize is the size, in bytes, of keyreg verifier
const KeyStoreRootSize = 64

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

// MerkleVerifier is a state proof
type MerkleVerifier [KeyStoreRootSize]byte

const microAlgoConversionFactor = 1e6

// ToAlgos converts amount in microAlgos to Algos
func (microalgos MicroAlgos) ToAlgos() float64 {
	return float64(microalgos) / microAlgoConversionFactor
}

// ToMicroAlgos converts amount in Algos to microAlgos
func ToMicroAlgos(algos float64) MicroAlgos {
	return MicroAlgos(math.Round(algos * microAlgoConversionFactor))
}

// FromBase64String converts a base64 string to a SignedTxn
func (signedTxn *SignedTxn) FromBase64String(b64string string) error {
	txnBytes, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		return err
	}
	err = msgpack.Decode(txnBytes, &signedTxn)
	if err != nil {
		return err
	}
	return nil
}

// FromBase64String converts a base64 string to a Block
func (block *Block) FromBase64String(b64string string) error {
	txnBytes, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		return err
	}
	err = msgpack.Decode(txnBytes, &block)
	if err != nil {
		return err
	}
	return nil
}

// DigestFromString converts a string to a Digest
func DigestFromString(str string) (d Digest, err error) {
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(str)
	if err != nil {
		return d, err
	}
	if len(decoded) != len(d) {
		msg := fmt.Sprintf(`Attempted to decode a string which was not a Digest: "%v"`, str)
		return d, errors.New(msg)
	}
	copy(d[:], decoded[:])
	return d, err
}
