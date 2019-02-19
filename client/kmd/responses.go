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
type VersionsResponse struct {
	_struct  struct{} `codec:",omitempty,omitemptyarray"`
	Versions []string `json:"versions"`
}

// ListWalletsResponse is the response to `GET /v1/wallets`
type ListWalletsResponse struct {
	APIV1ResponseEnvelope
	Wallets []APIV1Wallet `json:"wallets"`
}

// CreateWalletResponse is the response to `POST /v1/wallet`
type CreateWalletResponse struct {
	APIV1ResponseEnvelope
	Wallet APIV1Wallet `json:"wallet"`
}

// InitWalletHandleResponse is the response to `POST /v1/wallet/init`
type InitWalletHandleResponse struct {
	APIV1ResponseEnvelope
	WalletHandleToken string `json:"wallet_handle_token"`
}

// ReleaseWalletHandleResponse is the response to `POST /v1/wallet/release`
type ReleaseWalletHandleResponse struct {
	APIV1ResponseEnvelope
}

// RenewWalletHandleResponse is the response to `POST /v1/wallet/renew`
type RenewWalletHandleResponse struct {
	APIV1ResponseEnvelope
	WalletHandle APIV1WalletHandle `json:"wallet_handle"`
}

// RenameWalletResponse is the response to `POST /v1/wallet/rename`
type RenameWalletResponse struct {
	APIV1ResponseEnvelope
	Wallet APIV1Wallet `json:"wallet"`
}

// GetWalletResponse is the response to `POST /v1/wallet/info`
type GetWalletResponse struct {
	APIV1ResponseEnvelope
	WalletHandle APIV1WalletHandle `json:"wallet_handle"`
}

// ExportMasterDerivationKeyResponse is the reponse to `POST /v1/master-key/export`
type ExportMasterDerivationKeyResponse struct {
	APIV1ResponseEnvelope
	MasterDerivationKey types.MasterDerivationKey `json:"master_derivation_key"`
}

// ImportKeyResponse is the repsonse to `POST /v1/key/import`
type ImportKeyResponse struct {
	APIV1ResponseEnvelope
	Address string `json:"address"`
}

// ExportKeyResponse is the reponse to `POST /v1/key/export`
type ExportKeyResponse struct {
	APIV1ResponseEnvelope
	PrivateKey ed25519.PrivateKey `json:"private_key"`
}

// GenerateKeyResponse is the response to `POST /v1/key`
type GenerateKeyResponse struct {
	APIV1ResponseEnvelope
	Address string `json:"address"`
}

// DeleteKeyResponse is the response to `DELETE /v1/key`
type DeleteKeyResponse struct {
	APIV1ResponseEnvelope
}

// ListKeysResponse is the response to `POST /v1/key/list`
type ListKeysResponse struct {
	APIV1ResponseEnvelope
	Addresses []string `json:"addresses"`
}

// SignTransactionResponse is the repsonse to `POST /v1/transaction/sign`
type SignTransactionResponse struct {
	APIV1ResponseEnvelope
	SignedTransaction []byte `json:"signed_transaction"`
}

// ListMultisigResponse is the response to `POST /v1/multisig/list`
type ListMultisigResponse struct {
	APIV1ResponseEnvelope
	Addresses []string `json:"addresses"`
}

// ImportMultisigResponse is the response to `POST /v1/multisig/import`
type ImportMultisigResponse struct {
	APIV1ResponseEnvelope
	Address string `json:"address"`
}

// ExportMultisigResponse is the response to `POST /v1/multisig/export`
type ExportMultisigResponse struct {
	APIV1ResponseEnvelope
	Version   uint8               `json:"multisig_version"`
	Threshold uint8               `json:"threshold"`
	PKs       []ed25519.PublicKey `json:"pks"`
}

// DeleteMultisigResponse is the response to POST /v1/multisig/delete`
type DeleteMultisigResponse struct {
	APIV1ResponseEnvelope
}

// SignMultisigTransactionResponse is the repsonse to `POST /v1/multisig/sign`
type SignMultisigTransactionResponse struct {
	APIV1ResponseEnvelope
	Multisig types.MultisigSig `json:"multisig"`
}
