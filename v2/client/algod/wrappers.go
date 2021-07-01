package algod

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/algod/models"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Status retrieves the StatusResponse from the running node
// the StatusResponse includes data like the consensus version and current round
func (client Client) Status(headers ...*Header) (response models.NodeStatus, err error) {
	err = client.get(&response, "/status", nil, headers)
	return
}

// HealthCheck does a health check on the the potentially running node,
// returning an error if the API is down
func (client Client) HealthCheck(headers ...*Header) error {
	return client.get(nil, "/health", nil, headers)
}

// StatusAfterBlock waits for a block to occur then returns the StatusResponse after that block
// blocks on the node end
func (client Client) StatusAfterBlock(blockNum uint64, headers ...*Header) (response models.NodeStatus, err error) {
	err = client.get(&response, fmt.Sprintf("/status/wait-for-block-after/%d", blockNum), nil, headers)
	return
}

type pendingTransactionsParams struct {
	Max uint64 `url:"max"`
}

// GetPendingTransactions asks algod for a snapshot of current pending txns on the node, bounded by maxTxns.
// If maxTxns = 0, fetches as many transactions as possible.
func (client Client) GetPendingTransactions(maxTxns uint64, headers ...*Header) (response models.PendingTransactions, err error) {
	err = client.get(&response, fmt.Sprintf("/transactions/pending"), pendingTransactionsParams{maxTxns}, headers)
	return
}

// Versions retrieves the VersionResponse from the running node
// the VersionResponse includes data like version number and genesis ID
func (client Client) Versions(headers ...*Header) (response models.Version, err error) {
	err = client.get(&response, "/versions", nil, headers)
	return
}

// LedgerSupply gets the supply details for the specified node's Ledger
func (client Client) LedgerSupply(headers ...*Header) (response models.Supply, err error) {
	err = client.get(&response, "/ledger/supply", nil, headers)
	return
}

type transactionsByAddrParams struct {
	FirstRound uint64 `url:"firstRound,omitempty"`
	LastRound  uint64 `url:"lastRound,omitempty"`
	FromDate   string `url:"fromDate,omitempty"`
	ToDate     string `url:"toDate,omitempty"`
	Max        uint64 `url:"max,omitempty"`
}

// TransactionsByAddr returns all transactions for a PK [addr] in the [first,
// last] rounds range.
func (client Client) TransactionsByAddr(addr string, first, last uint64, headers ...*Header) (response models.TransactionList, err error) {
	params := transactionsByAddrParams{FirstRound: first, LastRound: last}
	err = client.get(&response, fmt.Sprintf("/account/%s/transactions", addr), params, headers)
	return
}

// TransactionsByAddrLimit returns the last [limit] number of transaction for a PK [addr].
func (client Client) TransactionsByAddrLimit(addr string, limit uint64, headers ...*Header) (response models.TransactionList, err error) {
	params := transactionsByAddrParams{Max: limit}
	err = client.get(&response, fmt.Sprintf("/account/%s/transactions", addr), params, headers)
	return
}

// TransactionsByAddrForDate returns all transactions for a PK [addr] in the [first,
// last] date range. Dates are of the form "2006-01-02".
func (client Client) TransactionsByAddrForDate(addr string, first, last string, headers ...*Header) (response models.TransactionList, err error) {
	params := transactionsByAddrParams{FromDate: first, ToDate: last}
	err = client.get(&response, fmt.Sprintf("/account/%s/transactions", addr), params, headers)
	return
}

// AccountInformation also gets the AccountInformationResponse associated with the passed address
func (client Client) AccountInformation(address string, headers ...*Header) (response models.Account, err error) {
	err = client.get(&response, fmt.Sprintf("/account/%s", address), nil, headers)
	return
}

// AssetInformation also gets the AssetInformationResponse associated with the passed asset creator and index
func (client Client) AssetInformation(index uint64, headers ...*Header) (response models.AssetParams, err error) {
	err = client.get(&response, fmt.Sprintf("/asset/%d", index), nil, headers)
	return
}

// TransactionInformation gets information about a specific transaction involving a specific account
// it will only return information about transactions submitted to the node queried
func (client Client) TransactionInformation(accountAddress, transactionID string, headers ...*Header) (response models.Transaction, err error) {
	transactionID = stripTransaction(transactionID)
	err = client.get(&response, fmt.Sprintf("/account/%s/transaction/%s", accountAddress, transactionID), nil, headers)
	return
}

// PendingTransactionInformation gets information about a recently issued
// transaction.  There are several cases when this might succeed:
//
// - transaction committed (CommittedRound > 0)
// - transaction still in the pool (CommittedRound = 0, PoolError = "")
// - transaction removed from pool due to error (CommittedRound = 0, PoolError != "")
//
// Or the transaction may have happened sufficiently long ago that the
// node no longer remembers it, and this will return an error.
func (client Client) PendingTransactionInformation(transactionID string, headers ...*Header) (response models.Transaction, err error) {
	transactionID = stripTransaction(transactionID)
	err = client.get(&response, fmt.Sprintf("/transactions/pending/%s", transactionID), nil, headers)
	return
}

// TransactionByID gets a transaction by its ID. Works only if the indexer is enabled on the node
// being queried.
func (client Client) TransactionByID(transactionID string, headers ...*Header) (response models.Transaction, err error) {
	transactionID = stripTransaction(transactionID)
	err = client.get(&response, fmt.Sprintf("/transaction/%s", transactionID), nil, headers)
	return
}

// SuggestedFee gets the recommended transaction fee from the node
func (client Client) SuggestedFee(headers ...*Header) (response models.TransactionFee, err error) {
	err = client.get(&response, "/transactions/fee", nil, headers)
	return
}

// SuggestedParams gets the suggested transaction parameters
func (client Client) SuggestedParams(headers ...*Header) (response models.TransactionParams, err error) {
	err = client.get(&response, "/transactions/params", nil, headers)
	return
}

// BuildSuggestedParams gets the suggested transaction parameters and
// builds a types.SuggestedParams to pass to transaction builders (see package future)
func (client Client) BuildSuggestedParams(headers ...*Header) (response types.SuggestedParams, err error) {
	var httpResponse models.TransactionParams
	err = client.get(&httpResponse, "/transactions/params", nil, headers)
	response.FlatFee = false
	response.Fee = types.MicroAlgos(httpResponse.Fee)
	response.GenesisID = httpResponse.GenesisID
	response.GenesisHash = httpResponse.GenesisHash
	response.FirstRoundValid = types.Round(httpResponse.LastRound)
	response.LastRoundValid = types.Round(httpResponse.LastRound + 1000)
	response.ConsensusVersion = httpResponse.ConsensusVersion
	return
}

// SendRawTransaction gets the bytes of a SignedTxn and broadcasts it to the network
func (client Client) SendRawTransaction(stx []byte, headers ...*Header) (response models.TransactionID, err error) {
	// Set default Content-Type, if not the user didn't specify it.
	addContentType := true
	for _, header := range headers {
		if strings.ToLower(header.Key) == "content-type" {
			addContentType = false
			break
		}
	}
	if addContentType {
		headers = append(headers, &Header{"Content-Type", "application/x-binary"})
	}
	err = client.post(&response, "/transactions", stx, headers)
	return
}

// Block gets the block info for the given round
func (client Client) Block(round uint64, headers ...*Header) (response models.Block, err error) {
	err = client.get(&response, fmt.Sprintf("/block/%d", round), nil, headers)
	return
}

func responseReadAll(resp *http.Response, maxContentLength int64) (body []byte, err error) {
	if resp.ContentLength > 0 {
		// more efficient path if we know the ContentLength
		if maxContentLength > 0 && resp.ContentLength > maxContentLength {
			return nil, errors.New("Content too long")
		}
		body = make([]byte, resp.ContentLength)
		_, err = io.ReadFull(resp.Body, body)
		return
	}

	return ioutil.ReadAll(resp.Body)
}

// BlockRaw gets the raw block msgpack bytes for the given round
func (client Client) BlockRaw(round uint64, headers ...*Header) (blockbytes []byte, err error) {
	var resp *http.Response
	request := struct {
		Raw string `url:"raw"`
	}{Raw: "1"}
	resp, err = client.submitFormRaw(fmt.Sprintf("/block/%d", round), request, "GET", false, headers)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// Current blocks are about 1MB. 10MB should be a safe backstop.
	return responseReadAll(resp, 10000000)
}

func (client Client) doGetWithQuery(ctx context.Context, path string, queryArgs map[string]string) (result string, err error) {
	queryURL := client.serverURL
	queryURL.Path = path

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	for k, v := range queryArgs {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set(authHeader, client.apiToken)
	for _, header := range client.headers {
		req.Header.Add(header.Key, header.Value)
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = extractError(resp)
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = string(bytes)
	return
}
