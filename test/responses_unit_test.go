package test

import (
	"context"
	baseJson "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/encoding/json"
)

var algodC *algod.Client
var indexerC *indexer.Client
var baselinePath string
var expectedStatus int
var response interface{}

func mockHttpResponsesInLoadedFromWithStatus(
	jsonfile, loadedFrom /* generated_responses*/ string, status int) error {
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
			response = models.LookupAccountByIDResponse{
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
			response = models.LookupAssetByIDResponse{
				CurrentRound: round,
				Asset:        something.(models.Asset),
			}
		case "searchForAssets":
			round, something, err = indexerC.SearchForAssets().Do(context.Background())
			response = models.AssetsResponse{
				Assets:       something.([]models.Asset),
				CurrentRound: round,
			}
		case "lookupAccountTransactions":
			response, err = indexerC.LookupAccountTransactions("").Do(context.Background())
		case "lookupAssetTransactions":
			response, err = indexerC.LookupAssetTransactions(10).Do(context.Background())
		case "searchForTransactions":
			response, err = indexerC.SearchForTransactions().Do(context.Background())
		default:
			err = fmt.Errorf("unknown endpoint: %s", endpoint)
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
	var err error

	responseJson := string(json.Encode(response))

	jsonfile, err := os.Open(baselinePath)
	if err != nil {
		return err
	}
	fileBytes, err := ioutil.ReadAll(jsonfile)
	if err != nil {
		return err
	}
	fmt.Printf("sss_sss_E____%s\n ", responseJson)
	fmt.Printf("sss_sss_baseline____%s\n ", string(fileBytes))

	ans, err := EqualJson(string(fileBytes), responseJson)

	fmt.Printf("sss_sss_F____ %v\n", ans)

	return err
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

func EqualJson(j1, j2 string) (ans bool, err error) {

	var o1 interface{}
	var o2 interface{}
	err = baseJson.Unmarshal([]byte(j1), &o1)
	if err != nil {
		return false, err
	}
	err = baseJson.Unmarshal([]byte(j2), &o2)
	if err != nil {
		return false, err
	}

	//	reflect.
	return reflect.DeepEqual(o1, o2), nil
}
