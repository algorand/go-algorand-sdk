package models

// TransactionSignaturePQsig structure holding a post-quantum signature.
// Definition:
// data/transactions/pqsig.go : PQSig
type TransactionSignaturePQsig struct {
	// PublicKey (pk)
	PublicKey []byte `json:"public-key"`

	// Salt (slt) a single byte, added to ensure the hashed address is not an Ed25519
	// curve point
	Salt uint64 `json:"salt,omitempty"`

	// Scheme (sch) identifies the internal signature scheme.
	Scheme string `json:"scheme"`

	// Signature (sig)
	Signature []byte `json:"signature"`
}
