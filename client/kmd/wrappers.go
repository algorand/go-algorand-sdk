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

// GenerateKey wraps APIV1POSTKeyRequest
func (kcl Client) GenerateKey(walletHandle []byte) (resp APIV1POSTKeyResponse, err error) {
	req := APIV1POSTKeyRequest{
		WalletHandleToken: string(walletHandle),
		DisplayMnemonic:   false,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// CreateWallet wraps APIV1POSTWalletRequest
func (kcl Client) CreateWallet(walletName []byte, walletDriverName string, walletPassword []byte, walletMDK types.MasterDerivationKey) (resp APIV1POSTWalletResponse, err error) {
	req := APIV1POSTWalletRequest{
		WalletName:          string(walletName),
		WalletDriverName:    walletDriverName,
		WalletPassword:      string(walletPassword),
		MasterDerivationKey: walletMDK,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// InitWallet wraps APIV1POSTWalletInitRequest
func (kcl Client) InitWallet(walletID []byte, walletPassword []byte) (resp APIV1POSTWalletInitResponse, err error) {
	req := APIV1POSTWalletInitRequest{
		WalletID:       string(walletID),
		WalletPassword: string(walletPassword),
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ReleaseWalletHandle wraps APIV1POSTWalletReleaseRequest
func (kcl Client) ReleaseWalletHandle(walletHandle []byte) (resp APIV1POSTWalletReleaseResponse, err error) {
	req := APIV1POSTWalletReleaseRequest{
		WalletHandleToken: string(walletHandle),
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ListKeys wraps APIV1POSTKeysListRequest
func (kcl Client) ListKeys(walletHandle []byte) (resp APIV1POSTKeysListResponse, err error) {
	req := APIV1POSTKeysListRequest{
		WalletHandleToken: string(walletHandle),
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ListMultisigAddrs wraps APIV1POSTMultisigListRequest
func (kcl Client) ListMultisigAddrs(walletHandle []byte) (resp APIV1POSTMultisigListResponse, err error) {
	req := APIV1POSTMultisigListRequest{
		WalletHandleToken: string(walletHandle),
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ImportMultisigAddr wraps APIV1POSTMultisigImportRequest
func (kcl Client) ImportMultisigAddr(walletHandle []byte, version, threshold uint8, pks []ed25519.PublicKey) (resp APIV1POSTMultisigImportResponse, err error) {
	req := APIV1POSTMultisigImportRequest{
		WalletHandleToken: string(walletHandle),
		Version:           version,
		Threshold:         threshold,
		PKs:               pks,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// MultisigSignTransaction wraps APIV1POSTMultisigTransactionSignRequest
func (kcl Client) MultisigSignTransaction(walletHandle, pw []byte, tx []byte, pk ed25519.PublicKey, partial types.MultisigSig) (resp APIV1POSTMultisigTransactionSignResponse, err error) {
	req := APIV1POSTMultisigTransactionSignRequest{
		WalletHandleToken: string(walletHandle),
		WalletPassword:    string(pw),
		Transaction:       tx,
		PublicKey:         pk,
		PartialMsig:       partial,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// RenewWalletHandle wraps APIV1POSTKeysListRequest
func (kcl Client) RenewWalletHandle(walletHandle []byte) (resp APIV1POSTWalletRenewResponse, err error) {
	req := APIV1POSTWalletRenewRequest{
		WalletHandleToken: string(walletHandle),
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ImportKey wraps APIV1POSTKeyImportRequest
func (kcl Client) ImportKey(walletHandle []byte, secretKey ed25519.PrivateKey) (resp APIV1POSTKeyImportResponse, err error) {
	req := APIV1POSTKeyImportRequest{
		WalletHandleToken: string(walletHandle),
		PrivateKey:        secretKey,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// ExportMasterDerivationKey wraps APIV1POSTMasterKeyExportRequest
func (kcl Client) ExportMasterDerivationKey(walletHandle []byte, walletPassword []byte) (resp APIV1POSTMasterKeyExportResponse, err error) {
	req := APIV1POSTMasterKeyExportRequest{
		WalletHandleToken: string(walletHandle),
		WalletPassword:    string(walletPassword),
	}
	err = kcl.DoV1Request(req, &resp)
	return
}

// SignTransaction wraps APIV1POSTTransactionSignRequest
func (kcl Client) SignTransaction(walletHandle, pw []byte, tx types.Transaction) (resp APIV1POSTTransactionSignResponse, err error) {
	txBytes := msgpack.Encode(tx)
	req := APIV1POSTTransactionSignRequest{
		WalletHandleToken: string(walletHandle),
		WalletPassword:    string(pw),
		Transaction:       txBytes,
	}
	err = kcl.DoV1Request(req, &resp)
	return
}
