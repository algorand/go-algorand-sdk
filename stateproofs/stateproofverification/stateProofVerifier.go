package stateproofverification

import (
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification/stateproof"
	"github.com/algorand/go-algorand-sdk/stateproofs/transactionverificationtypes"
)

const strengthTarget = uint64(256)

type StateProofVerifier struct {
	stateProofVerifier *stateproof.Verifier
}

func InitializeVerifier(votersCommitment transactionverificationtypes.GenericDigest, lnProvenWeight uint64) *StateProofVerifier {
	return &StateProofVerifier{stateProofVerifier: stateproof.MkVerifierWithLnProvenWeight(votersCommitment,
		lnProvenWeight, strengthTarget)}
}

func (v *StateProofVerifier) VerifyStateProofMessage(stateProof *transactionverificationtypes.EncodedStateProof, message transactionverificationtypes.Message) error {
	messageHash := message.IntoStateProofMessageHash()

	var decodedStateProof stateproof.StateProof
	err := msgpack.Decode(*stateProof, &decodedStateProof)
	if err != nil {
		return err
	}

	return v.stateProofVerifier.Verify(message.LastAttestedRound, messageHash, &decodedStateProof)
}
