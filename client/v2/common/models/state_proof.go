package models

// StateProof represents a state proof and its corresponding message
type StateProof struct {
	// Message represents the message that the state proofs are attesting to.
	Message StateProofMessage `json:"Message"`

	// Stateproof the encoded StateProof for the message.
	Stateproof []byte `json:"StateProof"`
}
