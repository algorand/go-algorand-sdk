package test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	modelsV2 "github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/cucumber/godog"
)

func AlgodClientV2Context(s *godog.Suite) {
	s.Step(`^mock http responses in "([^"]*)" loaded from "([^"]*)"$`, mockHttpResponsesInLoadedFrom)
	s.Step(`^expect error string to contain "([^"]*)"$`, expectErrorStringToContain)
	s.Step(`^we make any Pending Transaction Information call$`, weMakeAnyPendingTransactionInformationCall)
	s.Step(`^the parsed Pending Transaction Information response should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Pending Transactions Information call$`, weMakeAnyPendingTransactionsInformationCall)
	s.Step(`^the parsed Pending Transactions Information response should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Send Raw Transaction call$`, weMakeAnySendRawTransactionCall)
	s.Step(`^the parsed Send Raw Transaction response should have txid "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Pending Transactions By Address call$`, weMakeAnyPendingTransactionsByAddressCall)
	s.Step(`^the parsed Pending Transactions By Address response should contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Node Status call$`, weMakeAnyNodeStatusCall)
	s.Step(`^the parsed Node Status response should have a last round of (\d+)$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Ledger Supply call$`, weMakeAnyLedgerSupplyCall)
	s.Step(`^the parsed Ledger Supply response should have totalMoney (\d+) onlineMoney (\d+) on round (\d+)$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Status After Block call$`, weMakeAnyStatusAfterBlockCall)
	s.Step(`^the parsed Status After Block response should have a last round of (\d+)$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Account Information call$`, weMakeAnyAccountInformationCall)
	s.Step(`^the parsed Account Information response should have address "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Get Block call$`, weMakeAnyGetBlockCall)
	s.Step(`^the parsed Get Block response should have rewards pool "([^"]*)"$`, theParsedGetBlockResponseShouldHaveRewardsPool)
	s.Step(`^we make any Suggested Transaction Parameters call$`, weMakeAnySuggestedTransactionParametersCall)
	s.Step(`^the parsed Suggested Transaction Parameters response should have first round valid of (\d+)$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^expect the path used to be "([^"]*)"$`, expectThePathUsedToBe)
	s.Step(`^we make a Pending Transaction Information against txid "([^"]*)" with max (\d+)$`, weMakeAPendingTransactionInformationAgainstTxidWithMax)
	s.Step(`^we make a Pending Transactions By Address call against account "([^"]*)" and max (\d+)$`, weMakeAPendingTransactionsByAddressCallAgainstAccountAndMax)
	s.Step(`^we make a Status after Block call with round (\d+)$`, weMakeAStatusAfterBlockCallWithRound)
	s.Step(`^we make an Account Information call against account "([^"]*)"$`, weMakeAnAccountInformationCallAgainstAccount)
	s.Step(`^we make a Get Block call against block number (\d+)$`, weMakeAGetBlockCallAgainstBlockNumber)
	s.Step(`^the parsed Pending Transactions Information response should contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make a Pending Transaction Information against txid "([^"]*)" with format "([^"]*)"$`, weMakeAPendingTransactionInformationAgainstTxidWithFormat)
	s.Step(`^we make a Pending Transaction Information with max (\d+) and format "([^"]*)"$`, weMakeAPendingTransactionInformationWithMaxAndFormat)
	s.Step(`^we make a Pending Transactions By Address call against account "([^"]*)" and max (\d+) and format "([^"]*)"$`, weMakeAPendingTransactionsByAddressCallAgainstAccountAndMaxAndFormat)
	s.Step(`^we make a Get Block call against block number (\d+) with format "([^"]*)"$`, weMakeAGetBlockCallAgainstBlockNumberWithFormat)
	s.Step(`^we make any Dryrun call$`, weMakeAnyDryrunCall)
	s.Step(`^the parsed Dryrun Response should have global delta "([^"]*)" with (\d+)$`, parsedDryrunResponseShouldHave)
	s.BeforeScenario(func(interface{}) {
		globalErrForExamination = nil
	})
}

func weMakeAnyPendingTransactionInformationCall() error {
	return weMakeAnyCallTo("algod", "PendingTransactionInformation")
}

func weMakeAnyPendingTransactionsInformationCall() error {
	return weMakeAnyCallTo("algod", "GetPendingTransactions")
}

func weMakeAnySendRawTransactionCall() error {
	return weMakeAnyCallTo("algod", "RawTransaction")
}

func weMakeAnyPendingTransactionsByAddressCall() error {
	return weMakeAnyCallTo("algod", "GetPendingTransactionsByAddress")
}

func weMakeAnyNodeStatusCall() error {
	return weMakeAnyCallTo("algod", "GetStatus")
}

func weMakeAnyLedgerSupplyCall() error {
	return weMakeAnyCallTo("algod", "GetSupply")
}

func weMakeAnyStatusAfterBlockCall() error {
	return weMakeAnyCallTo("algod", "WaitForBlock")
}

func weMakeAnyAccountInformationCall() error {
	return weMakeAnyCallTo("algod", "GetAccountInformation")
}

var blockResponse types.Block

func weMakeAnyGetBlockCall() error {
	// This endpoint requires some sort of base64 decoding for the verification step and it isn't working properly.
	//return weMakeAnyCallTo("algod", "GetBlock")
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	blockResponse, globalErrForExamination = algodClient.Block(0).Do(context.Background())
	return nil
}

func theParsedGetBlockResponseShouldHaveRewardsPool(pool string) error {
	blockResponseRewardsPoolBytes := [32]byte(blockResponse.RewardsPool)
	poolBytes, err := base64.StdEncoding.DecodeString(pool)
	if err != nil {
		return err
	}
	if !bytes.Equal(poolBytes, blockResponseRewardsPoolBytes[:]) {
		return fmt.Errorf("response pool %v mismatched expected pool %v", blockResponseRewardsPoolBytes, poolBytes)
	}
	return nil
}

func weMakeAnySuggestedTransactionParametersCall() error {
	return weMakeAnyCallTo("algod", "TransactionParams")
}

func weMakeAPendingTransactionInformationAgainstTxidWithMax(txid string, max int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = algodClient.PendingTransactionInformation(txid).Do(context.Background())
	return nil
}

func weMakeAPendingTransactionsByAddressCallAgainstAccountAndMax(account string, max int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = algodClient.PendingTransactionsByAddress(account).Max(uint64(max)).Do(context.Background())
	return nil
}

func weMakeAStatusAfterBlockCallWithRound(round int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = algodClient.StatusAfterBlock(uint64(round)).Do(context.Background())
	return nil
}

func weMakeAnAccountInformationCallAgainstAccount(account string) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = algodClient.AccountInformation(account).Do(context.Background())
	return nil
}

func weMakeAGetBlockCallAgainstBlockNumber(blocknum int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = algodClient.Block(uint64(blocknum)).Do(context.Background())
	return nil
}

func weMakeAPendingTransactionInformationAgainstTxidWithFormat(txid, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	return weMakeAPendingTransactionInformationAgainstTxidWithMax(txid, 0)
}

func weMakeAPendingTransactionInformationWithMaxAndFormat(max int, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = algodClient.PendingTransactions().Max(uint64(max)).Do(context.Background())
	return nil
}

func weMakeAPendingTransactionsByAddressCallAgainstAccountAndMaxAndFormat(account string, max int, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	return weMakeAPendingTransactionsByAddressCallAgainstAccountAndMax(account, max)
}

func weMakeAGetBlockCallAgainstBlockNumberWithFormat(blocknum int, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	return weMakeAGetBlockCallAgainstBlockNumber(blocknum)
}

var dryrunResponse models.DryrunResponse

func weMakeAnyDryrunCall() (err error) {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return
	}

	dryrunResponse, err = algodClient.TealDryrun(modelsV2.DryrunRequest{}).Do(context.Background())
	return
}

func parsedDryrunResponseShouldHave(key string, action int) error {
	if len(dryrunResponse.Txns) != 1 {
		return fmt.Errorf("expected 1 txn in result but got %d", len(dryrunResponse.Txns))
	}
	gd := dryrunResponse.Txns[0].GlobalDelta
	if len(gd) != 1 {
		return fmt.Errorf("expected 1 global delta in result but got %d", len(gd))
	}
	if gd[0].Key != key {
		return fmt.Errorf("key %s != %s", key, gd[0].Key)
	}
	if int(gd[0].Value.Action) != action {
		return fmt.Errorf("action %d != %d", action, int(gd[0].Value.Action))
	}
	return nil
}

func mockHttpResponsesInLoadedFrom(jsonfiles, directory string) error {
	return mockHttpResponsesInLoadedFromWithStatus(jsonfiles, directory, 200)
}
