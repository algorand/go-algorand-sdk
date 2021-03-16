package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type SuggestedParams struct {
	c *Client
}

func (s *SuggestedParams) Do(ctx context.Context, headers ...*common.Header) (params types.SuggestedParams, err error) {
	var response models.TransactionParams
	err = s.c.get(ctx, &response, "/v2/transactions/params", nil, headers)
	params = types.SuggestedParams{
		Fee:              types.MicroAlgos(response.Fee),
		GenesisID:        response.GenesisID,
		GenesisHash:      response.Genesishash,
		FirstRoundValid:  types.Round(response.LastRound),
		LastRoundValid:   types.Round(response.LastRound + 1000),
		ConsensusVersion: response.ConsensusVersion,
		MinFee:           response.MinFee,
	}
	return
}
