package models

// StateProof represents a state proof and its corresponding message
type StateProof struct {
	// Message the encoded message.
	Message []byte `json:"Message"`

	// Stateproof the encoded StateProof for the message.
	Stateproof []byte `json:"StateProof"`
}
