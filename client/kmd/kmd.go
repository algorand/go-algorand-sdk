package kmd

import (
	"bytes"
	"fmt"
        "net/http"
        "time"

	"github.com/algorand/go-algorand-sdk/encoding/json"
)

const (
        timeoutSecs = 120
	kmdTokenHeader = "X-KMD-API-Token"
)

// Client is the client used to interact with the kmd API
type Client struct {
        httpClient http.Client
        apiToken   string
        address    string
}

func makeHTTPClient() http.Client {
        client := http.Client{
                Timeout: timeoutSecs * time.Second,
        }
        return client
}

// MakeClient instantiates a Client for the given address and apiToken
func MakeClient(address string, apiToken string) (Client, error) {
        kcl := Client{
                httpClient: makeHTTPClient(),
                apiToken:   apiToken,
                address:    address,
        }
        return kcl, nil
}

// DoV1Request accepts a request from kmdapi/requests and
func (kcl Client) DoV1Request(req APIV1Request, resp APIV1Response) error {
	var body []byte

	// Get the path and method for this request type
	reqPath, reqMethod, err := getPathAndMethod(req)
	if err != nil {
		return err
	}

	// Encode the request
	body = json.Encode(req)
	fullPath := fmt.Sprintf("%s/%s", kcl.address, reqPath)
	hreq, err := http.NewRequest(reqMethod, fullPath, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// Add the auth token
	hreq.Header.Add(kmdTokenHeader, kcl.apiToken)

	// Send the request
	hresp, err := kcl.httpClient.Do(hreq)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(hresp.Body)
	err = decoder.Decode(resp)
	hresp.Body.Close()
	if err != nil {
		return err
	}

	// Check if this was an error response
	err = resp.GetError()
	if err != nil {
		return err
	}

	return nil
}

// getPathAndMethod infers the request path and method from the request type
func getPathAndMethod(req APIV1Request) (reqPath string, reqMethod string, err error) {
	switch req.(type) {
	default:
		err = fmt.Errorf("unknown request type")
	case VersionsRequest:
		reqPath = "versions"
		reqMethod = "GET"
	case APIV1GETWalletsRequest:
		reqPath = "v1/wallets"
		reqMethod = "GET"
	case APIV1POSTWalletRequest:
		reqPath = "v1/wallet"
		reqMethod = "POST"
	case APIV1POSTWalletInitRequest:
		reqPath = "v1/wallet/init"
		reqMethod = "POST"
	case APIV1POSTWalletReleaseRequest:
		reqPath = "v1/wallet/release"
		reqMethod = "POST"
	case APIV1POSTWalletRenewRequest:
		reqPath = "v1/wallet/renew"
		reqMethod = "POST"
	case APIV1POSTWalletRenameRequest:
		reqPath = "v1/wallet/rename"
		reqMethod = "POST"
	case APIV1POSTWalletInfoRequest:
		reqPath = "v1/wallet/info"
		reqMethod = "POST"
	case APIV1POSTMasterKeyExportRequest:
		reqPath = "v1/master_key/export"
		reqMethod = "POST"
	case APIV1POSTKeyImportRequest:
		reqPath = "v1/key/import"
		reqMethod = "POST"
	case APIV1POSTKeyExportRequest:
		reqPath = "v1/key/export"
		reqMethod = "POST"
	case APIV1POSTKeyRequest:
		reqPath = "v1/key"
		reqMethod = "POST"
	case APIV1DELETEKeyRequest:
		reqPath = "v1/key"
		reqMethod = "DELETE"
	case APIV1POSTKeysListRequest:
		reqPath = "v1/keys/list"
		reqMethod = "POST"
	case APIV1POSTTransactionSignRequest:
		reqPath = "v1/transaction/sign"
		reqMethod = "POST"
	case APIV1POSTMultisigListRequest:
		reqPath = "v1/multisig/list"
		reqMethod = "POST"
	case APIV1POSTMultisigImportRequest:
		reqPath = "v1/multisig/import"
		reqMethod = "POST"
	case APIV1POSTMultisigExportRequest:
		reqPath = "v1/multisig/export"
		reqMethod = "POST"
	case APIV1POSTMultisigTransactionSignRequest:
		reqPath = "v1/multisig/sign"
		reqMethod = "POST"
	}
	return
}
