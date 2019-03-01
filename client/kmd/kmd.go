package kmd

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/algorand/go-algorand-sdk/encoding/json"
)

const (
	timeoutSecs    = 120
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
	case ListWalletsRequest:
		reqPath = "v1/wallets"
		reqMethod = "GET"
	case CreateWalletRequest:
		reqPath = "v1/wallet"
		reqMethod = "POST"
	case InitWalletHandleRequest:
		reqPath = "v1/wallet/init"
		reqMethod = "POST"
	case ReleaseWalletHandleRequest:
		reqPath = "v1/wallet/release"
		reqMethod = "POST"
	case RenewWalletHandleRequest:
		reqPath = "v1/wallet/renew"
		reqMethod = "POST"
	case RenameWalletRequest:
		reqPath = "v1/wallet/rename"
		reqMethod = "POST"
	case GetWalletRequest:
		reqPath = "v1/wallet/info"
		reqMethod = "POST"
	case ExportMasterDerivationKeyRequest:
		reqPath = "v1/master-key/export"
		reqMethod = "POST"
	case ImportKeyRequest:
		reqPath = "v1/key/import"
		reqMethod = "POST"
	case ExportKeyRequest:
		reqPath = "v1/key/export"
		reqMethod = "POST"
	case GenerateKeyRequest:
		reqPath = "v1/key"
		reqMethod = "POST"
	case DeleteKeyRequest:
		reqPath = "v1/key"
		reqMethod = "DELETE"
	case ListKeysRequest:
		reqPath = "v1/key/list"
		reqMethod = "POST"
	case SignTransactionRequest:
		reqPath = "v1/transaction/sign"
		reqMethod = "POST"
	case ListMultisigRequest:
		reqPath = "v1/multisig/list"
		reqMethod = "POST"
	case ImportMultisigRequest:
		reqPath = "v1/multisig/import"
		reqMethod = "POST"
	case ExportMultisigRequest:
		reqPath = "v1/multisig/export"
		reqMethod = "POST"
	case SignMultisigTransactionRequest:
		reqPath = "v1/multisig/sign"
		reqMethod = "POST"
	}
	return
}
