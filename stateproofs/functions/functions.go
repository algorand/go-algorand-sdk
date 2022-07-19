package functions

import (
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/stateproofs/datatypes"
	"github.com/algorand/go-algorand/crypto/stateproof"
)

const strengthTarget = uint64(256)

type StateProofVerifier struct {
	stateProofVerifier *stateproof.Verifier
}

func initializeVerifier(genesisVotersCommitment datatypes.GenericDigest, genesisLnProvenWeight uint64) *StateProofVerifier {
	return &StateProofVerifier{stateProofVerifier: stateproof.MkVerifierWithLnProvenWeight([]byte(genesisVotersCommitment),
		genesisLnProvenWeight, strengthTarget)}
}

func (v *StateProofVerifier) advanceVerifier(message datatypes.Message) {
	v.stateProofVerifier = stateproof.MkVerifierWithLnProvenWeight(message.VotersCommitment, message.LnProvenWeight, strengthTarget)
}

func (v *StateProofVerifier) verifyStateProofMessage(stateProof *datatypes.EncodedStateProof, message datatypes.Message) error {
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

func (v *StateProofVerifier) AdvanceState(stateProof *datatypes.EncodedStateProof, message datatypes.Message) error {
	err := v.verifyStateProofMessage(stateProof, message)
	if err != nil {
		return err
	}

	v.advanceVerifier(message)
	return nil
}
