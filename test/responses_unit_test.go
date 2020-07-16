package test

import (
	"context"
	"fmt"
	"path"
	
	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"	
)

var algodC *algod.Client
var indexerC *indexer.Client

func mockHttpResponsesInLoadedFromWithStatus(
	jsonfile, loadedFrom /* generated_responses*/ string, status int) error {
	directory := path.Join("./features/resources/", loadedFrom)
	var err error
	err = mockHttpResponsesInLoadedFromHelper(jsonfile, directory)
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

var response interface{}

func weMakeAnyCallTo(client /* algod/indexer */, endpoint string) (err error) {
	switch client {
	case "indexer":
		switch endpoint {
		case "lookupAccountByID":
			_, response, err = indexerC.LookupAccountByID("").Do(context.Background())
		case "searchForAccounts":
			response, err = indexerC.SearchAccounts().Do(context.Background())
		case "lookupApplicationByID":
			response, err = indexerC.LookupApplicationByID(10).Do(context.Background())
		case "searchForApplications":
			response, err = indexerC.SearchForApplications().Do(context.Background())
		case "lookupAssetBalances":
			response, err = indexerC.LookupAssetBalances(10).Do(context.Background())
		case "lookupAssetByID":
			_, response, err = indexerC.LookupAssetByID(10).Do(context.Background())
		case "searchForAssets":
			_, response, err = indexerC.SearchForAssets().Do(context.Background())
		case "lookupAccountTransactions":
			response, err = indexerC.LookupAccountTransactions("").Do(context.Background())
		case "lookupAssetTransactions":
			response, err = indexerC.LookupAssetTransactions(10).Do(context.Background())
		case "searchForTransactions":
			response, err = indexerC.SearchForTransactions().Do(context.Background())
		default:
			fmt.Errorf("unknown endpoint: %s", endpoint)
		}
		/*	case "algod":
		switch endpoint {
		case "GetStatus                
		case "WaitForBlock             
		case "TealCompile              
		case "RawTransaction           
		case "GetSupply                
		case "TransactionParams        
		case "GetApplicationByID       
		case "GetAssetByID
		default:
		}*/
	}
	return err
}

func theParsedResponseShouldEqualTheMockResponse() error {
	return godog.ErrPending
}

func ResponsesContext(s *godog.Suite) {
	s.Step(`^mock http responses in "([^"]*)" loaded from "([^"]*)" with status (\d+)\.$`, mockHttpResponsesInLoadedFromWithStatus)
	s.Step(`^we make any "([^"]*)" call to "([^"]*)"\.$`, weMakeAnyCallTo)
	s.Step(`^the parsed response should equal the mock response\.$`, theParsedResponseShouldEqualTheMockResponse)
}
/*
      | indexer_applications_AccountResponse_0.json             | 200    | indexer |  |
      | indexer_applications_AccountsResponse_0.json            | 200    | indexer |  |
      | indexer_applications_ApplicationResponse_0.json         | 200    | indexer |  |
      | indexer_applications_ApplicationsResponse_0.json        | 200    | indexer |  |
      | indexer_applications_AssetBalancesResponse_0.json       | 200    | indexer |  |
      | indexer_applications_AssetResponse_0.json               | 200    | indexer |  |
      | indexer_applications_AssetsResponse_0.json              | 200    | indexer |  |
      | indexer_applications_TransactionsResponse_0.json        | 200    | indexer |  |
      | indexer_applications_TransactionsResponse_0.json        | 200    | indexer |  |
      | indexer_applications_TransactionsResponse_0.json        | 200    | indexer |  |
      | indexer_applications_ErrorResponse_0.json               | 500    | indexer |  |
      | algod_applications_NodeStatusResponse_0.json            | 200    | algod   |  |
      | algod_applications_NodeStatusResponse_0.json            | 200    | algod   |  |
      | algod_applications_CompileResponse_0.json               | 200    | algod   |  |
      | algod_applications_PostTransactionsResponse_0.json      | 200    | algod   |  |
      | algod_applications_SupplyResponse_0.json                | 200    | algod   |  |
      | algod_applications_TransactionParametersResponse_0.json | 200    | algod   |  |
      | algod_applications_ApplicationResponse_0.json           | 200    | algod   |  |
      | algod_applications_AssetResponse_0.json                 | 200    | algod   |  |
*/
