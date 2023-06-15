package test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	modelsV2 "github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/types"

	"github.com/cucumber/godog"
)

func AlgodClientV2Context(s *godog.Suite) {
	s.Step(`^mock http responses in "([^"]*)" loaded from "([^"]*)"$`, mockHttpResponsesInLoadedFrom)
	s.Step(`^expect error string to contain "([^"]*)"$`, expectErrorStringToContain)
	s.Step(`^we make any Pending Transaction Information call$`, weMakeAnyPendingTransactionInformationCall)
	s.Step(`^the parsed Pending Transaction Information response should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make any Pending Transactions Information call$`, weMakeAnyPendingTransactionsInformationCall)
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
	s.Step(`^expect the request to be "([^"]*)" "([^"]*)"$`, expectTheRequestToBe)
	s.Step(`^we make a Status after Block call with round (\d+)$`, weMakeAStatusAfterBlockCallWithRound)
	s.Step(`^we make an Account Information call against account "([^"]*)"$`, weMakeAnAccountInformationCallAgainstAccount)
	s.Step(`^the parsed Pending Transactions Information response should contain an array of len (\d+) and element number (\d+) should have sender "([^"]*)"$`, theParsedResponseShouldEqualTheMockResponse)
	s.Step(`^we make a Pending Transaction Information against txid "([^"]*)" with format "([^"]*)"$`, weMakeAPendingTransactionInformationAgainstTxidWithFormat)
	s.Step(`^we make a Pending Transaction Information with max (\d+) and format "([^"]*)"$`, weMakeAPendingTransactionInformationWithMaxAndFormat)
	s.Step(`^we make a Pending Transactions By Address call against account "([^"]*)" and max (\d+) and format "([^"]*)"$`, weMakeAPendingTransactionsByAddressCallAgainstAccountAndMaxAndFormat)
	s.Step(`^we make a Get Block call against block number (\d+) with format "([^"]*)"$`, weMakeAGetBlockCallAgainstBlockNumberWithFormat)
	s.Step(`^we make any Dryrun call$`, weMakeAnyDryrunCall)
	s.Step(`^the parsed Dryrun Response should have global delta "([^"]*)" with (\d+)$`, parsedDryrunResponseShouldHave)
	s.Step(`^we make an Account Information call against account "([^"]*)" with exclude "([^"]*)"$`, weMakeAnAccountInformationCallAgainstAccountWithExclude)
	s.Step(`^we make an Account Asset Information call against account "([^"]*)" assetID (\d+)$`, weMakeAnAccountAssetInformationCallAgainstAccountAssetID)
	s.Step(`^we make an Account Application Information call against account "([^"]*)" applicationID (\d+)$`, weMakeAnAccountApplicationInformationCallAgainstAccountApplicationID)
	s.Step(`^we make a GetApplicationBoxByName call for applicationID (\d+) with encoded box name "([^"]*)"$`, weMakeAGetApplicationBoxByNameCall)
	s.Step(`^we make a GetApplicationBoxes call for applicationID (\d+) with max (\d+)$`, weMakeAGetApplicationBoxesCall)
	s.Step(`^we make a GetLightBlockHeaderProof call for round (\d+)$`, weMakeAGetLightBlockHeaderProofCallForRound)
	s.Step(`^we make a GetStateProof call for round (\d+)$`, weMakeAGetStateProofCallForRound)
	s.Step(`^we make a GetTransactionProof call for round (\d+) txid "([^"]*)" and hashtype "([^"]*)"$`, weMakeAGetTransactionProofCallForRoundTxidAndHashtype)
	s.Step(`^we make a Lookup Block Hash call against round (\d+)$`, weMakeALookupBlockHashCallAgainstRound)
	s.Step(`^we make a Ready call$`, weMakeAReadyCall)
	s.Step(`^we make a SetSyncRound call against round (\d+)$`, weMakeASetSyncRoundCallAgainstRound)
	s.Step(`^we make a GetSyncRound call$`, weMakeAGetSyncRoundCall)
	s.Step(`^we make a UnsetSyncRound call$`, weMakeAUnsetSyncRoundCall)
	s.Step(`^we make a SetBlockTimeStampOffset call against offset (\d+)$`, weMakeASetBlockTimeStampOffsetCallAgainstOffset)
	s.Step(`^we make a GetBlockTimeStampOffset call$`, weMakeAGetBlockTimeStampOffsetCall)
	s.Step(`^we make a GetLedgerStateDelta call against round (\d+)$`, weMakeAGetLedgerStateDeltaCallAgainstRound)
	s.Step(`^we make a LedgerStateDeltaForTransactionGroupResponse call for ID "([^"]*)"$`, weMakeALedgerStateDeltaForTransactionGroupResponseCallForID)
	s.Step(`^we make a TransactionGroupLedgerStateDeltaForRoundResponse call for round (\d+)$`, weMakeATransactionGroupLedgerStateDeltaForRoundResponseCallForRound)

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
	return weMakeAnyCallTo("algod", "AccountInformation")
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

func weMakeAPendingTransactionInformationAgainstTxidWithFormat(txid, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = algodClient.PendingTransactionInformation(txid).Do(context.Background())
	return nil
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
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, _, globalErrForExamination = algodClient.PendingTransactionsByAddress(account).Max(uint64(max)).Do(context.Background())
	return nil
}

func weMakeAGetBlockCallAgainstBlockNumberWithFormat(blocknum int, format string) error {
	if format != "msgpack" {
		return fmt.Errorf("this sdk does not support format %s", format)
	}
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	_, globalErrForExamination = algodClient.Block(uint64(blocknum)).Do(context.Background())
	return nil
}

var dryrunResponse modelsV2.DryrunResponse

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

func weMakeAnAccountInformationCallAgainstAccountWithExclude(account, exclude string) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.AccountInformation(account).Exclude(exclude).Do(context.Background())
	return nil
}

func weMakeAnAccountAssetInformationCallAgainstAccountAssetID(account string, assetID int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.AccountAssetInformation(account, uint64(assetID)).Do(context.Background())
	return nil
}

func weMakeAnAccountApplicationInformationCallAgainstAccountApplicationID(account string, appID int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.AccountApplicationInformation(account, uint64(appID)).Do(context.Background())
	return nil
}

func weMakeAGetApplicationBoxByNameCall(appId int, encodedBoxName string) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	decodedBoxNames, err := parseAppArgs(encodedBoxName)
	if err != nil {
		return err
	}
	algodClient.GetApplicationBoxByName(uint64(appId), decodedBoxNames[0]).Do(context.Background())
	return nil
}

func weMakeAGetApplicationBoxesCall(appId int, max int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetApplicationBoxes(uint64(appId)).Max(uint64(max)).Do(context.Background())
	return nil
}

func weMakeAGetLightBlockHeaderProofCallForRound(round int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetLightBlockHeaderProof(uint64(round)).Do(context.Background())
	return nil
}

func weMakeAGetStateProofCallForRound(round int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetStateProof(uint64(round)).Do(context.Background())
	return nil
}

func weMakeAGetTransactionProofCallForRoundTxidAndHashtype(round int, txid, hashtype string) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetTransactionProof(uint64(round), txid).Hashtype(hashtype).Do(context.Background())
	return nil
}

func weMakeALookupBlockHashCallAgainstRound(round int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetBlockHash(uint64(round)).Do(context.Background())
	return nil
}

func weMakeAReadyCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetReady().Do(context.Background())
	return nil
}

func weMakeASetSyncRoundCallAgainstRound(round int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.SetSyncRound(uint64(round)).Do(context.Background())
	return nil
}

func weMakeAGetSyncRoundCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetSyncRound().Do(context.Background())
	return nil
}

func weMakeAUnsetSyncRoundCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.UnsetSyncRound().Do(context.Background())
	return nil
}
func weMakeASetBlockTimeStampOffsetCallAgainstOffset(offset int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.SetBlockTimeStampOffset(uint64(offset)).Do(context.Background())
	return nil
}

func weMakeAGetBlockTimeStampOffsetCall() error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetBlockTimeStampOffset().Do(context.Background())
	return nil
}

func weMakeAGetLedgerStateDeltaCallAgainstRound(round int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetLedgerStateDelta(uint64(round)).Do(context.Background())
	return nil
}

func weMakeALedgerStateDeltaForTransactionGroupResponseCallForID(id string) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetLedgerStateDeltaForTransactionGroup(id).Do(context.Background())
	return nil
}

func weMakeATransactionGroupLedgerStateDeltaForRoundResponseCallForRound(round int) error {
	algodClient, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	algodClient.GetTransactionGroupLedgerStateDeltasForRound(uint64(round)).Do(context.Background())
	return nil
}
