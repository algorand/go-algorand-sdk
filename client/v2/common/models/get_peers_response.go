package models

// GetPeersResponse response containing the network peers of the node
type GetPeersResponse struct {
	// Peers
	Peers []PeerStatus `json:"Peers"`
}
