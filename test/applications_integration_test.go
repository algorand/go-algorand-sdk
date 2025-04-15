package test

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/algorand/go-algorand-sdk/v2/transaction"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/v2/abi"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	sdkJson "github.com/algorand/go-algorand-sdk/v2/encoding/json"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

var algodV2client *algod.Client
var indexerV2client *indexer.Client
var tx types.Transaction
var transientAccount crypto.Account
var applicationId uint64
var applicationIds []uint64
var txComposerMethodResults []transaction.ABIMethodResult
var lastTxConfirmedRound uint64

func anAlgodVClientConnectedToPortWithToken(v int, host string, port int, token string) error {
	var err error
	portAsString := strconv.Itoa(port)
	fullHost := "http://" + host + ":" + portAsString
	algodV2client, err = algod.MakeClient(fullHost, token)
	aclv2 = algodV2client
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

	funder, err := fundingAccount(algodV2client, uint64(microalgos))
	if err != nil {
		return err
	}

	ltxn, err := transaction.MakePaymentTxn(funder, transientAccount.Address.String(), uint64(microalgos), note, close, params)
	if err != nil {
		return err
	}
	lsk, err := kcl.ExportKey(handle, walletPswd, funder)
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
	_, err = transaction.WaitForConfirmation(algodV2client, ltxid, 1, context.Background())
	if err != nil {
		return err
	}
	return nil
}

func iFundTheCurrentApplicationsAddress(microalgos int) error {
	address := crypto.GetApplicationAddress(applicationId)

	params, err := algodV2client.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}

	funder, err := fundingAccount(algodV2client, uint64(microalgos))
	if err != nil {
		return err
	}
	txn, err := transaction.MakePaymentTxn(funder, address.String(), uint64(microalgos), nil, "", params)
	if err != nil {
		return err
	}

	res, err := kcl.SignTransaction(handle, walletPswd, txn)
	if err != nil {
		return err
	}

	txid, err = algodV2client.SendRawTransaction(res.SignedTransaction).Do(context.Background())
	if err != nil {
		return err
	}

	_, err = transaction.WaitForConfirmation(algodV2client, txid, 1, context.Background())
	return err
}

func readTealProgram(fileName string) ([]byte, error) {
	fileContents, err := loadResource(fileName)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(fileName, ".teal") {
		// need to compile TEAL source file
		response, err := algodV2client.TealCompile(fileContents).Do(context.Background())
		if err != nil {
			return nil, err
		}
		return base64.StdEncoding.DecodeString(response.Result)
	}

	return fileContents, nil
}

func iBuildAnApplicationTransaction(
	operation, approvalProgram, clearProgram string,
	globalBytes, globalInts, localBytes, localInts int,
	appArgs, foreignApps, foreignAssets, appAccounts string,
	extraPages int, boxes string) error {

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
		approvalP, err = readTealProgram(approvalProgram)
		if err != nil {
			return err
		}
	}

	if clearProgram != "" {
		clearP, err = readTealProgram(clearProgram)
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

	staticBoxes, err := parseBoxes(boxes)
	if err != nil {
		return err
	}

	gSchema := types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}
	lSchema := types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)}
	switch operation {
	case "create":
		tx, err = transaction.MakeApplicationCreateTxWithBoxes(false, approvalP, clearP,
			gSchema, lSchema, uint32(extraPages), args, accs, fApp, fAssets, staticBoxes,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "create_optin":
		tx, err = transaction.MakeApplicationCreateTxWithBoxes(true, approvalP, clearP,
			gSchema, lSchema, uint32(extraPages), args, accs, fApp, fAssets, staticBoxes,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "update":
		tx, err = transaction.MakeApplicationUpdateTxWithBoxes(applicationId, args, accs, fApp, fAssets, staticBoxes,
			approvalP, clearP,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "call":
		tx, err = transaction.MakeApplicationNoOpTxWithBoxes(applicationId, args, accs,
			fApp, fAssets, staticBoxes,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "optin":
		tx, err = transaction.MakeApplicationOptInTxWithBoxes(applicationId, args, accs, fApp, fAssets, staticBoxes,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "clear":
		tx, err = transaction.MakeApplicationClearStateTxWithBoxes(applicationId, args, accs, fApp, fAssets, staticBoxes,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "closeout":
		tx, err = transaction.MakeApplicationCloseOutTxWithBoxes(applicationId, args, accs, fApp, fAssets, staticBoxes,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	case "delete":
		tx, err = transaction.MakeApplicationDeleteTxWithBoxes(applicationId, args, accs, fApp, fAssets, staticBoxes,
			suggestedParams, transientAccount.Address, nil, types.Digest{}, [32]byte{}, types.Address{})
	default:
		return fmt.Errorf("unsupported tx type: %s", operation)
	}

	return err
}

func iSignAndSubmitTheTransactionSavingTheTxidIfThereIsAnErrorItIs(expectedErr string) error {
	var err error
	var lstx []byte

	txid, lstx, err = crypto.SignTransaction(transientAccount.PrivateKey, tx)
	if err != nil {
		return err
	}
	txid, err = algodV2client.SendRawTransaction(lstx).Do(context.Background())
	if err != nil && len(expectedErr) != 0 {
		if strings.Contains(err.Error(), expectedErr) {
			return nil
		}
	}
	return err
}

func iWaitForTheTransactionToBeConfirmed() error {
	res, err := transaction.WaitForConfirmation(algodV2client, txid, 1, context.Background())
	if err != nil {
		return err
	}
	lastTxConfirmedRound = res.ConfirmedRound
	return nil
}

func iRememberTheNewApplicationID() error {
	response, _, err := algodV2client.PendingTransactionInformation(txid).Do(context.Background())
	applicationId = response.ApplicationIndex
	applicationIds = append(applicationIds, applicationId)
	return err
}

func iResetTheArrayOfApplicationIDsToRemember() error {
	applicationIds = nil
	return nil
}

func iGetTheAccountAddressForTheCurrentApp() error {
	actual := crypto.GetApplicationAddress(applicationId)

	prefix := []byte("appID")
	preimage := make([]byte, len(prefix)+8) // 8 = number of bytes in a uint64
	copy(preimage, prefix)
	binary.BigEndian.PutUint64(preimage[len(prefix):], applicationId)

	expected := types.Address(sha512.Sum512_256(preimage))

	if expected != actual {
		return fmt.Errorf("Addresses do not match. Expected %s, got %s", expected.String(), actual.String())
	}
	return nil
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
		case "b64":
			d, err := (base64.StdEncoding.DecodeString(typeArg[1]))
			if err != nil {
				return nil, fmt.Errorf("failed to b64 decode arg = %s", arg)
			}
			resp[idx] = d
		default:
			return nil, fmt.Errorf("Applications doesn't currently support argument of type %s", typeArg[0])
		}
	}
	return resp, err
}

func parseBoxes(boxesStr string) (staticBoxes []types.AppBoxReference, err error) {
	if boxesStr == "" {
		return make([]types.AppBoxReference, 0), nil
	}

	boxesArray := strings.Split(boxesStr, ",")

	for i := 0; i < len(boxesArray); i += 2 {
		appID, err := strconv.ParseUint(boxesArray[i], 10, 64)
		if err != nil {
			return nil, err
		}

		var nameBytes []byte
		nameArray := strings.Split(boxesArray[i+1], ":")
		if nameArray[0] == "str" {
			nameBytes = []byte(nameArray[1])
		} else {
			nameBytes, err = base64.StdEncoding.DecodeString(nameArray[1])
			if err != nil {
				return nil, err
			}
		}

		if len(nameBytes) == 0 {
			nameBytes = nil
		}

		staticBoxes = append(staticBoxes,
			types.AppBoxReference{
				AppID: appID,
				Name:  nameBytes,
			})
	}

	return
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
		return fmt.Errorf("Could not find key '%s'", key)
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
	txComposerResult, err := txComposer.Execute(algodV2client, context.Background(), 10)
	if err != nil {
		return err
	}
	txComposerMethodResults = txComposerResult.MethodResults
	return nil
}

func iSimulateTheCurrentTransactionGroupWithTheComposer() error {
	result, err := txComposer.Simulate(context.Background(), algodV2client, models.SimulateRequest{})
	if err != nil {
		return err
	}
	simulateResponse = result.SimulateResponse
	txComposerMethodResults = result.MethodResults
	return nil
}

func theAppShouldHaveReturned(commaSeparatedB64Results string) error {
	b64ExpectedResults := strings.Split(commaSeparatedB64Results, ",")

	if len(b64ExpectedResults) != len(txComposerMethodResults) {
		return fmt.Errorf("length of expected results doesn't match actual: %d != %d", len(b64ExpectedResults), len(txComposerMethodResults))
	}

	for i, b64ExpectedResult := range b64ExpectedResults {
		expectedResult, err := base64.StdEncoding.DecodeString(b64ExpectedResult)
		if err != nil {
			return err
		}

		actualResult := txComposerMethodResults[i]
		method := actualResult.Method

		if actualResult.DecodeError != nil {
			return actualResult.DecodeError
		}

		if !bytes.Equal(actualResult.RawReturnValue, expectedResult) {
			return fmt.Errorf("Actual result does not match expected result. Actual: %s\n", base64.StdEncoding.EncodeToString(actualResult.RawReturnValue))
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

func theAppShouldHaveReturnedABITypes(colonSeparatedExpectedTypeStrings string) error {
	expectedTypeStrings := strings.Split(colonSeparatedExpectedTypeStrings, ":")

	if len(expectedTypeStrings) != len(txComposerMethodResults) {
		return fmt.Errorf("length of expected results doesn't match actual: %d != %d", len(expectedTypeStrings), len(txComposerMethodResults))
	}

	for i, expectedTypeString := range expectedTypeStrings {
		actualResult := txComposerMethodResults[i]

		if actualResult.DecodeError != nil {
			return actualResult.DecodeError
		}

		if expectedTypeString == abi.VoidReturnType {
			if len(actualResult.RawReturnValue) != 0 {
				return fmt.Errorf("No return bytes were expected, but some are present")
			}
			continue
		}

		expectedType, err := abi.TypeOf(expectedTypeString)
		if err != nil {
			return err
		}

		decoded, err := expectedType.Decode(actualResult.RawReturnValue)
		if err != nil {
			return err
		}
		reencoded, err := expectedType.Encode(decoded)
		if err != nil {
			return err
		}

		if !bytes.Equal(reencoded, actualResult.RawReturnValue) {
			return fmt.Errorf("The round trip result does not match the original result")
		}
	}

	return nil
}

func checkAtomicResultAgainstValue(resultIndex int, path, expectedValue string) error {
	keys := strings.Split(path, ".")
	info := txComposerMethodResults[resultIndex].TransactionInfo

	jsonBytes := sdkJson.Encode(&info)

	var genericJson interface{}
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.UseNumber()
	err := decoder.Decode(&genericJson)
	if err != nil {
		return err
	}

	for i, key := range keys {
		var value interface{}

		intKey, err := strconv.Atoi(key)
		if err == nil {
			// key is an array index
			genericArray, ok := genericJson.([]interface{})
			if !ok {
				return fmt.Errorf("Path component %d is an array index (%d), but the parent is not an array. Parent type: %s", i, intKey, reflect.TypeOf(genericJson))
			}
			if intKey < 0 || intKey >= len(genericArray) {
				return fmt.Errorf("Path component %d is an array index (%d) outside of the parent array's range. Parent length: %d", i, intKey, len(genericArray))
			}
			value = genericArray[intKey]
		} else {
			// key is an object field
			genericObject, ok := genericJson.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Path component %d is an object field ('%s'), but the parent is not an object. Parent type: %s", i, key, reflect.TypeOf(genericJson))
			}
			value, ok = genericObject[key]
			if !ok {
				var parentFields []string
				for field := range genericObject {
					parentFields = append(parentFields, "'"+field+"'")
				}
				return fmt.Errorf("Path component %d is an object field ('%s'), but the parent does not contain the field. Parent fields are: %s", i, key, strings.Join(parentFields, ","))
			}
		}

		genericJson = value
	}

	// we have arrived at the target object
	switch actual := genericJson.(type) {
	case string:
		if actual != expectedValue {
			return fmt.Errorf("String values not equal. Expected '%s', got '%s'", expectedValue, actual)
		}
	case json.Number:
		actualParsed, err := strconv.ParseUint(actual.String(), 10, 64)
		if err != nil {
			return err
		}
		expectedParsed, err := strconv.ParseUint(expectedValue, 10, 64)
		if err != nil {
			return err
		}
		if actualParsed != expectedParsed {
			return fmt.Errorf("Int values not equal. Expected %d, got %d", expectedParsed, actualParsed)
		}
	default:
		return fmt.Errorf("The final path element does not resolve to a support type. Type is %s", reflect.TypeOf(genericJson))
	}

	return nil
}

func checkInnerTxnGroupIDs(colonSeparatedPathsString string) error {
	paths := [][]int{}

	commaSeparatedPathStrings := strings.Split(colonSeparatedPathsString, ":")
	for _, commaSeparatedPathString := range commaSeparatedPathStrings {
		pathOfStrings := strings.Split(commaSeparatedPathString, ",")
		path := make([]int, len(pathOfStrings))
		for i, stringComponent := range pathOfStrings {
			intComponent, err := strconv.Atoi(stringComponent)
			if err != nil {
				return err
			}
			path[i] = intComponent
		}
		paths = append(paths, path)
	}

	var txInfosToCheck []models.PendingTransactionResponse

	for _, path := range paths {
		var current models.PendingTransactionResponse
		for pathIndex, innerTxnIndex := range path {
			if pathIndex == 0 {
				current = models.PendingTransactionResponse(txComposerMethodResults[innerTxnIndex].TransactionInfo)
			} else {
				current = current.InnerTxns[innerTxnIndex]
			}
		}
		txInfosToCheck = append(txInfosToCheck, current)
	}

	var group types.Digest
	for i, txInfo := range txInfosToCheck {
		if i == 0 {
			group = txInfo.Transaction.Txn.Group
		}

		if group != txInfo.Transaction.Txn.Group {
			return fmt.Errorf("Group hashes differ: %s != %s", group, txInfo.Transaction.Txn.Group)
		}
	}

	return nil
}

func checkSpinResult(resultIndex int, method, r string) error {
	if method != "spin()" {
		return fmt.Errorf("Incorrect method name, expected 'spin()', got '%s'", method)
	}

	result := txComposerMethodResults[resultIndex]
	decodedResult := result.ReturnValue.([]interface{})

	spin := decodedResult[0].([]interface{})
	spinBytes := make([]byte, len(spin))
	for i, value := range spin {
		spinBytes[i] = value.(byte)
	}

	matched, err := regexp.Match(r, spinBytes)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("Result did not match regex. Spin result: %s", string(spinBytes))
	}

	return nil
}

func sha512_256AsUint64(preimage []byte) uint64 {
	hashed := sha512.Sum512_256(preimage)
	return binary.BigEndian.Uint64(hashed[:8])
}

func checkRandomIntResult(resultIndex, input int) error {
	result := txComposerMethodResults[resultIndex]
	decodedResult := result.ReturnValue.([]interface{})

	randInt := decodedResult[0].(uint64)

	witness := decodedResult[1].([]interface{})
	witnessBytes := make([]byte, len(witness))
	for i, value := range witness {
		witnessBytes[i] = value.(byte)
	}

	x := sha512_256AsUint64(witnessBytes)
	quotient := x % uint64(input)
	if quotient != randInt {
		return fmt.Errorf("Unexpected result: quotient is %d and randInt is %d", quotient, randInt)
	}

	return nil
}

func checkRandomElementResult(resultIndex int, input string) error {
	result := txComposerMethodResults[resultIndex]
	decodedResult := result.ReturnValue.([]interface{})

	randElt := decodedResult[0].(byte)

	witness := decodedResult[1].([]interface{})
	witnessBytes := make([]byte, len(witness))
	for i, value := range witness {
		witnessBytes[i] = value.(byte)
	}

	x := sha512_256AsUint64(witnessBytes)
	quotient := x % uint64(len(input))
	if input[quotient] != randElt {
		return fmt.Errorf("Unexpected result: input[quotient] is %d and randElt is %d", input[quotient], randElt)
	}

	return nil
}

func theContentsOfTheBoxWithNameShouldBeIfThereIsAnErrorItIs(fromClient, encodedBoxName, boxContents, errStr string) error {
	var box models.Box
	decodedBoxNames, err := parseAppArgs(encodedBoxName)
	if err != nil {
		return err
	}
	decodedBoxName := decodedBoxNames[0]
	if fromClient == "algod" {
		box, err = algodV2client.GetApplicationBoxByName(applicationId, decodedBoxName).Do(context.Background())
	} else if fromClient == "indexer" {
		box, err = indexerV2client.LookupApplicationBoxByIDAndName(applicationId, decodedBoxName).Do(context.Background())
	} else {
		err = fmt.Errorf("expecting algod or indexer, got " + fromClient)
	}
	if err != nil {
		// If the expected error string is not empty, check if it is a substring of the actual error string.
		// Note that if the expected error string is empty, then the second condition will always return true.
		if len(errStr) != 0 && strings.Contains(err.Error(), errStr) {
			return nil
		}
		return err
	}
	if len(errStr) != 0 {
		return errors.New("expected an error but none was reported")
	}

	b64Value := base64.StdEncoding.EncodeToString(box.Value)
	if b64Value != boxContents {
		return fmt.Errorf("expected box value %s is not equal to actual box value %s", boxContents, b64Value)
	}

	return nil
}

func currentApplicationShouldHaveFollowingBoxes(fromClient, encodedBoxesRaw string) error {
	var expectedNames [][]byte
	if len(encodedBoxesRaw) > 0 {
		encodedBoxes := strings.Split(encodedBoxesRaw, ":")
		expectedNames = make([][]byte, len(encodedBoxes))
		for i, b := range encodedBoxes {
			expected, err := base64.StdEncoding.DecodeString(b)
			if err != nil {
				return err
			}
			expectedNames[i] = expected
		}
	}
	sort.Slice(expectedNames, func(i, j int) bool {
		return bytes.Compare(expectedNames[i], expectedNames[j]) < 0
	})

	var r models.BoxesResponse
	var err error

	if fromClient == "algod" {
		var boxResponse *models.BoxesResponse
		boxResponse, err = advanceChainAndWaitForBoxesToBeAvailable(len(expectedNames))
		if err == nil {
			r = *boxResponse
		}
	} else if fromClient == "indexer" {
		r, err = indexerV2client.SearchForApplicationBoxes(applicationId).Do(context.Background())
	} else {
		err = fmt.Errorf("expecting algod or indexer, got " + fromClient)
	}
	if err != nil {
		return err
	}

	actualNames := make([][]byte, len(r.Boxes))
	for i, b := range r.Boxes {
		actualNames[i] = b.Name
	}
	sort.Slice(actualNames, func(i, j int) bool {
		return bytes.Compare(actualNames[i], actualNames[j]) < 0
	})

	err = sliceOfBytesEqual(expectedNames, actualNames)
	if err != nil {
		return fmt.Errorf("expected and actual box names do not match: %w", err)
	}

	return nil
}

func advanceChainAndWaitForBoxesToBeAvailable(expectedBoxLength int) (*models.BoxesResponse, error) {
	maxRoundsToAdvance := 50
	for i := 0; i < maxRoundsToAdvance; i++ {
		params, err := algodV2client.SuggestedParams().Do(context.Background())
		if err != nil {
			return nil, err
		}

		txn, err := transaction.MakePaymentTxn(transientAccount.Address.String(), transientAccount.Address.String(),
			uint64(0), nil, "", params)
		if err != nil {
			return nil, err
		}

		_, lstx, err := crypto.SignTransaction(transientAccount.PrivateKey, txn)
		if err != nil {
			return nil, err
		}

		txid, err = algodV2client.SendRawTransaction(lstx).Do(context.Background())
		if err != nil {
			return nil, err
		}

		_, err = transaction.WaitForConfirmation(algodV2client, txid, 1, context.Background())

		// MaxAccountLookback is 4, we start checking for our boxes after the 5th round from create transaction
		if i > 5 {
			r, err := algodV2client.GetApplicationBoxes(applicationId).Do(context.Background())
			if err != nil {
				return nil, err
			}

			if len(r.Boxes) == expectedBoxLength {
				return &r, nil
			}
			// Sleep for 1 second and try again
			time.Sleep(time.Second)
		}
	}

	return nil, errors.New("timeout waiting for boxes to become available")
}

func waitForIndexerToCatchUp() error {
	const maxAttempts = 30

	roundToWaitFor := lastTxConfirmedRound
	indexerRound := uint64(0)
	attempts := 0

	for {
		res, err := indexerV2client.HealthCheck().Do(context.Background())
		if err != nil {
			return err
		}
		indexerRound = res.Round
		if indexerRound >= roundToWaitFor {
			// Success
			break
		}

		// Sleep for 1 second and try again
		time.Sleep(time.Second)
		attempts++

		if attempts > maxAttempts {
			// Failsafe to prevent infinite loop
			return fmt.Errorf("timeout waiting for indexer to catch up to round %d. It is currently on %d", roundToWaitFor, indexerRound)
		}
	}

	return nil
}

func currentApplicationShouldHaveBoxNum(fromClient string, max int, expectedNum int) error {
	var r models.BoxesResponse
	var err error

	if fromClient == "algod" {
		r, err = algodV2client.GetApplicationBoxes(applicationId).Max(uint64(max)).Do(context.Background())
	} else if fromClient == "indexer" {
		r, err = indexerV2client.SearchForApplicationBoxes(applicationId).Limit(uint64(max)).Do(context.Background())
	} else {
		err = fmt.Errorf("expecting algod or indexer, got " + fromClient)
	}
	if err != nil {
		return err
	}

	if len(r.Boxes) != expectedNum {
		return fmt.Errorf("expected and actual box number do not match: %v != %v", expectedNum, len(r.Boxes))
	}
	return nil
}

func indexerSaysCurrentAppShouldHaveTheseBoxes(max int, next string, encodedBoxesRaw string) error {
	r, err := indexerV2client.SearchForApplicationBoxes(applicationId).Limit(uint64(max)).Next(next).Do(context.Background())
	if err != nil {
		return err
	}
	actualNames := make([][]byte, len(r.Boxes))
	for i, b := range r.Boxes {
		actualNames[i] = b.Name
	}
	sort.Slice(actualNames, func(i, j int) bool {
		return bytes.Compare(actualNames[i], actualNames[j]) < 0
	})

	var expectedNames [][]byte
	if len(encodedBoxesRaw) > 0 {
		encodedBoxes := strings.Split(encodedBoxesRaw, ":")
		expectedNames = make([][]byte, len(encodedBoxes))
		for i, b := range encodedBoxes {
			expected, err := base64.StdEncoding.DecodeString(b)
			if err != nil {
				return err
			}
			expectedNames[i] = expected
		}
	}
	sort.Slice(expectedNames, func(i, j int) bool {
		return bytes.Compare(expectedNames[i], expectedNames[j]) < 0
	})

	err = sliceOfBytesEqual(expectedNames, actualNames)
	if err != nil {
		return fmt.Errorf("expected and actual box names do not match: %w", err)
	}

	return nil
}

// @applications.verified and @applications.boxes
func ApplicationsContext(s *godog.ScenarioContext) {
	s.Step(`^an algod v(\d+) client connected to "([^"]*)" port (\d+) with token "([^"]*)"$`, anAlgodVClientConnectedToPortWithToken)
	s.Step(`^I create a new transient account and fund it with (\d+) microalgos\.$`, iCreateANewTransientAccountAndFundItWithMicroalgos)
	s.Step(`^I build an application transaction with the transient account, the current application, suggested params, operation "([^"]*)", approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), app-args "([^"]*)", foreign-apps "([^"]*)", foreign-assets "([^"]*)", app-accounts "([^"]*)", extra-pages (\d+), boxes "([^"]*)"$`, iBuildAnApplicationTransaction)
	s.Step(`^I sign and submit the transaction, saving the txid\. If there is an error it is "([^"]*)"\.$`, iSignAndSubmitTheTransactionSavingTheTxidIfThereIsAnErrorItIs)
	s.Step(`^I wait for the transaction to be confirmed\.$`, iWaitForTheTransactionToBeConfirmed)
	s.Step(`^I remember the new application ID\.$`, iRememberTheNewApplicationID)
	s.Step(`^I reset the array of application IDs to remember\.$`, iResetTheArrayOfApplicationIDsToRemember)
	s.Step(`^I get the account address for the current application and see that it matches the app id\'s hash$`, iGetTheAccountAddressForTheCurrentApp)
	s.Step(`^The transient account should have the created app "([^"]*)" and total schema byte-slices (\d+) and uints (\d+), the application "([^"]*)" state contains key "([^"]*)" with value "([^"]*)"$`,
		theTransientAccountShouldHave)
	s.Step(`^suggested transaction parameters from the algod v2 client$`, suggestedParamsAlgodV2)
	s.Step(`^I add the current transaction with signer to the composer\.$`, iAddTheCurrentTransactionWithSignerToTheComposer)
	s.Step(`^I clone the composer\.$`, iCloneTheComposer)
	s.Step(`^I execute the current transaction group with the composer\.$`, iExecuteTheCurrentTransactionGroupWithTheComposer)
	s.Step(`^I simulate the current transaction group with the composer$`, iSimulateTheCurrentTransactionGroupWithTheComposer)
	s.Step(`^The app should have returned "([^"]*)"\.$`, theAppShouldHaveReturned)
	s.Step(`^The app should have returned ABI types "([^"]*)"\.$`, theAppShouldHaveReturnedABITypes)

	s.Step(`^I fund the current application\'s address with (\d+) microalgos\.$`, iFundTheCurrentApplicationsAddress)

	s.Step(`^I can dig the (\d+)th atomic result with path "([^"]*)" and see the value "([^"]*)"$`, checkAtomicResultAgainstValue)
	s.Step(`^I dig into the paths "([^"]*)" of the resulting atomic transaction tree I see group ids and they are all the same$`, checkInnerTxnGroupIDs)
	s.Step(`^The (\d+)th atomic result for "([^"]*)" satisfies the regex "([^"]*)"$`, checkSpinResult)

	s.Step(`^The (\d+)th atomic result for randomInt\((\d+)\) proves correct$`, checkRandomIntResult)
	s.Step(`^The (\d+)th atomic result for randElement\("([^"]*)"\) proves correct$`, checkRandomElementResult)

	s.Step(`^according to "([^"]*)", the contents of the box with name "([^"]*)" in the current application should be "([^"]*)"\. If there is an error it is "([^"]*)"\.$`, theContentsOfTheBoxWithNameShouldBeIfThereIsAnErrorItIs)
	s.Step(`^according to "([^"]*)", the current application should have the following boxes "([^"]*)"\.$`, currentApplicationShouldHaveFollowingBoxes)
	s.Step(`^according to "([^"]*)", with (\d+) being the parameter that limits results, the current application should have (\d+) boxes\.$`, currentApplicationShouldHaveBoxNum)
	s.Step(`^according to indexer, with (\d+) being the parameter that limits results, and "([^"]*)" being the parameter that sets the next result, the current application should have the following boxes "([^"]*)"\.$`, indexerSaysCurrentAppShouldHaveTheseBoxes)
	s.Step(`^I wait for indexer to catch up to the round where my most recent transaction was confirmed\.$`, waitForIndexerToCatchUp)
}
