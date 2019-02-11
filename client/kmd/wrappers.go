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

// ListWallets wraps APIV1GETWalletsRequest
func (kcl Client) ListWallets() (resp APIV1GETWalletsResponse, err error) {
	req := APIV1GETWalletsRequest{}
	err = kcl.DoV1Request(req, &resp)
	return
}

// CreateWallet wraps APIV1POSTWalletRequest
func (kcl Client) CreateWallet(walletName, walletDriverName, walletPassword string, walletMDK types.MasterDerivationKey) (resp APIV1POSTWalletResponse, err error) {
	req := APIV1POSTWalletRequest{
		WalletName:          walletName,
		WalletDriverName:    walletDriverName,
		WalletPassword:      walletPassword,
		MasterDerivationKey: walletMDK,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// InitWallet wraps APIV1POSTWalletInitRequest
func (kcl Client) InitWallet(walletID, walletPassword string) (resp APIV1POSTWalletInitResponse, err error) {
	req := APIV1POSTWalletInitRequest{
		WalletID:       walletID,
		WalletPassword: walletPassword,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ReleaseWalletHandle wraps APIV1POSTWalletReleaseRequest
func (kcl Client) ReleaseWalletHandle(walletHandle string) (resp APIV1POSTWalletReleaseResponse, err error) {
	req := APIV1POSTWalletReleaseRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// RenewWalletHandle wraps APIV1POSTKeysListRequest
func (kcl Client) RenewWalletHandle(walletHandle string) (resp APIV1POSTWalletRenewResponse, err error) {
	req := APIV1POSTWalletRenewRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// RenameWallet wraps APIV1POSTWalletRenameRequest
func (kcl Client) RenameWallet(walletID, walletPassword, newWalletName string) (resp APIV1POSTWalletRenameResponse, err error) {
	req := APIV1POSTWalletRenameRequest{
		WalletID: walletID,
		WalletPassword: walletPassword,
		NewWalletName: newWalletName,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// WalletInfo wraps APIV1POSTWalletInfoRequest
func (kcl Client) WalletInfo(walletHandle string) (resp APIV1POSTWalletInfoResponse, err error) {
	req := APIV1POSTWalletInfoRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ExportMasterDerivationKey wraps APIV1POSTMasterKeyExportRequest
func (kcl Client) ExportMasterDerivationKey(walletHandle, walletPassword string) (resp APIV1POSTMasterKeyExportResponse, err error) {
	req := APIV1POSTMasterKeyExportRequest{
		WalletHandleToken: walletHandle,
		WalletPassword:    walletPassword,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ImportKey wraps APIV1POSTKeyImportRequest
func (kcl Client) ImportKey(walletHandle string, secretKey ed25519.PrivateKey) (resp APIV1POSTKeyImportResponse, err error) {
	req := APIV1POSTKeyImportRequest{
		WalletHandleToken: walletHandle,
		PrivateKey:        secretKey,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ExportKey wraps APIV1POSTKeyExportRequest
func (kcl Client) ExportKey(walletHandle, walletPassword, addr string) (resp APIV1POSTKeyExportResponse, err error) {
	req := APIV1POSTKeyExportRequest{
		WalletHandleToken: walletHandle,
		WalletPassword: walletPassword,
		Address: addr,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// GenerateKey wraps APIV1POSTKeyRequest
func (kcl Client) GenerateKey(walletHandle string) (resp APIV1POSTKeyResponse, err error) {
	req := APIV1POSTKeyRequest{
		WalletHandleToken: walletHandle,
		DisplayMnemonic:   false,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// DeleteKey wraps APIV1DELETEKeyRequest
func (kcl Client) DeleteKey(walletHandle, walletPassword, addr string) (resp APIV1DELETEKeyResponse, err error) {
	req := APIV1DELETEKeyRequest{
		WalletHandleToken: walletHandle,
		WalletPassword: walletPassword,
		Address: addr,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ListKeys wraps APIV1POSTKeysListRequest
func (kcl Client) ListKeys(walletHandle string) (resp APIV1POSTKeysListResponse, err error) {
	req := APIV1POSTKeysListRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// SignTransaction wraps APIV1POSTTransactionSignRequest
func (kcl Client) SignTransaction(walletHandle, pw string, tx types.Transaction) (resp APIV1POSTTransactionSignResponse, err error) {
	txBytes := msgpack.Encode(tx)
	req := APIV1POSTTransactionSignRequest{
		WalletHandleToken: walletHandle,
		WalletPassword:    pw,
		Transaction:       txBytes,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ListMultisigAddrs wraps APIV1POSTMultisigListRequest
func (kcl Client) ListMultisigAddrs(walletHandle string) (resp APIV1POSTMultisigListResponse, err error) {
	req := APIV1POSTMultisigListRequest{
		WalletHandleToken: walletHandle,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ImportMultisigAddr wraps APIV1POSTMultisigImportRequest
func (kcl Client) ImportMultisigAddr(walletHandle string, version, threshold uint8, pks []ed25519.PublicKey) (resp APIV1POSTMultisigImportResponse, err error) {
	req := APIV1POSTMultisigImportRequest{
		WalletHandleToken: walletHandle,
		Version:           version,
		Threshold:         threshold,
		PKs:               pks,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ExportMultisigAddr wraps APIV1POSTMultisigExportRequest
func (kcl Client) ExportMultisigAddr(walletHandle, walletPassword, addr string) (resp APIV1POSTMultisigExportResponse, err error) {
	req := APIV1POSTMultisigExportRequest{
		WalletHandleToken: walletHandle,
		WalletPassword: walletPassword,
		Address: addr,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// MultisigSignTransaction wraps APIV1POSTMultisigTransactionSignRequest
func (kcl Client) MultisigSignTransaction(walletHandle, pw string, tx []byte, pk ed25519.PublicKey, partial types.MultisigSig) (resp APIV1POSTMultisigTransactionSignResponse, err error) {
	req := APIV1POSTMultisigTransactionSignRequest{
		WalletHandleToken: walletHandle,
		WalletPassword:    pw,
		Transaction:       tx,
		PublicKey:         pk,
		PartialMsig:       partial,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}
