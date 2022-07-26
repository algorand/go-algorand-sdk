package stateproofverification

import (
	"github.com/algorand/go-algorand/crypto/stateproof"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
)

const strengthTarget = uint64(256)

type StateProofVerifier struct {
	stateProofVerifier *stateproof.Verifier
}

func InitializeVerifier(votersCommitment stateprooftypes.GenericDigest, lnProvenWeight uint64) *StateProofVerifier {
	return &StateProofVerifier{stateProofVerifier: stateproof.MkVerifierWithLnProvenWeight([]byte(votersCommitment),
		lnProvenWeight, strengthTarget)}
}

func (v *StateProofVerifier) verifyStateProofMessage(stateProof *stateprooftypes.EncodedStateProof, message stateprooftypes.Message) error {
	messageHash := message.IntoStateProofMessageHash()

	var decodedStateProof stateproof.StateProof
	err := msgpack.Decode(*stateProof, &decodedStateProof)
	if err != nil {
		return err
	}

	var stateProofMessageHash stateproof.MessageHash
	copy(stateProofMessageHash[:], messageHash[:])
	return v.stateProofVerifier.Verify(message.LastAttestedRound, stateProofMessageHash, &decodedStateProof)
}
