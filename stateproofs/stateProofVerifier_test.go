package stateproofs

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/encoding/json"
	"github.com/algorand/go-algorand-sdk/types"
)

//go:embed "prevStateProof.json" "newStateProof.json"
var res embed.FS

func readJsonFileFromRes(filePath string, target interface{}, assertions *require.Assertions) {
	contents, err := res.ReadFile(filePath)
	assertions.NoError(err)

	err = json.Decode(contents, &target)
	assertions.NoError(err)
}

func TestStateProofVerification(t *testing.T) {
	a := require.New(t)

	prevStateProofFileName := "prevStateProof.json"
	newStateProofFileName := "newStateProof.json"

	var prevStateProof models.StateProof
	var newStateProof models.StateProof

	readJsonFileFromRes(prevStateProofFileName, &prevStateProof, a)
	readJsonFileFromRes(newStateProofFileName, &newStateProof, a)

	message := types.Message{
		BlockHeadersCommitment: newStateProof.Message.Blockheaderscommitment,
		VotersCommitment:       newStateProof.Message.Voterscommitment,
		LnProvenWeight:         newStateProof.Message.Lnprovenweight,
		FirstAttestedRound:     newStateProof.Message.Firstattestedround,
		LastAttestedRound:      newStateProof.Message.Lastattestedround,
	}
	encodedStateProof := types.EncodedStateProof(newStateProof.Stateproof)

	verifier := InitializeVerifier(prevStateProof.Message.Voterscommitment, prevStateProof.Message.Lnprovenweight)
	err := verifier.Verify(&encodedStateProof, &message)
	a.NoError(err)
}
