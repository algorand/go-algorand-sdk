package crypto

import (
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// inMemoryEd25519Signer is a simplistic implementation of an Ed25519Signer that
// keeps the public/private key pair in memory
type inMemoryEd25519Signer struct {
	sk ed25519.PrivateKey
	pk Ed25519PublicKey
}

// Ed25519Sign signs the given bytes with an ed25519 signature
func (sgnr inMemoryEd25519Signer) Ed25519Sign(toBeSigned []byte) ([]byte, error) {
	return ed25519.Sign(sgnr.sk, toBeSigned), nil
}

// Ed25519PublicKey returns the public key that should be used to verify
// the signatures performed by this signer
func (sgnr inMemoryEd25519Signer) Ed25519PublicKey() Ed25519PublicKey {
	return sgnr.pk
}

// GenerateAddressFromSK take a secret key and returns the corresponding Address
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged
func GenerateAddressFromSK(sk []byte) (types.Address, error) {
	edsk := ed25519.PrivateKey(sk)

	var a types.Address
	pk := edsk.Public()
	n := copy(a[:], []byte(pk.(ed25519.PublicKey)))
	if n != ed25519.PublicKeySize {
		return [32]byte{}, fmt.Errorf("generated public key has the wrong size, expected %d, got %d", ed25519.PublicKeySize, n)
	}
	return a, nil
}

// SKToInMemorySigner wraps an ed25519 private key with an Ed25519Signer implementation
//
// Note: usage of in-memory cryptographic APIs is discouraged
func SKToInMemorySigner(sk ed25519.PrivateKey) (Ed25519Signer, error) {
	var pk Ed25519PublicKey
	n := copy(pk[:], sk.Public().(ed25519.PublicKey))
	if n != ed25519.PublicKeySize {
		return inMemoryEd25519Signer{}, fmt.Errorf("generated public key has the wrong size, expected %d, got %d", ed25519.PublicKeySize, n)
	}

	return inMemoryEd25519Signer{sk: sk, pk: Ed25519PublicKey(pk)}, nil
}

// SignTransaction accepts a private key and a transaction, and returns the
// bytes of a signed transaction ready to be broadcasted to the network
// If the SK's corresponding address is different than the txn sender's, the SK's
// corresponding address will be assigned as AuthAddr
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519SignTransaction instead
func SignTransaction(sk ed25519.PrivateKey, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519SignTransaction(sgnr, tx)
}

// SignBytes signs the bytes and returns the signature
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519SignBytes instead
func SignBytes(sk ed25519.PrivateKey, bytesToSign []byte) (signature []byte, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519SignBytes(sgnr, bytesToSign)
}

// SignBid accepts a private key and a bid, and returns the signature of the
// bid under that key
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519SignBid instead
func SignBid(sk ed25519.PrivateKey, bid types.Bid) (signedBid []byte, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519SignBid(sgnr, bid)
}

// SignMultisigTransaction signs the given transaction, and multisig preimage, with the
// private key, returning the bytes of a signed transaction with the multisig field
// partially populated, ready to be passed to other multisig signers to sign or broadcast.
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519SignMultisigTransaction instead
func SignMultisigTransaction(sk ed25519.PrivateKey, ma MultisigAccount, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519SignMultisigTransaction(sgnr, ma, tx)
}

// AppendMultisigToLogicSig adds a new signature to multisigned LogicSig
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519AppendMultisigToLogicSig instead
func AppendMultisigToLogicSig(lsig *types.LogicSig, sk ed25519.PrivateKey) error {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return err
	}

	return Ed25519AppendMultisigToLogicSig(lsig, sgnr)
}

// AppendMultisigTransaction appends the signature corresponding to the given private key,
// returning an encoded signed multisig transaction including the signature.
// While we could compute the multisig preimage from the multisig blob, we ask the caller
// to pass it back in, to explicitly check that they know who they are signing as.
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519AppendMultisigTransaction instead
func AppendMultisigTransaction(sk ed25519.PrivateKey, ma MultisigAccount, preStxBytes []byte) (txid string, stxBytes []byte, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519AppendMultisigTransaction(sgnr, ma, preStxBytes)
}

// TealSign creates a signature compatible with ed25519verify opcode from contract address
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519TealSign instead
func TealSign(sk ed25519.PrivateKey, data []byte, contractAddress types.Address) (rawSig types.Signature, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519TealSign(sgnr, data, contractAddress)
}

// TealSignFromProgram creates a signature compatible with ed25519verify opcode from raw program bytes
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519TealSignFromProgram instead
func TealSignFromProgram(sk ed25519.PrivateKey, data []byte, program []byte) (rawSig types.Signature, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519TealSignFromProgram(sgnr, data, program)
}
