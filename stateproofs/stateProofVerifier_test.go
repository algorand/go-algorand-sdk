package stateproofverification

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/json"
	"github.com/algorand/go-algorand-sdk/types"
)

func readJsonFile(filePath string, target interface{}, assertions *require.Assertions) {
	contents, err := os.ReadFile(filePath)
	assertions.NoError(err)

	err = json.Decode(contents, &target)
	assertions.NoError(err)
}

func TestStateProofVerificationSanity(t *testing.T) {
	a := require.New(t)

	prevStateProofFileName := "proof_20133889_to_20134144"
	currentStateProofFileName := "proof_20134145_to_20134400"

	var oldStateProofData models.StateProof
	var newStateProofData models.StateProof

	readJsonFile(prevStateProofFileName, &oldStateProofData, a)
	readJsonFile(currentStateProofFileName, &newStateProofData, a)

	message := types.Message{
		BlockHeadersCommitment: newStateProofData.Message.Blockheaderscommitment,
		VotersCommitment:       newStateProofData.Message.Voterscommitment,
		LnProvenWeight:         newStateProofData.Message.Lnprovenweight,
		FirstAttestedRound:     newStateProofData.Message.Firstattestedround,
		LastAttestedRound:      newStateProofData.Message.Lastattestedround,
	}

	encodedStateProof := types.EncodedStateProof(newStateProofData.Stateproof)
	verifier := InitializeVerifier(oldStateProofData.Message.Voterscommitment, oldStateProofData.Message.Lnprovenweight)
	err := verifier.VerifyStateProofMessage(&encodedStateProof, &message)
	a.NoError(err)
}
