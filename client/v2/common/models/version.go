package models

// Version algod version information.
type Version struct {
	// Build
	Build BuildVersion `json:"build"`

	// GenesisHash
	GenesisHash []byte `json:"genesis_hash_b64"`

	// GenesisID
	GenesisID string `json:"genesis_id"`

	// Versions
	Versions []string `json:"versions"`
}
