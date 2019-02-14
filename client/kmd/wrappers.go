package kmd

import (
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

// Version wraps VersionsRequest
func (kcl Client) Version() (resp VersionsResponse, err error) {
	req := VersionsRequest{}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ListWallets wraps ListWalletsRequest
func (kcl Client) ListWallets() (resp ListWalletsResponse, err error) {
	req := ListWalletsRequest{}
	err = kcl.DoV1Request(req, &resp)
	return
}

// CreateWallet wraps CreateWalletRequest
func (kcl Client) CreateWallet(walletName, walletPassword, walletDriverName string, walletMDK types.MasterDerivationKey) (resp CreateWalletResponse, err error) {
	req := CreateWalletRequest{
		WalletName:          walletName,
		WalletDriverName:    walletDriverName,
		WalletPassword:      walletPassword,
		MasterDerivationKey: walletMDK,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// InitWalletHandle wraps InitWalletHandleRequest
func (kcl Client) InitWalletHandle(walletID, walletPassword string) (resp InitWalletHandleResponse, err error) {
	req := InitWalletHandleRequest{
		WalletID:       walletID,
		WalletPassword: walletPassword,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ReleaseWalletHandleHandle wraps ReleaseWalletHandleRequest
func (kcl Client) ReleaseWalletHandleHandle(walletHandle string) (resp ReleaseWalletHandleResponse, err error) {
	req := ReleaseWalletHandleRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// RenewWalletHandle wraps RenewWalletHandleRequest
func (kcl Client) RenewWalletHandle(walletHandle string) (resp RenewWalletHandleResponse, err error) {
	req := RenewWalletHandleRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// RenameWallet wraps RenameWalletRequest
func (kcl Client) RenameWallet(walletID, walletPassword, newWalletName string) (resp RenameWalletResponse, err error) {
	req := RenameWalletRequest{
		WalletID: walletID,
		WalletPassword: walletPassword,
		NewWalletName: newWalletName,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// GetWallet wraps GetWalletRequest
func (kcl Client) GetWallet(walletHandle string) (resp GetWalletResponse, err error) {
	req := GetWalletRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ExportMasterDerivationKey wraps ExportMasterDerivationKeyRequest
func (kcl Client) ExportMasterDerivationKey(walletHandle, walletPassword string) (resp ExportMasterDerivationKeyResponse, err error) {
	req := ExportMasterDerivationKeyRequest{
		WalletHandleToken: walletHandle,
		WalletPassword:    walletPassword,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ImportKey wraps ImportKeyRequest
func (kcl Client) ImportKey(walletHandle string, secretKey ed25519.PrivateKey) (resp ImportKeyResponse, err error) {
	req := ImportKeyRequest{
		WalletHandleToken: walletHandle,
		PrivateKey:        secretKey,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ExportKey wraps ExportKeyRequest
func (kcl Client) ExportKey(walletHandle, walletPassword, addr string) (resp ExportKeyResponse, err error) {
	req := ExportKeyRequest{
		WalletHandleToken: walletHandle,
		WalletPassword: walletPassword,
		Address: addr,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// GenerateKey wraps GenerateKeyRequest
func (kcl Client) GenerateKey(walletHandle string) (resp GenerateKeyResponse, err error) {
	req := GenerateKeyRequest{
		WalletHandleToken: walletHandle,
		DisplayMnemonic:   false,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// DeleteKey wraps DeleteKeyRequest
func (kcl Client) DeleteKey(walletHandle, walletPassword, addr string) (resp DeleteKeyResponse, err error) {
	req := DeleteKeyRequest{
		WalletHandleToken: walletHandle,
		WalletPassword: walletPassword,
		Address: addr,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ListKeys wraps ListKeysRequest
func (kcl Client) ListKeys(walletHandle string) (resp ListKeysResponse, err error) {
	req := ListKeysRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// SignTransaction wraps SignTransactionRequest
func (kcl Client) SignTransaction(walletHandle, pw string, tx types.Transaction) (resp SignTransactionResponse, err error) {
	txBytes := msgpack.Encode(tx)
	req := SignTransactionRequest{
		WalletHandleToken: walletHandle,
		WalletPassword:    pw,
		Transaction:       txBytes,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ListMultisig wraps ListMultisigRequest
func (kcl Client) ListMultisig(walletHandle string) (resp ListMultisigResponse, err error) {
	req := ListMultisigRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ImportMultisig wraps ImportMultisigRequest
func (kcl Client) ImportMultisig(walletHandle string, version, threshold uint8, pks []ed25519.PublicKey) (resp ImportMultisigResponse, err error) {
	req := ImportMultisigRequest{
		WalletHandleToken: walletHandle,
		Version:           version,
		Threshold:         threshold,
		PKs:               pks,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ExportMultisig wraps ExportMultisigRequest
func (kcl Client) ExportMultisig(walletHandle, walletPassword, addr string) (resp ExportMultisigResponse, err error) {
	req := ExportMultisigRequest{
		WalletHandleToken: walletHandle,
		WalletPassword: walletPassword,
		Address: addr,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// MultisigSignTransaction wraps SignMultisigTransactionRequest
func (kcl Client) MultisigSignTransaction(walletHandle, pw string, tx []byte, pk ed25519.PublicKey, partial types.MultisigSig) (resp SignMultisigTransactionResponse, err error) {
	req := SignMultisigTransactionRequest{
		WalletHandleToken: walletHandle,
		WalletPassword:    pw,
		Transaction:       tx,
		PublicKey:         pk,
		PartialMsig:       partial,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}
