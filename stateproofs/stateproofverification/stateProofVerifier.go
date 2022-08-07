package stateproofverification

import (
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
)

const strengthTarget = uint64(256)

type StateProofVerifier struct {
	stateProofVerifier *Verifier
}

func InitializeVerifier(votersCommitment stateprooftypes.GenericDigest, lnProvenWeight uint64) *StateProofVerifier {
	return &StateProofVerifier{stateProofVerifier: MkVerifierWithLnProvenWeight(votersCommitment,
		lnProvenWeight, strengthTarget)}
}

func (v *StateProofVerifier) VerifyStateProofMessage(stateProof *stateprooftypes.EncodedStateProof, message stateprooftypes.Message) error {
	messageHash := message.IntoStateProofMessageHash()

	var decodedStateProof StateProof
	err := msgpack.Decode(*stateProof, &decodedStateProof)
	if err != nil {
		return err
	}

	var stateProofMessageHash stateprooftypes.MessageHash
	copy(stateProofMessageHash[:], messageHash[:])
	return v.stateProofVerifier.Verify(message.LastAttestedRound, stateProofMessageHash, &decodedStateProof)
}
