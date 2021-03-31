package models

// TransactionSignatureMultisig (msig) structure holding multiple subsignatures.
// Definition:
// crypto/multisig.go : MultisigSig
type TransactionSignatureMultisig struct {
	// Subsignature (subsig) holds pairs of public key and signatures.
	Subsignature []TransactionSignatureMultisigSubsignature `json:"subsignature,omitempty"`

	// Threshold (thr)
	Threshold uint64 `json:"threshold,omitempty"`

	// Version (v)
	Version uint64 `json:"version,omitempty"`
}
