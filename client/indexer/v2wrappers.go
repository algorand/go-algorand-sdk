package indexer

import "fmt"

type AssetBalancesResponse struct {
	// A simplified version of AssetHolding
	Balances []MiniAssetHolding `json:"balances"`

	// Round at which the results are valid.
	Round uint64 `json:"round"`
}

type GetAssetHoldingParams struct {
	Limit   uint64 `url:"limit,omitempty"`
	Offset  uint64 `url:"offset,omitempty"`
	Round   uint64 `url:"round,omitempty"`
	CurrencyGreaterThan     uint64 `url:"currency-greater-than,omitempty"`
	CurrencyLessThan        uint64 `url:"currency-less-than,omitempty"`
}

type MiniAssetHolding struct {
	Address  string `json:"address"`
	Amount   uint64 `json:"amount"`
	IsFrozen bool   `json:"isFrozen"`
}

func (client Client) GetAssetHoldings(assetIndex uint64, params getAssetHoldingParams, headers ...*Header) (validRound uint64, holders []models.MiniAssetHolding, err) {
	var response models.GetAssetHoldingResponse
	err = client.get(&response, fmt.Sprintf("/assets/%d/balances", assetIndex), params, headers)
	validRound = response.Round
	holders = response.Balances
	return
}
