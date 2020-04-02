package algod

import (
	"context"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

type SuggestedParamsService struct {
	c *Client
}

func (s *SuggestedParamsService) Do(ctx context.Context, headers ...*common.Header) (params types.SuggestedParams, err error) {
	var response models.TransactionParams
	err = s.c.get(ctx, &response, "/transactions/params", nil, headers)
	params = types.SuggestedParams{
		Fee:              types.MicroAlgos(response.Fee),
		GenesisID:        response.GenesisID,
		GenesisHash:      response.Genesishash,
		FirstRoundValid:  types.Round(response.LastRound),
		LastRoundValid:   types.Round(response.LastRound + 1000),
		ConsensusVersion: response.ConsensusVersion,
	}
	return
}
