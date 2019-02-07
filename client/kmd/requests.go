package kmd

import (
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/types"
)

// APIV1Request is the interface that all API V1 requests must satisfy
type APIV1Request interface{}

// APIV1RequestEnvelope is a common envelope that all API V1 requests must embed
type APIV1RequestEnvelope struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
}

// VersionsRequest is the request for `GET /versions`
//
// swagger:model VersionsRequest
type VersionsRequest struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
}

// APIV1GETWalletsRequest is the request for `GET /v1/wallets`
//
// swagger:model ListWalletsRequest
type APIV1GETWalletsRequest struct {
	APIV1RequestEnvelope
}

// APIV1POSTWalletRequest is the request for `POST /v1/wallet`
//
// swagger:model CreateWalletRequest
type APIV1POSTWalletRequest struct {
	APIV1RequestEnvelope
	WalletName          string                    `json:"wallet_name"`
	WalletDriverName    string                    `json:"wallet_driver_name"`
	WalletPassword      string                    `json:"wallet_password"`
	MasterDerivationKey types.MasterDerivationKey `json:"master_derivation_key"`
}

// APIV1POSTWalletInitRequest is the request for `POST /v1/wallet/init`
//
// swagger:model InitWalletHandleTokenRequest
type APIV1POSTWalletInitRequest struct {
	APIV1RequestEnvelope
	WalletID       string `json:"wallet_id"`
	WalletPassword string `json:"wallet_password"`
}

// APIV1POSTWalletReleaseRequest is the request for `POST /v1/wallet/release`
//
// swagger:model ReleaseWalletHandleTokenRequest
type APIV1POSTWalletReleaseRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// APIV1POSTWalletRenewRequest is the request for `POST /v1/wallet/renew`
//
// swagger:model RenewWalletHandleTokenRequest
type APIV1POSTWalletRenewRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// APIV1POSTWalletRenameRequest is the request for `POST /v1/wallet/rename`
//
// swagger:model RenameWalletRequest
type APIV1POSTWalletRenameRequest struct {
	APIV1RequestEnvelope
	WalletID       string `json:"wallet_id"`
	WalletPassword string `json:"wallet_password"`
	NewWalletName  string `json:"wallet_name"`
}

// APIV1POSTWalletInfoRequest is the request for `POST /v1/wallet/info`
//
// swagger:model WalletInfoRequest
type APIV1POSTWalletInfoRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// APIV1POSTMasterKeyExportRequest is the request for `POST /v1/master_key/export`
//
// swagger:model ExportMasterKeyRequest
type APIV1POSTMasterKeyExportRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	WalletPassword    string `json:"wallet_password"`
}

// APIV1POSTKeyImportRequest is the request for `POST /v1/key/import`
//
// swagger:model ImportKeyRequest
type APIV1POSTKeyImportRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string             `json:"wallet_handle_token"`
	PrivateKey        ed25519.PrivateKey `json:"private_key"`
}

// APIV1POSTKeyExportRequest is the request for `POST /v1/key/export`
//
// swagger:model ExportKeyRequest
type APIV1POSTKeyExportRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Address           string `json:"address"`
	WalletPassword    string `json:"wallet_password"`
}

// APIV1POSTKeyRequest is the request for `POST /v1/key`
//
// swagger:model GenerateKeyRequest
type APIV1POSTKeyRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	DisplayMnemonic   bool   `json:"display_mnemonic"`
}

// APIV1DELETEKeyRequest is the request for `DELETE /v1/key`
//
// swagger:model DeleteKeyRequest
type APIV1DELETEKeyRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Address           string `json:"address"`
	WalletPassword    string `json:"wallet_password"`
}

// APIV1POSTKeysListRequest is the request for `POST /v1/keys/list`
//
// swagger:model ListKeysRequest
type APIV1POSTKeysListRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// APIV1POSTTransactionSignRequest is the request for `POST /v1/transaction/sign`
//
// swagger:model SignTransactionRequest
type APIV1POSTTransactionSignRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Transaction       []byte `json:"transaction"`
	WalletPassword    string `json:"wallet_password"`
}

// APIV1POSTMultisigListRequest is the request for `POST /v1/multisig/list`
//
// swagger:model ListMultisigRequest
type APIV1POSTMultisigListRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// APIV1POSTMultisigImportRequest is the request for `POST /v1/multisig/import`
//
// swagger:model ImportMultisigRequest
type APIV1POSTMultisigImportRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string              `json:"wallet_handle_token"`
	Version           uint8               `json:"multisig_version"`
	Threshold         uint8               `json:"threshold"`
	PKs               []ed25519.PublicKey `json:"pks"`
}

// APIV1POSTMultisigExportRequest is the request for `POST /v1/multisig/export`
//
// swagger:model ExportMultisigRequest
type APIV1POSTMultisigExportRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Address           string `json:"address"`
	WalletPassword    string `json:"wallet_password"`
}

// APIV1POSTMultisigDeleteRequest is the request for `POST /v1/multisig/delete`
//
// swagger:model DeleteMultisigRequest
type APIV1POSTMultisigDeleteRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
	Address           string `json:"address"`
	WalletPassword    string `json:"wallet_password"`
}

// APIV1POSTMultisigTransactionSignRequest is the request for `POST /v1/multisig/sign`
//
// swagger:model SignMultisigRequest
type APIV1POSTMultisigTransactionSignRequest struct {
	APIV1RequestEnvelope
	WalletHandleToken string            `json:"wallet_handle_token"`
	Transaction       []byte            `json:"transaction"`
	PublicKey         ed25519.PublicKey `json:"public_key"`
	PartialMsig       types.MultisigSig `json:"partial_multisig"`
	WalletPassword    string            `json:"wallet_password"`
}
