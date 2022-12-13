package models

// TxLease
type TxLease struct {
	// Expiration round that the lease expires
	Expiration uint64 `json:"expiration"`

	// Lease lease data
	Lease []byte `json:"lease"`

	// Sender address of the lease sender
	Sender string `json:"sender"`
}
