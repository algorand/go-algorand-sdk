package models

// GetBlockTimeStampOffsetResponse response containing the timestamp offset in
// seconds
type GetBlockTimeStampOffsetResponse struct {
	// Offset timestamp offset in seconds.
	Offset uint64 `json:"offset"`
}
