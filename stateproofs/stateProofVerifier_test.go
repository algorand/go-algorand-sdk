package stateproofs

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/json"
	"github.com/algorand/go-algorand-sdk/types"
)

//go:embed "prevStateProof.json"
var prevStateProofData []byte

//go:embed "newStateProof.json"
var newStateProofData []byte

func TestStateProofVerification(t *testing.T) {
	a := require.New(t)

	var prevStateProof models.StateProof
	var newStateProof models.StateProof

	err := json.Decode(prevStateProofData, &prevStateProof)
	a.NoError(err)
	err = json.Decode(newStateProofData, &newStateProof)
	a.NoError(err)

	message := types.Message{
		BlockHeadersCommitment: newStateProof.Message.Blockheaderscommitment,
		VotersCommitment:       newStateProof.Message.Voterscommitment,
		LnProvenWeight:         newStateProof.Message.Lnprovenweight,
		FirstAttestedRound:     newStateProof.Message.Firstattestedround,
		LastAttestedRound:      newStateProof.Message.Lastattestedround,
	}
	encodedStateProof := types.EncodedStateProof(newStateProof.Stateproof)

	verifier := InitializeVerifier(prevStateProof.Message.Voterscommitment, prevStateProof.Message.Lnprovenweight)
	err = verifier.Verify(&encodedStateProof, &message)
	a.NoError(err)
}
