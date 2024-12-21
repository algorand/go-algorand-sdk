package models

// HbProofFields (hbprf) HbProof is a signature using HeartbeatAddress's partkey,
// thereby showing it is online.
type HbProofFields struct {
	// HbPk (p) Public key of the heartbeat message.
	HbPk []byte `json:"hb-pk,omitempty"`

	// HbPk1sig (p1s) Signature of OneTimeSignatureSubkeyOffsetID(PK, Batch, Offset)
	// under the key PK2.
	HbPk1sig []byte `json:"hb-pk1sig,omitempty"`

	// HbPk2 (p2) Key for new-style two-level ephemeral signature.
	HbPk2 []byte `json:"hb-pk2,omitempty"`

	// HbPk2sig (p2s) Signature of OneTimeSignatureSubkeyBatchID(PK2, Batch) under the
	// master key (OneTimeSignatureVerifier).
	HbPk2sig []byte `json:"hb-pk2sig,omitempty"`

	// HbSig (s) Signature of the heartbeat message.
	HbSig []byte `json:"hb-sig,omitempty"`
}
