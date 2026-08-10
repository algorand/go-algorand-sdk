package models

// PeerStatus the status of a connected peer in the P2P network
type PeerStatus struct {
	// ConnectionType connection type
	ConnectionType string `json:"connection-type,omitempty"`

	// NetworkAddress network address of the peer
	NetworkAddress string `json:"network-address"`

	// NetworkType network type
	NetworkType string `json:"network-type,omitempty"`
}
