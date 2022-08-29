package stateproofverification

import (
	"github.com/algorand/go-stateproof-verification/stateproof"
	"github.com/algorand/go-stateproof-verification/stateprooftypes"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

const strengthTarget = uint64(256)

type StateProofVerifier struct {
	stateProofVerifier *stateproof.Verifier
}

func InitializeVerifier(votersCommitment types.GenericDigest, lnProvenWeight uint64) *StateProofVerifier {
	return &StateProofVerifier{stateProofVerifier: stateproof.MkVerifierWithLnProvenWeight(stateprooftypes.GenericDigest(votersCommitment),
		lnProvenWeight, strengthTarget)}
}

func (v *StateProofVerifier) VerifyStateProofMessage(stateProof *types.EncodedStateProof, message *types.Message) error {
	messageHash := message.IntoStateProofMessageHash()

	var decodedStateProof stateproof.StateProof
	err := msgpack.Decode(*stateProof, &decodedStateProof)
	if err != nil {
		return err
	}

	return v.stateProofVerifier.Verify(message.LastAttestedRound, stateprooftypes.MessageHash(messageHash), &decodedStateProof)
}
