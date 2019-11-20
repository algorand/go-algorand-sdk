package types

// Transaction describes a transaction that can appear in a block.
type Transaction struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Type of transaction
	Type TxType `codec:"type"`

	// Common fields for all types of transactions
	Header

	// Fields for different types of transactions
	KeyregTxnFields
	PaymentTxnFields
	AssetConfigTxnFields
	AssetTransferTxnFields
	AssetFreezeTxnFields
}

// SignedTxn wraps a transaction and a signature. The encoding of this struct
// is suitable to broadcast on the network
type SignedTxn struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Sig  Signature   `codec:"sig"`
	Msig MultisigSig `codec:"msig"`
	Lsig LogicSig    `codec:"lsig"`
	Txn  Transaction `codec:"txn"`
}

// KeyregTxnFields captures the fields used for key registration transactions.
type KeyregTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	VotePK          VotePK `codec:"votekey"`
	SelectionPK     VRFPK  `codec:"selkey"`
	VoteFirst       Round  `codec:"votefst"`
	VoteLast        Round  `codec:"votelst"`
	VoteKeyDilution uint64 `codec:"votekd"`
}

// PaymentTxnFields captures the fields used by payment transactions.
type PaymentTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Receiver Address    `codec:"rcv"`
	Amount   MicroAlgos `codec:"amt"`

	// When CloseRemainderTo is set, it indicates that the
	// transaction is requesting that the account should be
	// closed, and all remaining funds be transferred to this
	// address.
	CloseRemainderTo Address `codec:"close"`
}

// AssetConfigTxnFields captures the fields used for asset
// allocation, re-configuration, and destruction.
type AssetConfigTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// ConfigAsset is the asset being configured or destroyed.
	// A zero value means allocation.
	ConfigAsset AssetIndex `codec:"caid"`

	// AssetParams are the parameters for the asset being
	// created or re-configured.  A zero value means destruction.
	AssetParams AssetParams `codec:"apar"`
}

// AssetTransferTxnFields captures the fields used for asset transfers.
type AssetTransferTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	XferAsset AssetIndex `codec:"xaid"`

	// AssetAmount is the amount of asset to transfer.
	// A zero amount transferred to self allocates that asset
	// in the account's Assets map.
	AssetAmount uint64 `codec:"aamt"`

	// AssetSender is the sender of the transfer.  If this is not
	// a zero value, the real transaction sender must be the Clawback
	// address from the AssetParams.  If this is the zero value,
	// the asset is sent from the transaction's Sender.
	AssetSender Address `codec:"asnd"`

	// AssetReceiver is the recipient of the transfer.
	AssetReceiver Address `codec:"arcv"`

	// AssetCloseTo indicates that the asset should be removed
	// from the account's Assets map, and specifies where the remaining
	// asset holdings should be transferred.  It's always valid to transfer
	// remaining asset holdings to the creator account.
	AssetCloseTo Address `codec:"aclose"`
}

// AssetFreezeTxnFields captures the fields used for freezing asset slots.
type AssetFreezeTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// FreezeAccount is the address of the account whose asset
	// slot is being frozen or un-frozen.
	FreezeAccount Address `codec:"fadd"`

	// FreezeAsset is the asset ID being frozen or un-frozen.
	FreezeAsset AssetIndex `codec:"faid"`

	// AssetFrozen is the new frozen value.
	AssetFrozen bool `codec:"afrz"`
}

// Header captures the fields common to every transaction type.
type Header struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Sender      Address    `codec:"snd"`
	Fee         MicroAlgos `codec:"fee"`
	FirstValid  Round      `codec:"fv"`
	LastValid   Round      `codec:"lv"`
	Note        []byte     `codec:"note"`
	GenesisID   string     `codec:"gen"`
	GenesisHash Digest     `codec:"gh"`

	// Group specifies that this transaction is part of a
	// transaction group (and, if so, specifies the hash
	// of a TxGroup).
	Group Digest `codec:"grp"`

	// Lease enforces mutual exclusion of transactions.  If this field is
	// nonzero, then once the transaction is confirmed, it acquires the
	// lease identified by the (Sender, Lease) pair of the transaction until
	// the LastValid round passes.  While this transaction possesses the
	// lease, no other transaction specifying this lease can be confirmed.
	Lease [32]byte `codec:"lx"`
}

// TxGroup describes a group of transactions that must appear
// together in a specific order in a block.
type TxGroup struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// TxGroupHashes specifies a list of hashes of transactions that must appear
	// together, sequentially, in a block in order for the group to be
	// valid.  Each hash in the list is a hash of a transaction with
	// the `Group` field omitted.
	TxGroupHashes []Digest `codec:"txlist"`
}

// AddLease adds the passed lease (see types/transaction.go) to the header of the passed transaction
// and updates fee accordingly
// - lease: the [32]byte lease to add to the header
// - feePerByte: the new feePerByte
func (tx *Transaction) AddLease(lease [32]byte, feePerByte uint64) {
	copy(tx.Header.Lease[:], lease[:])
	// normally we would use estimateSize,
	// and set fee = feePerByte * estimateSize,
	// but this would cause a circular import.
	// we know we are adding 32 bytes (+ a few bytes to hold the 32), so increase fee accordingly.
	tx.Header.Fee = tx.Header.Fee + MicroAlgos(37*feePerByte)
}

// AddLeaseWithFlatFee adds the passed lease (see types/transaction.go) to the header of the passed transaction
// and updates fee accordingly
// - lease: the [32]byte lease to add to the header
// - flatFee: the new flatFee
func (tx *Transaction) AddLeaseWithFlatFee(lease [32]byte, flatFee uint64) {
	tx.Header.Lease = lease
	tx.Header.Fee = MicroAlgos(flatFee)
}
