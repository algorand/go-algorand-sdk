package test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

var algodV2client *algod.Client
var tx types.Transaction
var transientAccount crypto.Account
var applicationId uint64
var txids []string
var confirmedRound uint64
var results []future.ABIResult

func anAlgodVClientConnectedToPortWithToken(v int, host string, port int, token string) error {
	var err error
	portAsString := strconv.Itoa(port)
	fullHost := "http://" + host + ":" + portAsString
	algodV2client, err = algod.MakeClient(fullHost, token)
	gh = []byte("MLWBXKMRJ5W3USARAFOHPQJAF4DN6KY3ZJVPIXKODKNN5ZXSZ2DQ")

	return err
}

func iCreateANewTransientAccountAndFundItWithMicroalgos(microalgos int) error {
	var err error

	transientAccount = crypto.GenerateAccount()

	params, err := algodV2client.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}

	params.Fee = types.MicroAlgos(fee)
	ltxn, err := future.MakePaymentTxn(accounts[1], transientAccount.Address.String(), uint64(microalgos), note, close, params)
	if err != nil {
		return err
	}
	lsk, err := kcl.ExportKey(handle, walletPswd, accounts[1])
	if err != nil {
		return err
	}
	ltxid, lstx, err := crypto.SignTransaction(lsk.PrivateKey, ltxn)
	if err != nil {
		return err
	}
	_, err = algodV2client.SendRawTransaction(lstx).Do(context.Background())
	if err != nil {
		return err
	}
	_, err = future.WaitForConfirmation(algodV2client, ltxid, 5, context.Background())
	if err != nil {
		return err
	}
	return nil
}

func iBuildAnApplicationTransaction(
	operation, approvalProgram, clearProgram string,
	globalBytes, globalInts, localBytes, localInts int,
	appArgs, foreignApps, foreignAssets, appAccounts string, extraPages int) error {

	var clearP []byte
	var approvalP []byte
	var err error

	var ghbytes [32]byte
	copy(ghbytes[:], gh)

	var suggestedParams types.SuggestedParams
	suggestedParams, err = algodV2client.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}

	if approvalProgram != "" {
		approvalP, err = ioutil.ReadFile("features/resources/" + approvalProgram)
		if err != nil {
			return err
		}
	}

	if clearProgram != "" {
		clearP, err = ioutil.ReadFile("features/resources/" + clearProgram)
		if err != nil {
			return err
		}
	}
	args, err := parseAppArgs(appArgs)
	if err != nil {
		return err
	}
	var accs []string
	if appAccounts != "" {
		accs = strings.Split(appAccounts, ",")
	}
	fApp, err := splitUint64(foreignApps)
	if err != nil {
		return err
	}

	fAssets, err := splitUint64(foreignAssets)
	if err != nil {
		return err
	}

	gSchema := types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}
	lSchema := types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)}
	switch operation {
	case "create":
		if extraPages > 0 {
			tx, err = future.MakeApplicationCreateTxWithExtraPages(false, approvalP, clearP,
				gSchema, lSchema, args, accs, fApp, fAssets,
				suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{}, uint32(extraPages))
		} else {
			tx, err = future.MakeApplicationCreateTx(false, approvalP, clearP,
				gSchema, lSchema, args, accs, fApp, fAssets,
				suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		}

		if err != nil {
			return err
		}

	case "create_optin":
		if extraPages > 0 {
			tx, err = future.MakeApplicationCreateTxWithExtraPages(true, approvalP, clearP,
				gSchema, lSchema, args, accs, fApp, fAssets,
				suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{}, uint32(extraPages))
		} else {
			tx, err = future.MakeApplicationCreateTx(true, approvalP, clearP,
				gSchema, lSchema, args, accs, fApp, fAssets,
				suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		}

		if err != nil {
			return err
		}
	case "update":
		tx, err = future.MakeApplicationUpdateTx(applicationId, args, accs, fApp, fAssets,
			approvalP, clearP,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "call":
		tx, err = future.MakeApplicationCallTx(applicationId, args, accs,
			fApp, fAssets, types.NoOpOC, approvalP, clearP, gSchema, lSchema,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "optin":
		tx, err = future.MakeApplicationOptInTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "clear":
		tx, err = future.MakeApplicationClearStateTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "closeout":
		tx, err = future.MakeApplicationCloseOutTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}

	case "delete":
		tx, err = future.MakeApplicationDeleteTx(applicationId, args, accs, fApp, fAssets,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported tx type: %s", operation)
	}

	return nil
}

func iSignAndSubmitTheTransactionSavingTheTxidIfThereIsAnErrorItIs(err string) error {
	var e error
	var lstx []byte

	txid, lstx, e = crypto.SignTransaction(transientAccount.PrivateKey, tx)
	if e != nil {
		return e
	}
	txid, e = algodV2client.SendRawTransaction(lstx).Do(context.Background())
	if e != nil {
		if strings.Contains(e.Error(), err) {
			return nil
		}
	}
	return e
}

func iWaitForTheTransactionToBeConfirmed() error {
	_, err := future.WaitForConfirmation(algodV2client, txid, 5, context.Background())
	if err != nil {
		return err
	}
	return nil
}

func iRememberTheNewApplicationID() error {
	response, _, err := algodV2client.PendingTransactionInformation(txid).Do(context.Background())
	applicationId = response.ApplicationIndex
	appId = applicationId
	return err
}

func parseAppArgs(appArgsString string) (appArgs [][]byte, err error) {
	if appArgsString == "" {
		return make([][]byte, 0), nil
	}
	argsArray := strings.Split(appArgsString, ",")
	resp := make([][]byte, len(argsArray))

	for idx, arg := range argsArray {
		typeArg := strings.Split(arg, ":")
		switch typeArg[0] {
		case "str":
			resp[idx] = []byte(typeArg[1])
		case "int":
			intval, _ := strconv.ParseUint(typeArg[1], 10, 64)

			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.BigEndian, intval)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to bytes", arg)
			}
			resp[idx] = buf.Bytes()
		default:
			return nil, fmt.Errorf("Applications doesn't currently support argument of type %s", typeArg[0])
		}
	}
	return resp, err
}

func splitUint64(uint64s string) ([]uint64, error) {
	if uint64s == "" {
		return make([]uint64, 0), nil
	}
	uintarr := strings.Split(uint64s, ",")
	resp := make([]uint64, len(uintarr))
	var err error
	for idx, val := range uintarr {
		resp[idx], err = strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func theTransientAccountShouldHave(appCreated string, byteSlices, uints int,
	applicationState, key, value string) error {

	acct, err := algodV2client.AccountInformation(transientAccount.Address.String()).Do(context.Background())
	if err != nil {
		return err
	}
	if acct.AppsTotalSchema.NumByteSlice != uint64(byteSlices) {
		return fmt.Errorf("expected NumByteSlices %d, got %d",
			byteSlices, acct.AppsTotalSchema.NumByteSlice)
	}
	if acct.AppsTotalSchema.NumUint != uint64(uints) {
		return fmt.Errorf("expected NumUnit %d, got %d",
			uints, acct.AppsTotalSchema.NumUint)
	}
	created, err := strconv.ParseBool(appCreated)
	if err != nil {
		return err
	}

	// If we don't expect the app to exist, verify that it isn't there and exit.
	if !created {
		for _, ap := range acct.CreatedApps {
			if ap.Id == applicationId {
				return fmt.Errorf("AppId is not expected to be in the account")
			}
		}
		return nil
	}

	appIdFound := false
	for _, ap := range acct.CreatedApps {
		if ap.Id == applicationId {
			appIdFound = true
			break
		}
	}
	if !appIdFound {
		return fmt.Errorf("AppId %d is not found in the account", applicationId)
	}

	// If there is no key to check, we're done.
	if key == "" {
		return nil
	}

	// Verify the key-value is set
	found := false
	var keyValues []models.TealKeyValue

	switch strings.ToLower(applicationState) {
	case "local":
		count := 0
		for _, als := range acct.AppsLocalState {
			if als.Id == applicationId {
				keyValues = als.KeyValue
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("Expected only one matching AppsLocalState, found %d",
				count)
		}
	case "global":
		count := 0
		for _, als := range acct.CreatedApps {
			if als.Id == applicationId {
				keyValues = als.Params.GlobalState
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("Expected only one matching CreatedApps, found %d",
				count)
		}
	default:
		return fmt.Errorf("Unknown application state: %s", applicationState)
	}

	if len(keyValues) == 0 {
		return fmt.Errorf("Expected keyvalues length to be greater than 0")
	}
	for _, kv := range keyValues {
		if kv.Key == key {
			if kv.Value.Type == 1 {
				if kv.Value.Bytes != value {
					return fmt.Errorf("Value mismatch: expected: '%s', got '%s'",
						value, kv.Value.Bytes)
				}
			} else if kv.Value.Type == 0 {
				val, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					return err
				}
				if kv.Value.Uint != val {
					return fmt.Errorf("Value mismatch: expected: '%s', got '%s'",
						value, kv.Value.Bytes)
				}
			}
			found = true
		}
	}
	if !found {
		fmt.Errorf("Could not find key '%s'", key)
	}
	return nil
}

func theUnconfirmedPendingTransactionByIDShouldHaveNoApplyDataFields() error {
	status, _, err := algodV2client.PendingTransactionInformation(txid).Do(context.Background())
	if err != nil {
		return err
	}
	if status.ConfirmedRound == 0 {
		if len(status.GlobalStateDelta) != 0 {
			return fmt.Errorf("unexpected global state delta, there should be none: %v", status.GlobalStateDelta)
		}
		if len(status.LocalStateDelta) != 0 {
			return fmt.Errorf("unexpected local state delta, there should be none: %v", status.LocalStateDelta)
		}
	}
	return nil
}

func getAccountDelta(addr string, data []models.AccountStateDelta) []models.EvalDeltaKeyValue {
	for _, v := range data {
		if v.Address == addr {
			return v.Delta
		}
	}
	return nil
}
func theConfirmedPendingTransactionByIDShouldHaveAStateChangeForToIndexerShouldAlsoConfirmThis(stateLocation, key, newValue string, indexer int) error {
	status, _, err := algodV2client.PendingTransactionInformation(txid).Do(context.Background())
	if err != nil {
		return err
	}

	c1 := make(chan models.Transaction, 1)

	go func() {
		for true {
			indexerResponse, _ := indexerClients[indexer].SearchForTransactions().TXID(txid).Do(context.Background())
			if len(indexerResponse.Transactions) == 1 {
				c1 <- indexerResponse.Transactions[0]
			}
			time.Sleep(time.Second)
		}
	}()

	var indexerTx models.Transaction
	select {
	case res := <-c1:
		indexerTx = res
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for indexer trasaction")
	}

	var algodKeyValues []models.EvalDeltaKeyValue
	var indexerKeyValues []models.EvalDeltaKeyValue

	switch stateLocation {
	case "local":
		addr := indexerTx.Sender
		algodKeyValues = getAccountDelta(addr, status.LocalStateDelta)
		indexerKeyValues = getAccountDelta(addr, indexerTx.LocalStateDelta)
	case "global":
		algodKeyValues = status.GlobalStateDelta
		indexerKeyValues = indexerTx.GlobalStateDelta
	default:
		return fmt.Errorf("unknown location: " + stateLocation)
	}

	// algod
	if len(algodKeyValues) != 1 {
		return fmt.Errorf("expected 1 key value, found: %d", len(algodKeyValues))
	}
	if algodKeyValues[0].Key != key {
		return fmt.Errorf("wrong key in algod: %s != %s", algodKeyValues[0].Key, key)
	}

	// indexer
	if len(indexerKeyValues) != 1 {
		return fmt.Errorf("expected 1 key value, found: %d", len(indexerKeyValues))
	}
	if indexerKeyValues[0].Key != key {
		return fmt.Errorf("wrong key in indexer: %s != %s", indexerKeyValues[0].Key, key)
	}

	if indexerKeyValues[0].Value.Action != algodKeyValues[0].Value.Action {
		return fmt.Errorf("action mismatch between algod and indexer")
	}

	switch algodKeyValues[0].Value.Action {
	case uint64(1):
		// bytes
		if algodKeyValues[0].Value.Bytes != newValue {
			return fmt.Errorf("algod value mismatch: %s != %s", algodKeyValues[0].Value.Bytes, newValue)
		}
		if indexerKeyValues[0].Value.Bytes != newValue {
			return fmt.Errorf("indexer value mismatch: %s != %s", indexerKeyValues[0].Value.Bytes, newValue)
		}
	case uint64(2):
		// int
		newValueInt, err := strconv.ParseUint(newValue, 10, 64)
		if err != nil {
			return fmt.Errorf("problem parsing new int value: %s", newValue)
		}

		if algodKeyValues[0].Value.Uint != newValueInt {
			return fmt.Errorf("algod value mismatch: %d != %s", algodKeyValues[0].Value.Uint, newValue)
		}
		if indexerKeyValues[0].Value.Uint != newValueInt {
			return fmt.Errorf("indexer value mismatch: %d != %s", indexerKeyValues[0].Value.Uint, newValue)
		}
	default:
		return fmt.Errorf("unexpected action: %d", algodKeyValues[0].Value.Action)
	}

	return nil
}

func suggestedParamsAlgodV2() error {
	var err error
	sugParams, err = aclv2.SuggestedParams().Do(context.Background())
	return err
}

func iAddTheCurrentTransactionWithSignerToTheComposer() error {
	err := txComposer.AddTransaction(accountTxAndSigner)
	return err
}

func iCloneTheComposer() error {
	txComposer = txComposer.Clone()
	return nil
}

func iExecuteTheCurrentTransactionGroupWithTheComposer() error {
	var err error
	confirmedRound, txids, results, err = txComposer.Execute(algodV2client, context.Background(), 10)
	return err
}

func theAppShouldHaveReturned(commaSeparatedB64Results string) error {
	b64ExpectedResults := strings.Split(commaSeparatedB64Results, ",")

	if len(b64ExpectedResults) != len(results) {
		return fmt.Errorf("length of expected results doesn't match actual: %d != %d", len(b64ExpectedResults), len(results))
	}

	if len(txComposerMethods) != len(results) {
		return fmt.Errorf("length of composer's methods doesn't match results: %d != %d", len(txComposerMethods), len(results))
	}

	for i, b64ExpectedResult := range b64ExpectedResults {
		expectedResult, err := base64.StdEncoding.DecodeString(b64ExpectedResult)
		if err != nil {
			return err
		}

		actualResult := results[i]
		method := txComposerMethods[i]

		if actualResult.DecodeError != nil {
			return actualResult.DecodeError
		}

		if !bytes.Equal(actualResult.RawReturnValue, expectedResult) {
			return fmt.Errorf("Actual result does not match expected result. Actual: %s\n", base64.RawStdEncoding.EncodeToString(actualResult.RawReturnValue))
		}

		if method.Returns.IsVoid() {
			if len(expectedResult) > 0 {
				return fmt.Errorf("Expected result should be empty")
			}
			continue
		}

		abiReturnType, err := method.Returns.GetTypeObject()
		if err != nil {
			return err
		}

		expectedValue, err := abiReturnType.Decode(expectedResult)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(expectedValue, actualResult.ReturnValue) {
			return fmt.Errorf("the decoded expected value doesn't match the actual result")
		}
	}

	return nil
}

//@applications.verified
func ApplicationsContext(s *godog.Suite) {
	s.Step(`^an algod v(\d+) client connected to "([^"]*)" port (\d+) with token "([^"]*)"$`, anAlgodVClientConnectedToPortWithToken)
	s.Step(`^I create a new transient account and fund it with (\d+) microalgos\.$`, iCreateANewTransientAccountAndFundItWithMicroalgos)
	s.Step(`^I build an application transaction with the transient account, the current application, suggested params, operation "([^"]*)", approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), app-args "([^"]*)", foreign-apps "([^"]*)", foreign-assets "([^"]*)", app-accounts "([^"]*)", extra-pages (\d+)$`, iBuildAnApplicationTransaction)
	s.Step(`^I sign and submit the transaction, saving the txid\. If there is an error it is "([^"]*)"\.$`, iSignAndSubmitTheTransactionSavingTheTxidIfThereIsAnErrorItIs)
	s.Step(`^I wait for the transaction to be confirmed\.$`, iWaitForTheTransactionToBeConfirmed)
	s.Step(`^I remember the new application ID\.$`, iRememberTheNewApplicationID)
	s.Step(`^The transient account should have the created app "([^"]*)" and total schema byte-slices (\d+) and uints (\d+), the application "([^"]*)" state contains key "([^"]*)" with value "([^"]*)"$`,
		theTransientAccountShouldHave)
	s.Step(`^the unconfirmed pending transaction by ID should have no apply data fields\.$`, theUnconfirmedPendingTransactionByIDShouldHaveNoApplyDataFields)
	s.Step(`^the confirmed pending transaction by ID should have a "([^"]*)" state change for "([^"]*)" to "([^"]*)", indexer (\d+) should also confirm this\.$`, theConfirmedPendingTransactionByIDShouldHaveAStateChangeForToIndexerShouldAlsoConfirmThis)
	s.Step(`^suggested transaction parameters from the algod v2 client$`, suggestedParamsAlgodV2)
	s.Step(`^I add the current transaction with signer to the composer\.$`, iAddTheCurrentTransactionWithSignerToTheComposer)
	s.Step(`^I clone the composer\.$`, iCloneTheComposer)
	s.Step(`^I execute the current transaction group with the composer\.$`, iExecuteTheCurrentTransactionGroupWithTheComposer)
	s.Step(`^The app should have returned "([^"]*)"\.$`, theAppShouldHaveReturned)
}
