// Package signing holds the signing primitives shared by the crypto and
// transaction packages.
//
// It lets the transaction package's signers and the crypto package's
// deprecated in-memory signing functions share one implementation without
// crypto having to export scheme-specific signing functions.
package signing

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"errors"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Ed25519PublicKeySize is the size in bytes of an ed25519 public key
const Ed25519PublicKeySize = 32

// Ed25519PublicKey represents a 32 byte ed25519 public key.
//
// crypto.Ed25519PublicKey is an alias of this type, so that crypto.Ed25519Signer
// and Ed25519Signer share a method set.
type Ed25519PublicKey [Ed25519PublicKeySize]byte

// Ed25519Signer mirrors crypto.Ed25519Signer, see its documentation.
type Ed25519Signer interface {
	Ed25519Sign(toBeSigned []byte) ([]byte, error)
	Ed25519PublicKey() Ed25519PublicKey
}

// PQSigner mirrors crypto.PQSigner, see its documentation.
type PQSigner interface {
	PQSign(toBeSigned []byte) ([]byte, error)
	PQPublicKey() []byte
	PQScheme() types.PQScheme
}

// ErrInvalidSignatureReturned is returned when an Ed25519Signer produces a
// signature whose length is not ed25519.SignatureSize.
var ErrInvalidSignatureReturned = errors.New("ed25519 library returned an invalid signature")

// txidPrefix is prepended to a transaction when computing its txid
var txidPrefix = []byte("TX")

// bytesPrefix is prepended to a message when signing
var bytesPrefix = []byte("MX")

// programPrefix is prepended to a logic program when computing a hash
var programPrefix = []byte("Program")

// msigProgramPrefix is prepended to a logic program when computing a hash for a program signed by multisig
var msigProgramPrefix = []byte("MsigProgram")

// pqProgramPrefix is prepended to a logic program when computing the bytes a
// post-quantum scheme signs for a delegated LogicSig.
var pqProgramPrefix = []byte("PQProgram")

// programDataPrefix is prepended to teal sign data
var programDataPrefix = []byte("ProgData")

// TransactionBytesToSign returns the byte form of the tx that we actually sign
// and compute txID from: the canonical msgpack encoding of tx, prefixed with
// "TX" for domain separation.
func TransactionBytesToSign(tx types.Transaction) []byte {
	return bytes.Join([][]byte{txidPrefix, msgpack.Encode(tx)}, nil)
}

// TxIDFromBytesToSign computes a transaction id base32 string from the bytes
// returned by TransactionBytesToSign.
func TxIDFromBytesToSign(toBeSigned []byte) string {
	txidBytes := sha512.Sum512_256(toBeSigned)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])
}

// BytesToSign returns the bytes that are actually signed when signing an
// arbitrary message.
func BytesToSign(message []byte) []byte {
	return bytes.Join([][]byte{bytesPrefix, message}, nil)
}

// ProgramToSign returns the bytes signed when delegating a LogicSig to a single
// ed25519 account, which are also hashed to compute an escrow address.
func ProgramToSign(program []byte) []byte {
	return bytes.Join([][]byte{programPrefix, program}, nil)
}

// MsigProgramToSign returns the bytes signed when delegating a LogicSig to the
// multisig account at msigAddr.
func MsigProgramToSign(msigAddr types.Address, program []byte) []byte {
	return bytes.Join([][]byte{msigProgramPrefix, msigAddr[:], program}, nil)
}

// PQProgramToSign returns the bytes a post-quantum scheme signs when delegating
// a LogicSig to a PQ account: ("PQProgram" || address || program).
func PQProgramToSign(addr types.Address, program []byte) []byte {
	return bytes.Join([][]byte{pqProgramPrefix, addr[:], program}, nil)
}

// TealSignData returns the bytes signed by Ed25519TealSign
func TealSignData(data []byte, contractAddress types.Address) []byte {
	return bytes.Join([][]byte{programDataPrefix, contractAddress[:], data}, nil)
}

// EncodeSignedTxn encodes the SignedTxn, assigning the signing account as the
// AuthAddr when it is not the sender of the transaction.
func EncodeSignedTxn(stx types.SignedTxn, signerAddress types.Address) []byte {
	if stx.Txn.Sender != signerAddress {
		stx.AuthAddr = signerAddress
	}
	return msgpack.Encode(&stx)
}
