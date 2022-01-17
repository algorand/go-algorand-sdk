package future

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

// abiReturnHash is the 4-byte prefix for logged return values, from https://github.com/algorandfoundation/ARCs/blob/main/ARCs/arc-0004.md#standard-format
var abiReturnHash = []byte{0x15, 0x1f, 0x7c, 0x75}

// maxAppArgs is the maximum number of arguments for an application call transaction at the time
// ARC-4 was created
const maxAppArgs = 16

// The tuple threshold is maxAppArgs, minus 1 for the method selector in the first app arg,
// minus 1 for the final app argument becoming a tuple of the remaining method args
const methodArgsTupleThreshold = maxAppArgs - 2

// TransactionWithSigner represents an unsigned transactions and a signer that can authorize that
// transaction.
type TransactionWithSigner struct {
	// An unsigned transaction
	Txn types.Transaction
	// A transaction signer that can authorize the transaction
	Signer TransactionSigner
}

// Represents the output from a successful ABI method call.
type ABIResult struct {
	// The TxID of the transaction that invoked the ABI method call.
	TxID string
	// The raw bytes of the return value from the ABI method call. This will be empty if the method
	// does not return a value (return type "void").
	RawReturnValue []byte
	// The return value from the ABI method call. This will be nil if the method does not return
	// a value (return type "void"), or if the SDK was unable to decode the returned value.
	ReturnValue interface{}
	// If the SDK was unable to decode a return value, the error will be here. Make sure to check
	// this before examinging ReturnValue
	DecodeError error
}

// AddMethodCallParams contains the parameters for the method AtomicTransactionComposer.AddMethodCall
type AddMethodCallParams struct {
	// The ID of the smart contract to call. Set this to 0 to indicate an application creation call.
	AppID uint64
	// The method to call on the smart contract
	Method abi.Method
	// The arguments to include in the method call. If omitted, no arguments will be passed to the
	// method.
	MethodArgs []interface{}
	// The address of the sender of this application call
	Sender types.Address
	// Transactions params to use for this application call
	SuggestedParams types.SuggestedParams
	// The OnComplete action to take for this application call
	OnComplete types.OnCompletion
	// The approval program for this application call. Only set this if this is an application
	// creation call, or if onComplete is UpdateApplicationOC.
	ApprovalProgram []byte
	// The clear program for this application call. Only set this if this is an application creation
	// call, or if onComplete is UpdateApplicationOC.
	ClearProgram []byte
	// The global schema sizes. Only set this if this is an application creation call.
	GlobalSchema types.StateSchema
	// The local schema sizes. Only set this if this is an application creation call.
	LocalSchema types.StateSchema
	// The number of extra pages to allocate for the application's programs. Only set this if this
	// is an application creation call.
	ExtraPages uint32
	// The note value for this application call
	Note []byte
	// The lease value for this application call
	Lease [32]byte
	// If provided, the address that the sender will be rekeyed to at the conclusion of this application call
	RekeyTo types.Address
	// A transaction Signer that can authorize this application call from sender
	Signer TransactionSigner
}

// AtomicTransactionComposerStatus represents the status of an AtomicTransactionComposer
type AtomicTransactionComposerStatus = int

const (
	// The atomic group is still under construction.
	BUILDING AtomicTransactionComposerStatus = iota

	// The atomic group has been finalized, but not yet signed.
	BUILT

	// The atomic group has been finalized and signed, but not yet submitted to the network.
	SIGNED

	// The atomic group has been finalized, signed, and submitted to the network.
	SUBMITTED

	// The atomic group has been finalized, signed, submitted, and successfully committed to a block.
	COMMITTED
)

type transactionContext struct {
	// The main transaction.
	txn types.Transaction

	// The corresponding signer responsible for producing the signed transaction.
	signer TransactionSigner

	// The corresponding Method constructed from information passed into atc.AddMethodCall().
	method *abi.Method

	// The raw signed transaction populated after invocation of atc.GatherSignatures().
	stxBytes []byte

	// The txid of the transaction, empty until populated by a call to txContext.txID()
	txid string
}

func (txContext *transactionContext) txID() string {
	if txContext.txid == "" {
		txContext.txid = crypto.GetTxID(txContext.txn)
	}
	return txContext.txid
}

func (txContext *transactionContext) isMethodCallTx() bool {
	return txContext.method != nil
}

// The maximum size of an atomic transaction group.
const MaxAtomicGroupSize = 16

// AtomicTransactionComposer is a helper class used to construct and execute atomic transaction groups
type AtomicTransactionComposer struct {
	// The current status of the composer. The status increases monotonically.
	status AtomicTransactionComposerStatus

	// The transaction contexts in the group with their respective signers. If status is greater then
	// BUILDING then this slice cannot change.
	txContexts []transactionContext
}

// GetStatus returns the status of this composer's transaction group.
func (atc *AtomicTransactionComposer) GetStatus() AtomicTransactionComposerStatus {
	return atc.status
}

// Count returns the number of transactions currently in this atomic group.
func (atc *AtomicTransactionComposer) Count() int {
	return len(atc.txContexts)
}

// Clone creates a new composer with the same underlying transactions. The new composer's status
// will be BUILDING, so additional transactions may be added to it.
func (atc *AtomicTransactionComposer) Clone() AtomicTransactionComposer {
	newTxContexts := make([]transactionContext, len(atc.txContexts))
	copy(newTxContexts, atc.txContexts)
	for i := range newTxContexts {
		newTxContexts[i].txn.Group = types.Digest{}
	}

	if len(newTxContexts) == 0 {
		newTxContexts = nil
	}

	return AtomicTransactionComposer{
		status:     BUILDING,
		txContexts: newTxContexts,
	}
}

func (atc *AtomicTransactionComposer) validateTransaction(txn types.Transaction, expectedType string) error {
	emtpyGroup := types.Digest{}
	if txn.Group != emtpyGroup {
		return fmt.Errorf("expected empty group id")
	}

	if expectedType != abi.AnyTransactionType && expectedType != string(txn.Type) {
		return fmt.Errorf("expected transaction with type %s, but got type %s", expectedType, string(txn.Type))
	}

	return nil
}

// AddTransaction adds a transaction to this atomic group.
//
// An error will be thrown if the composer's status is not BUILDING, or if adding this transaction
// causes the current group to exceed MaxAtomicGroupSize.
func (atc *AtomicTransactionComposer) AddTransaction(txnAndSigner TransactionWithSigner) error {
	if atc.status != BUILDING {
		return errors.New("status must be BUILDING in order to add tranactions")
	}

	if atc.Count() == MaxAtomicGroupSize {
		return fmt.Errorf("reached max group size: %d", MaxAtomicGroupSize)
	}

	err := atc.validateTransaction(txnAndSigner.Txn, abi.AnyTransactionType)
	if err != nil {
		return err
	}

	txContext := transactionContext{
		txn:    txnAndSigner.Txn,
		signer: txnAndSigner.Signer,
	}
	atc.txContexts = append(atc.txContexts, txContext)
	return nil
}

// AddMethodCall adds a smart contract method call to this atomic group.
//
// An error will be thrown if the composer's status is not BUILDING, if adding this transaction
// causes the current group to exceed MaxAtomicGroupSize, or if the provided arguments are invalid
// for the given method.
func (atc *AtomicTransactionComposer) AddMethodCall(params AddMethodCallParams) error {
	if atc.status != BUILDING {
		return errors.New("status must be BUILDING in order to add transactions")
	}

	if len(params.MethodArgs) != len(params.Method.Args) {
		return fmt.Errorf("the incorrect number of arguments were provided: %d != %d", len(params.MethodArgs), len(params.Method.Args))
	}

	if atc.Count()+params.Method.GetTxCount() > MaxAtomicGroupSize {
		return fmt.Errorf("reached max group size: %d", MaxAtomicGroupSize)
	}

	if params.AppID == 0 {
		if len(params.ApprovalProgram) == 0 || len(params.ClearProgram) == 0 {
			return fmt.Errorf("ApprovalProgram and ClearProgram must be provided for an application creation call")
		}
	} else if params.OnComplete == types.UpdateApplicationOC {
		if len(params.ApprovalProgram) == 0 || len(params.ClearProgram) == 0 {
			return fmt.Errorf("ApprovalProgram and ClearProgram must be provided for an application update call")
		}
		if (params.GlobalSchema != types.StateSchema{}) || (params.LocalSchema != types.StateSchema{}) {
			return fmt.Errorf("GlobalSchema and LocalSchema must not be provided for an application update call")
		}
	} else if len(params.ApprovalProgram) != 0 || len(params.ClearProgram) != 0 || (params.GlobalSchema != types.StateSchema{}) || (params.LocalSchema != types.StateSchema{}) {
		return fmt.Errorf("ApprovalProgram, ClearProgram, GlobalSchema, and LocalSchema must not be provided for a non-creation call")
	}

	var txsToAdd []TransactionWithSigner
	var basicArgValues []interface{}
	var basicArgTypes []abi.Type
	var refArgValues []interface{}
	var refArgTypes []string
	refArgIndexToBasicArgIndex := make(map[int]int)
	for i, arg := range params.Method.Args {
		argValue := params.MethodArgs[i]

		if arg.IsTransactionArg() {
			txnAndSigner, ok := argValue.(TransactionWithSigner)
			if !ok {
				return fmt.Errorf("invalid arg type, expected transaction")
			}

			err := atc.validateTransaction(txnAndSigner.Txn, arg.Type)
			if err != nil {
				return err
			}
			txsToAdd = append(txsToAdd, txnAndSigner)
		} else {
			var abiType abi.Type
			var err error

			if arg.IsReferenceArg() {
				refArgIndexToBasicArgIndex[len(refArgTypes)] = len(basicArgTypes)
				refArgValues = append(refArgValues, argValue)
				refArgTypes = append(refArgTypes, arg.Type)

				// treat the reference as a uint8 for encoding purposes
				abiType, err = abi.TypeOf("uint8")
			} else {
				abiType, err = arg.GetTypeObject()
			}
			if err != nil {
				return err
			}

			basicArgValues = append(basicArgValues, argValue)
			basicArgTypes = append(basicArgTypes, abiType)
		}
	}

	var foreignAccounts []string
	var foreignApps []uint64
	var foreignAssets []uint64
	refArgsResolved, err := populateMethodCallReferenceArgs(params.Sender.String(), params.AppID, refArgTypes, refArgValues, &foreignAccounts, &foreignApps, &foreignAssets)
	if err != nil {
		return err
	}
	for i, resolved := range refArgsResolved {
		basicArgIndex := refArgIndexToBasicArgIndex[i]
		// use the foreign array index as the encoded argument value
		basicArgValues[basicArgIndex] = resolved
	}

	// Up to 16 app arguments can be passed to app call. First is reserved for method selector,
	// and the rest are for method call arguments. But if more than 15 method call arguments
	// are present, then the method arguments after the 14th are placed in a tuple in the last app
	// argument slot
	if len(basicArgValues) > maxAppArgs-1 {
		typesForTuple := make([]abi.Type, len(basicArgTypes)-methodArgsTupleThreshold)
		copy(typesForTuple, basicArgTypes[methodArgsTupleThreshold:])

		valueForTuple := make([]interface{}, len(basicArgValues)-methodArgsTupleThreshold)
		copy(valueForTuple, basicArgValues[methodArgsTupleThreshold:])

		tupleType, err := abi.MakeTupleType(typesForTuple)
		if err != nil {
			return err
		}

		basicArgValues = append(basicArgValues[:methodArgsTupleThreshold], valueForTuple)
		basicArgTypes = append(basicArgTypes[:methodArgsTupleThreshold], tupleType)
	}

	encodedAbiArgs := [][]byte{params.Method.GetSelector()}

	for i, abiArg := range basicArgValues {
		encodedArg, err := basicArgTypes[i].Encode(abiArg)
		if err != nil {
			return err
		}

		encodedAbiArgs = append(encodedAbiArgs, encodedArg)
	}

	tx, err := MakeApplicationCallTx(
		params.AppID,
		encodedAbiArgs,
		foreignAccounts,
		foreignApps,
		foreignAssets,
		params.OnComplete,
		params.ApprovalProgram,
		params.ClearProgram,
		params.GlobalSchema,
		params.LocalSchema,
		params.SuggestedParams,
		params.Sender,
		params.Note,
		types.Digest{},
		params.Lease,
		params.RekeyTo)
	if err != nil {
		return err
	}

	if params.ExtraPages != 0 {
		tx, err = MakeApplicationCallTxWithExtraPages(tx, params.ExtraPages)
		if err != nil {
			return err
		}
	}

	txAndSigner := TransactionWithSigner{
		Txn:    tx,
		Signer: params.Signer,
	}

	for _, txAndSigner := range txsToAdd {
		txContext := transactionContext{
			txn:    txAndSigner.Txn,
			signer: txAndSigner.Signer,
		}
		atc.txContexts = append(atc.txContexts, txContext)
	}

	methodCallTxContext := transactionContext{
		txn:    txAndSigner.Txn,
		signer: txAndSigner.Signer,
		method: &params.Method,
	}
	atc.txContexts = append(atc.txContexts, methodCallTxContext)
	return nil
}

func (atc *AtomicTransactionComposer) getFinalizedTxWithSigners() []TransactionWithSigner {
	txWithSigners := make([]TransactionWithSigner, len(atc.txContexts))
	for i, txContext := range atc.txContexts {
		txWithSigners[i] = TransactionWithSigner{
			Txn:    txContext.txn,
			Signer: txContext.signer,
		}
	}
	return txWithSigners
}

// BuildGroup finalizes the transaction group and returned the finalized transactions.
//
// The composer's status will be at least BUILT after executing this method.
func (atc *AtomicTransactionComposer) BuildGroup() ([]TransactionWithSigner, error) {
	if atc.status > BUILDING {
		return atc.getFinalizedTxWithSigners(), nil
	}

	if atc.Count() == 0 {
		return nil, fmt.Errorf("attempting to build group with zero transactions")
	}

	var txns []types.Transaction
	for _, txContext := range atc.txContexts {
		txns = append(txns, txContext.txn)
	}

	if len(txns) > 1 {
		gid, err := crypto.ComputeGroupID(txns)
		if err != nil {
			return nil, err
		}

		for i := range atc.txContexts {
			atc.txContexts[i].txn.Group = gid
		}
	}

	atc.status = BUILT
	return atc.getFinalizedTxWithSigners(), nil
}

func (atc *AtomicTransactionComposer) getRawSignedTxs() [][]byte {
	stxs := make([][]byte, len(atc.txContexts))
	for i, txContext := range atc.txContexts {
		stxs[i] = txContext.stxBytes
	}
	return stxs
}

// GatherSignatures obtains signatures for each transaction in this group. If signatures have
// already been obtained, this method will return cached versions of the signatures.
//
// The composer's status will be at least SIGNED after executing this method.
//
// An error will be thrown if signing any of the transactions fails. Otherwise, this will return an
// array of signed transactions.
func (atc *AtomicTransactionComposer) GatherSignatures() ([][]byte, error) {
	// if status is at least signed then return cached signed transactions
	if atc.status >= SIGNED {
		return atc.getRawSignedTxs(), nil
	}

	// retrieve built transactions and verify status is BUILT
	txsWithSigners, err := atc.BuildGroup()
	if err != nil {
		return nil, err
	}

	var txs []types.Transaction
	for _, txWithSigner := range txsWithSigners {
		txs = append(txs, txWithSigner.Txn)
	}

	visited := make([]bool, len(txs))
	rawSignedTxs := make([][]byte, len(txs))
	for i, txWithSigner := range txsWithSigners {
		if visited[i] {
			continue
		}

		var indexesToSign []int
		for j, other := range txsWithSigners {
			if !visited[j] && txWithSigner.Signer.Equals(other.Signer) {
				indexesToSign = append(indexesToSign, j)
				visited[j] = true
			}
		}

		if len(indexesToSign) == 0 {
			return nil, fmt.Errorf("invalid tx signer provided, isn't equal to self")
		}

		sigStxs, err := txWithSigner.Signer.SignTransactions(txs, indexesToSign)
		if err != nil {
			return nil, err
		}

		for i, index := range indexesToSign {
			rawSignedTxs[index] = sigStxs[i]
		}
	}

	for i, stxBytes := range rawSignedTxs {
		atc.txContexts[i].stxBytes = stxBytes
	}
	atc.status = SIGNED
	return rawSignedTxs, nil
}

func (atc *AtomicTransactionComposer) getTxIDs() []string {
	txIDs := make([]string, len(atc.txContexts))
	for i, txContext := range atc.txContexts {
		txIDs[i] = txContext.txID()
	}
	return txIDs
}

// Submit sends the transaction group to the network, but doesn't wait for it to be committed to a
// block. An error will be thrown if submission fails.
//
// The composer's status must be SUBMITTED or lower before calling this method. If submission is
// successful, this composer's status will update to SUBMITTED.
//
// Note: a group can only be submitted again if it fails.
//
// Returns a list of TxIDs of the submitted transactions.
func (atc *AtomicTransactionComposer) Submit(client *algod.Client, ctx context.Context) ([]string, error) {
	if atc.status > SUBMITTED {
		return nil, errors.New("status must be SUBMITTED or lower in order to call Submit()")
	}

	stxs, err := atc.GatherSignatures()
	if err != nil {
		return nil, err
	}

	var serializedStxs []byte
	for _, stx := range stxs {
		serializedStxs = append(serializedStxs, stx...)
	}

	_, err = client.SendRawTransaction(serializedStxs).Do(ctx)
	if err != nil {
		return nil, err
	}

	atc.status = SUBMITTED
	return atc.getTxIDs(), nil
}

// Execute sends the transaction group to the network and waits until it's committed to a block. An
// error will be thrown if submission or execution fails.
//
// The composer's status must be SUBMITTED or lower before calling this method, since execution is
// only allowed once. If submission is successful, this composer's status will update to SUBMITTED.
// If the execution is also successful, this composer's status will update to COMMITTED.
//
// Note: a group can only be submitted again if it fails.
//
// Returns the confirmed round for this transaction, the txIDs of the submitted transactions, and an
// ABIResult for each method call in this group.
func (atc *AtomicTransactionComposer) Execute(client *algod.Client, ctx context.Context, waitRounds uint64) (uint64, []string, []ABIResult, error) {
	if atc.status == COMMITTED {
		return 0, nil, nil, errors.New("status is already committed")
	}

	_, err := atc.Submit(client, ctx)
	if err != nil {
		return 0, nil, nil, err
	}
	atc.status = SUBMITTED

	indexToWaitFor := 0
	for i, txContext := range atc.txContexts {
		if txContext.isMethodCallTx() {
			// if there is a method call in the group, we need to query the
			// pending tranaction endpoint for it anyway, so as an optimization
			// we should wait for its TxID
			indexToWaitFor = i
			break
		}
	}

	txinfo, err := WaitForConfirmation(client, atc.txContexts[indexToWaitFor].txID(), waitRounds, ctx)
	if err != nil {
		return 0, nil, nil, err
	}
	atc.status = COMMITTED

	var returnValues []ABIResult
	for i, txContext := range atc.txContexts {
		// Verify method call is available. This may not be the case if the App Call Tx wasn't created
		// by AddMethodCall().
		if !txContext.isMethodCallTx() {
			continue
		}

		result := ABIResult{TxID: txContext.txID()}

		var methodCallInfo models.PendingTransactionInfoResponse
		if i == indexToWaitFor {
			methodCallInfo = txinfo
		} else {
			methodCallInfo, _, err = client.PendingTransactionInformation(result.TxID).Do(ctx)
			if err != nil {
				result.DecodeError = err
				returnValues = append(returnValues, result)
				continue
			}
		}

		if txContext.method.Returns.IsVoid() {
			result.RawReturnValue = []byte{}
			returnValues = append(returnValues, result)
			continue
		}

		if len(methodCallInfo.Logs) == 0 {
			result.DecodeError = errors.New("method call did not log a return value")
			returnValues = append(returnValues, result)
			continue
		}

		lastLog := methodCallInfo.Logs[len(methodCallInfo.Logs)-1]
		if !bytes.HasPrefix(lastLog, abiReturnHash) {
			result.DecodeError = errors.New("method call did not log a return value")
			returnValues = append(returnValues, result)
			continue
		}

		result.RawReturnValue = lastLog[len(abiReturnHash):]

		abiType, err := txContext.method.Returns.GetTypeObject()
		if err != nil {
			result.DecodeError = err
			returnValues = append(returnValues, result)
			break
		}

		result.ReturnValue, result.DecodeError = abiType.Decode(result.RawReturnValue)
		returnValues = append(returnValues, result)
	}

	return txinfo.ConfirmedRound, atc.getTxIDs(), returnValues, nil
}

// marshallAbiUint64 converts any value used to represent an ABI "uint64" into
// a golang uint64
func marshallAbiUint64(value interface{}) (uint64, error) {
	abiType, err := abi.TypeOf("uint64")
	if err != nil {
		return 0, err
	}
	encoded, err := abiType.Encode(value)
	if err != nil {
		return 0, err
	}
	decoded, err := abiType.Decode(encoded)
	if err != nil {
		return 0, err
	}
	marshalledValue, ok := decoded.(uint64)
	if !ok {
		err = fmt.Errorf("Decoded value is not a uint64")
	}
	return marshalledValue, err
}

// marshallAbiAddress converts any value used to represent an ABI "address" into
// a golang address string
func marshallAbiAddress(value interface{}) (string, error) {
	abiType, err := abi.TypeOf("address")
	if err != nil {
		return "", err
	}
	encoded, err := abiType.Encode(value)
	if err != nil {
		return "", err
	}
	decoded, err := abiType.Decode(encoded)
	if err != nil {
		return "", err
	}
	marshalledValue, ok := decoded.([]byte)
	if !ok || len(marshalledValue) != len(types.ZeroAddress) {
		err = fmt.Errorf("Decoded value is not a 32 length byte slice")
	}
	var addressValue types.Address
	copy(addressValue[:], marshalledValue)
	return addressValue.String(), err
}

// populateMethodCallReferenceArgs parses reference argument types and resolves them to an index
// into the appropriate foreign array. Their placement will be as compact as possible, which means
// values will be deduplicated and any value that is the sender or the current app will not be added
// to the foreign array.
func populateMethodCallReferenceArgs(sender string, currentApp uint64, types []string, values []interface{}, accounts *[]string, apps *[]uint64, assets *[]uint64) ([]int, error) {
	resolvedIndexes := make([]int, len(types))

	for i, value := range values {
		var resolved int

		switch types[i] {
		case abi.AccountReferenceType:
			address, err := marshallAbiAddress(value)
			if err != nil {
				return nil, err
			}
			if address == sender {
				resolved = 0
			} else {
				duplicate := false
				for j, account := range *accounts {
					if address == account {
						resolved = j + 1 // + 1 because 0 is the sender
						duplicate = true
						break
					}
				}
				if !duplicate {
					resolved = len(*accounts) + 1
					*accounts = append(*accounts, address)
				}
			}
		case abi.ApplicationReferenceType:
			appID, err := marshallAbiUint64(value)
			if err != nil {
				return nil, err
			}
			if appID == currentApp {
				resolved = 0
			} else {
				duplicate := false
				for j, app := range *apps {
					if appID == app {
						resolved = j + 1 // + 1 because 0 is the current app
						duplicate = true
						break
					}
				}
				if !duplicate {
					resolved = len(*apps) + 1
					*apps = append(*apps, appID)
				}
			}
		case abi.AssetReferenceType:
			assetID, err := marshallAbiUint64(value)
			if err != nil {
				return nil, err
			}
			duplicate := false
			for j, asset := range *assets {
				if assetID == asset {
					resolved = j
					duplicate = true
					break
				}
			}
			if !duplicate {
				resolved = len(*assets)
				*assets = append(*assets, assetID)
			}
		default:
			return nil, fmt.Errorf("Unknown reference type: %s", types[i])
		}

		resolvedIndexes[i] = resolved
	}

	return resolvedIndexes, nil
}
