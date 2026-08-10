package models

// TransactionSignature validation signature associated with some data. Only one of
// the signatures should be provided.
type TransactionSignature struct {
	// Logicsig (lsig) Programatic transaction signature.
	// Definition:
	// data/transactions/logicsig.go
	Logicsig TransactionSignatureLogicsig `json:"logicsig,omitempty"`

	// Multisig structure holding multiple subsignatures.
	// Definition:
	// crypto/multisig.go : MultisigSig
	Multisig TransactionSignatureMultisig `json:"multisig,omitempty"`

	// Pqsig structure holding a post-quantum signature.
	// Definition:
	// data/transactions/pqsig.go : PQSig
	Pqsig TransactionSignaturePQsig `json:"pqsig,omitempty"`

	// Sig (sig) Standard ed25519 signature.
	Sig []byte `json:"sig,omitempty"`
}
