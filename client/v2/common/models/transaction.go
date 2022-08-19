package models

// Transaction contains all fields common to all transactions and serves as an
// envelope to all transactions type. Represents both regular and inner
// transactions.
// Definition:
// data/transactions/signedtxn.go : SignedTxn
// data/transactions/transaction.go : Transaction
type Transaction struct {
	// ApplicationTransaction fields for application transactions.
	// Definition:
	// data/transactions/application.go : ApplicationCallTxnFields
	ApplicationTransaction TransactionApplication `json:"application-transaction,omitempty"`

	// AssetConfigTransaction fields for asset allocation, re-configuration, and
	// destruction.
	// A zero value for asset-id indicates asset creation.
	// A zero value for the params indicates asset destruction.
	// Definition:
	// data/transactions/asset.go : AssetConfigTxnFields
	AssetConfigTransaction TransactionAssetConfig `json:"asset-config-transaction,omitempty"`

	// AssetFreezeTransaction fields for an asset freeze transaction.
	// Definition:
	// data/transactions/asset.go : AssetFreezeTxnFields
	AssetFreezeTransaction TransactionAssetFreeze `json:"asset-freeze-transaction,omitempty"`

	// AssetTransferTransaction fields for an asset transfer transaction.
	// Definition:
	// data/transactions/asset.go : AssetTransferTxnFields
	AssetTransferTransaction TransactionAssetTransfer `json:"asset-transfer-transaction,omitempty"`

	// AuthAddr (sgnr) this is included with signed transactions when the signing
	// address does not equal the sender. The backend can use this to ensure that auth
	// addr is equal to the accounts auth addr.
	AuthAddr string `json:"auth-addr,omitempty"`

	// CloseRewards (rc) rewards applied to close-remainder-to account.
	CloseRewards uint64 `json:"close-rewards,omitempty"`

	// ClosingAmount (ca) closing amount for transaction.
	ClosingAmount uint64 `json:"closing-amount,omitempty"`

	// ConfirmedRound round when the transaction was confirmed.
	ConfirmedRound uint64 `json:"confirmed-round,omitempty"`

	// CreatedApplicationIndex specifies an application index (ID) if an application
	// was created with this transaction.
	CreatedApplicationIndex uint64 `json:"created-application-index,omitempty"`

	// CreatedAssetIndex specifies an asset index (ID) if an asset was created with
	// this transaction.
	CreatedAssetIndex uint64 `json:"created-asset-index,omitempty"`

	// Fee (fee) Transaction fee.
	Fee uint64 `json:"fee"`

	// FirstValid (fv) First valid round for this transaction.
	FirstValid uint64 `json:"first-valid"`

	// GenesisHash (gh) Hash of genesis block.
	GenesisHash []byte `json:"genesis-hash,omitempty"`

	// GenesisId (gen) genesis block ID.
	GenesisId string `json:"genesis-id,omitempty"`

	// GlobalStateDelta (gd) Global state key/value changes for the application being
	// executed by this transaction.
	GlobalStateDelta []EvalDeltaKeyValue `json:"global-state-delta,omitempty"`

	// Group (grp) Base64 encoded byte array of a sha512/256 digest. When present
	// indicates that this transaction is part of a transaction group and the value is
	// the sha512/256 hash of the transactions in that group.
	Group []byte `json:"group,omitempty"`

	// Id transaction ID
	Id string `json:"id,omitempty"`

	// InnerTxns inner transactions produced by application execution.
	InnerTxns []Transaction `json:"inner-txns,omitempty"`

	// IntraRoundOffset offset into the round where this transaction was confirmed.
	IntraRoundOffset uint64 `json:"intra-round-offset,omitempty"`

	// KeyregTransaction fields for a keyreg transaction.
	// Definition:
	// data/transactions/keyreg.go : KeyregTxnFields
	KeyregTransaction TransactionKeyreg `json:"keyreg-transaction,omitempty"`

	// LastValid (lv) Last valid round for this transaction.
	LastValid uint64 `json:"last-valid"`

	// Lease (lx) Base64 encoded 32-byte array. Lease enforces mutual exclusion of
	// transactions. If this field is nonzero, then once the transaction is confirmed,
	// it acquires the lease identified by the (Sender, Lease) pair of the transaction
	// until the LastValid round passes. While this transaction possesses the lease, no
	// other transaction specifying this lease can be confirmed.
	Lease []byte `json:"lease,omitempty"`

	// LocalStateDelta (ld) Local state key/value changes for the application being
	// executed by this transaction.
	LocalStateDelta []AccountStateDelta `json:"local-state-delta,omitempty"`

	// Logs (lg) Logs for the application being executed by this transaction.
	Logs [][]byte `json:"logs,omitempty"`

	// Note (note) Free form data.
	Note []byte `json:"note,omitempty"`

	// PaymentTransaction fields for a payment transaction.
	// Definition:
	// data/transactions/payment.go : PaymentTxnFields
	PaymentTransaction TransactionPayment `json:"payment-transaction,omitempty"`

	// ReceiverRewards (rr) rewards applied to receiver account.
	ReceiverRewards uint64 `json:"receiver-rewards,omitempty"`

	// RekeyTo (rekey) when included in a valid transaction, the accounts auth addr
	// will be updated with this value and future signatures must be signed with the
	// key represented by this address.
	RekeyTo string `json:"rekey-to,omitempty"`

	// RoundTime time when the block this transaction is in was confirmed.
	RoundTime uint64 `json:"round-time,omitempty"`

	// Sender (snd) Sender's address.
	Sender string `json:"sender"`

	// SenderRewards (rs) rewards applied to sender account.
	SenderRewards uint64 `json:"sender-rewards,omitempty"`

	// Signature validation signature associated with some data. Only one of the
	// signatures should be provided.
	Signature TransactionSignature `json:"signature,omitempty"`

	// StateProofTransaction fields for a state proof transaction.
	// Definition:
	// data/transactions/stateproof.go : StateProofTxnFields
	StateProofTransaction TransactionStateProof `json:"state-proof-transaction,omitempty"`

	// Type (type) Indicates what type of transaction this is. Different types have
	// different fields.
	// Valid types, and where their fields are stored:
	// * (pay) payment-transaction
	// * (keyreg) keyreg-transaction
	// * (acfg) asset-config-transaction
	// * (axfer) asset-transfer-transaction
	// * (afrz) asset-freeze-transaction
	// * (appl) application-transaction
	// * (stpf) state-proof-transaction
	Type string `json:"tx-type,omitempty"`
}
