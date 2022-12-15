package models

// KvDelta a single Delta containing the key, the previous value and the current
// value for a single round.
type KvDelta struct {
	// Key the key, base64 encoded.
	Key []byte `json:"key,omitempty"`

	// Value the new value of the KV store entry, base64 encoded.
	Value []byte `json:"value,omitempty"`
}
