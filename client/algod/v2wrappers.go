package algod

import (
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
)

// TODO ejr talk to indexer team about Shutdown() impl
// TODO ejr talk to indexer team about RegisterParticipationKeys impl

// GetPendingTransactions asks algod for a snapshot of current pending txns on the node, bounded by maxTxns.
// If maxTxns = 0, fetches as many transactions as possible.
func (client Client) GetPendingTransactions(maxTxns uint64, headers ...*Header) (response models.PendingTransactions, err error) {
	err = client.get(&response, fmt.Sprintf("/transactions/pending"), pendingTransactionsParams{maxTxns}, headers)
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
// TODO ejr add Max query
func (client Client) PendingTransactionInformation(transactionID string, headers ...*Header) (response models.Transaction, err error) {
	transactionID = stripTransaction(transactionID)
	err = client.get(&response, fmt.Sprintf("/transactions/pending/%s", transactionID), nil, headers)
	return
}

// SendRawTransaction gets the bytes of a SignedTxn and broadcasts it to the network
func (client Client) SendRawTransaction(stx []byte, headers ...*Header) (response models.TransactionID, err error) {
	err = client.post(&response, "/transactions", stx, headers)
	return
}

// TODO ejr these params may be out of date
type transactionsByAddrParams struct {
	FirstRound uint64 `url:"firstRound,omitempty"`
	LastRound  uint64 `url:"lastRound,omitempty"`
	FromDate   string `url:"fromDate,omitempty"`
	ToDate     string `url:"toDate,omitempty"`
	Max        uint64 `url:"max,omitempty"`
}

// TransactionsByAddrLimit returns the last [limit] number of transaction for a PK [addr].
func (client Client) TransactionsByAddrLimit(addr string, limit uint64, headers ...*Header) (response models.TransactionList, err error) {
	params := transactionsByAddrParams{Max: limit}
	err = client.get(&response, fmt.Sprintf("/account/%s/transactions", addr), params, headers)
	return
}

// Status retrieves the StatusResponse from the running node
// the StatusResponse includes data like the consensus version and current round
func (client Client) Status(headers ...*Header) (response models.NodeStatus, err error) {
	err = client.get(&response, "/status", nil, headers)
	return
}

// LedgerSupply gets the supply details for the specified node's Ledger
func (client Client) LedgerSupply(headers ...*Header) (response models.Supply, err error) {
	err = client.get(&response, "/ledger/supply", nil, headers)
	return
}

// StatusAfterBlock waits for a block to occur then returns the StatusResponse after that block
// blocks on the node end
func (client Client) StatusAfterBlock(blockNum uint64, headers ...*Header) (response models.NodeStatus, err error) {
	err = client.get(&response, fmt.Sprintf("/status/wait-for-block-after/%d", blockNum), nil, headers)
	return
}

// AccountInformation also gets the AccountInformationResponse associated with the passed address
func (client Client) AccountInformation(address string, headers ...*Header) (response models.Account, err error) {
	err = client.get(&response, fmt.Sprintf("/account/%s", address), nil, headers)
	return
}

// Block gets the block info for the given round // TODO ejr add Raw query
func (client Client) Block(round uint64, headers ...*Header) (response models.Block, err error) {
	err = client.get(&response, fmt.Sprintf("/block/%d", round), nil, headers)
	return
}

// SuggestedParams gets the suggested transaction parameters
// TODO ejr reflect go suggested params change here once changes have settled
func (client Client) SuggestedParams(headers ...*Header) (response models.TransactionParams, err error) {
	err = client.get(&response, "/transactions/params", nil, headers)
	return
}
