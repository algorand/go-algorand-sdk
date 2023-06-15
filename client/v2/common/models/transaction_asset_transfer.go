package models

// TransactionAssetTransfer fields for an asset transfer transaction.
// Definition:
// data/transactions/asset.go : AssetTransferTxnFields
type TransactionAssetTransfer struct {
	// Amount (aamt) Amount of asset to transfer. A zero amount transferred to self
	// allocates that asset in the account's Assets map.
	Amount uint64 `json:"amount"`

	// AssetId (xaid) ID of the asset being transferred.
	AssetId uint64 `json:"asset-id"`

	// CloseAmount number of assets transferred to the close-to account as part of the
	// transaction.
	CloseAmount uint64 `json:"close-amount,omitempty"`

	// CloseTo (aclose) Indicates that the asset should be removed from the account's
	// Assets map, and specifies where the remaining asset holdings should be
	// transferred. It's always valid to transfer remaining asset holdings to the
	// creator account.
	CloseTo string `json:"close-to,omitempty"`

	// Receiver (arcv) Recipient address of the transfer.
	Receiver string `json:"receiver"`

	// Sender (asnd) The effective sender during a clawback transactions. If this is
	// not a zero value, the real transaction sender must be the Clawback address from
	// the AssetParams.
	Sender string `json:"sender,omitempty"`
}
