package models

// Version algod version information.
type Version struct {
	// Build
	Build BuildVersion `json:"build,omitempty"`

	// GenesisHash
	GenesisHash []byte `json:"genesis_hash_b64,omitempty"`

	// GenesisID
	GenesisID string `json:"genesis_id,omitempty"`

	// Versions
	Versions []string `json:"versions,omitempty"`
}
