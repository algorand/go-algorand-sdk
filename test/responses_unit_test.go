package test

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/encoding/json"
	"github.com/algorand/go-algorand-sdk/types"
)

var algodC *algod.Client
var indexerC *indexer.Client
var baselinePath string
var expectedStatus int
var response interface{}

// @unit
// @unit.responses

func mockHttpResponsesInLoadedFromWithStatus(jsonfile, loadedFrom string, status int) error {
	directory := path.Join("./features/resources/", loadedFrom)
	baselinePath = path.Join(directory, jsonfile)
	var err error
	expectedStatus = status
	err = mockHttpResponsesInLoadedFromHelper(jsonfile, directory, status)
	if err != nil {
		return err
	}
	algodC, err = algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	indexerC, err = indexer.MakeClient(mockServer.URL, "")
	return err
}

func weMakeAnyCallTo(client /* algod/indexer */, endpoint string) (err error) {
	var round uint64
	var something interface{}

	switch client {
	case "indexer":
		switch endpoint {
		case "lookupAccountByID":
			round, something, err = indexerC.LookupAccountByID("").Do(context.Background())
			response = models.AccountResponse{
				CurrentRound: round,
				Account:      something.(models.Account),
			}
		case "searchForAccounts":
			response, err = indexerC.SearchAccounts().Do(context.Background())
		case "lookupApplicationByID":
			response, err = indexerC.LookupApplicationByID(10).Do(context.Background())
		case "searchForApplications":
			response, err = indexerC.SearchForApplications().Do(context.Background())
		case "lookupAssetBalances":
			response, err = indexerC.LookupAssetBalances(10).Do(context.Background())
		case "lookupAssetByID":
			round, something, err = indexerC.LookupAssetByID(10).Do(context.Background())
			response = models.AssetResponse{
				CurrentRound: round,
				Asset:        something.(models.Asset),
			}
		case "searchForAssets":
			response, err = indexerC.SearchForAssets().Do(context.Background())
		case "lookupAccountTransactions":
			response, err = indexerC.LookupAccountTransactions("").Do(context.Background())
		case "lookupAssetTransactions":
			response, err = indexerC.LookupAssetTransactions(10).Do(context.Background())
		case "searchForTransactions":
			response, err = indexerC.SearchForTransactions().Do(context.Background())
		case "lookupBlock":
			response, err = indexerC.LookupBlock(10).Do(context.Background())
		case "any":
			// This is an error case
			// pickup the error as the response
			_, response = indexerC.SearchForTransactions().Do(context.Background())
		default:
			err = fmt.Errorf("unknown indexer endpoint: %s", endpoint)
		}
	case "algod":
		switch endpoint {
		case "GetStatus":
			response, err = algodC.Status().Do(context.Background())
		case "GetBlock":
			response, err = algodC.Block(10).Do(context.Background())
		case "WaitForBlock":
			response, err = algodC.StatusAfterBlock(10).Do(context.Background())
		case "TealCompile":
			response, err = algodC.TealCompile([]byte{}).Do(context.Background())
		case "RawTransaction":
			var returnedTxid string
			returnedTxid, err = algodC.SendRawTransaction([]byte{}).Do(context.Background())
			response = txidresponse{TxID: returnedTxid}
		case "GetSupply":
			response, err = algodC.Supply().Do(context.Background())
		case "TransactionParams":
			var sParams types.SuggestedParams
			sParams, err = algodC.SuggestedParams().Do(context.Background())
			response = models.TransactionParametersResponse{
				ConsensusVersion: sParams.ConsensusVersion,
				Fee:              uint64(sParams.Fee),
				GenesisId:        sParams.GenesisID,
				GenesisHash:      sParams.GenesisHash,
				LastRound:        uint64(sParams.FirstRoundValid),
				MinFee:           sParams.MinFee,
			}
		case "GetAccountInformation":
			response, err = algodC.AccountInformation("acct").Do(context.Background())
		case "GetApplicationByID":
			response, err = algodC.GetApplicationByID(10).Do(context.Background())
		case "GetAssetByID":
			response, err = algodC.GetAssetByID(10).Do(context.Background())
		case "PendingTransactionInformation":
			response, _, err = algodC.PendingTransactionInformation("transaction").Do(context.Background())
		case "GetPendingTransactions":
			var total uint64
			var top []types.SignedTxn
			total, top, err = algodC.PendingTransactions().Do(context.Background())
			response = models.PendingTransactionsResponse{
				TopTransactions:   top,
				TotalTransactions: total,
			}
		case "GetPendingTransactionsByAddress":
			var total uint64
			var top []types.SignedTxn
			total, top, err = algodC.PendingTransactionsByAddress("address").Do(context.Background())
			response = models.PendingTransactionsResponse{
				TopTransactions:   top,
				TotalTransactions: total,
			}
		case "DryRun":
			response, err = algodC.TealDryrun(models.DryrunRequest{}).Do(context.Background())
		case "Proof":
			response, err = algodC.GetProof(10, "asdf").Do(context.Background())
		case "any":
			// This is an error case
			// pickup the error as the response
			_, response = indexerC.SearchForTransactions().Do(context.Background())
		default:
			err = fmt.Errorf("unknown algod endpoint: %s", endpoint)
		}
	}
	return err
}

type txidresponse struct {
	TxID string `json:"txId"`
}

func theParsedResponseShouldEqualTheMockResponse() error {
	var responseJson string

	if expectedStatus != 200 {
		responseJson = response.(error).Error()
		// The error message is not a well formed Json.
		// Verify the expected status code, and remove the json corrupting message
		statusCode := fmt.Sprintf("%d", expectedStatus)
		if !strings.Contains(responseJson, statusCode) {
			return fmt.Errorf("Expected error code: %d, got otherwise", expectedStatus)
		}
		parts := strings.SplitAfterN(responseJson, ":", 2)
		responseJson = parts[1]
	} else {
		responseJson = string(json.Encode(response))
	}

	return VerifyResponse(baselinePath, responseJson)
}

func ResponsesContext(s *godog.Suite) {
	s.Step(`^mock http responses in "([^"]*)" loaded from "([^"]*)" with status (\d+)\.$`, mockHttpResponsesInLoadedFromWithStatus)
	s.Step(`^we make any "([^"]*)" call to "([^"]*)"\.$`, weMakeAnyCallTo)
	s.Step(`^the parsed response should equal the mock response\.$`, theParsedResponseShouldEqualTheMockResponse)
}
