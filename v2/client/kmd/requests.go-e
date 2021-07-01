package kmd

import (
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/types"
)

// DefaultWalletDriver is the wallet backend that kmd will use by default
const DefaultWalletDriver = "sqlite"

// APIV1Request is the interface that all API V1 requests must satisfy
type APIV1Request interface{}

// APIV1RequestEnvelope is a common envelope that all API V1 requests must embed
type APIV1RequestEnvelope struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
}

// VersionsRequest is the request for `GET /versions`
type VersionsRequest struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
}

// ListWalletsRequest is the request for `GET /v1/wallets`
type ListWalletsRequest struct {
	APIV1RequestEnvelope
}

// CreateWalletRequest is the request for `POST /v1/wallet`
type CreateWalletRequest struct {
	APIV1RequestEnvelope
	WalletName          string                    `json:"wallet_name"`
	WalletDriverName    string                    `json:"wallet_driver_name"`
	WalletPassword      string                    `json:"wallet_password"`
	MasterDerivationKey types.MasterDerivationKey `json:"master_derivation_key"`
}

// InitWalletHandleRequest is the request for `POST /v1/wallet/init`
type InitWalletHandleRequest struct {
	APIV1RequestEnvelope
	WalletID       string `json:"wallet_id"`
	WalletPassword string `json:"wallet_password"`
}

// ReleaseWalletHandleRequest is the request for `POST /v1/wallet/release`
type ReleaseWalletHandleRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// RenewWalletHandleRequest is the request for `POST /v1/wallet/renew`
type RenewWalletHandleRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// RenameWalletRequest is the request for `POST /v1/wallet/rename`
type RenameWalletRequest struct {
	APIV1RequestEnvelope
	WalletID       string `json:"wallet_id"`
	WalletPassword string `json:"wallet_password"`
	NewWalletName  string `json:"wallet_name"`
}

// GetWalletRequest is the request for `POST /v1/wallet/info`
type GetWalletRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// ExportMasterDerivationKeyRequest is the request for `POST /v1/master_key/export`
type ExportMasterDerivationKeyRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	WalletPassword    string `json:"wallet_password"`
}

// ImportKeyRequest is the request for `POST /v1/key/import`
type ImportKeyRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string             `json:"wallet_handle_token"`
	PrivateKey        ed25519.PrivateKey `json:"private_key"`
}

// ExportKeyRequest is the request for `POST /v1/key/export`
type ExportKeyRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Address           string `json:"address"`
	WalletPassword    string `json:"wallet_password"`
}

// GenerateKeyRequest is the request for `POST /v1/key`
type GenerateKeyRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	DisplayMnemonic   bool   `json:"display_mnemonic"`
}

// DeleteKeyRequest is the request for `DELETE /v1/key`
type DeleteKeyRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Address           string `json:"address"`
	WalletPassword    string `json:"wallet_password"`
}

// ListKeysRequest is the request for `POST /v1/keys/list`
type ListKeysRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// SignTransactionRequest is the request for `POST /v1/transaction/sign`
type SignTransactionRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string            `json:"wallet_handle_token"`
	Transaction       []byte            `json:"transaction"`
	WalletPassword    string            `json:"wallet_password"`
	PublicKey         ed25519.PublicKey `json:"public_key"`
}

// ListMultisigRequest is the request for `POST /v1/multisig/list`
type ListMultisigRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// ImportMultisigRequest is the request for `POST /v1/multisig/import`
type ImportMultisigRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string              `json:"wallet_handle_token"`
	Version           uint8               `json:"multisig_version"`
	Threshold         uint8               `json:"threshold"`
	PKs               []ed25519.PublicKey `json:"pks"`
}

// ExportMultisigRequest is the request for `POST /v1/multisig/export`
type ExportMultisigRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Address           string `json:"address"`
	WalletPassword    string `json:"wallet_password"`
}

// DeleteMultisigRequest is the request for `POST /v1/multisig/delete`
type DeleteMultisigRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Address           string `json:"address"`
	WalletPassword    string `json:"wallet_password"`
}

// SignMultisigTransactionRequest is the request for `POST /v1/multisig/sign`
type SignMultisigTransactionRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string            `json:"wallet_handle_token"`
	Transaction       []byte            `json:"transaction"`
	PublicKey         ed25519.PublicKey `json:"public_key"`
	PartialMsig       types.MultisigSig `json:"partial_multisig"`
	WalletPassword    string            `json:"wallet_password"`
}
