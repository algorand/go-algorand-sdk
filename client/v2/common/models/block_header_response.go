package models

import "github.com/algorand/go-algorand-sdk/v2/types"

// BlockHeaderResponse block header.
type BlockHeaderResponse struct {
	// Blockheader block header data.
	Blockheader types.Block `json:"blockHeader"`
}
