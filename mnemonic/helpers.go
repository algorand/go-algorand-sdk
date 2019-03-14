package mnemonic

import (
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/types"
)

// FromPrivateKey is a helper that converts an ed25519 private key to a
// human-readable mnemonic
func FromPrivateKey(sk ed25519.PrivateKey) (string, error) {
	seed := sk.Seed()
	return FromKey(seed)
}

// ToPrivateKey is a helper that converts a mnemonic directly to an ed25519
// private key
func ToPrivateKey(mnemonic string) (sk ed25519.PrivateKey, err error) {
	seedBytes, err := ToKey(mnemonic)
	if err != nil {
		return
	}
	return ed25519.NewKeyFromSeed(seedBytes), nil
}

// FromMasterDerivationKey is a helper that converts an MDK to a human-readable
// mnemonic
func FromMasterDerivationKey(mdk types.MasterDerivationKey) (string, error) {
	return FromKey(mdk[:])
}

// ToMasterDerivationKey is a helper that converts a mnemonic directly to a
// master derivation key
func ToMasterDerivationKey(mnemonic string) (mdk types.MasterDerivationKey, err error) {
	mdkBytes, err := ToKey(mnemonic)
	if err != nil {
		return
	}
	if len(mdkBytes) != len(mdk) {
		panic("recovered mdk is wrong length")
	}
	copy(mdk[:], mdkBytes)
	return
}
