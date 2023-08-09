package models

// BlockTxidsResponse top level transaction IDs in a block.
type BlockTxidsResponse struct {
	// Blocktxids block transaction IDs.
	Blocktxids []string `json:"blockTxids"`
}
