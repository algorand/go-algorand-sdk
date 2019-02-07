package kmd

import (
	"github.com/algorand/go-algorand-sdk/types"
)

// APIV1Wallet is the API's representation of a wallet
type APIV1Wallet struct {
	ID                    string         `json:"id"`
	Name                  string         `json:"name"`
	DriverName            string         `json:"driver_name"`
	DriverVersion         uint32         `json:"driver_version"`
	SupportsMnemonicUX    bool           `json:"mnemonic_ux"`
	SupportedTransactions []types.TxType `json:"supported_txs"`
}

// APIV1WalletHandle includes the wallet the handle corresponds to
// and the number of number of seconds to expiration
type APIV1WalletHandle struct {
	Wallet         APIV1Wallet `json:"wallet"`
	ExpiresSeconds int64       `json:"expires_seconds"`
}
