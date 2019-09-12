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
}

// SignedTxn wraps a transaction and a signature. The encoding of this struct
// is suitable to broadcast on the network
type SignedTxn struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Sig  Signature   `codec:"sig"`
	Msig MultisigSig `codec:"msig"`
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
	// A zero value (including the creator address) means allocation.
	ConfigAsset AssetID `codec:"caid"`

	// AssetParams are the parameters for the asset being
	// created or re-configured.  A zero value means destruction.
	AssetParams AssetParams `codec:"apar"`
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
}
