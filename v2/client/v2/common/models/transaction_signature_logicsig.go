package models

// TransactionSignatureLogicsig (lsig) Programatic transaction signature.
// Definition:
// data/transactions/logicsig.go
type TransactionSignatureLogicsig struct {
	// Args (arg) Logic arguments, base64 encoded.
	Args [][]byte `json:"args,omitempty"`

	// Logic (l) Program signed by a signature or multi signature, or hashed to be the
	// address of ana ccount. Base64 encoded TEAL program.
	Logic []byte `json:"logic"`

	// MultisigSignature (msig) structure holding multiple subsignatures.
	// Definition:
	// crypto/multisig.go : MultisigSig
	MultisigSignature TransactionSignatureMultisig `json:"multisig-signature,omitempty"`

	// Signature (sig) ed25519 signature.
	Signature []byte `json:"signature,omitempty"`
}
