package models

// HoldingRef holdingRef names a holding by referring to an Address and Asset it
// belongs to.
type HoldingRef struct {
	// Address (d) Address in access list, or the sender of the transaction.
	Address string `json:"address"`

	// Asset (s) Asset ID for asset in access list.
	Asset uint64 `json:"asset"`
}
