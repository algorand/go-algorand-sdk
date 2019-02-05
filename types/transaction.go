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

	VotePK      VotePK `codec:"votekey"`
	SelectionPK VRFPK  `codec:"selkey"`
}

// PaymentTxnFields captures the fields used by payment transactions.
type PaymentTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Receiver Address `codec:"rcv"`
	Amount   Algos   `codec:"amt"`
}

// Header captures the fields common to every transaction type.
type Header struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Sender     Address `codec:"snd"`
	Fee        Algos   `codec:"fee"`
	FirstValid Round   `codec:"fv"`
	LastValid  Round   `codec:"lv"`
	Note       []byte  `codec:"note"`
}
