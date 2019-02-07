package kmd

import (
	"errors"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/types"
)

// APIV1Response is the interface that all API V1 responses must satisfy
type APIV1Response interface {
	GetError() error
}

// APIV1ResponseEnvelope is a common envelope that all API V1 responses must embed
type APIV1ResponseEnvelope struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
	Error   bool     `json:"error"`
	Message string   `json:"message"`
}

// GetError allows VersionResponse to satisfy the APIV1Response interface, even
// though it can never return an error and is not versioned
func (r VersionsResponse) GetError() error {
	return nil
}

// GetError allows responses that embed an APIV1ResponseEnvelope to satisfy the
// APIV1Response interface
func (r APIV1ResponseEnvelope) GetError() error {
	if r.Error {
		return errors.New(r.Message)
	}
	return nil
}

// VersionsResponse is the response to `GET /versions`
// friendly:VersionsResponse
type VersionsResponse struct {
	_struct  struct{} `codec:",omitempty,omitemptyarray"`
	Versions []string `json:"versions"`
}

// APIV1GETWalletsResponse is the response to `GET /v1/wallets`
// friendly:ListWalletsResponse
type APIV1GETWalletsResponse struct {
	APIV1ResponseEnvelope
	Wallets []APIV1Wallet `json:"wallets"`
}

// APIV1POSTWalletResponse is the response to `POST /v1/wallet`
// friendly:CreateWalletResponse
type APIV1POSTWalletResponse struct {
	APIV1ResponseEnvelope
	Wallet APIV1Wallet `json:"wallet"`
}

// APIV1POSTWalletInitResponse is the response to `POST /v1/wallet/init`
// friendly:InitWalletHandleTokenResponse
type APIV1POSTWalletInitResponse struct {
	APIV1ResponseEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// APIV1POSTWalletReleaseResponse is the response to `POST /v1/wallet/release`
// friendly:ReleaseWalletHandleTokenResponse
type APIV1POSTWalletReleaseResponse struct {
	APIV1ResponseEnvelope
}

// APIV1POSTWalletRenewResponse is the response to `POST /v1/wallet/renew`
// friendly:RenewWalletHandleTokenResponse
type APIV1POSTWalletRenewResponse struct {
	APIV1ResponseEnvelope
	WalletHandle APIV1WalletHandle `json:"wallet_handle"`
}

// APIV1POSTWalletRenameResponse is the response to `POST /v1/wallet/rename`
// friendly:RenameWalletResponse
type APIV1POSTWalletRenameResponse struct {
	APIV1ResponseEnvelope
	Wallet APIV1Wallet `json:"wallet"`
}

// APIV1POSTWalletInfoResponse is the response to `POST /v1/wallet/info`
// friendly:WalletInfoResponse
type APIV1POSTWalletInfoResponse struct {
	APIV1ResponseEnvelope
	WalletHandle APIV1WalletHandle `json:"wallet_handle"`
}

// APIV1POSTMasterKeyExportResponse is the reponse to `POST /v1/master_key/export`
// friendly:ExportMasterKeyResponse
type APIV1POSTMasterKeyExportResponse struct {
	APIV1ResponseEnvelope
	MasterDerivationKey types.MasterDerivationKey `json:"master_derivation_key"`
}

// APIV1POSTKeyImportResponse is the repsonse to `POST /v1/key/import`
// friendly:ImportKeyResponse
type APIV1POSTKeyImportResponse struct {
	APIV1ResponseEnvelope
	Address string `json:"address"`
}

// APIV1POSTKeyExportResponse is the reponse to `POST /v1/key/export`
// friendly:ExportKeyResponse
type APIV1POSTKeyExportResponse struct {
	APIV1ResponseEnvelope
	PrivateKey ed25519.PrivateKey `json:"private_key"`
}

// APIV1POSTKeyResponse is the response to `POST /v1/key`
// friendly:GenerateKeyResponse
type APIV1POSTKeyResponse struct {
	APIV1ResponseEnvelope
	Address string `json:"address"`
}

// APIV1DELETEKeyResponse is the response to `DELETE /v1/key`
// friendly:DeleteKeyResponse
type APIV1DELETEKeyResponse struct {
	APIV1ResponseEnvelope
}

// APIV1POSTKeysListResponse is the response to `POST /v1/keys/list`
// friendly:ListKeysResponse
type APIV1POSTKeysListResponse struct {
	APIV1ResponseEnvelope
	Addresses []string `json:"addresses"`
}

// APIV1POSTTransactionSignResponse is the repsonse to `POST /v1/transaction/sign`
// friendly:SignTransactionResponse
type APIV1POSTTransactionSignResponse struct {
	APIV1ResponseEnvelope
	SignedTransaction []byte `json:"signed_transaction"`
}

// APIV1POSTMultisigListResponse is the response to `POST /v1/multisig/list`
// friendly:ListMultisigResponse
type APIV1POSTMultisigListResponse struct {
	APIV1ResponseEnvelope
	Addresses []string `json:"addresses"`
}

// APIV1POSTMultisigImportResponse is the response to `POST /v1/multisig/import`
// friendly:ImportMultisigResponse
type APIV1POSTMultisigImportResponse struct {
	APIV1ResponseEnvelope
	Address string `json:"address"`
}

// APIV1POSTMultisigExportResponse is the response to `POST /v1/multisig/export`
// friendly:ExportMultisigResponse
type APIV1POSTMultisigExportResponse struct {
	APIV1ResponseEnvelope
	Version   uint8               `json:"multisig_version"`
	Threshold uint8               `json:"threshold"`
	PKs       []ed25519.PublicKey `json:"pks"`
}

// APIV1POSTMultisigDeleteResponse is the response to POST /v1/multisig/delete`
// friendly:DeleteMultisigResponse
type APIV1POSTMultisigDeleteResponse struct {
	APIV1ResponseEnvelope
}

// APIV1POSTMultisigTransactionSignResponse is the repsonse to `POST /v1/multisig/sign`
// friendly:SignMultisigResponse
type APIV1POSTMultisigTransactionSignResponse struct {
	APIV1ResponseEnvelope
	Multisig types.MultisigSig `json:"multisig"`
}
