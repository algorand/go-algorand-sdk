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
