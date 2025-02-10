package models

// BlockHeaderResponse block header.
type BlockHeaderResponse struct {
	// Blockheader block header data.
	Blockheader *map[string]interface{} `json:"blockHeader"`
}
