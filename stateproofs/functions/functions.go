package functions

import (
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/algorand/go-algorand/crypto/stateproof"
)

func MkVerifierWithLnProvenWeight(partcom datatypes.GenericDigest, lnProvenWt uint64, strengthTarget uint64) *stateproof.Verifier {
	return stateproof.MkVerifierWithLnProvenWeight([]byte(partcom), lnProvenWt, strengthTarget)
}

func Verify(verifier *stateproof.Verifier, round types.Round, messageHash datatypes.MessageHash, stateProof *datatypes.EncodedStateProof) error {
	var decodedStateProof stateproof.StateProof
	err := msgpack.Decode(*stateProof, &decodedStateProof)
	if err != nil {
		return err
	}

	var stateProofMessageHash stateproof.MessageHash
	copy(stateProofMessageHash[:], messageHash[:])
	return verifier.Verify(uint64(round), stateProofMessageHash, &decodedStateProof)
}
