package stateproof

import (
	crypto "github.com/algorand/go-algorand-sdk/algod-crypto"
	"github.com/algorand/go-algorand-sdk/algod-crypto/merklearray"
)

// Verify is a wrapper method for algod-crypto Verify method. This verifies a state proof
func Verify(root crypto.GenericDigest, elems map[uint64]crypto.Hashable, proof *merklearray.Proof) error {
	return merklearray.Verify(root,elems,proof)
}
