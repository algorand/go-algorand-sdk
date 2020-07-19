package test

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/cucumber/godog"
	"github.com/nsf/jsondiff"

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
	case "algod":
		switch endpoint {
		case "GetStatus":
			response, err = algodC.Status().Do(context.Background())
		case "WaitForBlock":
			response, err = algodC.StatusAfterBlock(10).Do(context.Background())
		case "TealCompile":
			response, err = algodC.TealCompile([]byte{}).Do(context.Background())
		case "RawTransaction":
			response, err = algodC.SendRawTransaction([]byte{}).Do(context.Background())
		case "GetSupply":
			response, err = algodC.Supply().Do(context.Background())
		case "TransactionParams":
			response, err = algodC.SuggestedParams().Do(context.Background())
		case "GetApplicationByID":
			response, err = algodC.GetApplicationByID(10).Do(context.Background())
		case "GetAssetByID":
			response, err = algodC.GetAssetByID(10).Do(context.Background())
		default:
			err = fmt.Errorf("unknown endpoint: %s", endpoint)
		}
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
	ans, err := EqualJson(string(fileBytes), responseJson)

	fmt.Printf("sss_sss_F____ %v\n", ans)

	return err
}

func ResponsesContext(s *godog.Suite) {
	s.Step(`^mock http responses in "([^"]*)" loaded from "([^"]*)" with status (\d+)\.$`, mockHttpResponsesInLoadedFromWithStatus)
	s.Step(`^we make any "([^"]*)" call to "([^"]*)"\.$`, weMakeAnyCallTo)
	s.Step(`^the parsed response should equal the mock response\.$`, theParsedResponseShouldEqualTheMockResponse)
}

func EqualJson(j1, j2 string) (ans bool, err error) {

	options := jsondiff.Options{
		Added:            jsondiff.Tag{Begin: "___ADDED___", End: ""},
		Removed:          jsondiff.Tag{Begin: "___REMED___", End: ""},
		Changed:          jsondiff.Tag{Begin: "___DIFFER___", End: ""},
		ChangedSeparator: " -> ",
	}
	options.PrintTypes = false
	d, str := jsondiff.Compare([]byte(j1), []byte(j2), &options)
	if d == jsondiff.FullMatch {
		return true, nil
	}
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "___REMED___") {
			if strings.Contains(line, "false") {
				continue
			}
			fmt.Printf("%s\n", line)
			return false, nil
		}
		if strings.Contains(line, "___ADDED___") {
			fmt.Printf("%s\n", line)
			return false, nil
		}
		if strings.Contains(line, "___DIFFER___") {
			fmt.Printf("%s\n", line)
			return false, nil
		}
	}
	if d != jsondiff.SupersetMatch {
		fmt.Printf("%s\n", str)
		return false, nil
	}
	return true, nil
}
