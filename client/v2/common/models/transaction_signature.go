package models

// TransactionSignature validation signature associated with some data. Only one of
// the signatures should be provided.
type TransactionSignature struct {
	// Logicsig (lsig) Programatic transaction signature.
	// Definition:
	// data/transactions/logicsig.go
	Logicsig TransactionSignatureLogicsig `json:"logicsig,omitempty"`

	// Multisig (msig) structure holding multiple subsignatures.
	// Definition:
	// crypto/multisig.go : MultisigSig
	Multisig TransactionSignatureMultisig `json:"multisig,omitempty"`

	// Sig (sig) Standard ed25519 signature.
	Sig []byte `json:"sig,omitempty"`
}
